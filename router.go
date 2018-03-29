package web

/*
	Router 负责把所有的URL映射出去
	1.提供全局变量操作
	2.提供
*/

import (
	"context"
	_template "html/template"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
	"vectors/logger"
	"vectors/utils"
	"vectors/web/template"
)

const (
	ROUTER_VER = "1.2.0"
)

const (
	CommomRoute RouteType = iota // extenion route
	HookBeforeRoute
	HookAfterRoute
	ReplaceRoute // the route replace orgin

)

type (
	RouteType byte

	// 路由节点绑定的 方法(Handler)
	TMethodType struct {
		Name      string         // 名称
		Func      reflect.Value  // 方法本体
		FuncType  reflect.Type   // 方法类型
		ArgType   []reflect.Type // 参数组类型
		ReplyType []reflect.Type //TODO 返回多结果
	}

	// TRoute 路,表示一个Link 连接地址"../webgo/"
	// 提供基础数据参数供Handler处理
	TRoute struct {
		ID       int64  // 在Tree里的Id号
		Path     string // 网络路径
		FilePath string // 短存储路径
		Model    string // 模型/对象/模块名称 Tmodule/Tmodel, "Model.Action", "404"
		Action   string // 动作名称[包含模块名，动作名] "Model.Action", "/index.html","/filename.png"
		//TemplatePath string // 废弃
		FileName string
		Type     RouteType // Route 类型 决定合并的形式

		// ReverseProxy
		//Scheme         string
		//Host           string
		Host           *url.URL
		isReverseProxy bool //# 是反向代理
		isDynRoute     bool // 是否*动态路由   /base/*.html

		MainCtrl TMethodType   // 主控制器 每个Route都会有一个主要的Ctrl,其他为Hook的Ctrl
		Ctrls    []TMethodType // 最终控制器 合并主控制器+次控制器
		//HookCtrl map[string][]TMethodType // 次控制器 map[*][]TMethodType 匹配所有  Hook的Ctrl会在主的Ctrl执行完后执行
		//Ctrls    map[string][]TMethodType // 最终控制器 合并主控制器+次控制器
	}

	// TRouter 路由,路的集合.负责管理TRoute
	TRouter struct {
		I18n        *TI18n
		Server      *TServer
		Logger      *logger.TLogger
		Template    *template.TTemplateSet
		TempleteVar map[string]interface{} // ！！！储存模板全局变量.将应用到说有模板中去
		GVar        map[string]interface{} // ！！！全局变量. 需改进
		//Sock        net.Listener           //保存套接字以供关闭

		show_route bool

		//Routes      []TRoute //废弃
		tree       *TTree
		middleware *TMiddlewareManager // 中间件
		//beforeRoute reflect.Value //废弃 动作处理器

		lock              sync.RWMutex
		handlerPool       sync.Pool
		proxy_handlerPool sync.Pool
		respPool          sync.Pool
		actionPool        *TPool
	}
)

// 创建Router
func NewRouter() *TRouter {
	lRouter := &TRouter{
		TempleteVar: map[string]interface{}{},
		GVar:        map[string]interface{}{},
		tree:        NewRouteTree(),
	}

	//
	lRouter.middleware = NewMiddlewareManager()

	// inite HandlerPool New function
	lRouter.handlerPool.New = func() interface{} {
		return NewHandler()
	}
	//lRouter.handlerPool, _ = cache.NewCacher("memory", `{"interval":180,"expired":3600}`)
	//lRouter.handlerPool.New(func() interface{} {
	//	return NewHandler()
	//})

	lRouter.proxy_handlerPool.New = func() interface{} {
		return NewProxyHandler()
	}

	// inite HandlerPool New function
	lRouter.respPool.New = func() interface{} {
		return NewResponser()
	}
	//lRouter.respPool, _ = cache.NewCacher("memory", `{"interval":180,"expired":3600}`)
	//lRouter.respPool.New(func() interface{} {
	//	return NewResponser()
	//})

	lRouter.actionPool = NewPool()

	// Router 默认[全局变量]
	lRouter.GVar["Version"] = ROUTER_VER
	lRouter.GVar["StratDateTime"] = time.Now().UTC()

	return lRouter
}

/*
func (self *TRoute) IsNil() bool {
	return !(len(self.Controllers) == 0)
}
*/

