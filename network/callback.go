package network

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/yenole/sugar/packet"
)

type callback struct {
	dict map[int]func(*packet.Response)
}

func (c *callback) Done(resp *packet.Response) error {
	if fn, ok := c.dict[resp.ID]; ok {
		fn(resp)
	}
	return nil
}

func (c *callback) Write(w io.Writer, req *packet.Request, resp interface{}) error {
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

	err := req.Write(w)
	if err != nil {
		return err
	}

	select {
	case err := <-done:
		return err

	case <-ctx.Done():
		delete(c.dict, req.ID)
		return ctx.Err()
	}
}
