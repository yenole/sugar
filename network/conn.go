package network

import (
	"context"
	"errors"
	"io"
	"net"
	"time"

	"github.com/yenole/sugar/packet"
)

type Writer interface {
	Write(io.Writer) error
}

type Conn interface {
	net.Conn

	ReadPack() (*packet.Request, error)

	WritePack(Writer, ...interface{}) error
}

func Wrap(cnn net.Conn) *connrsp {
	return &connrsp{Conn: cnn}
}

type connrsp struct {
	net.Conn
	dict map[int]func(*packet.Response)
}

func (c *connrsp) Done(resp *packet.Response) error {
	if fn, ok := c.dict[resp.ID]; ok {
		fn(resp)
	}
	return nil
}
func (c *connrsp) ReadPack() (*packet.Request, error) {
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
		c.Done(&resp)
		goto loop
	}

	var req packet.Request
	return &req, req.Read(c)
}

func (c *connrsp) WritePack(w Writer, resp ...interface{}) error {
	if v, ok := w.(*packet.Request); ok && len(resp) > 0 {
		return c.WriteWithRsp(v, resp[0])
	}
	return w.Write(c)
}

func (c *connrsp) WriteWithRsp(req *packet.Request, resp interface{}) error {
	if c.dict == nil {
		c.dict = make(map[int]func(*packet.Response))
	}

	done := make(chan error)
	defer close(done)

	req.ID = int(time.Now().UnixMilli())
	c.dict[req.ID] = func(rsp *packet.Response) {
		if rsp.Error != "" {
			done <- errors.New(rsp.Error)
			return
		}

		done <- rsp.UnPack(resp)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()

	err := req.Write(c)
	if err != nil {
		delete(c.dict, req.ID)
		return err
	}

	select {
	case err := <-done:
		delete(c.dict, req.ID)
		return err

	case <-ctx.Done():
		delete(c.dict, req.ID)
		return ctx.Err()
	}
}
