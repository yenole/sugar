package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yenole/sugar/handler"
	"github.com/yenole/sugar/unit"
)

func main() {
	uni := unit.New("calc", &unit.Option{
		// Protocol: "tcp://35.220.206.100:7899",
		Protocol: "tcp://:8081",
	})
	err := uni.Run()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	uni.Handle("add", func(ctx *handler.Context) interface{} {
		var nums []uint
		if err := ctx.Bind(&nums); err != nil {
			return err
		}
		fmt.Printf("nums: %v\n", nums)
		return nums[0] + nums[1]
	})

	go func() {
		time.Sleep(time.Second * 5)
		var resp json.RawMessage
		err := uni.Call("platform", "login", map[string]interface{}{
			"token": "xusir92@gmail.com:e10adc3949ba59abbe56e057f20f883e",
		}, &resp)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		fmt.Printf("string(resp): %v\n", string(resp))
	}()

	sign := make(chan struct{})
	<-sign
}
