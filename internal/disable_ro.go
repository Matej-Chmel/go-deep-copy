package internal

// Code in this file was taken from
// https://github.com/brunoga/deep/blob/92c699d4e2e304e7c4a4d4138817a8c96e8abb72/disable_ro.go
// License notice available in file brunoga_LICENSE

import (
	"reflect"
	"unsafe"
)

var flagOffset uintptr = 0

const (
	readOnlyFlag uintptr = 0x60
)

// Disable read-only mode for Value
func disableRO(v *reflect.Value) {
	// Get pointer to flags
	flags := (*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + flagOffset))

	// Clear the read-only flags
	*flags &^= readOnlyFlag
}

// Computes the offset of the read-only flag or does nothing
// if it's already computed
func initFlagOffset() {
	if flagOffset != 0 {
		return
	}

	defaultValueType := reflect.TypeOf(reflect.Value{})

	if defaultValueType.Kind() != reflect.Struct {
		panic("deep: reflect.Value is not a struct")
	}

	for i := 0; i < defaultValueType.NumField(); i++ {
		field := defaultValueType.Field(i)

		if field.Name == "flag" {
			if field.Type.Kind() != reflect.Uintptr {
				panic("deep: reflect.Value.flag is not a uintptr")
			}

			flagOffset = field.Offset
			return
		}
	}

	panic("deep: reflect.Value has no flag field")
}
