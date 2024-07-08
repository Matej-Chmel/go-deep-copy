package internal

import (
	r "reflect"

	gs "github.com/Matej-Chmel/go-generic-stack"
)

// Struct for creating deep copies
type CopyWriter struct {
	// Stack of staged or final copies
	products gs.Stack[*r.Value]
	// Stack of original values
	work gs.Stack[*item]
}

// Constructs
func NewCopyWriter(v *r.Value) CopyWriter {
	w := CopyWriter{
		products: gs.Stack[*r.Value]{},
		work:     gs.Stack[*item]{},
	}
	w.pushWork(None, v)
	return w
}

// Copy array or slice
func (w *CopyWriter) copyArray(it *item, isArray bool) {
	val := it.val

	if it.flagOrIndex == None {
		var n r.Value

		if elemType := val.Type().Elem(); isArray {
			// Create new array
			elem := r.ArrayOf(val.Len(), elemType)
			n = r.New(elem).Elem()
		} else {
			// Create new slice
			elem := r.SliceOf(elemType)
			n = r.MakeSlice(elem, val.Len(), val.Cap())
		}

		it.flagOrIndex = 0
		w.products.Push(&n)
	}

	if it.flagOrIndex > 0 {
		// Set top product to previous index
		prev := w.products.PopAndReturn()
		w.products.Top().Index(it.flagOrIndex - 1).Set(*prev)
	}

	if it.flagOrIndex == val.Len() {
		// End of array, pop work item
		w.work.Pop()
	} else {
		// Push element at current index onto the stack
		elem := val.Index(it.flagOrIndex)
		it.flagOrIndex++
		w.pushWork(None, &elem)
	}
}

// Copy a bool value
func (w *CopyWriter) copyBool(val *r.Value) {
	w.products.Top().SetBool(val.Bool())
}

// Attempts to copy a built-in type.
// Returns true if val is of a built-in type.
func (w *CopyWriter) copyBuiltInType(val *r.Value) {
	switch kind := val.Kind(); kind {
	case r.Bool, r.Complex64, r.Complex128, r.Float32, r.Float64,
		r.Int, r.Int8, r.Int16, r.Int32, r.Int64, r.String, r.Uint,
		r.Uint8, r.Uint16, r.Uint32, r.Uint64, r.Uintptr, r.UnsafePointer:

		// Copy all data from the original Value
		aCopy := *val
		w.products.Push(&aCopy)
		return
	}

	// This item cannot be copied,
	// product is equal to the original work item
	w.products.Push(val)
}

// Attempts to copy a composite type.
// Returns true if it.val is of a composite type.
// Composite types are Array, Map, Pointer, Slice and Struct.
func (w *CopyWriter) copyCompositeType(it *item) bool {
	switch kind := it.kind(); kind {
	case r.Array:
		w.copyArray(it, true)
	case r.Map:
		w.copyMap(it)
	case r.Pointer:
		w.copyPointer(it)
	case r.Slice:
		w.copyArray(it, false)
	case r.Struct:
		w.copyStruct(it)
	default:
		return false
	}

	return true
}

// Copy a work item
func (w *CopyWriter) copyItem(it *item) {
	if w.copyCompositeType(it) {
		return
	}

	// This item will be processed in a single pass,
	// pop it from stack
	w.work.Pop()
	w.copyBuiltInType(it.val)
}

// Copy a map
func (w *CopyWriter) copyMap(it *item) {
	val := it.val

	if it.flagOrIndex == None {
		// Push new map onto the stack
		w.products.Push(newValue(val))
		it.keys = w.products.Top().MapKeys()
		it.flagOrIndex = 0
	}

	if it.flagOrIndex > 0 {
		// Set previous key-value pair
		key := w.products.PopAndReturn()
		val := w.products.PopAndReturn()
		w.products.Top().SetMapIndex(*key, *val)
	}

	if it.flagOrIndex == val.Len() {
		// End of map, pop work item
		w.work.Pop()
	} else {
		// Push key and value onto the stack
		key := it.keys[it.flagOrIndex]
		val := val.MapIndex(key)
		it.flagOrIndex++
		w.pushWork(None, &val)
		w.pushWork(None, &key)
	}
}

// Copy a pointer
func (w *CopyWriter) copyPointer(it *item) {
	val := it.val

	if it.flagOrIndex == None {
		// This item wasn't processed yet

		if val.IsNil() {
			// Nil pointer is a new product
			product := r.Zero(val.Type())

			// No further processing needed
			w.work.Pop()
			w.products.Push(&product)
		} else {
			// Push underlying value onto the stack
			elem := val.Elem()
			it.flagOrIndex = ValNext
			w.products.Push(NewPointer(&elem))
			w.pushWork(None, &elem)
		}
	} else if it.flagOrIndex == ValNext {
		// Set underlying value to the created copy
		w.work.Pop()
		data := w.products.PopAndReturn()
		w.products.Top().Elem().Set(*data)
	}
}

// Copy a struct
func (w *CopyWriter) copyStruct(it *item) {
	val := it.val

	if it.flagOrIndex == None {
		// This item wasn't processed yet, initialize flags
		initFlagOffset()

		// Push an instance with default values onto the stack
		w.products.Push(newValue(val))
		it.flagOrIndex = 0
	}

	if it.flagOrIndex > 0 {
		// Set previous field
		prevData := w.products.PopAndReturn()
		prevField := w.products.Top().Field(it.flagOrIndex - 1)

		// Disable read-only mode to set a potentionally unexported field
		disableRO(&prevField)
		prevField.Set(*prevData)
	}

	if it.flagOrIndex == val.NumField() {
		// All fields were processed
		w.work.Pop()
		return
	}

	field := val.Field(it.flagOrIndex)

	// Disable read-only mode to read a potentionally unexported field
	disableRO(&field)
	w.pushWork(None, &field)
	it.flagOrIndex++
}

// Copy all work items and return the single remaining product
func (w *CopyWriter) CopyWork() *r.Value {
	for w.work.HasItems() {
		w.copyItem(w.work.Top())
	}

	return w.products.PopAndReturn()
}

// Push new work item on top of the stack
func (w *CopyWriter) pushWork(flagOrIndex int, val *r.Value) {
	w.work.Push(newItem(flagOrIndex, val))
}
