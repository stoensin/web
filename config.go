package web

import (
	//	"os"
	//	"path"
	"path/filepath"

	"github.com/VectorsOrigin/utils"
	"github.com/go-ini/ini"
)

/*
	Config 负责服务器配置
	@使用INI作为配置格式是为了方便使用该框架的其他程序可以向同一个配置文件添加扩展选项
*/

type (
	TConfig struct {
		*ini.File `ini:"-"`
		fileName  string // 文件名称
		filePath  string // 文件路径
		//FilePath    string //设置文件的路径
		//LastModTime int64  //最后修改时间
		//RootPath       string // 服务器硬盘地址
		DebugMode             bool   `ini:"debug_mode"`
		LoggerLevel           int    `ini:"logger_level"` // 日志等级
		RecoverPanic          bool   `ini:"enabled_recover_panic"`
		PrintRouterTree       bool   `ini:"enabled_print_router_tree"`
		Host                  string `ini:"host"` //端口
		Port                  int    `ini:"port"` //端口
		EnabledTLS            bool   `ini:"enabled_tls"`
		TLSCertFile           string `ini:"tls_cert_file"`
		TLSKeyFile            string `ini:"tls_key_file"`
		CookieSecret          string
		DefaultDateFormat     string `ini:default_date_format`
		DefaultDateTimeFormat string `ini:default_date_time_format`

		/*
			ModuleDir             string `ini:"module_dir"` //模块,程序块目录
			TemplateDir           string `ini:"template_dir"`
			StaticDir             string `ini:"static_dir"`
			CssDir                string `ini:"css_dir"`
			JsDir                 string `ini:"js_dir"`
			ImgDir                string `ini:"img_dir"`
		*/
	}
)

const (
	CONFIG_FILE_NAME = "config.ini"
	MODULE_DIR       = "module" // # 模块文件夹名称
	DATA_DIR         = "data"
	STATIC_DIR       = "static"
	TEMPLATE_DIR     = "template"
	CSS_DIR          = "css"
	JS_DIR           = "js"
	IMG_DIR          = "img"
)

var (
	/*	SERVER_ROOT, _ = os.Getwd()
		STATIC_ROOT    = path.Join(SERVER_ROOT, "/static")   // 静态文件物理路径
		TEMPLATES_ROOT = path.Join(SERVER_ROOT, "/template") // 模板路径
		ModulePath     = path.Join(SERVER_ROOT, "/module")   // 模板路径
	*/

	cfg = ini.Empty()
	// 固定变量
	// App settings.
	AppVer      string // #程序版本
	AppName     string // #名称
	AppUrl      string //
	AppSubUrl   string //
	AppPath     string // #程序文件夹
	AppFilePath string // #程序绝对路径
	AppDir      string // # 文件夹名称

/*	DefaultDateFormat     string = "2006-01-02"
	DefaultDateTimeFormat string = "2006-01-02 15:04:05"
	ConfigFileName        string
*/
/*
	// debug
	DebugMode bool //  调试模式

	// logger
	LoggerLevel     int  // 日志等级
	RecoverPanic    bool // recover 时 panic
	PrintRouterTree bool
	// server
	Addr string //端口
	Port int    //端口
	// path
	ModuleDir   string //模块,程序块目录
	TemplateDir string
	StaticDir   string
	CssDir      string
	JsDir       string
	ImgDir      string

	FilePath string //设置文件的路径

	RootPath string // 服务器硬盘地址

	CookieSecret string
*/
)

func init() {
	AppFilePath = utils.AppFilePath()
	AppPath = filepath.Dir(AppFilePath)
	AppDir = filepath.Base(AppPath)
}

// 新建一个配置类
// 指定文件名时自动加载 不给名字手动加载
func NewConfig(file_name ...string) *TConfig {
	config := &TConfig{
		File:                  cfg,
		DebugMode:             false,
		LoggerLevel:           4,
		RecoverPanic:          true,
		PrintRouterTree:       true,
		Host:                  "localhost",
		Port:                  16888,
		EnabledTLS:            false,
		TLSCertFile:           "",
		TLSKeyFile:            "",
		DefaultDateFormat:     "2006-01-02",
		DefaultDateTimeFormat: "2006-01-02 15:04:05",
	}

	if len(file_name) != 0 {
		config.LoadFromFile(file_name[0])
		config.MapTo(config)
	}

	/*
			section := config.Section("logger")
			LoggerLevel = section.Key("level").MustInt(4)                  // 日志等级
			RecoverPanic = section.Key("enabled_recover_panic").MustBool() // recover 时 panic
			PrintRouterTree = section.Key("enabled_print_router_tree").MustBool()

		// path
		section := config.Section("server")
		DebugMode = section.Key("debug_mode").MustBool(false) // debug
		Addr = section.Key("addr").MustString("0.0.0.0")
		Port = section.Key("port").MustInt(16888)
		ModuleDir = section.Key("module_dir").MustString("module") //模块,程序块目录
		TemplateDir = section.Key("template_dir").MustString("template")
		StaticDir = section.Key("static_dir").MustString("static")
		CssDir = section.Key("css_dir").MustString("css")
		JsDir = section.Key("js_dir").MustString("js")
		ImgDir = section.Key("img_dir").MustString("img")
	*/
	return config
}

// 初始化
func (self *TConfig) Init() {
	if self.File == nil {
		self.LoadFromFile(CONFIG_FILE_NAME)
	}
}

func (self *TConfig) LoadFromFile(file_name string) {
	if file_name == "" {
		file_name = CONFIG_FILE_NAME
	}

	// STEP:保存数据
	self.fileName = file_name
	self.filePath = filepath.Join(AppPath, file_name)
	err := self.File.Append(self.filePath)
	if err != nil {
		logger.Err("Load ")
	}
	//self.File, err = ini.Load(self.filePath)

}

func (self *TConfig) Reload() bool {
	/*	fileinfo, _ := os.Stat(self.FilePath) //获取文件信息
		if fileinfo.ModTime().Unix() > self.LastModTime {
			self.LoadFromFile(self.FilePath)
			return true
		}
	*/
	return false
}

func (self *TConfig) Save() error {
	logger.Dbg("save", self.filePath)
	/*section := self.Section("logger")
	LoggerLevel = section.Key("level").SetValue(LoggerLevel)                   // 日志等级
	RecoverPanic = section.Key("enabled_recover_panic").SetValue(RecoverPanic) // recover 时 panic
	PrintRouterTree = section.Key("enabled_print_router_tree").SetValue(PrintRouterTree)

		// path
		section = self.Section("server")
		DebugMode = section.Key("debug_mode").SetValue(DebugMode) // debug
		Addr = section.Key("addr").SetValue(Addr)
		Port = section.Key("port").SetValue(Port)
		ModuleDir = section.Key("module_dir").SetValue(ModuleDir)
		TemplateDir = section.Key("template_dir").SetValue(TemplateDir)
		StaticDir = section.Key("static_dir").SetValue(StaticDir)
		CssDir = section.Key("css_dir").SetValue(CssDir)
		JsDir = section.Key("js_dir").SetValue(JsDir)
		ImgDir = section.Key("img_dir").SetValue(ImgDir)
	*/
	return self.SaveTo(self.filePath)
}

func (self *TConfig) FileName() string {
	return self.fileName
}

func (self *TConfig) FilePath() string {
	return self.filePath
}
