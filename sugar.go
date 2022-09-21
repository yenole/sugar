package sugar

import (
	"net/http"
	"sync"

	"github.com/yenole/sugar/group"
	"github.com/yenole/sugar/logger"
	"github.com/yenole/sugar/network/async"
)

type Option struct {
	Listen  string
	Gateway string
}

type Sugar struct {
	mux   sync.RWMutex
	glist map[string]*group.Group

	logger logger.Logger

	g *async.Async
}

func New(logger logger.Logger) *Sugar {
	s := &Sugar{
		glist:  make(map[string]*group.Group),
		g:      async.New(),
		logger: logger,
	}
	return s
}

func (s *Sugar) onRevUnit(un *group.Unit) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	if _, ok := s.glist[un.Name]; !ok {
		s.glist[un.Name] = group.NewGroup(s.logger)
	}

	g := s.glist[un.Name]
	if err := g.HandleRevUnit(un, s.done(un)); err != nil {
		s.logger.Errorf("rev sugar unit fail:%v", err.Error())
		return err
	}
	return nil
}

func (s *Sugar) Run(o *Option) {
	s.Listen(o.Gateway)

	s.logger.Infof("sugar listen %v\n", o.Listen)
	http.ListenAndServe(o.Listen, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.mux.RLock()
		defer s.mux.RUnlock()
		for _, g := range s.glist {
			if h := g.Match(r); h != nil {
				h.ServeHTTP(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}