// TODO 管理Ctrl 顺序 before center after
// 根据不同Action 名称合并Ctrls
func (self *TRoute) CombineController(aFrom *TRoute) {
	switch aFrom.Type {
	/*	case CommomRoute:
		// 普通Tree合并
		{
			self.MainCtrl = append(self.MainCtrl, aFrom.MainCtrl...)
		}*/
	case HookBeforeRoute:
		/*
			// 静态：Url xxx\xxx\action 将直接插入
			// 动态：根据Action 插入到Map 叠加
			if aFrom.isDynRoute {
				self.HookCtrl[aFrom.Action] = append(self.HookCtrl[aFrom.Action], aFrom.MainCtrl)
			} else {
				self.MainCtrl = append(self.MainCtrl, aFrom.MainCtrl...)
			}
		*/

		self.Ctrls = []TMethodType{self.MainCtrl}
		self.Ctrls = append(self.Ctrls, self.Ctrls...)
	case HookAfterRoute:
		self.Ctrls = append(self.Ctrls, aFrom.MainCtrl)
	default:
		// 替换路由会直接替换 主控制器 但不会影响其他Hook 进来的控制器
		self.MainCtrl = aFrom.MainCtrl
		self.Ctrls = []TMethodType{self.MainCtrl}
		//logger.Dbg("CombineController", self.MainCtrl, aFrom.MainCtrl)
	}
}

/*
// 从handler 获得中间件名称列表
// 例：Struct.Handle(Handler *THandler,Action Struct)的Struct
func (self *TRoute) mapMiddleware(aHandlerType reflect.Type) int64 {
	self.Middleware = map[string]bool{}
	//Warn("Middleware", aHandlerType)
	// NumIn returns a function type's input parameter count.  // It panics if the type's Kind is not Func.
	for i := 0; i < aHandlerType.NumIn(); i++ { //参数
		lParm := aHandlerType.In(i)
		//Warn("Middleware", lParm.Name(), lParm.Kind())
		if (i == 0 || i == 2) && lParm.Kind() == reflect.Struct {
			//Warn("Middleware", lParm.NumField())
			for k := 0; k < lParm.NumField(); k++ { //成员
				lTypeName := strings.Split(lParm.Field(k).Type.String(), ".")
				self.Middleware[lTypeName[len(lTypeName)-1]] = true
				//Warn("Middleware", lParm.NumField(), lParm.Field(k).Name)
			}
		}
	}
	return int64(len(self.Middleware))
}
*/

/*
 初始化所有加载工作
*/
func (self *TRouter) Init() {

	// 创建并初始化[国际化]
	if self.Server.Config.UseI18N {
		self.Logger.Info("Use I18N")
		self.I18n = NewI18n("XWeb") // I18N Name
		lPath := filepath.Join(self.Server.Config.RootPath, self.Server.Config.LocaleDir)
		if err := self.I18n.Init(lPath, self.Server.Config.LangCode); err != nil {
			self.Logger.ErrLn("I18N(Init):", err)
		}

		lTemplateFuncs := map[string]interface{}{
			"trans": func(aText string) _template.HTML {
				return _template.HTML(self.I18n.Translate(aText))
			},
		}

		self.Template.AddFuncs(lTemplateFuncs)
	}

	//self.RegisterModules(admin.Admin)
	if self.Server.Config.PrintRouterTree {
		self.tree.PrintTrees()
	}

}

