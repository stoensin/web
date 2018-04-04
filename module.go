package web

/*
	Router 负责把所有的URL映射出去
	1.提供全局变量操作
	2.提供
*/

import (
	urls "net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	//"vectors/logger"
	"vectors/utils"
)

type (
	// 提供普通接口
	IModule interface {
		// 返回Module所有Routes 理论上只需被调用一次
		GetRoutes() *TTree
		GetPath() string
		GetFilePath() string
	}

	// 提供注册接口
	IModuleRegister interface {
		Register() // 注册信息到App管理器提供展示安装
	}

	// 提供注册接口
	IModuleInstaller interface {
		Install()   // - 装载套件上的Module
		Uninstall() // - 卸载套件上的Module
	}

	//模块类-模块可以在任何地方创建,并最终归集到Router里.
	TModule struct {
		Parent *TModule
		Tree   *TTree

		// attrs
		Name        string // module name <dir name> or <named>
		Summary     string
		Description string
		Author      string
		Website     string
		Category    string //# Categories can be used to filter modules in modules listing
		Version     string
		Depends     string //# any module necessary for this one to work correctly
		Demo        string //# only loaded in demonstration mode

		// unconfirm
		Path     string // URL 路径
		FilePath string // 短文件系统路径-当前文件夹名称
		// lock     sync.RWMutex

		//beforeRoute reflect.Value // 废弃动 作处理器
		//afterRoute  reflect.Value // 废弃 动作处理器

		//beforeActionRoute reflect.Value //废弃 动作处理器
		//afterActionRoute  reflect.Value //废弃 动作处理器
		//##########新特新等待优化################
		Data []string //存储注册时导入的数据文件路径

	}
)

/*    """
Setup an import-hook to be able to import OpenERP addons from the different
addons paths.

This ensures something like ``import crm`` (or even
``import openerp.addons.crm``) works even if the addons are not in the
PYTHONPATH.
"""*/

func __initialize_sys_path() {

	/*
	   global ad_paths
	   global hooked

	   dd = tools.config.addons_data_dir
	   if dd not in ad_paths:
	       ad_paths.append(dd)

	   for ad in tools.config['addons_path'].split(','):
	       ad = os.path.abspath(tools.ustr(ad.strip()))
	       if ad not in ad_paths:
	           ad_paths.append(ad)

	   # add base module path
	   base_path = os.path.abspath(os.path.join(os.path.dirname(os.path.dirname(__file__)), 'addons'))
	   if base_path not in ad_paths:
	       ad_paths.append(base_path)

	   if not hooked:
	       sys.meta_path.append(AddonsImportHook())
	       hooked = True
	*/
}

/*    """Return the path of the given module.

Search the addons paths and return the first path where the given
module is found. If downloaded is True, return the default addons
path if nothing else is found.

"""*/
func GetModulePath(module string, downloaded bool, display_warning bool) (res string) {

	// initialize_sys_path()
	// for adp in ad_paths:
	//      if os.path.exists(opj(adp, module)) or os.path.exists(opj(adp, '%s.zip' % module)):
	//         return opj(adp, module)
	res = filepath.Join(AppPath, MODULE_DIR)
	//if _, err := os.Stat(res); err == nil {
	//	return res
	//}
	return

	// if downloaded:
	//    return opj(tools.config.addons_data_dir, module)
	if display_warning {
		logger.Warn(`module %s: module not found`, module)
	}

	return ""
}

/*
   """Return the full path of a resource of the given module.

   :param module: module name
   :param list(str) args: resource path components within module

   :rtype: str
   :return: absolute path to the resource

   TODO make it available inside on osv object (self.get_resource_path)
   """*/

func GetResourcePath(module_src_path string) (res string) {
	//filepath.SplitList(module_src_path)
	mod_path := GetModulePath("", false, true)

	res = filepath.Join(mod_path, module_src_path)

	if _, err := os.Stat(res); err == nil {
		return
	}

	/*
	   if  res!=="" return False
	   resource_path = opj(mod_path, *args)
	   if os.path.isdir(mod_path):
	       # the module is a directory - ignore zip behavior
	       if os.path.exists(resource_path):
	           return resource_path
	*/
	return ""
}

// 创建[模块]
// @name 为空时所有Routes 不加任何Path 为/时添加当前目录名称为Path
func NewModule(parent *TModule, name ...string) *TModule {
	m := &TModule{
		Tree: NewRouteTree(),
		/*&TModuleNameSpace{
			set: make(map[string]*TModule),
		},*/
		Parent: parent,
		//Name:     utils.Trim(name),
		FilePath: cur_path(), // 磁盘路径
	}
	//logger.Dbg("NewModule", cur_path(), cur_dir_name())

	if len(name) > 0 {
		if utils.Trim(name[0]) == "/" {
			m.Path = cur_dir_name()
			m.Name = m.Path
		} else {
			m.Name = utils.Trim(name[0])
		}
	}

	// 组合URL路径
	if parent != nil {
		m.Name = path.Join(parent.Name, m.Name) //废弃修改
		m.Path = path.Join(parent.Path, m.Path)
	}

	return m
}

func (self *TModule) GetRoutes() *TTree {
	return self.Tree
}

func (self *TModule) GetPath() string {
	return self.Path
}

