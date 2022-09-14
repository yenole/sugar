package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yenole/sugar/unit"
)

var (
	gateway = flag.String("gateway", "tcp://35.220.206.100:7899", "sugar gateway address")
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("args is error")
		os.Exit(1)
	}

	cli := unit.New("cli", &unit.Option{
		Protocol: *gateway,
	})

	err := cli.Run()
	if err != nil {
		fmt.Println("cli run fail:", err.Error())
		os.Exit(1)
	}
	time.Sleep(time.Second * 2)

	var rsp json.RawMessage
	err = cli.Call(os.Args[1], os.Args[2], json.RawMessage(os.Args[3]), &rsp)
	if err != nil {
		fmt.Printf("call %s:%s fail:%s\n", os.Args[1], os.Args[2], err.Error())
		os.Exit(1)
	}

	var buffer bytes.Buffer
	json.Indent(&buffer, rsp, "", "\t")
	fmt.Println(buffer.String())
}
