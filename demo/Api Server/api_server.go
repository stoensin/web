package main

import (
	"fmt"
	"strings"
	"vectors/web"
)

func Post(hd *web.THandler) {
	lParams := hd.MethodParams()
	lStr := lParams.AsString("test")

	hd.RespondByJson(lStr)
}

func Get(hd *web.THandler) {
	hd.RespondByJson("hello world")
}

func main() {
	lR := web.NewServer("")
	lR.Logger.SetLevel(4)
	lR.Get("/api/post", Post)
	lR.Get("/api/get", Get)
	lR.Listen(":8888")
	//<-make(chan int) //需要一个循环防止程序推出
}
