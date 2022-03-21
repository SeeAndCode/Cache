package csp

const (
	CmdSetReq  = 00001
	CmdSetResp = 10001

	CmdGetReq  = 00002
	CmdGetResp = 10002

	CmdGetRange = 00003

	CmdGetSet = 00004

	CmdGetBit = 00005

	CmdMGet = 00006

	CmdSetBit = 00007
)

// Request 请求，id表示命令字，data表示具体数据
type Request struct {
	Id   int
	Data interface{}
}

// Response 响应，code表示错误码，message表示基本错误信息
type Response struct {
	Id      int
	Code    int
	Message string
	Data    interface{}
}

// SetReqData Set 00001
type SetReqData struct {
	Key   string
	Value []byte
}

// SetRespData Set 10001
type SetRespData struct {
}

// GetReqData Get 00002
type GetReqData struct {
	Key string
}

// GetRespData Get 10002
type GetRespData struct {
	Value []byte
}

// GetRangeReqData GetRange 00003
type GetRangeReqData struct {
	Key   string
	Start int
	End   int
}

// GetRangeRespData GetRange 00003
type GetRangeRespData struct {
	Value []byte
}

// GetSetReqData GetSet 00004
type GetSetReqData struct {
	Key   string
	Value []byte
}

// GetSetRespData GetSet 10004
type GetSetRespData struct {
	Value []string
}

// GetBitReqData GetBit 00005
type GetBitReqData struct {
	Key    string
	Offset int
}

// GetBitRespData GetBit 10005
type GetBitRespData struct {
	Value string
}

// MGetReqData MGet 00006
type MGetReqData struct {
	Keys []string
}

// MGetRespData MGet 10006
type MGetRespData struct {
	Values [][]byte
}

// SetBitReqData SetBit 00007
type SetBitReqData struct {
	Key    string
	Value  []byte
	Offset int
}

// SetBitRespData SetBit 10007
type SetBitRespData struct {
}
