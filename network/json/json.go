package json

import (
	"encoding/binary"
	"encoding/json"
	"io"
)

type BytesPack struct {
	ID      uint        `json:"id,omitempty"`
	SN      string      `json:"sn,omitempty"`
	Method  string      `json:"method"`
	Pararms interface{} `json:"params"`
}

func (b *BytesPack) Wirte(w io.Writer, prefix ...byte) error {
	byts, err := json.Marshal(b)
	if err != nil {
		return err
	}
	sbyts := make([]byte, 4)
	binary.BigEndian.PutUint32(sbyts, uint32(len(byts)))
	_, err = w.Write(append(prefix, append(sbyts, byts...)...))
	return err
}

type BytesPackRsp struct {
}

func (b *BytesPackRsp) Write(w io.Writer) error {
	return nil
}

type BytesUnPack struct {
	ID     uint            `json:"id"`
	Error  string          `json:"error,omitempty"`
	Result json.RawMessage `json:"result"`
}

func (b *BytesUnPack) UnPack(resp interface{}) error {
	return json.Unmarshal(b.Result, resp)
}

func UnPack(r io.Reader) (interface{}, error) {
	byts := make([]byte, 2048)
	_, err := r.Read(byts[:5])
	if err != nil {
		return nil, err
	}
	size := binary.BigEndian.Uint32(byts[1:5])
	// TODO check byts size
	_, err = r.Read(byts[5 : size+5])
	if err != nil {
		return nil, err
	}
	if byts[0] == 0 {
		return byts[5 : size+5], nil
	}
	var pack BytesUnPack
	return &pack, json.Unmarshal(byts[5:][:size], &pack)
}
