package main

import (
	"fmt"

	"github.com/VectorsOrigin/web"
)

type (
	ctrls struct {
		//这里写中间件
	}
)

func (self ctrls) hello_world(hd *web.THandler) {
	hd.RespondString("Hello Webgo World!")
}

func main() {
	srv := web.NewServer("")
	srv.Get("/hello", ctrls.hello_world)
	srv.Get("/hello2", func(c *web.THandler) {
		c.RespondString("Hello, World")
		return
	})

	srv.Get("/hello3", func(c *web.THandler) {
		c.RenderTemplate("hello_world.html", map[string]interface{}{"static": "youpath"})
		fmt.Println("b", c.Route.FilePath)
		return
	})

	srv.Listen(":8080")
}
