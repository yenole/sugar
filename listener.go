package sugar

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/yenole/sugar/network"
)

var (
	hello = []byte(fmt.Sprintf(`sugar://%s`, Version))
)

func (s *Sugar) Listen() {
	s.logger.Infof("sugar reigst listen %v", *gListen)
	ln, err := net.Listen("tcp", *gListen)
	if err != nil {
		s.logger.Errorf("err: %v\n", err)
		return
	}
	go func() {
		for {
			rev, err := ln.Accept()
			if err != nil {
				s.logger.Errorf("err: %v\n", err)
				break
			}
			s.logger.Debugf("on has rev in %v\n", rev.RemoteAddr())
			go s.onAcceptRevAuth(rev)
		}
	}()
}

func (s *Sugar) onAcceptRevAuth(cnn net.Conn) {
	byts := make([]byte, 1024)
	n, err := cnn.Read(byts)
	if err != nil {
		s.logger.Errorf("accept auth fail: %v", err)
		return
	}
	if !bytes.HasPrefix(byts[:n], hello) {
		s.logger.Errorf("accept is not sugar uint")
		cnn.Close()
		return
	}

	url, err := url.Parse(string(byts[:n]))
	if err != nil || url.Path == "/" || url.Path == "" {
		s.logger.Errorf("accept is not sugar unit")
		cnn.Close()
		return
	}

	svr := &Server{
		sid:  PickSId(cnn.RemoteAddr()),
		plot: url.Query().Get("plot"),
		name: strings.TrimPrefix(url.Path, "/"),
		conn: network.Wrap(cnn),
	}

	err = s.revUnit(svr)
	if err != nil {
		s.logger.Errorf("err: %v\n", err)
		cnn.Close()
		return
	}
	// go s.svrs[sid].onReceive(s)
}

func CalsMD5(v string) string {
	w := md5.New()
	io.WriteString(w, v)
	return fmt.Sprintf("%x", w.Sum(nil))
}

func PickSId(addr net.Addr) string {
	return CalsMD5(fmt.Sprintf("%s-%v", addr, time.Now().Unix()))
}
