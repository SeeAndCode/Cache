package client

import (
	"Cache/internal/pkg/csp"
	"bufio"
	"fmt"
	"net"
)

type cspClient struct {
	net.Conn
	r *bufio.Reader
}

// return a client which use csp to
func newCSPClient(srvAddr string) *cspClient {
	c, err := net.Dial("tcp", srvAddr)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(c)
	return &cspClient{c, r}
}

// API based on csp
// Most function will return error when failed to construct request, or to send request, or to receive response.
// They will return error when response code is not ok.

// Every API function will be like below:
// func API() {
//     construct request
//     send request and receive response
//     check response
//     return
// }

// Set map key to value
func (c *cspClient) Set(key string, value []byte) error {
	data := &csp.SetReqData{
		Key:   key,
		Value: value,
	}
	req := &csp.Request{
		Id:   csp.CmdSetReq,
		Data: data,
	}

	resp, err := c.sendAndReceive(req)
	if err != nil {
		return err
	}

	if resp.Code != csp.CodeOK {
		return fmt.Errorf(resp.Message)
	}

	return nil
}

// Get get the value confirmed by the key
func (c *cspClient) Get(key string) ([]byte, error) {
	data := &csp.GetReqData{Key: key}
	req := &csp.Request{
		Id:   csp.CmdGetReq,
		Data: data,
	}

	resp, err := c.sendAndReceive(req)
	if err != nil {
		return nil, err
	}

	if resp.Code != csp.CodeOK {
		return nil, fmt.Errorf(resp.Message)
	}

	d, _ := resp.Data.(*csp.GetRespData)
	return d.Value, nil
}

func (c *cspClient) sendAndReceive(req *csp.Request) (*csp.Response, error) {
	// send request
	bs, err := csp.SerializeRequest(req)
	if err != nil {
		return nil, err
	}
	_, err = c.Write(bs)
	if err != nil {
		return nil, err
	}
	// receive response
	resp, err := csp.UnserializeResponse(c.r)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
