package web

import (
	"reflect"
	"sync"
)

/*
	本池不会自动生成任何东西只暂存保管
*/

type (
	TPool struct {
		pools map[reflect.Type]sync.Pool
	}
)

func NewPool() *TPool {
	return &TPool{
		pools: make(map[reflect.Type]sync.Pool), //@@@ 改进改为接口 String
	}
}

//@@@ 改进改为接口
// Resul:nil 当取不到时直接返回Nil 方便外部判断
// TODO:优化速度
func (self *TPool) Get(object reflect.Type) (val reflect.Value) {
	if object == nil {
		return
	}

	if pool, ok := self.pools[object]; ok {
		itf := pool.Get()
		if itf == nil {
			return reflect.New(object).Elem()
		}

		return itf.(reflect.Value)
	} else {
		return reflect.New(object).Elem()
	}
}

func (self *TPool) Put(typ reflect.Type, val reflect.Value) {
	if pool, ok := self.pools[typ]; ok {
		pool.Put(val)
	} else {
		var lPool sync.Pool
		lPool.Put(val)
		self.pools[typ] = lPool
	}
}
