package godeepcopy

import (
	"fmt"
	r "reflect"

	ite "github.com/Matej-Chmel/go-deep-copy/internal"
)

// Creates a deep copy of T. This copy will live on a new memory address
// and its every field will be unreachable from the original.
// If the copy fails, a default instance of T is returned.
func DeepCopy[T any](data T) T {
	val := r.ValueOf(data)
	aCopy, err := DeepCopyValue[T](&val)

	if err != nil {
		var defaultInstance T
		return defaultInstance
	}

	return aCopy
}

// Creates a deep copy from a reflect.Value val as a new variable of type T.
// This copy will live on a new memory address
// and its every field will be unreachable from the original.
func DeepCopyValue[T any](val *r.Value) (T, error) {
	w := ite.NewCopyWriter(val)
	aCopy := w.CopyWork()
	res, ok := aCopy.Interface().(T)

	if ok {
		return res, nil
	}

	var defaultRes T
	err := fmt.Errorf("Generic type %T doesn't match type %s",
		defaultRes, aCopy.Type().String())
	return defaultRes, err
}

// Returns true if a deep copy can be created from data.
// Only channels, functions and empty interface{} cannot be copied.
func IsCopyable[T any](data T) bool {
	val := r.ValueOf(data)
	kind := val.Kind()

	switch kind {
	case r.Chan, r.Func, r.Interface, r.Invalid:
		return false
	}

	return true
}
