package unit

import (
	"github.com/yenole/sugar/logger"
	"github.com/yenole/sugar/network"
	"github.com/yenole/sugar/network/async"
	"github.com/yenole/sugar/packet"
	"github.com/yenole/sugar/route"
	"github.com/yenole/sugar/unit/handler"
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
