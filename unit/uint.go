package unit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/yenole/sugar"
	"github.com/yenole/sugar/handler"
	"github.com/yenole/sugar/network"
	"github.com/yenole/sugar/packet"
	"github.com/yenole/sugar/route"
)

type Params map[string]interface{}

type Option struct {
	Plot     uint8
	Protocol string
	Router   *route.Option
}

type Unit struct {
	name   string
	option *Option
	conn   network.Conn

	handler *handler.Handler
}

func New(name string, option *Option) *Unit {
	return &Unit{
		name:    name,
		option:  option,
		handler: &handler.Handler{},
	}
}

func (u *Unit) Handle(way string, fn handler.HandlerFunc) {
	u.handler.Handle(way, fn)
}

func (u *Unit) Run() error {
	// TODO check protocol
	if u.option.Protocol != "" {
		u.dailer(u.option.Protocol)
	}
	return nil
}

func (u *Unit) Call(sn string, m string, req, resp interface{}) error {
	r := packet.NewRequest(sn, m, req)
	return u.conn.WritePack(r, resp)
}

func (u *Unit) dailer(addr string) {
	addrs := strings.Split(addr, "://")
	switch addrs[0] {
	case "tcp":
		cnn, err := net.Dial("tcp", addrs[1])
		if err != nil {
			fmt.Printf("err: %v\n", err)
			time.AfterFunc(time.Second*5, func() { u.dailer(u.option.Protocol) })
			return
		}
		u.onDailerSugar(network.Wrap(cnn))
	}
}

func (u *Unit) onDailerSugar(cnn network.Conn) {
	hello := fmt.Sprintf("sugar://%s/%s?plot=%d", sugar.Version, u.name, u.option.Plot)
	cnn.Write([]byte(hello))

	byts := make([]byte, 1024)
	n, err := cnn.Read(byts)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	if !bytes.Equal(byts[:n], []byte("sugar://welcome")) {
		cnn.Close()
		return
	}

	go u.onRev(cnn)

	u.conn = cnn
	if u.option.Router != nil {
		var ret json.RawMessage
		err := u.Call("", "route", u.option.Router, &ret)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		fmt.Printf("string(ret): %v\n", string(ret))
	}
}

func (u *Unit) onRev(cnn network.Conn) {
	defer u.dailer(u.option.Protocol)

	for {
		req, err := cnn.Request()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}

		go func(req *packet.Request) {
			ret, err := u.handler.Handler(req.Method, req.Params)
			if req.ID == 0 {
				fmt.Printf("req: %v params:%v\n", req.Method, string(req.Params))
				return
			}

			if err != nil {
				resp := packet.Response{
					ID:    req.ID,
					Error: err.Error(),
				}
				cnn.WritePack(&resp)

				fmt.Printf("req: %v params:%v error:%v\n", req.Method, string(req.Params), err.Error())
				return
			}
			resp := packet.NewResponse(req.ID, ret)
			cnn.WritePack(resp)
			fmt.Printf("req: %v params:%v result:%v\n", req.Method, string(req.Params), ret)
		}(req)

	}
}
