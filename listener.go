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
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	go func() {
		for {
			rev, err := ln.Accept()
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
			fmt.Printf("on has rev in %v\n", rev.RemoteAddr())
			go s.onAcceptRevAuth(rev)
		}
	}()
}

func (s *Sugar) onAcceptRevAuth(cnn net.Conn) {
	byts := make([]byte, 1024)
	n, err := cnn.Read(byts)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	if !bytes.HasPrefix(byts[:n], hello) {
		fmt.Println("accept is not sugar unit")
		cnn.Close()
		return
	}
	fmt.Printf("string(byts): %v\n", string(byts[:n]))
	url, err := url.Parse(string(byts[:n]))
	if err != nil || url.Path == "/" || url.Path == "" {
		fmt.Println("accept is not sugar unit")
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
		fmt.Printf("err: %v\n", err)
		cnn.Close()
		return
	}
	// go s.svrs[sid].onReceive(s)
}

func CalsMD5(v string) string {
	fmt.Printf("v: %v\n", v)
	w := md5.New()
	io.WriteString(w, v)
	return fmt.Sprintf("%x", w.Sum(nil))
}

func PickSId(addr net.Addr) string {
	return CalsMD5(fmt.Sprintf("%s-%v", addr, time.Now().Unix()))
}
