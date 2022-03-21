package csp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

/*
下面是使用ABNF描述的协议的请求部分

request = Data-command; 第一种请求是数据操作命令

bytes-array = length SP content
Key = bytes-array; 字符串格式的键
value = bytes-array
Key-Value = length SP length SP content content; 字符串格式的键和值
length = 1 * DIGIT; 长度，至少一位数值
content = *OCTET; 具体内容

Data-command = set; 数据操作命令Set
set = "00001" SP Key-Value; set Value with Key

Data-command =/ get; 数据操作命令Get
get = "00002" SP Key; get Value of Key

Data-command =/ getrange; 数据操作命令GetRange
getrange = "00003" SP Start SP End SP Key
Start = 1 * DIGIT
End = 1 * DIGIT

Data-command =/ getset; 数据操作命令GetSet
getset = "0004" SP Key-Value

Data-command =/ getbit; 数据操作命令GetBit
getbit = "00005" SP Offset SP Key
Offset = 1 * DIGIT

Data-command =/ mget; 数据操作命令MGet
mget = "00006" SP num SP 1 * Key

Data-command =/ setbit; 数据操作命令SetBit
setbit = "00007" SP Offset SP Key-Value
*/

// UnserializeRequest read and unserialize request from the buffer
// r should be created by bufio.newReader(net.conn)
//
func UnserializeRequest(r *bufio.Reader) (*Request, error) {
	// 读取命令id
	id, err := readInt(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read Id, %s", err.Error())
	}
	// 构造请求
	req := new(Request)

	switch id {
	case CmdSetReq:
		req.Data, err = readSetReqData(r)
	case CmdGetReq:
		req.Data, err = readGetReqData(r)
	case CmdGetRange:
		req.Data, err = readGetRangeReqData(r)
	case CmdGetSet:
		req.Data, err = readGetSetReqData(r)
	case CmdGetBit:
		req.Data, err = readGetBitReqData(r)
	case CmdMGet:
		req.Data, err = readMGetReqData(r)
	case CmdSetBit:
		req.Data, err = readSetBitReqData(r)
	default:
		return nil, fmt.Errorf("unknown cmd Id %d", id)
	}
	if err != nil {
		return nil, err
	}

	req.Id = id
	return req, nil
}

func SerializeRequest(req *Request) ([]byte, error) {
	id := req.Id
	bf := []byte(fmt.Sprintf("%d ", id))

	switch id {
	case CmdSetReq:
		if data, ok := req.Data.(*SetReqData); !ok {
			return nil, errors.New("invalid request")
		} else {
			toBytesSetReqData(data, &bf)
		}
	case CmdGetReq:
		if data, ok := req.Data.(*GetReqData); !ok {
			return nil, errors.New("invalid request")
		} else {
			toBytesGetReqData(data, &bf)
		}
	case CmdGetRange:
		if data, ok := req.Data.(*GetRangeReqData); !ok {
			return nil, errors.New("invalid request")
		} else {
			writeGetRangeReqData(data, &bf)
		}
	case CmdGetSet:
		if data, ok := req.Data.(*GetSetReqData); !ok {
			return nil, errors.New("invalid request")
		} else {
			writeGetSetReqData(data, &bf)
		}
	case CmdGetBit:
		if data, ok := req.Data.(*GetBitReqData); !ok {
			return nil, errors.New("invalid request")
		} else {
			writeGetBitReqData(data, &bf)
		}
	case CmdMGet:
		if data, ok := req.Data.(*MGetReqData); !ok {
			return nil, errors.New("invalid request")
		} else {
			writeMGetReqData(data, &bf)
		}
	case CmdSetBit:
		if data, ok := req.Data.(*SetBitReqData); !ok {
			return nil, errors.New("invalid request")
		} else {
			writeSetBitReqData(data, &bf)
		}
	default:
		return nil, fmt.Errorf("unknown cmd Id %d", id)
	}
	return bf, nil
}

// "00001" SP Key-Value
func readSetReqData(r *bufio.Reader) (*SetReqData, error) {
	key, value, err := readKeyValue(r)
	if err != nil {
		return nil, err
	}
	return &SetReqData{key, value}, nil
}

