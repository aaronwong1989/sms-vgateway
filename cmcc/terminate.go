package cmcc

func NewTerminate() *MessageHeader {
	t := MessageHeader{}
	t.TotalLength = 12
	t.SequenceId = uint32(RequestSeq.NextVal())
	t.CommandId = CMPP_TERMINATE
	return &t
}

func NewTerminateResp(seq uint32) *MessageHeader {
	t := MessageHeader{}
	t.TotalLength = 12
	t.SequenceId = seq
	t.CommandId = CMPP_TERMINATE_RESP
	return &t
}
