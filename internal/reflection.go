package internal

import r "reflect"

// Create a new pointer to a new value of the same type as val
func NewPointer(val *r.Value) *r.Value {
	n := r.New(val.Type())
	return &n
}

// Create a new value of the same type as val
func newValue(val *r.Value) *r.Value {
	n := NewPointer(val).Elem()
	return &n
}
