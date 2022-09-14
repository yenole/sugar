package sugar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/yenole/sugar/packet"
)

type Matcher interface {
	Match(*http.Request) http.Handler
}

type Sugar struct {
	mux    sync.RWMutex
	svrs   map[string]*units
	routes []Matcher

	logger Logger
}

func New(logger Logger) *Sugar {
	s := &Sugar{
		logger: logger,
		svrs:   make(map[string]*units),
	}
	return s
}

func (s *Sugar) revUnit(unit *Server) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if list, ok := s.svrs[unit.name]; ok {
		if err := list.Rev(s, unit); err != nil {
			return err
		}
	} else {
		s.svrs[unit.name] = newUnits(s, unit)
	}
	return nil
}

func (s *Sugar) revRoute(units *units) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.routes = append(s.routes, units)
	return nil
}

func (s *Sugar) handleRev(r *packet.Request, unit *Server) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	s.logger.Debugf("%s <==> %s:%s(%s)", unit.name, r.SN, r.Method, string(r.Params))
	if svr, ok := s.svrs[r.SN]; ok {
		req := packet.NewRequest("", r.Method, r.Params)
		if r.ID == 0 {
			svr.WritePack(req, nil)
		} else {
			var resp json.RawMessage
			err := svr.WritePack(req, &resp)
			if err != nil {
				rsp := &packet.Response{
					ID:    r.ID,
					Error: err.Error(),
				}
				unit.conn.WritePack(rsp)
				return
			}
			rsp := packet.NewResponse(r.ID, resp)
			unit.conn.WritePack(rsp)
		}
	} else {
		rsp := &packet.Response{
			ID:    r.ID,
			Error: fmt.Sprintf("%s not found", r.SN),
		}
		unit.conn.WritePack(rsp)
	}
}

func (s *Sugar) Run() {
	s.Listen()

	s.logger.Infof("sugar listen %v\n", *listen)
	http.ListenAndServe(*listen, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.mux.RLock()
		defer s.mux.RUnlock()
		for _, m := range s.routes {
			if h := m.Match(r); h != nil {
				h.ServeHTTP(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}
