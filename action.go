package web

/*

提供单handler 多Action 顺序/并非执行

*/
import (
	"sync"
)

type (
	TAction struct {
		listeners map[string][]func(interface{}, func(bool))
		lock      *sync.RWMutex
	}
)

var Events *TAction = NewAction()

func NewAction() *TAction {
	return &TAction{
		listeners: make(map[string][]func(interface{}, func(bool))),
		lock:      new(sync.RWMutex),
	}
}

//并发执行事件
func GoEvent(eventName string, session interface{}, next func(bool)) {
	Events.GoExecute(eventName, session, next)
}

//顺序执行事件
func Event(eventName string, session interface{}, next func(bool)) {
	Events.Execute(eventName, session, next)
}

//添加事件
func AddEvent(eventName string, handler func(interface{}, func(bool))) {
	Events.Register(eventName, handler)
}

/*
注册事件
[Examle:]
Events.Register("AfterResponse", func(session interface{}, next func(bool)) {
	log.Println("Got AfterResponse event!")
	isSuccess := true
	next(isSuccess) //这里的next函数无论什么情况下必须执行。
})
采用不同的方式执行事件时，此处的next函数的作用也是不同的：
1、在并发执行事件的时候，next函数的作用是通知程序我已经执行完了(不理会这一步是否执行成功)；
2、在顺序执行事件的时候，next函数的作用是通知程序是否继续执行下一步，next(true)是继续执行下一步，next(false)是终止执行下一步
*/
func (self *TAction) Register(eventName string, handler func(interface{}, func(bool))) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.listeners == nil {
		self.listeners = make(map[string][]func(interface{}, func(bool)))
	}
	_, ok := self.listeners[eventName]
	if !ok {
		self.listeners[eventName] = make([]func(interface{}, func(bool)), 0)
	}
	self.listeners[eventName] = append(self.listeners[eventName], handler)
}

/*
并发执行事件
[Examle 1:]
Events.GoExecute("AfterHandler", session, func(_ bool) {//此匿名函数在本事件的最后执行
	session.Response.Send()
	session.Response.Close()
})
[Examle 2:]
Events.Execute("AfterResponse", session, func(_ bool) {})
*/
func (self *TAction) GoExecute(eventName string, session interface{}, next func(bool)) {
	if self.listeners == nil {
		next(true)
		return
	}
	c := make(chan int)
	n := 0
	self.lock.RLock()
	defer self.lock.RUnlock()
	if l, ok := self.listeners[eventName]; ok {
		if len(l) > 0 {
			for _, h := range l {
				n++
				//h 的原型为 func(interface{}, func(bool))
				go h(session, func(_ bool) {
					c <- 1
				})
			}
		}
	}
	for n > 0 {
		i := <-c
		if i == 1 {
			n--
		}
	}
	next(true)
}

/**
 * 顺序执行事件
 */
func (self *TAction) Execute(eventName string, session interface{}, next func(bool)) {
	if self.listeners == nil {
		next(true)
		return
	}
	self.lock.RLock()
	defer self.lock.RUnlock()
	var nextStep bool = false
	if l, ok := self.listeners[eventName]; ok {
		if len(l) > 0 {
			for _, h := range l {
				h(session, func(ok bool) {
					nextStep = ok
				})
				//一旦传入false，后面的全部忽略执行
				if !nextStep {
					break
				}
			}
		} else {
			nextStep = true
		}
	} else {
		nextStep = true
	}
	next(nextStep)
}
