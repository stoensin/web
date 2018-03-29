package xnet

import (
	"log"
	"net"
	"net/http"
	"os"
)

type (
	// 服务对象
	// registered to serve a particular path or subtree
	// in the HTTP server.
	//
	// ServeHTTP should write reply headers and data to the ResponseWriter
	// and then return.  Returning signals that the request is finished
	// and that the HTTP server can move on to the next request on
	// the connection.
	TServeObject interface {
		ServeHTTP(http.ResponseWriter, *http.Request)
	}

	THttpServer struct {
		*http.Server
		Sock net.Listener //保存套接字以供关闭
	}
)

func NewHttpServer() *THttpServer {
	return &THttpServer{
		Server: &http.Server{
		//ReadTimeout:  c.ReadTimeout,
		//WriteTimeout: c.WriteTimeout,
		//Addr:         c.Address,
		},
	}
}

func (self *THttpServer) Listen(kind, laddr string) {
	var err error
	self.Sock, err = net.Listen(kind, laddr)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
		//self.Logger.Fatal("Router.Listen:", err)
	}
}

// Serve 会对每次请求,使用Router接口ServeHTTP调用处理对应的Route携带的Handler
func (self *THttpServer) Serve(object TServeObject) {
	err := self.Serve(self.Sock, object) // 接受一个Sock和一个带接口ServeHTTP的Router,handler..
	if err != nil {
		//self.Logger.Fatal("Router.Listen:", err)

		log.Println(err)
		os.Exit(-1)
	}
}

//Stops the web server
func (self *THttpServer) Close() {
	if self.Sock != nil {
		self.Sock.Close()
	}
}
