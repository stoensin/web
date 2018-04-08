# webgo
a golang web framework for vectors ERP system

框架提供关于Http服务器端最精简逻辑实现,理论上可以兼容大部分其他框架中间件(需要小量修改)

##QQ Group 151120790

	服务器目录树
	App
	│
	├─module 应用模块目录
	│  ├─web 模块目录
	│  │  ├─static 静态资源目录
	│  │  │   ├─uploads 上传根目录
	│  │  │   ├─lib 资源库文件目录(常用作前端框架库)
	│  │  │   └─src 资源文件
	│  │  │      ├─js 资源Js文件目录
	│  │  │      ├─img 资源图片文件目录
	│  │  │      └─css 资源Css文件
	│  │  ├─model 模型目录
	│  │  ├─template 视图文件目录
	│  │  ├─data 数据目录
	│  │  ├─model 模型目录
	│  │  └─controller.go 控制器
	│  │
	│  ├─base 模块目录
	│  │
	│  └─... 扩展的可装卸功能模块或插件
	│
	├─static 静态资源目录
	│  ├─uploads 上传根目录
	│  ├─lib 资源库文件目录(常用作前端框架库)
	│  └─src 资源文件
	│     ├─js 资源Js文件目录
	│     ├─img 资源图片文件目录
	│     └─css 资源Css文件
	├─template 视图文件目录
	├─deploy 部署文件目录
	│
	├─main.go 主文件
	└─main.ini 配置文件

## hello world demo

	package main

	import (
		"fmt"

		"github.com/VectorsOrigin/web"
	)

	type (
		ctrls struct {
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