/*
func (self *TRouter) RegisterModules(obj interface{}) {
	self.lock.Lock()
	defer self.lock.Unlock()

	Trace("RegisterModules:", reflect.TypeOf(obj))
	if m, ok := obj.(*TModule); ok { //检查是否有IRestCreator.Create(hd *TRestHandle)
		self.Routes = append(self.Routes, m.Routes...) // 注意要加[省略号] !!!暂时有重复合并问题
		//self.Routes = MergeMaps(self.Routes, m.Routes) // 合并两个Maps安全点
		//#创建文件夹
		//os.Mkdir("./modules/aa", 0700) //>>>>>>>>>>
		mdlPath := filepath.Join(self.Server.Config.ModulesDir, m.Path)
		err := os.Mkdir(mdlPath, 0700)
		if os.IsExist(err) {
		}
		os.Mkdir(filepath.Join(mdlPath, self.Server.Config.StaticDir), 0700)
		os.Mkdir(filepath.Join(mdlPath, self.Server.Config.TemplatesDir), 0700)

		//Trace("OK")
	} else {
		Debug("RegisterModules:Object is not useable!")
	}
	//log.Println("r", self.Routes)
}
*/
func (self *TRouter) RegisterModule(aMd IModule, build_path ...bool) {
	if aMd == nil {
		self.Logger.WarnLn("RegisterModule is nil")
		return
	}

	// 执行注册器接口
	if a, ok := aMd.(IModuleRegister); ok {
		a.Register()
	}

	///lRoutes := aMd.GetRoutes()
	lModuleFilePath := utils.Trim(aMd.GetFilePath())

	//self.Logger.("RegisterModules:", reflect.TypeOf(aMd))

	self.lock.Lock() //<-锁
	self.tree.Conbine(aMd.GetRoutes())
	///self.Routes = append(self.Routes, lRoutes...) // 注意要加[省略号] !!!暂时有重复合并问题
	//self.Routes = MergeMaps(self.Routes, m.Routes) // 合并两个Maps安全点

	self.lock.Unlock() //<-

	//#创建文件夹
	//os.Mkdir("./modules/aa", 0700) //>>>>>>>>>>

	// The Path must be not blank.
	// <待优化静态路径管理>必须不是空白路径才能组合正确
	if len(build_path) > 0 && build_path[0] && len(lModuleFilePath) > 0 {
		lModuleFilePath := filepath.Join(self.Server.Config.ModulesDir, lModuleFilePath)
		err := os.Mkdir(lModuleFilePath, 0700)
		if err != nil {
			os.Mkdir(filepath.Join(lModuleFilePath, self.Server.Config.StaticDir), 0700)
			os.Mkdir(filepath.Join(lModuleFilePath, self.Server.Config.TemplatesDir), 0700)
		}
	}

}

// 注册中间件
func (self *TRouter) RegisterMiddleware(aMd ...IMiddleware) {
	for _, m := range aMd {
		lType := reflect.TypeOf(m)
		if lType.Kind() == reflect.Ptr {
			lType = lType.Elem()
		}
		lName := lType.String()
		self.middleware.Add(lName, m)
		//lName := strings.Split(lType.String(), ".")
		//self.middleware.Add(lName[len(lName)-1], m)
		//Trace("RegisterMiddleware:", lName)
	}

}

func (self *TRouter) AddVar(name string, value interface{}) {
	self.GVar[name] = value
}

func (self *TRouter) GetVar(name string) interface{} {
	return self.GVar[name]
}

func (self *TRouter) DelVar(name string) {
	delete(self.GVar, name)
}

//安全调用Handle 函数(resp []reflect.Value, e interface{})
// TODO 有待优化
func (self *TRouter) safelyCall(function reflect.Value, args []reflect.Value, hd *THandler, aActionValue reflect.Value) {
	// 错误处理
	defer func() {
		if err := recover(); err != nil {
			if self.Server.Config.RecoverPanic { //是否绕过错误处理直接关闭程序
				self.routePanic(hd, aActionValue)
				//e = err
				//resp = nil
				self.Logger.ErrLn("Handler crashed with error : ", err)
				for i := 1; ; i++ {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					self.Logger.ErrLn(file, line)
				}
			} else {
				// go back to panic
				self.Logger.WarnLn("Panic:SafelyCall")
				panic(err)
			}
		}
	}()
	//调用函数 >>>输出HTML数据

	function.Call(args)
}

