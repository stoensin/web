package web

/*
	#模板全局默认变量
		App
		│
		├─module 应用模块目录
		│  ├─web 模块目录
		│  │  ├─staric 静态资源目录
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
		├─staric 静态资源目录
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
*/

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"vectors/utils"
	"vectors/web/template"

	log "github.com/VectorsOrigin/logger"
)

type (
	TServer struct {
		TModule
		GVar map[string]interface{} //全局变量. 废弃

		Config *TConfig // 配置类
		Router *TRouter // 路由类
		//Logger   *logger.TLogger        // 日志类
		Template *template.TTemplateSet // 模板类

		debugMode bool
	}
)

func Version() string {
	return "0.0.1.161210"
}

// 新建一个记录器
func NewServer(configFile string) *TServer {
	srv := &TServer{
		TModule: *NewModule(nil, ""),
		Config:  NewConfig(configFile),
		//Logger:   logger.NewLogger(""),
		Router:   NewRouter(),
		Template: template.NewTemplateSet(),
	}
	// 初始化服务器资源路径为APP当前路径
	srv.TModule.Path = ""     //utils.AppDir()
	srv.TModule.FilePath = "" //utils.AppDir() // 替换文件层次

	// 传递
	srv.Router.Server = srv // 传递服务器指针
	//srv.Router.Logger = srv.Logger
	srv.Router.Template = srv.Template

	return srv
}

// 调试模式
// 关闭所有缓存
func (self *TServer) Debug(debug_mode bool) {
	self.debugMode = debug_mode
	if debug_mode {
		logger.SetLevel(log.LevelDebug)
		//self.Router.
		self.Template.Cacheable = false // 关闭模板缓存
	} else {
		logger.SetLevel(log.LevelInfo)
		//self.Router.
		self.Template.Cacheable = true // 关闭模板缓存
	}

}

func (self *TServer) _fmtAddr(aAddr []string) string {
	// 如果已经配置了端口则不使用
	if self.Config.Port < 10 {
		self.Config.Addr = ""
		self.Config.Port = 8000

		if len(aAddr) != 0 {
			lAddrSplitter := strings.Split(aAddr[0], ":")
			if len(lAddrSplitter) != 2 {
				logger.Err("Address %s is unavailable!", aAddr[0])

			} else {
				self.Config.Addr = lAddrSplitter[0]
				self.Config.Port = utils.StrToInt64(lAddrSplitter[1])
			}
		}
	}

	return fmt.Sprintf("%s:%d", self.Config.Addr, self.Config.Port)
}
func (self *TServer) Listen(aAddr ...string) {
	//注册主Route
	self.Router.RegisterModule(self)
	self.Router.Init()
	lAddr := self._fmtAddr(aAddr)

	// 阻塞监听
	err := http.ListenAndServe(lAddr, self.Router)
	if err != nil {
		logger.Panic("start server faild : %s", err)
	}

	// 显示系统信息
	logger.Info("Listening on http %s", lAddr)
}

func (self *TServer) ListenTLS(certFile, keyFile string, aAddr ...string) {
	//注册主Route
	self.Router.RegisterModule(self)

	self.Router.Init()
	lAddr := self._fmtAddr(aAddr)

	// 阻塞监听
	err := http.ListenAndServeTLS(lAddr, certFile, keyFile, self.Router)
	if err != nil {
		logger.Panic("start server faild : %s", err)
	}

	// 显示系统信息
	logger.Info("listening on https %s", lAddr)
}

func (self *TServer) ShowRoute(sw bool) {
	self.Router.ShowRoute(sw)
}

func (self *TServer) ShowRouter() {
	self.Config.PrintRouterTree = true
}

func (self *TServer) SetVar(name string, value interface{}) {
	//self.GVar[name] = value
}

func (self *TServer) AddVar(name string, value interface{}) {
	self.Router.AddVar(name, value)
}

func (self *TServer) DelVar(name string) {
	self.Router.DelVar(name)
}

func (self *TServer) RegisterModule(obj IModule) {
	self.Router.RegisterModule(obj)
}

// 注册中间件
// 中间件可以使用在Conntroller，全局Object 上
func (self *TServer) RegisterMiddleware(obj ...IMiddleware) {
	self.Router.RegisterMiddleware(obj...)

}

func (self *TServer) LoadConfigFile(filepath string) {
	self.Config.LoadFromFile(filepath)

}

func SetStaticPath(url string, path string) {
	//log.Print(path)
	//log.Print(http.Dir(path))
	http.Handle(url+"/", http.StripPrefix(url, http.FileServer(http.Dir(path))))
	//log.Print(http.StripPrefix(url, http.FileServer(http.Dir(path))))
	return
}

// 获得[服务器静态]文件地址
func (self *TServer) GetStaticPath() string {
	return path.Join(AppPath, STATIC_DIR)
}

// 获得[服务器模板]文件地址
func (self *TServer) GetTemplatesPath() string {
	return path.Join(AppPath, TEMPLATE_DIR)
}

//废弃
func (self *TServer) __SetLoggerLevel(aLevel int) {
	logger.SetLevel(aLevel)
}

//Stops the web server
func Close() {
	os.Exit(0)
}
