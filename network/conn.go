package network

import (
	"io"
	"net"

	"github.com/yenole/sugar/packet"
)

type Writer interface {
	Write(io.Writer) error
}

type Conn interface {
	net.Conn

	Request() (*packet.Request, error)

	WritePack(Writer, ...interface{}) error
}

type conn struct {
	net.Conn

	resp *callback
}

func Wrap(cnn net.Conn) *conn {
	return &conn{Conn: cnn, resp: &callback{}}
}

func (c *conn) Request() (*packet.Request, error) {
loop:
	byts := make([]byte, 1)
	_, err := c.Read(byts)
	if err != nil {
		return nil, err
	}

	if byts[0] == 1 {
		var resp packet.Response
		err := resp.Read(c)
		if err != nil {
			return nil, err
		}
		c.resp.Done(&resp)
		goto loop
	}

	var req packet.Request
	return &req, req.Read(c)
}

func (c *conn) WritePack(w Writer, resp ...interface{}) error {
	if v, ok := w.(*packet.Request); ok && len(resp) > 0 {
		return c.resp.Write(c, v, resp[0])
	}
	return w.Write(c)
}
