// Package sugar provides ...
package sugar

import (
	"encoding/json"
	"fmt"

	"github.com/yenole/sugar/group"
	"github.com/yenole/sugar/packet"
)

type done struct {
	s  *Sugar
	un *group.Unit
}

func (d *done) Do(r *packet.Request) {
	d.s.g.Do(d.s.handleRevRequest(r, d.un))
}

func (s *Sugar) done(un *group.Unit) *done {
	return &done{s: s, un: un}
}
func (s *Sugar) handleRevRequest(r *packet.Request, un *group.Unit) func() {
	return func() {
		var rsp interface{}
		if g := s.group(r.SN); g != nil {
			rsp = g.HandlePack(r, un)
		} else {
			rsp = s.HandlePack(r, un)
		}

		if r.ID != 0 {
			err := un.WritePack(packet.NewRsp(r.ID, rsp))
			if err != nil {
				s.logger.Errorf("%s <=> %s:%s(%v) err:%v", un.Name, r.SN, r.Method, string(r.Params), err.Error())
				return
			}
			if err, ok := rsp.(error); ok {
				s.logger.Debugf("%s <=> %s:%s(%v) err:%v", un.Name, r.SN, r.Method, string(r.Params), err.Error())
			} else if raw, ok := rsp.(*json.RawMessage); ok {
				s.logger.Debugf("%s <=> %s:%s(%v) result:%v", un.Name, r.SN, r.Method, string(r.Params), string(*raw))
			} else {
				s.logger.Debugf("%s <=> %s:%s(%v) result:%v", un.Name, r.SN, r.Method, string(r.Params), rsp)
			}
			return
		}
		s.logger.Debugf("%s <=> %s:%s(%v) success", un.Name, r.SN, r.Method, string(r.Params))
	}
}

func (s *Sugar) group(name string) *group.Group {
	if name == "" {
		return nil
	}

	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.glist[name]
}

func (s *Sugar) HandlePack(r *packet.Request, un *group.Unit) interface{} {
	switch r.Method {
	case "route":
		return s.handleRevRoute(r.Params, un)

	case "state":
		return s.handleRevState(r.Params, un)
	}
	return nil
}

func (s *Sugar) handleRevRoute(raw []byte, un *group.Unit) interface{} {
	if g := s.group(un.Name); g != nil {
		return g.HandleRouter(raw, un)
	}
	return fmt.Errorf("not found %s group", un.Name)
}

func (s *Sugar) handleRevState(raw []byte, un *group.Unit) interface{} {
	return nil
}
