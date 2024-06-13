package godeepcopy

import (
	r "reflect"

	ite "github.com/Matej-Chmel/go-deep-copy/internal"
	gs "github.com/Matej-Chmel/go-generic-stack"
)

func DeepCopy[T any](data T) T {
	val := r.ValueOf(data)
	return DeepCopyValue[T](&val)
}

func DeepCopyValue[T any](v *r.Value) T {
	p := ite.NewProcessor(v)
	p.Run()
	val := p.Result()
	return val.Interface().(T)
}

func IsCopyable[T any](data T) bool {
	val := r.ValueOf(data)
	kind := val.Kind()

	switch kind {
	case r.Chan, r.Func, r.Interface, r.Invalid:
		return false
	case r.Struct:
		return fullyExported(&val)
	}

	return true
}

func IsFullyExported[T any](data T) bool {
	val := r.ValueOf(data)
	return fullyExported(&val)
}

func fullyExported(val *r.Value) bool {
	if !hasFields(val) {
		return true
	}

	stack := gs.Stack[*r.Value]{}
	stack.Push(val)

	for stack.HasItems() {
		top, _ := stack.Pop()
		top = makeAddressable(top)
		nf := top.NumField()

		for i := 0; i < nf; i++ {
			field := top.Field(i)

			if !field.CanInterface() {
				return false
			}

			if hasFields(&field) {
				stack.Push(&field)
			}
		}
	}

	return true
}

func hasFields(val *r.Value) bool {
	kind := val.Kind()
	return kind == r.Struct || kind == r.Pointer && val.Elem().Kind() == r.Struct
}

func makeAddressable(val *r.Value) *r.Value {
	var elem r.Value
	kind := val.Kind()

	if kind == r.Pointer {
		elem = val.Elem()
	} else {
		tmp := ite.NewPtr(val)
		tmp.Elem().Set(*val)
		elem = tmp.Elem()
	}

	return &elem
}