func (self *TModule) GetFilePath() string {
	return self.FilePath
}

// get current file path without file name
func cur_path() string {
	_, file, _, _ := runtime.Caller(2) // level 3
	path, _ := path.Split(file)
	return path
}

// 随文件引用层次而变
// get current file dir name
func cur_dir_name() string {
	_, file, _, _ := runtime.Caller(3) // level 3
	path, _ := path.Split(file)
	return filepath.Base(path)
}

/*
func (self *TModule) SetParent(parent *TModule) {
	self.Parent = parent

	// 组合路径
	if parent != nil {
		self.Name = path.Join(self.Parent.Name, self.Name)
	}
}
*/

/*
pos: true 为插入Before 反之After
*/
func (self *TModule) url(rote_type RouteType, aMethod []string, url string, controller interface{}, scheme string, host string) *TRoute {
	if rote_type != ProxyRoute && controller == nil {
		logger.Panic("the route must binding a controller!")
	}

	/* 有代商议废除
	lUrl := ""

	// 组合父系URL路径
	// 路径：/Self.Name or Parent.Name/Self.Name
	if self.Parent == nil {
		lUrl = self.Name
	} else {
		lUrl = strings.Trim(path.Join(self.Parent.Name, self.Name), " ")
	}

	url = utils.Trim(url)
	if url == "" {
		logger.Err("the route must have a path")
	} else if url == "/" { // 添加Index页面
		lUrl = lUrl //  相当于 /admin
	} else {
		lUrl = utils.JoinURL(lUrl, url)
	}
	*/

	// 修整Url > /+Url
	if !strings.HasPrefix(url, "/") && url != "/" {
		url = "/" + url
	}

	route := &TRoute{
		Path:     url,
		FilePath: self.FilePath,
		Model:    self.Name,
		Action:   "", //
		Type:     rote_type,
		//HookCtrl: make([]TMethodType, 0),
		Ctrls: make([]TMethodType, 0),
		//Host:     host,
		//Scheme:   scheme,

	}
	// # is it proxy route
	if scheme != "" && host != "" {
		route.Host = &urls.URL{
			Scheme: scheme,
			Host:   host,
		}
		route.isReverseProxy = true
	}

	lValueType := TMethodType{
		FuncType: reflect.TypeOf(controller)}

	//handler URL函数
	if fv, ok := controller.(reflect.Value); ok { //****得到函数参数
		lValueType.Func = fv
	} else {
		lValueType.Func = reflect.ValueOf(controller)
	}

	//route.MainCtrl = append(route.MainCtrl, lValueType)
	route.MainCtrl = lValueType
	route.Ctrls = append(route.Ctrls, route.MainCtrl)
	//logger.Dbg("url", route.MainCtrl, aMethod)
	for _, m := range aMethod {
		self.Tree.AddRoute(m, url, route)
	}

	return route
}

func (self *TModule) Get(url string, controller interface{}) *TRoute {
	return self.url(CommomRoute, []string{"GET", "HEAD"}, url, controller, "", "")
}

func (self *TModule) Post(url string, controller interface{}) *TRoute {
	/*//#1 获得Dir 名称
	if utils.Trim(self.FilePath) == "" {
		self.FilePath = utils.CurDirName()
	}*/

	return self.url(CommomRoute, []string{"POST"}, url, controller, "", "")
}

func (self *TModule) Head(url string, controller interface{}) *TRoute {
	return self.url(CommomRoute, []string{"HEAD"}, url, controller, "", "")
}

func (self *TModule) Options(url string, controller interface{}) *TRoute {
	return self.url(CommomRoute, []string{"OPTIONS"}, url, controller, "", "")
}

func (self *TModule) Trace(url string, controller interface{}) *TRoute {
	return self.url(CommomRoute, []string{"TRACE"}, url, controller, "", "")
}

func (self *TModule) Patch(url string, controller interface{}) *TRoute {
	return self.url(CommomRoute, []string{"PATCH"}, url, controller, "", "")
}

func (self *TModule) Delete(url string, controller interface{}) *TRoute {
	return self.url(CommomRoute, []string{"DELETE"}, url, controller, "", "")
}

func (self *TModule) Put(url string, controller interface{}) *TRoute {
	return self.url(CommomRoute, []string{"PUT"}, url, controller, "", "")
}

// 重组添加模块[URL]
func (self *TModule) Url(url string, controller interface{}) *TRoute {
	return self.url(CommomRoute, HttpMethods, url, controller, "", "")
}

//
// @ Proxy(nil, "/", "https", "www.baidu.com")
func (self *TModule) Proxy(methods []string, url string, scheme string, host string) *TRoute {
	if host == "" || scheme == "" || url == "" {
		logger.Panic("all the args must correct and not blank!")
	}

	if methods == nil {
		methods = HttpMethods
	}

	return self.url(ProxyRoute, methods, url, nil, scheme, host)
}

/*
  Hook 钩子
	xxx/xxx/id
	xxx/xxx/(:name)
	xxx/xxx/*
*/
func (self *TModule) HookBefore(methods []string, url string, controller interface{}) {
	self.url(HookBeforeRoute, methods, url, controller, "", "")
}

func (self *TModule) HookAfter(methods []string, url string, controller interface{}) {
	self.url(HookAfterRoute, methods, url, controller, "", "")
}
