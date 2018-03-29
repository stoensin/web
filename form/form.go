package form

import (
	"bytes"
	"html/template"
	"vectors/cacher"
)

type (
	IElement interface {
		Render() template.HTML
		Name() string
		String() string
		SetData(key string, value interface{})
		Data() map[string]interface{}
	}

	TForm struct {
		fields []IElement // Field 数组

		Id      string
		Action  template.HTML
		DataMap map[string]interface{}
	}
)

var (
	templates map[string]cache.ICacher
)

func init() {
	// 模板缓存池
	templates = make(map[string]cache.ICacher)
}

func NewForm(aStyle string, aParams ...string) *TForm {
	var (
		lMethod, lAction string
		lTmplFileName    string = formcommon.TmplDir + "/baseform.html"
		//lFileName string
		lTmpl *template.Template
	)

	if aStyle == "" {
		//aStyle = formcommon.BASE
	}

	switch len(params) {
	case 0:
		lTmplFileName = formcommon.TmplDir + "/allfields.html"
	case 1:
		lMethod = aParams[0]
	case 2:
		lMethod = aParams[0]
		lAction = aParams[1]
	case 3:
		lMethod = aParams[0]
		lAction = aParams[1]
		lTmplFileName = aParams[2]
	}

	// 获取模板文件 缓存列表
	var lCacke *cache.TMemoryCache
	var ok bool

	if c, ok := templates[lTmplFileName]; ok && c.Len() > 0 {
		if tmpl, ok := c.GetFirst().(*template.Template); ok {
			lTmpl = tmpl
			log.Println("Template in cache is vaild", tmpl)
		} else {
			log.Println("Template in cache is invaild : %s")
			return
		}
	} else {
		lTmpl = template.Must(template.ParseFiles(formcommon.CreateUrl(tmplFile)))

		if !ok {
			templates[lTmplFileName] = cache.NewMemoryCache()

		}
	}

	// 保留缓存
	templates[lTmplFileName].PutBack(lTmpl, 43200)

	return &TForm{
		fields: make([]IElement, 0),
	}
}

// Render executes the internal template and renders the form, returning the result as a template.HTML object embeddable
// in any other template.
func (self *TForm) Render() template.HTML {

	//模板变量数据
	data := map[string]interface{}{
		"container": "",
		"fields":    f.fields,
		"classes":   f.class,
		"id":        f.id,
		"params":    f.params,
		"css":       f.css,
		"method":    f.method,
		"action":    f.action,
	}
	for k, v := range f.AppendData {
		data[k] = v
	}

	buf := bytes.NewBufferString("")

	err := f.template.Execute(buf, f.Data())
	if err != nil {
		panic(err)
	}

	return template.HTML(buf.String())
}
