package main

import (
	"vectors/web"
)

type (
	ctrls struct {
	}
)

func (self ctrls) hello_world(hd *web.THandler) {
	hd.RespondString("Hello Webgo World!")
}

func main() {
	srv1 := web.NewServer("server1")
	srv1.ShowRoute(true)
	srv1.Get("/hello", ctrls.hello_world)
	srv1.Get("/hello2", func(c *web.THandler) {
		c.RespondString("Hello, World")
		return
	})

	go srv1.Listen(":8080")

	srv2 := web.NewServer("")
	srv2.ShowRoute(true)
	//srv2.Logger.SetLevel(4)
	srv2.Proxy(nil, "/hello", "http", "localhost:8080")
	srv2.Proxy(nil, "/s", "https", "www.baidu.com")
	srv2.Listen(":8888")
	//<-make(chan int) //需要一个循环防止程序推出
}
