package group

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/yenole/sugar/logger"
	"github.com/yenole/sugar/network"
	"github.com/yenole/sugar/packet"
)

const (
	ACT_ROUTER uint8 = iota
)

type Params map[string]interface{}

type Done interface {
	Do(r *packet.Request)
}

type Unit struct {
	Name string

	sid  string
	plot string
	conn network.Conn
}

func NewUnit(sn string, cnn network.Conn, o string) (*Unit, error) {
	url, err := url.Parse(o)
	if err != nil || url.Path == "/" || url.Path == "" {
		return nil, errors.New("unit scheme fail")
	}

	return &Unit{
		sid:  sn,
		plot: url.Query().Get("plot"),
		Name: strings.TrimPrefix(url.Path, "/"),
		conn: cnn,
	}, nil
}

func (u *Unit) HandlePack(r *packet.Request) interface{} {
	req := packet.NewRequest("", r.Method, r.Params)
	if r.ID == 0 {
		return u.conn.WritePack(req)
	}

	var rsp json.RawMessage
	err := u.conn.WritePack(req, &rsp)
	if err != nil {
		return err
	}
	return &rsp
}

func (u *Unit) WritePack(w network.Writer) error {
	return u.conn.WritePack(w)
}

func (u *Unit) onRev(logger logger.Logger, do Done, close func()) {
	defer close()

	logger.Infof("sn:%s sid:%s join suagr", u.Name, u.sid)
	for {
		req, err := u.conn.ReadPack()
		if err != nil {
			logger.Errorf("sn:%v sid:%s read pack err:%v", u.Name, u.sid, err.Error())
			return
		}

		do.Do(req)
	}
}
