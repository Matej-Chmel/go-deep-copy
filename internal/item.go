package internal

import r "reflect"

type Item struct {
	Flag   int
	Flag2  int
	Field  *r.Value
	Field2 *r.Value
	Keys   []r.Value
	Ix     int
	NewVal *r.Value
	Val    *r.Value
}

const (
	None = iota
	KeyNext
	Result
	StructData
	ValNext
)

func NewItem(flag int, field *r.Value, ix int, val *r.Value) *Item {
	return &Item{
		Flag:   flag,
		Flag2:  None,
		Field:  field,
		Field2: nil,
		Keys:   nil,
		NewVal: nil,
		Ix:     ix,
		Val:    val,
	}
}