// 问题1:
func (self *TRouter) routeStatic(req *http.Request, w *TResponseWriter) {
	var lFilePath string
	lPath, lFileName := filepath.Split(req.URL.Path) //products/js/base.js
	//urlPath := strings.Split(strings.Trim(req.URL.Path, `/`), `/`) // Split不能去除/products

	//根目录静态文件映射过滤
	if lPath == "/" {
		switch filepath.Ext(lFileName) {
		case ".txt", ".html", ".htm": // 目前只开放这种格式
			lFilePath = filepath.Join(lFileName)
		}

	} else {
		//urlPath := strings.Trim(req.URL.Path, `/`)

		//如果第一个是静态文件夹名则选用主静态文件夹,反之使用模块
		// /static/js/base.js
		// /ModuleName/static/js/base.js
		lDirs := strings.Split(lPath, "/")
		if strings.EqualFold(lDirs[1], self.Server.Config.StaticDir) {
			//if strings.HasPrefix(lPath, "/"+self.Server.Config.StaticDir) { // 如果请求是 /Static/js/base.js
			/* static_file = filepath.Join(
			self.Server.Config.RootPath,                           // c:\project\
			self.Server.Config.StaticDir,                          // c:\project\static\
			strings.Join(urlPath[1:], string(filepath.Separator)), // c:\project\static\js\base.js
			fileName)
			*/

			lFilePath = filepath.Join(req.URL.Path)

		} else { // 如果请求是 products/Static/js/base.js
			/* static_file = filepath.Join(
			self.Server.Config.RootPath,                           // c:\project\
			self.Server.Config.ModulesDir,                         // c:\project\Modules
			urlPath[0],                                            // c:\project\Modules\products\
			self.Server.Config.StaticDir,                          // c:\project\Modules\products\static\
			strings.Join(urlPath[1:], string(filepath.Separator)), // c:\project\Modules\products\static\js\base.js
			fileName)
			*/

			//Debug("lDirsD", lDirs, self.Server.Config.StaticDir, string(os.PathSeparator))
			// 再次检查 Module Name 后必须是 /static 目录
			if strings.EqualFold(lDirs[2], self.Server.Config.StaticDir) {
				lFilePath = filepath.Join(
					self.Server.Config.ModulesDir, // c:\project\Modules
					req.URL.Path)
			} else {
				http.NotFound(w, req)
				return

			}
		}
	}

	// 当模块路径无该文件时，改为程序static文件夹
	if !utils.FileExists(lFilePath) {
		lIndex := strings.Index(lFilePath, self.Server.Config.StaticDir)
		if lIndex != -1 {
			lFilePath = lFilePath[lIndex-1:]

		}
	}

	//Info("static_file", static_file)
	if req.Method == "GET" || req.Method == "HEAD" {
		lFilePath = filepath.Join(
			self.Server.Config.RootPath,
			lFilePath)
		// 当程序文件夹无该文件时
		if !utils.FileExists(lFilePath) {
			self.Logger.DbgLn("Not Found", lFilePath)
			http.NotFound(w, req)
			return
		}
		//Noted: ServeFile() can not accept "/AA.exe" string, only accepy "AA.exe" string.
		http.ServeFile(w, req, lFilePath) // func ServeFile(w ResponseWriter, r *Request, name string)
		//self.Server.Logger.Println("RouteFile:" + static_file)
		return
	}
	//Debug("RouteStatic", path, fileName)
	return
}

