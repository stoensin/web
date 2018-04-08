package web

/*
	Handler 负责把处理URL调用的过程程序
*/
import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"vectors/utils"
)

/*
  废弃 由中间件实现 保留代码
*/
type TI18n struct {
	Name            string //名称
	localeDir       string //目录
	defaultLanguage string //默认语言

	rmutex       sync.RWMutex
	mutex        sync.Mutex
	Locales      map[string]map[string]string //语言文件数据
	currentLocal string
	lastModTime  map[string]int64
}

// 创建[区域设置]
func NewI18n(name string) *TI18n {
	if name == "" {
		name = "Locale"
	}

	return &TI18n{
		Name:        name,
		Locales:     map[string]map[string]string{},
		lastModTime: map[string]int64{},
	}
}

// [区域]设置
func (self *TI18n) Init(path, lang string) error {
	if utils.DirExists(path) {
		self.localeDir = path
		self.defaultLanguage = lang

		return self.Load(self.defaultLanguage)
	} else {
		return errors.New("Dir not Exist")
	}
	return nil
}

// 载入[语言翻译文件]
func (self *TI18n) Load(lang string) error {
	self.rmutex.RLock()
	_, found := self.Locales[lang]
	oldLastModTime, mFound := self.lastModTime[lang] // 语言文件的修改时间
	self.rmutex.RUnlock()

	langFile := filepath.Join(self.localeDir, lang) // 组织路径
	newer := false
	dataFi, err := os.Stat(langFile) //检查并返回文件状态
	if err == nil {
		lastModTime := dataFi.ModTime().Unix() //获取最后修改时间
		if !mFound {                           // 如果该Lang没有修改时间 则添加
			self.rmutex.Lock()
			self.lastModTime[lang] = lastModTime
			self.rmutex.Unlock()
			newer = true
		} else { //如果有则对比更新
			if lastModTime > oldLastModTime {
				newer = true
			}
		}

		// 不是新文件 没必要更新
		if found && !newer {
			return nil
		}

		data, _ := ioutil.ReadFile(langFile) //读取文件
		m := map[string]string{}
		err = json.Unmarshal(data, &m) //转换文件到MAP[]
		if err == nil {
			self.mutex.Lock()
			self.Locales[lang] = m //更新
			self.mutex.Unlock()
		}
	}

	return err
}

func (self *TI18n) SetLocale(local string) {
	self.currentLocal = strings.ToLower(local) //通常小写
}

func (self *TI18n) Translate(aText string, aLocal ...string) string {
	var lLocal string
	if len(aLocal) == 0 {
		lLocal = self.currentLocal
	} else {
		lLocal = aLocal[0]
	}

	if ct, ok := self.Locales[lLocal]; ok {
		if v, o := ct[aText]; o {
			return v
		}
	}
	return aText
}
func (self *TI18n) GetLangMap() map[string]string {
	err := self.Load(self.currentLocal)
	if err != nil {
		self.currentLocal = self.defaultLanguage
	}

	self.rmutex.RLock()
	msgs := self.Locales[self.currentLocal]
	self.rmutex.RUnlock()

	return msgs
}
