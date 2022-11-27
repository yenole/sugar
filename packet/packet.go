package packet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
)

type Request struct {
	ID     int             `json:"id"`
	SN     string          `json:"sn"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

func NewRequest(sn, m string, req interface{}) *Request {
	byts, _ := json.Marshal(req)
	return &Request{
		SN:     sn,
		Method: m,
		Params: json.RawMessage(byts),
	}
}

func (r *Request) Read(rr io.Reader) error {
	var buffer bytes.Buffer
	_, err := io.CopyN(&buffer, rr, 4)
	if err != nil {
		return err
	}
	size := binary.BigEndian.Uint32(buffer.Bytes()[:4])
	_, err = io.CopyN(&buffer, rr, int64(size))
	if err != nil {
		return err
	}
	return json.Unmarshal(buffer.Bytes()[4:], r)
}

func (r *Request) Write(w io.Writer) error {
	byts, err := json.Marshal(r)
	if err != nil {
		return err
	}

	sbyts := make([]byte, 5)
	binary.BigEndian.PutUint32(sbyts[1:], uint32(len(byts)))
	_, err = w.Write(append(sbyts, byts...))
	return err
}

type Response struct {
	ID     int             `json:"id"`
	Error  string          `json:"error,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}

func NewRsp(id int, v interface{}) *Response {
	if err, ok := v.(error); ok {
		return &Response{ID: id, Error: err.Error()}
	}
	byts, _ := json.Marshal(v)
	return &Response{ID: id, Result: byts}
}

func NewResponse(id int, v interface{}) *Response {
	byts, _ := json.Marshal(v)
	return &Response{ID: id, Result: byts}
}

func (r *Response) Read(rr io.Reader) error {
	var buffer bytes.Buffer
	_, err := io.CopyN(&buffer, rr, 4)
	if err != nil {
		return err
	}
	size := binary.BigEndian.Uint32(buffer.Bytes())
	_, err = io.CopyN(&buffer, rr, int64(size))
	if err != nil {
		return err
	}
	return json.Unmarshal(buffer.Bytes()[4:], r)
}

func (r *Response) Write(w io.Writer) error {
	byts, err := json.Marshal(r)
	if err != nil {
		return err
	}

	sbyts := make([]byte, 5)
	sbyts[0] = 1
	binary.BigEndian.PutUint32(sbyts[1:], uint32(len(byts)))
	_, err = w.Write(append(sbyts, byts...))
	return err
}

func (r *Response) UnPack(resp interface{}) error {
	return json.Unmarshal(r.Result, resp)
}