// ServeHTTP
// 每个连接Route的入口
func (self *TRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Pool 提供TResponseWriter
	lResp := self.respPool.Get().(*TResponseWriter)
	lResp.connect(w)

	self.routeHandler(req, lResp)

	// Pool 回收TResponseWriter
	lResp.ResponseWriter = nil
	self.respPool.Put(lResp)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// 克隆interface 并复制里面的指针
func cloneInterfacePtrFeild(s interface{}) interface{} {
	lVal := reflect.Indirect(reflect.ValueOf(s)) //Indirect 等同 Elem()
	lType := reflect.TypeOf(s).Elem()            // 返回类型
	lrVal := reflect.New(lType)                  //创建某类型
	lrVal.Elem().Set(lVal)
	/*
		for i := 0; i < lVal.NumField(); i++ {
			lField := lVal.Field(i)
			Warn("jj", lField, lField.Kind())
			if lField.Kind() == reflect.Ptr {
				//fmt.Println("jj", lField, lField.Elem())
				//lrVal.Field(i).SetPointer(unsafe.Pointer(lField.Pointer()))
				lrVal.Elem().Field(i).Set(lField)
				//lrVal.FieldByName("Id").SetString("fasd")
			}
		}
	*/
	//fmt.Println(lrVal)
	//return reflect.Indirect(lrVal).Interface()
	return lrVal.Interface()
}

// TODO:过滤 _ 的中间件
func (self *TRouter) routeBefore(hd *THandler, aActionValue reflect.Value) {
	// Action结构外的其他中间件
	var (
		IsFound         bool
		lField, lMethod reflect.Value
		lType           reflect.Type
		lNew            interface{}
		ml              IMiddleware
	)
	//self.Server.Logger.DbgLn("Name  self.middleware.names", self.middleware.Names)
	for _, key := range self.middleware.Names {
		// @:直接返回 放弃剩下的Handler
		if hd.Response.Written() {
			break
		}
		IsFound = false
		for i := 0; i < aActionValue.NumField(); i++ { // Action结构下的中间件
			lField = aActionValue.Field(i) // 获得成员
			lType = lField.Type()

			if lType.Kind() == reflect.Ptr {
				lType = lType.Elem()
			}

			//self.Server.Logger.DbgLn("Name %s %s", key, lType.Name(), lField.Interface(), lField.Kind(), lField.String())
			ml = self.middleware.Get(key).(IMiddleware)
			if lType.String() == key {
				//Warn("lField.IsValid(),lField.IsNil()", lField.IsValid(), lField.IsNil())
				//if lField.IsValid() { // 存在该Filed

				if lField.Kind() == reflect.Struct {
					//	过滤继承的结构体
					//	type TAction struct {
					//		TEvent
					//	}
					ml.Request(aActionValue.Interface(), hd)
				} else if lField.IsNil() {
					//Warn("lField.IsNil()", key)
					lNew = cloneInterfacePtrFeild(ml)                        // 克隆
					lNew.(IMiddleware).Request(aActionValue.Interface(), hd) // 首先获得基本数据

					lMVal := reflect.ValueOf(lNew) // or reflect.ValueOf(lMiddleware).Convert(lField.Type())
					if lField.Kind() == lMVal.Kind() {
						lField.Set(lMVal) // 通过
					}
				} else {
					// 尝试获取方法
					lMethod = lField.MethodByName("Request")
					//Warn("routeBefore", key, lMethod.IsValid())
					if lMethod.IsValid() {
						lMethod.Call([]reflect.Value{aActionValue, reflect.ValueOf(hd)}) //执行方法
					}
				}
				//}

				// STEP:结束循环
				IsFound = true
				break
			}
		}

		// 更新非控制器中的中间件
		if !IsFound {
			//Warn(" routeBefore not IsFound", key, aActionValue.Interface(), hd)
			ml.Request(aActionValue.Interface(), hd)
		}
	}
}

func (self *TRouter) routeAfter(hd *THandler, aActionValue reflect.Value) {
	var (
		lField, lMethod reflect.Value
		lType           reflect.Type
		lNew            interface{}
	)
	for key, ml := range self.middleware.middlewares {
		for i := 0; i < aActionValue.NumField(); i++ { // Action结构下的中间件
			lField = aActionValue.Field(i) // 获得成员
			lType = lField.Type()

			if lType.Kind() == reflect.Ptr {
				lType = lType.Elem()
			}

			//self.Server.Logger.DbgLn("Name %s %s", key, lType.Name(), lField.Interface(), lField.Kind(), lField.String())

			if lType.String() == key {
				//Warn("lField.IsValid(),lField.IsNil()", lField.IsValid(), lField.IsNil())
				//if lField.IsValid() { // 存在该Filed
				//	过滤继承的结构体
				//	type TAction struct {
				//		TEvent
				//	}
				if lField.Kind() != reflect.Struct && lField.IsNil() {
					//Warn("!aActionValue.IsValid()", lrActionVal)
					lNew = cloneInterfacePtrFeild(ml)                         // 克隆
					lNew.(IMiddleware).Response(aActionValue.Interface(), hd) // 首先获得基本数据
					lMVal := reflect.ValueOf(lNew)                            // or reflect.ValueOf(lMiddleware).Convert(lField.Type())

					if lField.Kind() == lMVal.Kind() {
						lField.Set(lMVal) // 通过
					}
				} else {
					// 尝试获取方法
					lMethod = lField.MethodByName("Response")
					if lMethod.IsValid() {
						lMethod.Call([]reflect.Value{aActionValue, reflect.ValueOf(hd)}) //执行方法
					}
				}
				//}

				// STEP:结束循环
				break
			} else {
				//Warn(" routeBefore", key, aActionValue.Interface(), hd)
				ml.Response(aActionValue.Interface(), hd)
			}
		}
	}
	/*
		lNameLst := make(map[string]bool)
		for i := 0; i < aActionValue.NumField(); i++ { // Action结构下的中间件
			lField := aActionValue.Field(i) // 获得成员
			lType := lField.Type()

			if lField.Kind() == reflect.Struct {
				continue
			}

			if lType.Kind() == reflect.Ptr {
				lType = lType.Elem()
			}
			lFieldName := lType.Name() + lType.String()
			if self.middleware.Contain(lFieldName) {
				lNameLst[lFieldName] = true
				m := lField.MethodByName("Response")
				if m.IsValid() {
					lHdValue := reflect.ValueOf(hd)
					//self.Logger.Info("ccc", m, lHdValue, aActionValue)
					m.Call([]reflect.Value{aActionValue, lHdValue}) //执行方法
				}
			}
		}

		for key, ml := range self.middleware.middlewares { // Action结构下的中间件
			if !lNameLst[key] && ml != nil {
				//Warn("lNameLst", key, ml)
				ml.Response(aActionValue.Interface(), hd)
			}
		}*/

}

func (self *TRouter) routePanic(hd *THandler, aActionValue reflect.Value) {
	if aActionValue.IsValid() {
		lNameLst := make(map[string]bool)
		// @@@@@@@@@@有待优化 可以缓存For结果
		for i := 0; i < aActionValue.NumField(); i++ { // Action结构下的中间件
			lField := aActionValue.Field(i) // 获得成员
			lType := lField.Type()

			// 过滤继承结构的中间件
			if lField.Kind() == reflect.Struct {
				continue
			}

			if lType.Kind() == reflect.Ptr {
				lType = lType.Elem()
			}
			lFieldName := lType.Name() + lType.String()
			if self.middleware.Contain(lFieldName) {
				lNameLst[lFieldName] = true
				m := lField.MethodByName("Panic")
				if m.IsValid() {
					lHdValue := reflect.ValueOf(hd)
					//self.Logger.Info("", m, lHdValue, aActionValue)
					m.Call([]reflect.Value{aActionValue, lHdValue}) //执行方法
				}
			}
		}

		// 重复斌执行上面 遗漏的
		for key, ml := range self.middleware.middlewares { // Action结构下的中间件
			if !lNameLst[key] && ml != nil {
				//Warn("lNameLst", key, ml)
				ml.Panic(aActionValue.Interface(), hd)
			}
		}
	}
}

/*
// 处理URL 请求
// 优化处理
#Pool Route/ResponseWriter
*/
func (self *TRouter) routeHandler(req *http.Request, w *TResponseWriter) {
	lPath := req.URL.Path //获得的地址

	//ar lRoute *TRoute
	//ar lParam Params
	// # match route from tree
	lRoute, lParam := self.tree.Match(req.Method, lPath)
	if self.show_route {
		self.Logger.Info("[Path]%v [Route]%v", lPath, lRoute.FilePath)
	}

	self.Logger.Info("[Path]%v [Route]%v", lPath, lRoute)

	//opy(lParam, Param)
	if lRoute == nil {
		self.routeStatic(req, w) // # serve as a static file link
		return
	}

	if lRoute.isReverseProxy {
		self.routeProxy(lRoute, lParam, req, w)

		return
	}

	// # init Handler
	lHandler := self.handlerPool.Get().(*THandler)
	//lHandler, ok := self.handlerPool.Get().(*THandler)
	//if !ok {
	//	lHandler = NewHandler()
	//	self.Logger.Info("pool %v", self.handlerPool.Len())
	//}

	//lHandler := self.handlerPool.Get().(*THandler) // Pool 提供Handler
	lHandler.connect(w, req, self, lRoute)
	for _, param := range lParam {
		lHandler.setPathParams(param.Name, param.Value)
		//self.Logger.DbgLn("lParam", param.Name, param.Value)
	}

	var (
		args          []reflect.Value //handler参数
		lActionVal    reflect.Value
		lActionTyp    reflect.Type
		lParm         reflect.Type
		CtrlValidable bool
	)

	// TODO:将所有需要执行的Handler 存疑列表或者树-Node保存函数和参数
	//logger.Dbg("lParm %s %d:%d %p %p", lHandler.TemplateSrc, lRoute.Action, lRoute.MainCtrl, len(lRoute.Ctrls), lRoute.Ctrls)
	for index, ctrl := range lRoute.Ctrls {
		lHandler.CtrlIndex = index //index
		// STEP#: 获取<Ctrl.Func()>方法的参数
		for i := 0; i < ctrl.FuncType.NumIn(); i++ {
			lParm = ctrl.FuncType.In(i) // 获得参数

			//self.Logger.DbgLn("lParm%d:", i, lParm, lParm.Name())
			switch lParm { //arg0.Elem() { //获得Handler的第一个参数类型.
			case reflect.TypeOf(lHandler): // if is a pointer of THandler
				{
					//args = append(args, reflect.ValueOf(lHandler)) // 这里将传递本函数先前创建的handle 给请求函数
					args = append(args, lHandler.val) // 这里将传递本函数先前创建的handle 给请求函数
				}
			default:
				{
					//Trace("lParm->default")
					if i == 0 && lParm.Kind() == reflect.Struct { // 第一个 //第一个是方法的结构自己本身 例：(self TMiddleware) ProcessRequest（）的 self
						lActionTyp = lParm
						lActionVal = self.actionPool.Get(lParm)
						if !lActionVal.IsValid() {
							lActionVal = reflect.New(lParm).Elem() //由类生成实体值,必须指针转换而成才是Addressable  错误：lVal := reflect.Zero(aHandleType)
						}
						args = append(args, lActionVal) //插入该类型空值
						break
					}

					// STEP:如果是参数是 http.ResponseWriter 值
					if strings.EqualFold(lParm.String(), "http.ResponseWriter") { // Response 类
						//args = append(args, reflect.ValueOf(w.ResponseWriter))
						args = append(args, w.val)
						break
					}

					// STEP:如果是参数是 http.Request 值
					if lParm == reflect.TypeOf(req) { // request 指针
						args = append(args, reflect.ValueOf(req)) //TODO (同上简化reflect.ValueOf）
						break
					}

					// STEP#:
					args = append(args, reflect.Zero(lParm)) //插入该类型空值

				}
			}
		}

		CtrlValidable = lActionVal.IsValid()
		if CtrlValidable {
			self.Logger.Info("routeBefore")
			self.routeBefore(lHandler, lActionVal)
		}
		logger.Infof("safelyCall %v ,%v", lHandler.Response.Written(), args)
		if !lHandler.Response.Written() {
			//self.Logger.Info("safelyCall")
			// -- execute Handler or Panic Event
			self.safelyCall(ctrl.Func, args, lHandler, lActionVal) //传递参数给函数.<<<
		}

		if !lHandler.Response.Written() && CtrlValidable {
			// # after route
			self.routeAfter(lHandler, lActionVal)
		}

		if CtrlValidable {
			self.actionPool.Put(lActionTyp, lActionVal)
		}
	}

	if lHandler.finalCall.IsValid() {
		if f, ok := lHandler.finalCall.Interface().(func(*THandler)); ok {
			//f([]reflect.Value{reflect.ValueOf(lHandler)})
			f(lHandler)
			//			Trace("Handler Final Call")
		}
	}

	//##################
	//设置某些默认头
	//设置默认的 content-type
	//TODO 由Tree完成
	//tm := time.Now().UTC()
	lHandler.SetHeader(true, "Engine", "vectors web") //取当前时间
	//lHandler.SetHeader(true, "Date", WebTime(tm)) //
	//lHandler.SetHeader(true, "Content-Type", "text/html; charset=utf-8")
	if lHandler.TemplateSrc != "" {
		//添加[static]静态文件路径
		logger.Dbg(self.Server.Config.StaticDir, path.Join(utils.FilePathToPath(lRoute.FilePath), self.Server.Config.StaticDir))
		//	self.AddVar(self.Server.Config.StaticDir, path.Join(utils.FilePathToPath(lRoute.FilePath), self.Server.Config.StaticDir)) //添加[static]静态文件路径
		lHandler.RenderArgs[self.Server.Config.StaticDir] = path.Join(utils.FilePathToPath(lRoute.FilePath), self.Server.Config.StaticDir)
	}

	// 结束Route并返回内容
	lHandler.Apply()

	self.handlerPool.Put(lHandler) // Pool 回收Handler
	return
}

func (self *TRouter) routeProxy(route *TRoute, param Params, req *http.Request, rw *TResponseWriter) {
	lHandler := self.proxy_handlerPool.Get().(*TProxyHandler)
	lHandler.connect(rw, req, self, route)
	transport := lHandler.Transport

	ctx := req.Context()
	if cn, ok := rw.ResponseWriter.(http.CloseNotifier); ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()
		notifyChan := cn.CloseNotify()
		go func() {
			select {
			case <-notifyChan:
				cancel()
			case <-ctx.Done():
			}
		}()
	}

	outreq := new(http.Request)
	*outreq = *req // includes shallow copies of maps, but okay
	if req.ContentLength == 0 {
		outreq.Body = nil // Issue 16036: nil Body for http.Transport retries
	}
	outreq = outreq.WithContext(ctx)

	lHandler.Director(outreq)
	outreq.Close = false

	// We are modifying the same underlying map from req (shallow
	// copied above) so we only copy it if necessary.
	copiedHeaders := false

	// Remove hop-by-hop headers listed in the "Connection" header.
	// See RFC 2616, section 14.10.
	if c := outreq.Header.Get("Connection"); c != "" {
		for _, f := range strings.Split(c, ",") {
			if f = strings.TrimSpace(f); f != "" {
				if !copiedHeaders {
					outreq.Header = make(http.Header)
					copyHeader(outreq.Header, req.Header)
					copiedHeaders = true
				}
				outreq.Header.Del(f)
			}
		}
	}

	// Remove hop-by-hop headers to the backend. Especially
	// important is "Connection" because we want a persistent
	// connection, regardless of what the client sent to us.
	for _, h := range hopHeaders {
		if outreq.Header.Get(h) != "" {
			if !copiedHeaders {
				outreq.Header = make(http.Header)
				copyHeader(outreq.Header, req.Header)
				copiedHeaders = true
			}
			outreq.Header.Del(h)
		}
	}

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// If we aren't the first proxy retain prior
		// X-Forwarded-For information as a comma+space
		// separated list and fold multiple headers into one.
		if prior, ok := outreq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outreq.Header.Set("X-Forwarded-For", clientIP)
	}

	res, err := transport.RoundTrip(outreq)
	if err != nil {
		self.Logger.Err("http: proxy error: %v", err)
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	// Remove hop-by-hop headers listed in the
	// "Connection" header of the response.
	if c := res.Header.Get("Connection"); c != "" {
		for _, f := range strings.Split(c, ",") {
			if f = strings.TrimSpace(f); f != "" {
				res.Header.Del(f)
			}
		}
	}

	for _, h := range hopHeaders {
		res.Header.Del(h)
	}

	if lHandler.ModifyResponse != nil {
		if err := lHandler.ModifyResponse(res); err != nil {
			self.Logger.Err("http: proxy error: %v", err)
			rw.WriteHeader(http.StatusBadGateway)
			return
		}
	}

	copyHeader(rw.Header(), res.Header)

	// The "Trailer" header isn't included in the Transport's response,
	// at least for *http.Transport. Build it up from Trailer.
	if len(res.Trailer) > 0 {
		trailerKeys := make([]string, 0, len(res.Trailer))
		for k := range res.Trailer {
			trailerKeys = append(trailerKeys, k)
		}
		rw.Header().Add("Trailer", strings.Join(trailerKeys, ", "))
	}

	rw.WriteHeader(res.StatusCode)
	if len(res.Trailer) > 0 {
		// Force chunking if we saw a response trailer.
		// This prevents net/http from calculating the length for short
		// bodies and adding a Content-Length.
		if fl, ok := rw.ResponseWriter.(http.Flusher); ok {
			fl.Flush()
		}
	}
	lHandler.copyResponse(rw, res.Body)
	res.Body.Close() // close now, instead of defer, to populate res.Trailer
	copyHeader(rw.Header(), res.Trailer)

	self.proxy_handlerPool.Put(lHandler)
}

