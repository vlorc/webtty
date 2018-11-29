package peer

import (
	"github.com/json-iterator/go"
	"reflect"
	"strings"
	"sync"
)

type EventDriver struct {
	table  sync.Map
	global []func(Driver, string, []byte) error
	miss   []func(Driver, string, []byte) error
}

func NewEvent() Event {
	return &EventDriver{}
}

func (e *EventDriver) On(name string, fn ...interface{}) Event {
	switch name {
	case "*":
		for _, f := range fn {
			e.global = append(e.global, f.(func(Driver, string, []byte) error))
		}
	case "#":
		for _, f := range fn {
			e.miss = append(e.miss, f.(func(Driver, string, []byte) error))
		}
	default:
		e.__on(name, fn...)
	}
	return e
}

func (e *EventDriver) With(obj ...interface{}) Event {
	v := reflect.ValueOf(obj[0])
	t := reflect.TypeOf(obj[0])

	for i, l := 0, t.NumMethod(); i < l; i++ {
		m := t.Method(i)
		if !strings.HasPrefix(m.Name, "On") {
			continue
		}
		e.__append(strings.ToUpper(m.Name[2:]), v.Method(i))
	}
	return e
}

func (e *EventDriver) Emit(driver Driver, message []byte) error {
	iter := jsoniter.ConfigFastest.BorrowIterator(message)
	defer jsoniter.ConfigFastest.ReturnIterator(iter)

	return e.__emit(driver, iter)
}

func (e *EventDriver) __append(name string, fn reflect.Value) *Invoker {
	it, _ := e.table.LoadOrStore(name, &Invoker{})
	invoker := it.(*Invoker)
	if nil == invoker.New && fn.Type().NumIn() > 0 {
		invoker.New = func() interface{} {
			return reflect.New(fn.Type().In(0).Elem()).Interface()
		}
	}
	invoker.__append(fn)
	return invoker
}

func (e *EventDriver) __on(name string, fn ...interface{}) {
	it, _ := e.table.LoadOrStore(name, &Invoker{})
	invoker := it.(*Invoker)
	if nil == invoker.New {
		for _, f := range fn {
			if reflect.TypeOf(f).NumIn() > 0 {
				invoker.New = func() interface{} {
					return reflect.New(reflect.TypeOf(f).In(0).Elem()).Interface()
				}
			}
		}
	}
	invoker.Append(fn...)
}

func (e *EventDriver) __emit(driver Driver, iter *jsoniter.Iterator) error {
	var dest string
	var source string
	var command string
	var payload []byte
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, field string) bool {
		switch field {
		case "source":
			source = iter.ReadString()
		case "dest":
			dest = iter.ReadString()
		case "command":
			command = iter.ReadString()
		case "payload":
			payload = iter.SkipAndReturnBytes()
		default:
			iter.Skip()
		}
		return true
	})
	for _, f := range e.global {
		if err := f(driver, command, payload); nil != err {
			return err
		}
	}
	if it, ok := e.table.Load(command); ok {
		return e.__call(driver, string(source), it.(*Invoker), iter.ResetBytes(payload))
	}
	for _, f := range e.miss {
		if err := f(driver, command, payload); nil != err {
			return err
		}
	}
	return nil
}

func (e *EventDriver) __call(driver Driver, source string, inv *Invoker, iter *jsoniter.Iterator) error {
	val := inv.Get()
	defer inv.Put(val)

	if iter.ReadVal(val); nil != iter.Error {
		return iter.Error
	}
	return inv.Apply(driver, val, source)
}
