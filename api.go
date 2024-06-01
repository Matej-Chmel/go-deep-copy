package godeepcopy

import r "reflect"

func DeepCopy[T any](data T) T {
	val := r.ValueOf(data)
	return DeepCopyValue[T](&val)
}

func DeepCopyValue[T any](v *r.Value) T {
	p := newProcessor(v)
	p.run()
	val := p.result()
	return val.Interface().(T)
}

func IsCopyable[T any](data T) bool {
	val := r.ValueOf(data)
	kind := val.Kind()

	switch kind {
	case r.Chan, r.Func, r.Interface, r.Invalid:
		return false
	}

	return true
}