func toBytesSetReqData(data *SetReqData, bf *[]byte) {
	dataString := fmt.Sprintf("%d %d %s%s", len(data.Key), len(data.Value), data.Key, data.Value)
	*bf = append(*bf, []byte(dataString)...)
}

// "00002" SP Key
func readGetReqData(r *bufio.Reader) (*GetReqData, error) {
	key, err := readKey(r)
	if err != nil {
		return nil, err
	}
	return &GetReqData{key}, nil
}

func toBytesGetReqData(data *GetReqData, bf *[]byte) {
	dataString := fmt.Sprintf("%d %s", len(data.Key), data.Key)
	*bf = append(*bf, []byte(dataString)...)
}

// "00003" SP Start SP End SP Key
func readGetRangeReqData(r *bufio.Reader) (*GetRangeReqData, error) {
	start, err := readInt(r)
	if err != nil {
		return nil, err
	}
	end, err := readInt(r)
	if err != nil {
		return nil, err
	}
	key, err := readKey(r)
	if err != nil {
		return nil, err
	}
	return &GetRangeReqData{
		Key:   key,
		Start: start,
		End:   end,
	}, nil
}

func writeGetRangeReqData(data *GetRangeReqData, bf *[]byte) {
	dataString := fmt.Sprintf("%d %d %d %s", data.Start, data.End, len(data.Key), data.Key)
	*bf = append(*bf, []byte(dataString)...)
}

// "0004" SP Key-Value
func readGetSetReqData(r *bufio.Reader) (*GetSetReqData, error) {
	key, value, err := readKeyValue(r)
	if err != nil {
		return nil, err
	}
	return &GetSetReqData{key, value}, nil
}

func writeGetSetReqData(data *GetSetReqData, bf *[]byte) {
	dataString := fmt.Sprintf("%d %d %s %s", len(data.Key), len(data.Value), data.Key, data.Value)
	*bf = append(*bf, []byte(dataString)...)
}

// "00005" SP Offset SP Key
func readGetBitReqData(r *bufio.Reader) (*GetBitReqData, error) {
	offset, err := readInt(r)
	if err != nil {
		return nil, err
	}
	key, err := readKey(r)
	if err != nil {
		return nil, err
	}
	return &GetBitReqData{
		Key:    key,
		Offset: offset,
	}, nil
}

func writeGetBitReqData(data *GetBitReqData, bf *[]byte) {
	dataString := fmt.Sprintf("%d %d %s", data.Offset, len(data.Key), data.Key)
	*bf = append(*bf, []byte(dataString)...)
}

