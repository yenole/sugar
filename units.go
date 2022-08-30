package sugar

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/yenole/sugar/packet"
	"github.com/yenole/sugar/route"
)

type units struct {
	plot string

	mux  sync.RWMutex
	dict map[string]*Server

	route *route.Option
}

func newUnits(sg *Sugar, unit *Server) *units {
	s := &units{
		plot: unit.plot,
		dict: make(map[string]*Server),
	}
	s.Rev(sg, unit)
	return s
}

func (u *units) Match(r *http.Request) http.Handler {
	if len(u.dict) == 0 {
		return nil
	}

	fmt.Printf("r.URL.Host: %v\n", r.Host)
	if u.route.Host == r.Host {
		url, _ := url.Parse(u.route.Listen)
		return httputil.NewSingleHostReverseProxy(url)
	}
	return nil
}

func (u *units) WritePack(r *packet.Request, resp interface{}) error {
	for _, s := range u.dict {
		if resp == nil {
			return s.conn.WritePack(r)
		} else {

			return s.conn.WritePack(r, resp)
		}
	}
	return errors.New("not unit")
}

func (u *units) Rev(sg *Sugar, unit *Server) error {
	u.mux.Lock()
	defer u.mux.Unlock()

	if _, ok := u.dict[unit.sid]; ok || u.plot != unit.plot {
		return errors.New("uni not rev")
	}

	u.dict[unit.sid] = unit
	io.WriteString(unit.conn, "sugar://welcome")

	go unit.onReceive(func(r *packet.Request, s *Server) {
		u.handleRev(r, sg, unit)
	}, func() {
		fmt.Printf("name %s sid %s out\n", unit.name, unit.sid)
		u.mux.Lock()
		defer u.mux.Unlock()
		delete(u.dict, unit.sid)
	})
	return nil
}

func (u *units) handleRev(r *packet.Request, sg *Sugar, unit *Server) {
	fmt.Printf("rev sn:%v id:%v method:%v params:%v\n", unit.name, r.ID, r.Method, string(r.Params))

	if r.SN != "" {
		sg.handleRev(r, unit)
		return
	}

	var ret interface{}
	switch r.Method {
	case "route":
		ret = u.handleRouter(r.Params, sg)

	}

	if r.ID == 0 {
		return
	}

	if err, ok := ret.(error); ok {
		rsp := &packet.Response{
			ID:    r.ID,
			Error: err.Error(),
		}
		unit.conn.WritePack(rsp)
		return
	}

	rsp := packet.NewResponse(r.ID, ret)
	unit.conn.WritePack(rsp)
}

func (u *units) handleRouter(raw []byte, sg *Sugar) interface{} {
	var route route.Option
	err := json.Unmarshal(raw, &route)
	if err != nil {
		return err
	}
	defer func() {
		u.route = &route
	}()

	if u.route != nil {
		return true
	}
	return sg.revRoute(u)
}
