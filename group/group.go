package group

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/yenole/sugar/logger"
	"github.com/yenole/sugar/packet"
	"github.com/yenole/sugar/route"
)

type Group struct {
	plot string

	mux  sync.RWMutex
	dict map[string]*Unit

	route  *route.Option
	logger logger.Logger
}

func NewGroup(logger logger.Logger) *Group {
	return &Group{
		dict:   make(map[string]*Unit),
		logger: logger,
	}
}

func (g *Group) Match(r *http.Request) http.Handler {
	if len(g.dict) == 0 {
		return nil
	}

	if v := r.Header.Get("Sugar-Host"); v != "" {
		r.Host = v
	}
	if strings.Contains(g.route.Host, fmt.Sprintf(" %s ", r.Host)) {
		url, _ := url.Parse(g.route.Listen)
		return httputil.NewSingleHostReverseProxy(url)
	}
	return nil
}

func (g *Group) HandlePack(r *packet.Request) interface{} {
	for _, un := range g.dict {
		return un.HandlePack(r)
	}
	return errors.New("not found unit")
}

func (g *Group) HandleRevUnit(un *Unit, do Done) error {
	g.mux.Lock()
	defer g.mux.Unlock()

	if _, ok := g.dict[un.sid]; ok || g.plot != un.plot {
		return errors.New("uni not rev")
	}

	g.dict[un.sid] = un
	io.WriteString(un.conn, "sugar://welcome")

	go un.onRev(g.logger, do, func() {
		g.logger.Infof("sn:%s sid:%s quit", un.Name, un.sid)
		g.mux.Lock()
		defer g.mux.Unlock()
		delete(g.dict, un.sid)
	})
	return nil
}

func (g *Group) HandleRouter(raw []byte, un *Unit) interface{} {
	var route route.Option
	err := json.Unmarshal(raw, &route)
	if err != nil {
		return err
	}
	defer func() {
		if route.Host != "" {
			route.Host = fmt.Sprintf(" %s ", route.Host)
		}
		g.route = &route
	}()

	if g.route != nil {
		return true
	}
	return nil
}
