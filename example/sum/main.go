package main

import (
	"fmt"

	"github.com/yenole/sugar/handler"
	"github.com/yenole/sugar/unit"
)

func main() {
	uni := unit.New("calc", &unit.Option{
		Protocol: "tcp://127.0.0.1:8081",
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
	sign := make(chan struct{})
	<-sign
}
