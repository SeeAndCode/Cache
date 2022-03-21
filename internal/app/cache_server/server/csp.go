package server

import (
	"Cache/internal/app/cache_server/cache"
	"Cache/internal/pkg/csp"
	"bufio"
	"io"
	"log"
	"net"
)

type CspServer struct {
	cache.Cache
}

func newCSPServer(c cache.Cache) *CspServer {
	return &CspServer{c}
}

func (s *CspServer) Listen(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Println("csp server is ready on ", addr)
	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		go s.process(c)
	}
}

func (s *CspServer) process(c net.Conn) {
	defer c.Close()
	r := bufio.NewReader(c)

	for {
		req, err := csp.UnserializeRequest(r)
		if err != nil {
			if err != io.EOF {
				log.Println("close connection due to error:", err)
				return
			}
		}
		resp := s.handle(req)
		err = s.sendResp(c, resp)
		if err != nil {
			log.Println("send response failed due to err:", err)
		}
	}
}

// sendResp marshal and send resp to c
func (s *CspServer) sendResp(c net.Conn, resp *csp.Response) error {
	bs, err := csp.SerializeResponse(resp)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = c.Write(bs)
	return err
}

func (s *CspServer) handle(req *csp.Request) *csp.Response {
	id := req.Id
	resp := new(csp.Response)
	// handle request
	switch id {
	case csp.CmdSetReq:
		data, _ := req.Data.(*csp.SetReqData)
		s.handleSet(data, resp)
	case csp.CmdGetReq:
		data, _ := req.Data.(*csp.GetReqData)
		s.handleGet(data, resp)
	case csp.CmdGetRange:
		data, _ := req.Data.(*csp.GetRangeReqData)
		s.handleGetRange(data, resp)
	case csp.CmdGetSet:
		data, _ := req.Data.(*csp.GetSetReqData)
		s.handleGetSet(data, resp)
	case csp.CmdGetBit:
		data, _ := req.Data.(*csp.GetBitReqData)
		s.handleGetBit(data, resp)
	case csp.CmdMGet:
		data, _ := req.Data.(*csp.MGetReqData)
		s.handleMGet(data, resp)
	case csp.CmdSetBit:
		data, _ := req.Data.(*csp.SetBitReqData)
		s.handleSetBit(data, resp)
	}
	return resp
}

func (s *CspServer) handleSet(data *csp.SetReqData, resp *csp.Response) {
	resp.Id = csp.CmdSetResp

	key, value := data.Key, data.Value
	err := s.Set(key, value)
	if err != nil {
		resp.Code = csp.CodeUnknown
		resp.Message = csp.MsgUnknown
		return
	}

	resp.Code = csp.CodeOK
	resp.Message = csp.MsgOK
	resp.Data = &csp.SetRespData{}
}

func (s *CspServer) handleGet(data *csp.GetReqData, resp *csp.Response) {
	resp.Id = csp.CmdGetResp

	key := data.Key
	value, err := s.Get(key)
	if err != nil {
		resp.Code = csp.CodeUnknown
		resp.Message = csp.MsgUnknown
		return
	}

	resp.Code = csp.CodeOK
	resp.Message = csp.MsgOK
	resp.Data = &csp.GetRespData{
		Value: value,
	}
}

func (s *CspServer) handleGetRange(data *csp.GetRangeReqData, resp *csp.Response) {

}

func (s *CspServer) handleGetSet(data *csp.GetSetReqData, resp *csp.Response) {

}

func (s *CspServer) handleGetBit(data *csp.GetBitReqData, resp *csp.Response) {

}

func (s *CspServer) handleMGet(data *csp.MGetReqData, resp *csp.Response) {

}

func (s *CspServer) handleSetBit(data *csp.SetBitReqData, resp *csp.Response) {

}
