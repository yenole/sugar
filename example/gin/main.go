package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yenole/sugar/route"
	"github.com/yenole/sugar/unit"
)

func main() {
	rt := gin.New()

	uni := unit.New("api", &unit.Option{
		Protocol: "tcp://35.220.206.100:7899",
		Router: &route.Option{
			Type:   0,
			Host:   "ethsana.sana",
			Listen: "http://localhost:8082",
			Routes: map[string]string{
				`/hello`: http.MethodGet,
			},
		},
	})
	err := uni.Run()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	rt.GET("/hello", func(ctx *gin.Context) {
		var sum uint
		err := uni.Call("calc", "add", []interface{}{122, 100}, &sum)
		if err != nil {
			ctx.String(http.StatusOK, fmt.Sprintf("err:%s\n", err.Error()))
		} else {
			ctx.String(http.StatusOK, fmt.Sprintf("this is sugar unit resp , calc add result:%v\n", sum))
		}
	})

	rt.Run(":8082")
}
