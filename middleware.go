package web

/**
中间件控制器

TWebCtrl struct {
	event.TEvent
}

// 传递的必须为非指针的值(self TWebCtrl)
func (self TWebCtrl) Before(hd *web.THandler) {

}

*/

import (
	"fmt"
	"sync"
)

var (
	Middleware = TMiddlewareManager{}
)

type (
	IMiddleware interface {
		/*
			this will call before current ruote
			@act: the action interface which middleware bindding
			@hd: the Handler interface for controller
		*/
		Request(act interface{}, hd *THandler)

		/*
			this will call after current ruote
			@act: the action interface which middleware bindding
			@hd: the Handler interface for controller
		*/
		Response(act interface{}, hd *THandler)

		Panic(act interface{}, hd *THandler)
	}

	TMiddlewareManager struct {
		middlewares map[string]IMiddleware
		Names       []string     //
		lock        sync.RWMutex // 同步性不重要暂时不加锁
	}
)

func NewMiddlewareManager() *TMiddlewareManager {
	return &TMiddlewareManager{
		middlewares: make(map[string]IMiddleware),
	}

}

func (self *TMiddlewareManager) Contain(key string) bool {
	self.lock.RLock()
	_, ok := self.middlewares[key]
	self.lock.RUnlock()
	return ok
}

func (self *TMiddlewareManager) Add(key string, value IMiddleware) {
	self.lock.Lock()
	defer self.lock.Unlock()

	if _, exsit := self.middlewares[key]; !exsit {
		self.middlewares[key] = value
		self.Names = append(self.Names, key) // # 保存添加顺序
		//self.Names = append(self.Names, "gsdfgsf")
		//Warn("TMiddlewareManager", self.Names)
	} else {
		fmt.Println("key:" + key + " already exists")
	}
}

func (self *TMiddlewareManager) Set(key string, value IMiddleware) {
	self.lock.Lock()
	defer self.lock.Unlock()

	if _, ok := self.middlewares[key]; ok {
		self.middlewares[key] = value
	} else {
		fmt.Println("key:" + key + " does not exists")
	}
}

func (self *TMiddlewareManager) Get(key string) interface{} {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.middlewares[key]
}

func (self *TMiddlewareManager) Del(key string) {
	self.lock.Lock()
	defer self.lock.Unlock()

	delete(self.middlewares, key)
	for i, n := range self.Names {
		if n == key {
			self.Names = append(self.Names[:i], self.Names[i+1:]...)
			break
		}

	}
}

/*
// 根据Route记录的中间更新
func (self TMiddlewareManager) ProcessRequest(hd *THandler) {
	for key, val := range hd.Route.Middleware {
		if val {
			lMiddleware := self[key]
			if lMiddleware != nil {
				Warn("has lName[len(lName)-1:]", key)
				lMiddleware.Request(hd)
			}

		}
	}
}

// 中间件返回
func (self TMiddlewareManager) ProcessResponse(hd *THandler) {
		for key, val := range hd.Route.Middleware {
			if val {
				fmt.Println(self[key], hd.Route.Middleware)
				lMiddleware := self[key]
				if lMiddleware != nil {
					lMiddleware.Response(hd)
				}

			}
		}

}
*/
/*
// 中间件请求
func (self TMiddlewareManager) ProcessRequest(hd *IHandle) {
	for _, m := range self { //遍历所有中间件
		mType := reflect.TypeOf(m) //获得当前中间件的类型Type.
		//log.Println("ProcessRequest:", m, mType, mValue)           //获得当前中间件的值Value
		if method, found := mType.MethodByName("Request"); found { //获得当前中间件中是否有该方法
			mValue := reflect.ValueOf(m)
			methodType := method.Type      //获得该方法的类型Type
			hdValue := reflect.ValueOf(hd) //获得http.Request 的Value传送给方法
			//log.Println("Request1", mType, mValue, methodType, methodType.NumIn(), methodType.In(1), reqValue.Type())
			if methodType.NumIn() == 2 && methodType.In(1) == hdValue.Type() { // 目标方法不超过1个参数(NumIn()返回包括方法所属类名称)且类型是http.Request
				mValue.MethodByName("Request").Call([]reflect.Value{hdValue}) //执行方法
				//log.Println("Request3", mType, method.Func, methodType, reflect.TypeOf(m))
			}
		}
	}
}

// 中间件返回
func (self TMiddlewareManager) ProcessResponse(hd *IHandle) {
	for _, m := range self { //遍历所有中间件
		mType := reflect.TypeOf(m)                                  //获得当前中间件的类型Type.
		if method, found := mType.MethodByName("Response"); found { //获得当前中间件中是否有该方法
			mValue := reflect.ValueOf(m)
			methodType := method.Type                                          //获得该方法的类型Type
			hdValue := reflect.ValueOf(hd)                                     //获得http.ResponseWriter 的Value传送给方法
			if methodType.NumIn() == 2 && methodType.In(1) == hdValue.Type() { // 目标方法不超过1个参数(NumIn()返回包括方法所属类名称所以2)且类型是http.ResponseWriter
				mValue.MethodByName("Response").Call([]reflect.Value{hdValue}) //执行方法
			}
		}
	}
}
*/
