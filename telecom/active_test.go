package telecom

import (
	"testing"
)

func TestActiveTest(t *testing.T) {
	at := NewActiveTest()
	t.Logf("%T : %s", at, at)

	data := at.Encode()
	t.Logf("%T : %x", data, data)

	h := &MessageHeader{}
	_ = h.Decode(data)
	t.Logf("%T : %s", h, h)

	at2 := &ActiveTest{}
	_ = at2.Decode(h, data)
	t.Logf("%T : %s", at2, at2)

	resp := at.ToResponse(0).(*ActiveTestResp)
	t.Logf("%T : %s", resp, resp)

	data = resp.Encode()
	t.Logf("%T : %x", data, data)
	_ = h.Decode(data)

	resp2 := &ActiveTestResp{}
	_ = resp2.Decode(h, data)
	t.Logf("%T : %s", resp2, resp2)
}
