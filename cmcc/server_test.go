package main

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"github.com/stretchr/testify/assert"

	cmcc "sms-vgateway/cmcc/protocol"
)

var pool = goroutine.Default()

func TestClient(t *testing.T) {
	clients := 4
	senders := 1
	receivers := 1
	for i := 0; i < clients; i++ {
		runClient(t, senders, receivers)
	}
	time.Sleep(5 * time.Minute)
}

func runClient(t *testing.T, senders, receivers int) {
	go func(t *testing.T) {
		c, err := net.Dial("tcp", ":9000")
		if err != nil {
			log.Errorf("%v", err)
			return
		}
		defer func(c net.Conn) {
			err := c.Close()
			if err != nil {
				log.Errorf("%v", err)
			}
		}(c)

		sendConnect(t, c)
		time.Sleep(time.Millisecond * 10)

		for i := 0; i < senders; i++ {
			_ = pool.Submit(func() {
				for b := true; b; {
					b = sendMt(t, c)
				}
			})
		}

		for i := 0; i < receivers; i++ {
			_ = pool.Submit(func() {
				for b := true; b; {
					b = readResp(t, c)
				}
			})
		}
		time.Sleep(5 * time.Minute)
	}(t)
}

func sendConnect(t *testing.T, c net.Conn) {
	con := cmcc.NewConnect()
	// con.authenticatorSource = "000000" // 反例测试
	t.Logf(">>>: %s", con)
	i, _ := c.Write(con.Encode())
	assert.True(t, uint32(i) == con.TotalLength)
	resp := make([]byte, 33)
	i, _ = c.Read(resp)
	assert.True(t, i == 33)

	header := &cmcc.MessageHeader{}
	err := header.Decode(resp)
	if err != nil {
		return
	}
	rep := &cmcc.CmppConnectResp{}
	err = rep.Decode(header, resp[cmcc.HEAD_LENGTH:])
	if err != nil {
		return
	}
	t.Logf("<<<: %s", rep)
	assert.True(t, 0 == rep.Status)
}

func sendMt(t *testing.T, c net.Conn) bool {
	mts := cmcc.NewSubmit([]string{"13100001111"}, fmt.Sprintf("hello world! %d", rand.Uint64()))
	mt := mts[0]
	_, err := c.Write(mt.Encode())
	if err != nil {
		log.Errorf("%v", err)
		return false
	}
	t.Logf(">>> %s", mt)
	return true
}

func readResp(t *testing.T, c net.Conn) bool {
	bytes := make([]byte, 12)
	_, err := c.Read(bytes)
	if err != nil {
		log.Errorf("%v", err)
		return false
	}
	header := &cmcc.MessageHeader{}
	_ = header.Decode(bytes)
	l := int(header.TotalLength - 12)
	bytes = make([]byte, l)
	l, err = c.Read(bytes)
	if err != nil {
		log.Errorf("%v", err)
		return false
	}
	if header.CommandId == cmcc.CMPP_SUBMIT_RESP {
		csr := &cmcc.SubmitResp{}
		err := csr.Decode(header, bytes)
		if err != nil {
			log.Errorf("%v", err)
			return false
		} else {
			t.Logf("<<< %s", csr)
		}
	} else if header.CommandId == cmcc.CMPP_ACTIVE_TEST {
		at := cmcc.ActiveTest{MessageHeader: header}
		t.Logf("<<< %s", at)
		ats := at.ToResponse()
		_, err = c.Write(ats.Encode())
		if err != nil {
			log.Errorf("%v", err)
			return false
		} else {
			t.Logf(">>> %s", ats)
		}
	} else {
		t.Logf("<<< %s:%x", header, bytes)
	}
	return true
}
