package peer

import (
	"reflect"
	"sync"
)

type Invoker struct {
	sync.Pool
	callback []func([]reflect.Value) error
}

func (i *Invoker) Apply(d Driver, v interface{}, s string) error {
	param := []reflect.Value{
		reflect.ValueOf(v),
		reflect.ValueOf(d),
		reflect.ValueOf(s),
	}
	for _, f := range i.callback {
		if err := f(param[:]); nil != err {
			return err
		}
	}
	return nil
}

func (i *Invoker) Append(fn ...interface{}) *Invoker {
	for _, f := range fn {
		i.callback = append(i.callback, wrapper(f))
	}
	return i
}

func (i *Invoker) __append(fn ...reflect.Value) *Invoker {
	for _, f := range fn {
		i.callback = append(i.callback, __wrapper(f))
	}
	return i
}

func wrapper(fn interface{}) func([]reflect.Value) error {
	return __wrapper(reflect.ValueOf(fn))
}

func __wrapper(fn reflect.Value) func([]reflect.Value) error {
	proxy := __proxy(fn.Type())
	return func(param []reflect.Value) error {
		return proxy(param, fn.Call(param[:fn.Type().NumIn()]))
	}
}

func __proxy(fn reflect.Type) func([]reflect.Value, []reflect.Value) error {
	out := fn.NumOut()
	if 1 == out {
		return func(param, result []reflect.Value) error {
			if result[0].IsValid() && !result[0].IsNil() {
				return result[0].Interface().(error)
			}
			return nil
		}
	}
	if 3 == out {
		return func(param, result []reflect.Value) error {
			if result[2].IsValid() && !result[2].IsNil() {
				return result[2].Interface().(error)
			}
			if name := result[0].String(); "" != name {
				param[1].Interface().(Driver).EmitTo(param[2].String(), name, result[1].Interface())
			}
			return nil
		}
	}
	return func([]reflect.Value, []reflect.Value) error {
		return nil
	}
}