// "00006" SP num SP 1 * Key
func readMGetReqData(r *bufio.Reader) (*MGetReqData, error) {
	num, err := readInt(r)
	if err != nil {
		return nil, err
	}
	if num == 0 {
		return nil, errors.New("no Key to read")
	}
	var keys []string
	for i := 0; i < num; i++ {
		key, err := readKey(r)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return &MGetReqData{keys}, nil
}

func writeMGetReqData(data *MGetReqData, bf *[]byte) {
	*bf = append(*bf, []byte(fmt.Sprintf("%d ", len(data.Keys)))...)
	for _, key := range data.Keys {
		*bf = append(*bf, []byte(fmt.Sprintf("%d %s", len(key), key))...)
	}
}

// "00007" SP Offset SP Key-Value
func readSetBitReqData(r *bufio.Reader) (*SetBitReqData, error) {
	offset, err := readInt(r)
	if err != nil {
		return nil, err
	}
	key, value, err := readKeyValue(r)
	if err != nil {
		return nil, err
	}
	return &SetBitReqData{
		Key:    key,
		Value:  value,
		Offset: offset,
	}, nil
}

func writeSetBitReqData(data *SetBitReqData, bf *[]byte) {
	dataString := fmt.Sprintf("%d %d %d %s %s", data.Offset, len(data.Key), len(data.Value), data.Key, data.Key)
	*bf = append(*bf, []byte(dataString)...)
}

/*
下面是使用ABNF描述的协议的响应部分

response = data-command-resp

code = 1 * DIGIT
message = bytes-array

data-command-resp = set-resp
set-resp = "10001" SP code SP message

data-command-resp =/ get-resp
get-resp = "10002" SP code SP message value
*/

func UnserializeResponse(r *bufio.Reader) (*Response, error) {
	id, err := readInt(r)
	if err != nil {
		return nil, err
	}

	code, err := readInt(r)
	if err != nil {
		return nil, err
	}

	message, err := readBytesArray(r)
	if err != nil {
		return nil, err
	}
	resp := new(Response)

	switch id {
	case CmdSetResp:
		resp.Data, err = readSetRespData(r)
	case CmdGetResp:
		resp.Data, err = readGetRespData(r)
	default:
		return nil, errors.New("unknown cmd id")
	}
	if err != nil {
		return nil, err
	}
	resp.Id = id
	resp.Code = code
	resp.Message = string(message)
	return resp, nil
}

func SerializeResponse(resp *Response) ([]byte, error) {
	id := resp.Id
	bf := []byte(fmt.Sprintf("%d %d %d %s", id, resp.Code, len(resp.Message), resp.Message))
	switch id {
	case CmdSetResp:
		if data, ok := resp.Data.(*SetRespData); !ok {
			return nil, errors.New("invalid response")
		} else {
			toBytesSetRespData(data, &bf)
		}
	case CmdGetResp:
		if data, ok := resp.Data.(*GetRespData); !ok {
			return nil, errors.New("invalid response")
		} else {
			toBytesGetRespData(data, &bf)
		}
	}
	return bf, nil
}

// "10001" SP code SP message SP
func readSetRespData(r *bufio.Reader) (*SetRespData, error) {
	return &SetRespData{}, nil
}

func toBytesSetRespData(data *SetRespData, bf *[]byte) {
	return
}

// "10002" SP code SP message SP value
func readGetRespData(r *bufio.Reader) (*GetRespData, error) {
	value, err := readBytesArray(r)
	if err != nil {
		return nil, err
	}
	return &GetRespData{Value: value}, nil
}

func toBytesGetRespData(data *GetRespData, bf *[]byte) {
	dataString := fmt.Sprintf("%d %s", len(data.Value), data.Value)
	*bf = append(*bf, []byte(dataString)...)
}

// utils

// readInt 从流中读取下一个SP前的一个数字
// 1 * DIGIT
func readInt(r *bufio.Reader) (int, error) {
	tmp, err := r.ReadString(' ')
	if err != nil {
		return 0, err
	}
	if len(tmp) == 1 {
		return 0, errors.New("no integer to read")
	}
	value, err := strconv.Atoi(strings.TrimSpace(tmp))
	if err != nil {
		return 0, err
	}
	return value, nil
}

// readKey 从流中读取Key
// Key = bytes-array
func readKey(r *bufio.Reader) (string, error) {
	bs, err := readBytesArray(r)
	return string(bs), err
}

// readKeyValue 从流中读取key-Value
// Key-Value = length SP length SP content content
func readKeyValue(r *bufio.Reader) (string, []byte, error) {
	kLen, err := readInt(r)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read Key length, %s", err.Error())
	}
	vLen, err := readInt(r)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read Value length, %s", err.Error())
	}
	key, value := make([]byte, kLen), make([]byte, vLen)
	_, err = io.ReadFull(r, key)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read Key, %s", err.Error())
	}
	_, err = io.ReadFull(r, value)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read Value, %s", err.Error())
	}
	return string(key), value, nil
}

// readBytesArray 从流中读取bytes-array
// bytes-array = length SP content
func readBytesArray(r *bufio.Reader) ([]byte, error) {
	bsLen, err := readInt(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read bytes length, %s", err.Error())
	}
	bs := make([]byte, bsLen)
	_, err = io.ReadFull(r, bs)
	if err != nil {
		return nil, fmt.Errorf("failed to read bytes, %s", err.Error())
	}
	return bs, nil
}