func (self *TRouter) ShowRoute(sw bool) {
	self.show_route = sw
}

/*
// 返回 TRoute 的拷贝
//循环匹配URL
func (self *TRouter) matchRoute(aPath string) TRoute {
	for _, v := range self.Routes {
		if v.regexp.MatchString(aPath) { //正则表达式比对
			return v
		}
		continue
	}
	return TRoute{} // 空值
}
*/
/* 废弃
// Router提供Templete渲染的基本函数,支持String|File参数.
func (self *TRouter) RenderTemplate(html string, w http.ResponseWriter, data interface{}) {
	if data != nil {
		MergeMaps(self.GVar, data.(map[string]interface{})) // 添加Router的全局变量到Templete
		xtemplate.RenderToResponse(html, w, self.I18N.GetLangMap(), data)

	} else {
		xtemplate.RenderToResponse(html, w, self.I18N.GetLangMap(), self.GVar)
	}
}
*/

/*
func (self *TRouter) Listen() {
	//mux := http.NewServeMux()
	//mux.Handle("/", MainRouteServer) //匹配到所有
	//self.Logger.Printf("webgo serving %s\n", addr)

	var err error
	self.Sock, err = net.Listen("tcp", ":"+self.Config.Port) //服务器sock
	if err != nil {
		self.Logger.Fatal("Router.Listen:", err)
	}

	//Router := NewRouter()
	// Serve 会对每次请求,使用Router接口ServeHTTP调用处理对应的Route携带的Handler
	err = http.Serve(self.Sock, self) // 接受一个Sock和一个带接口ServeHTTP的Router,handler..
	if err != nil {
		self.Logger.Fatal("Router.Listen:", err)
	}

	self.Close() //self.Sock.Close()
}

//Stops the web server
func (self *TRouter) Close() {
	if self.Sock != nil {
		self.Logger.Println("Web Close")
		self.Sock.Close()
	}
}
*/
