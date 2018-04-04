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

// 用于标记服务器创建标记是否重复
var ServerNames []string

type (
	TServer struct {
		TModule
		GVar map[string]interface{} //全局变量. 废弃

		Config   *TConfig               // 配置类
		Router   *TRouter               // 路由类
		Template *template.TTemplateSet // 模板类
		//Logger   *logger.TLogger        // 日志类
		//debugMode bool
	}
)

func Version() string {
	return "0.0.1.161210"
}

// 新建一个记录器
// 服务器名称
func NewServer(name ...string) *TServer {
	// 确定服务器名称
	var server_name string
	if len(name) != 0 {
		server_name = strings.ToLower(name[0])
	}

	if server_name == "" {
		server_name = "server"
	}

	// 验校服务名称
	if utils.InStrings(server_name, ServerNames...) != -1 {
		logger.Panic("Server %s is existing in the list %v", server_name, ServerNames)
	}
	ServerNames = append(ServerNames, server_name)

	srv := &TServer{
		TModule: *NewModule(nil, server_name),
		Config:  NewConfig(),
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
	self.Config.DebugMode = debug_mode
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

func (self *TServer) parse_addr(addr []string) (host string, port int) {
	// 如果已经配置了端口则不使用
	if len(addr) != 0 {
		lAddrSplitter := strings.Split(addr[0], ":")
		if len(lAddrSplitter) != 2 {
			logger.Err("Address %s of server %s is unavailable!", addr[0], self.Name)
		} else {
			host = lAddrSplitter[0]
			port = utils.StrToInt(lAddrSplitter[1])
		}
	}

	return
}

func (self *TServer) Listen(addr ...string) {
	// 解析地址
	host, port := self.parse_addr(addr)
	logger.Dbg("", host, port)
	self.Config.LoadFromFile(CONFIG_FILE_NAME)
	// 确认配置已经被加载加载
	// 配置最终处理
	sec, err := self.Config.GetSection(self.Name)
	if err != nil {

		// 存储默认
		sec, err = self.Config.NewSection(self.Name)
		if err != nil {
			logger.Panic("creating ini' section faild! Name:%s Error:%s", self.Name, err.Error())
		}
		if host != "" {
			self.Config.Host = host
		}
		self.Config.Port = port
		sec.ReflectFrom(self.Config)
		//self.Config.Save()

	}

	// 映射到服务器配置结构里
	sec.MapTo(self.Config) // 加载
	self.Config.Save()     // 保存文件

	//注册主Route
	self.Router.RegisterModule(self)
	self.Router.Init()
	// 阻塞监听
	// 显示系统信息
	new_addr := fmt.Sprintf("%s:%d", self.Config.Host, self.Config.Port)
	logger.Info("Listening on address: %s", new_addr)
	if self.Config.EnabledTLS {
		if self.Config.TLSCertFile == "" || self.Config.TLSKeyFile == "" {
			logger.Panic("lost cert file or key file for TLS connection!")
		}

		err := http.ListenAndServeTLS(new_addr, self.Config.TLSCertFile, self.Config.TLSKeyFile, self.Router)
		if err != nil {
			logger.Panic("start server faild : %s", err)
		}

	} else {
		err := http.ListenAndServe(new_addr, self.Router)
		if err != nil {
			logger.Panic("start server faild : %s", err)
		}
	}
}

// 废弃
func (self *TServer) __ListenTLS(certFile, keyFile string, addr ...string) {
	/*	//注册主Route
		self.Router.RegisterModule(self)
		self.Router.Init()
		self.Config.Init()
		lAddr := self.parse_addr(addr)
		// 阻塞监听
		err := http.ListenAndServeTLS(lAddr, certFile, keyFile, self.Router)
		if err != nil {
			logger.Panic("start server faild : %s", err)
		}

		// 显示系统信息
		logger.Info("listening on https %s", lAddr)
	*/
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
