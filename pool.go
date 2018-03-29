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

/*
func (self *TPool) new(aType reflect.Type) interface{} {
	return self.pools[aType].New= func() interface{} {
		lActionVal := reflect.New(aType).Elem() //由类生成实体值,必须指针转换而成才是Addressable  错误：lVal := reflect.Zero(aHandleType)

		return lActionVal
	}
}
*/

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

	/*
		var lPoolName string

		lType := reflect.TypeOf(aName)
		//Debug("get", lType.Name(), lType.String(), lType.Kind())
		switch lType.Kind() {
		case reflect.Ptr: // 接受reflect.rtype，指针类型
			{
				lType = lType.Elem()

				// 如果是指针 可能是结构体/结构体指针 *struct/struct
				if lType.String() == "reflect.rtype" {
					lType = aName.(reflect.Type) //指针类型转换成类型

					// 接口reflect.rtype 可能是 指针/结构
					if lType.Kind() == reflect.Ptr {
						lType = lType.Elem()
					}
				}
				//Debug("getg", lType.Name(), lType.String(), lType.Kind())

				lPoolName = lType.String()
			}

		case reflect.Struct: // 接受 结构值
			{
				lPoolName = lType.String()
			}
		case reflect.String:
			{
				lPoolName = aName.(string)
			}
		}
	*/
	//Debug("lPoolName", lPoolName)
	if pool, ok := self.pools[object]; ok {
		itf := pool.Get()
		if itf == nil {
			return reflect.New(object).Elem()
		}

		//Debug("TPool.Get", ok, aType, len(self.pools))
		return itf.(reflect.Value)
	} else {
		//Debug("TPool.Get", ok, aType, len(self.pools))
		return reflect.New(object).Elem()
	}

}

func (self *TPool) Put(typ reflect.Type, val reflect.Value) {
	//if typ. == nil || val == nil {
	//	return
	//}
	/*
		var (
			lType     reflect.Type
			lVal      reflect.Value
			lPoolName string
		)

			lType = reflect.TypeOf(x)

			switch lType.String() {
			case "reflect.Value":
				{
					lVal = x.(reflect.Value)
					if lVal.IsValid() {
						lType = lVal.Type()
					}
					//Debug("TPool.Put", lType.Kind(), lType.Name(), lType.String(), lVal.Type(), lVal.Kind(), lVal.String())

				}
			default:
			}
			//lType2 := reflect.ValueOf(x)

			if lType.Kind() == reflect.Ptr {
				lType = lType.Elem() //指针类型转换成类型
			}
			lPoolName = lType.String()
	*/

	//Debug("TPool.Put", lType.Kind(), lType.Name())
	if pool, ok := self.pools[typ]; ok {
		//Debug("TPool.Put", pool, ok)
		pool.Put(val)
	} else {
		var lPool sync.Pool
		lPool.Put(val)
		self.pools[typ] = lPool
		//Debug("TPool.Put", lPool, x, lType, len(self.pools))
	}

}
