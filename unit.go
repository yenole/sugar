package sugar

import (
	"github.com/yenole/sugar/network"
	"github.com/yenole/sugar/packet"
)

const (
	ACT_ROUTER uint8 = iota
)

type Params map[string]interface{}

type Server struct {
	sid  string
	name string
	plot string
	conn network.Conn
}

func (s *Server) onReceive(logger Logger, h func(*packet.Request, *Server), close func()) {
	defer close()

	logger.Infof("name %s sid %s join sugar\n", s.name, s.sid)
	for {
		req, err := s.conn.Request()
		if err != nil {
			logger.Errorf("sn:%s sid:%s request err: %v\n", err)
			return
		}

		h(req, s)
	}
}
