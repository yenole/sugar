package sugar

import (
	"fmt"

	"github.com/yenole/sugar/network"
	"github.com/yenole/sugar/packet"
)

const (
	ACT_ROUTER uint8 = iota
)

type Server struct {
	sid  string
	name string
	plot string
	conn network.Conn
}

func (s *Server) onReceive(h func(*packet.Request, *Server), close func()) {
	fmt.Println(s.conn.RemoteAddr())
	fmt.Printf("name %s sid %s join sugar\n", s.name, s.sid)
	defer close()

	for {
		req, err := s.conn.Request()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}

		h(req, s)
	}
}
