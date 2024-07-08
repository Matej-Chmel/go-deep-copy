package internal

import r "reflect"

// Work item
type item struct {
	// Flag or index into an array, map, slice or struct
	flagOrIndex int
	// Saves an order of keys in a map
	keys []r.Value
	// Original value of the work item
	val *r.Value
}

const (
	// This item wasn't processed yet
	None int = -1
	// A pointer is waiting for a copy of its underlying value
	ValNext int = -2
)

// Constructs a new work item
func newItem(flagOrIndex int, val *r.Value) *item {
	return &item{
		flagOrIndex: flagOrIndex,
		keys:        nil,
		val:         val,
	}
}

// Returns the kind of the original value
func (it *item) kind() r.Kind {
	return it.val.Kind()
}
