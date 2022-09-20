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
	"github.com/yenole/sugar/logger"
	"github.com/yenole/sugar/network"
	"github.com/yenole/sugar/network/async"
	"github.com/yenole/sugar/packet"
	"github.com/yenole/sugar/route"
)

type Params map[string]interface{}

type Option struct {
	Plot     uint8
	Protocol string
	Logger   logger.Logger
	Router   *route.Option
}

type Unit struct {
	name   string
	option *Option
	conn   network.Conn

	handler *handler.Handler
	async   *async.Async
	logger  logger.Logger
}

func New(name string, option *Option) *Unit {
	if option.Logger == nil {
		option.Logger = newLogger()
	}

	return &Unit{
		name:    name,
		option:  option,
		handler: &handler.Handler{},

		async:  async.New(),
		logger: option.Logger,
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
			u.logger.Errorf("dialer %v fail:%v", addr[1], err.Error())
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
		u.logger.Errorf("read sugar bytes fail:%v", err.Error())
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
			u.logger.Errorf("regist router fail:%v", err.Error())
			return
		}
		u.logger.Debugf("reigst router result:%v", string(ret))
	}
}

func (u *Unit) onRevProc(cnn network.Conn, r *packet.Request) func() {
	return func() {
		rsp, err := u.handler.Handler(r.Method, r.Params)
		if r.ID == 0 {
			u.logger.Debugf("rev %v params:%v", r.Method, string(r.Params))
			return
		}

		if err != nil {
			cnn.WritePack(packet.NewRsp(r.ID, err))
			u.logger.Debugf("rev %v params:%v err:%v", r.Method, string(r.Params), err.Error())
		} else {
			cnn.WritePack(packet.NewRsp(r.ID, rsp))
			u.logger.Debugf("rev %v params:%v result:%v", r.Method, string(r.Params), rsp)
		}
	}
}
func (u *Unit) onRev(cnn network.Conn) {
	defer u.dailer(u.option.Protocol)

	for {
		req, err := cnn.ReadPack()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}

		u.async.Do(u.onRevProc(cnn, req))

	}
}
