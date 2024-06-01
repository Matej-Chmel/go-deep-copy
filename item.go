package godeepcopy

import "reflect"

type item struct {
	flag   int
	flag2  int
	field  *reflect.Value
	field2 *reflect.Value
	keys   []reflect.Value
	ix     int
	newVal *reflect.Value
	val    *reflect.Value
}

const (
	none = iota
	keyNext
	result
	structData
	valNext
)

func newItem(flag int, field *reflect.Value, ix int, val *reflect.Value) *item {
	return &item{
		flag:   flag,
		flag2:  none,
		field:  field,
		field2: nil,
		keys:   nil,
		newVal: nil,
		ix:     ix,
		val:    val,
	}
}
