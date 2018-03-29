package web

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"vectors/utils"
)

/*
	Config 负责设置服务器
*/

type (
	TConfig struct {
		LoggerLevel int    `json:"LOGGER_LEVEL"` // 日志等级
		FilePath    string //设置文件的路径
		LastModTime int64  //最后修改时间

		RootPath       string // 服务器硬盘地址
		ModulesDir     string `json:"ModulesDir"` //模块,程序块目录
		TemplatesDir   string `json:"TemplatesDir"`
		StaticDir      string `json:"StaticDir"`
		CssDir         string `json:"CssDir"`
		JsDir          string `json:"JsDir"`
		ImgDir         string `json:"ImgDir"`
		Addr           string `json:"Addr"` //端口
		Port           int64  `json:"PORT"` //端口
		CookieSecret   string
		RecoverPanic   bool
		UseSessions    bool   `json:"USE_SESSIONS"` //国际化
		UseMiddlewares bool   `json:"USE_MIDDLEWARES"`
		UseI18N        bool   `json:"USE_I18N"` //国际化
		LocaleDir      string `json:"LOCALE_DIR"`
		LangCode       string `json:"LANGUAGE_CODE"`
		EnableAdmin    bool   // flag of enable admin module to log every request info.

		DefaultDateFormat     string `json:DEFAULT_DATE_FORMAT`
		DefaultDateTimeFormat string `json:DEFAULT_DATETIME_FORMAT`

		PrintRouterTree bool
	}
)

const (
	DEFAULT_DATE_FORMAT     = "2006-01-02"
	DEFAULT_DATETIME_FORMAT = "2006-01-02 15:04:05"
	ErrNotFound             = "Page Not Found"
	ErrAccessDenied         = "Access Denied"
	ErrInternalServerError  = "Internal Server Error"
)

var (
	//memstats   = new(runtime.MemStats)
	//HandleType reflect.Type //
	//APP_FILE string

	SERVER_ROOT, _ = os.Getwd()
	STATIC_ROOT    = path.Join(SERVER_ROOT, "/static")   // 静态文件物理路径
	TEMPLATES_ROOT = path.Join(SERVER_ROOT, "/template") // 模板路径
	ModulePath     = path.Join(SERVER_ROOT, "/module")   // 模板路径

	// App settings.
	AppVer         string            // #程序版本
	AppName        string            // #名称
	AppUrl         string            //
	AppSubUrl      string            //
	AppPath        string            // #程序文件夹
	AppFilePath    string            // #程序绝对路径
	AppDir         string            // # 文件夹名称
	AppModuleDir   string = "module" // # 模块文件夹名称
	AppDataDir     string = "data"
	AppStaticDir   string = "static"
	AppTemplateDir string = "template"
)

func init() {
	AppFilePath = utils.AppFilePath()
	AppPath = filepath.Dir(AppFilePath)
	AppDir = filepath.Base(AppPath)
}

// 新建一个配置类
func NewConfig(file string) *TConfig {
	config := &TConfig{
		LoggerLevel: 0,

		DefaultDateFormat:     DEFAULT_DATE_FORMAT,
		DefaultDateTimeFormat: DEFAULT_DATETIME_FORMAT,

		RootPath:        SERVER_ROOT + "/", // 服务器路径
		ModulesDir:      "module",          // 原始模块
		TemplatesDir:    "template",        //静态文件地址modules/templates
		StaticDir:       "static",          //静态文件地址modules/static
		CssDir:          "css",             //静态文件地址modules/Static/css
		JsDir:           "js",              //静态文件地址modules/Static/js
		ImgDir:          "img",             //静态文件地址modules/Static/img
		Port:            0,                 //端口
		RecoverPanic:    true,
		UseSessions:     false,
		UseMiddlewares:  true,
		UseI18N:         true,
		EnableAdmin:     false,
		PrintRouterTree: false,
		LocaleDir:       "locale",
		LangCode:        "zh-cn"} //必须跟紧

	config.LoadFromFile(file)
	return config
}

func (self *TConfig) Init() {
}

func (self *TConfig) LoadFromFile(filepath string) {
	if filepath == "" {
		return
	}
	data, err := ioutil.ReadFile(filepath) //打开文件
	if err != nil {
		panic(err)
	}

	data = regexp.MustCompile(`#.*\n`).ReplaceAll(data, []byte("\n")) //替换掉注释

	err = json.Unmarshal(data, self) // 赋值到self里Json tag对应的变量
	if err != nil {
		panic(err)
	}
	/*
		self.UploadDirectory = c.StaticDirectory + c.UploadDirectory
			c.ThemeDirectory = c.ThemeDirectory + c.Theme + "/"
			c.StaticCssDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticCssDirectory
			c.StaticJsDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticJsDirectory
			c.StaticImgDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticImgDirectory
	*/
	self.FilePath = filepath
	fileinfo, _ := os.Stat(filepath)             //获取文件信息
	self.LastModTime = fileinfo.ModTime().Unix() //保存最后修改时间

}

func (self *TConfig) Reload() bool {
	var Result bool = false
	fileinfo, _ := os.Stat(self.FilePath) //获取文件信息
	if fileinfo.ModTime().Unix() > self.LastModTime {
		self.LoadFromFile(self.FilePath)
		/*
			data := c.format(self.FilePath)
			*c = NewConfig()
			json.Unmarshal(data, c)
			c.configPath = configPath
			c.configLastModTime = dataFi.ModTime().Unix()
			c.UploadDirectory = c.StaticDirectory + c.UploadDirectory
			c.ThemeDirectory = c.ThemeDirectory + c.Theme + "/"
			c.StaticCssDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticCssDirectory
			c.StaticJsDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticJsDirectory
			c.StaticImgDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticImgDirectory
		*/

		Result = true
	}

	return Result
}
