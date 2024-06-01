package godeepcopy

import (
	r "reflect"
	"unsafe"

	gs "github.com/Matej-Chmel/go-generic-stack"
)

type processor struct {
	stack gs.Stack[*item]
}

func newProcessor(v *r.Value) processor {
	p := processor{stack: gs.Stack[*item]{}}
	p.push(none, nil, 0, v)
	return p
}

func newPtr(val *r.Value) *r.Value {
	n := r.New(val.Type())
	return &n
}

func newVal(val *r.Value) *r.Value {
	n := newPtr(val).Elem()
	return &n
}

func (p *processor) finalize(it *item) {
	if p.stack.Len() == 1 {
		it.flag = result
	} else if p.stack.Len() > 1 {
		p.stack.Pop()

		top, _ := p.stack.Top()

		if it.flag2 == valNext {
			top.field2 = it.newVal
		} else {
			top.field = it.newVal
		}
	}
}

func (p *processor) push(flag int, field *r.Value, ix int, val *r.Value) {
	p.stack.Push(newItem(flag, field, ix, val))
}

func (p *processor) result() *r.Value {
	top, err := p.stack.Pop()

	if err != nil {
		return nil
	}

	return top.newVal
}

func (p *processor) run() {
	for p.stack.HasItems() {
		top, _ := p.stack.Top()

		if top.flag == result {
			break
		}

		p.processItem(top)
	}
}

func (p *processor) processItem(it *item) {
	kind := it.val.Kind()

	switch kind {
	case r.Array:
		p.processArray(it, true)
		return

	case r.Map:
		p.processMap(it)
		return

	case r.Pointer:
		p.processPtr(it)
		return

	case r.Slice:
		p.processArray(it, false)
		return

	case r.Struct:
		p.processStruct(it)
		return
	}

	it.newVal = newVal(it.val)

	switch kind {
	case r.Bool:
		p.processBool(it)
	case r.Complex64, r.Complex128:
		p.processComplex(it)
	case r.Float32, r.Float64:
		p.processFloat(it)
	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
		p.processInt(it)
	case r.String:
		p.processString(it)
	case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
		p.processUint(it)
	case r.Uintptr:
		p.processUintptr(it)
	case r.UnsafePointer:
		p.processUnsafePtr(it)
	default:
		it.newVal = it.val
	}

	p.finalize(it)
}

func (p *processor) processArray(it *item, isArray bool) {
	if it.flag == none {
		var n r.Value

		if isArray {
			elem := r.ArrayOf(it.val.Len(), it.val.Type().Elem())
			n = r.New(elem).Elem()
		} else {
			elem := r.SliceOf(it.val.Type().Elem())
			n = r.MakeSlice(elem, it.val.Len(), it.val.Cap())
		}

		it.newVal = &n
		it.flag = valNext
	} else if it.flag != valNext {
		return
	}

	if it.ix == it.val.Len() {
		p.finalize(it)
		return
	}

	if it.field == nil {
		field := it.val.Index(it.ix)
		p.push(none, nil, 0, &field)
	} else {
		it.newVal.Index(it.ix).Set(*it.field)
		it.field = nil
		it.ix++
	}
}

func (p *processor) processBool(it *item) {
	it.newVal.SetBool(it.val.Bool())
}

func (p *processor) processComplex(it *item) {
	it.newVal.SetComplex(it.val.Complex())
}

func (p *processor) processFloat(it *item) {
	it.newVal.SetFloat(it.val.Float())
}

func (p *processor) processInt(it *item) {
	it.newVal.SetInt(it.val.Int())
}

func (p *processor) processMap(it *item) {
	if it.flag == none {
		it.keys = it.val.MapKeys()
		it.newVal = newVal(it.val)
		it.flag = keyNext
	}

	if it.ix == it.val.Len() {
		p.finalize(it)
		return
	}

	if it.flag == keyNext {
		if it.field == nil {
			p.push(none, nil, 0, &it.keys[it.ix])
		} else {
			it.flag = valNext
		}
	} else if it.flag == valNext {
		if it.field2 == nil {
			v := it.val.MapIndex(it.keys[it.ix])
			n := newItem(none, nil, 0, &v)
			n.flag2 = valNext
		} else {
			it.newVal.SetMapIndex(*it.field, *it.field2)
			it.flag = keyNext
			it.ix++
		}
	}
}

func (p *processor) processPtr(it *item) {
	if it.flag == none {
		it.newVal = newPtr(it.val)
		p.push(none, nil, 0, it.val)
		it.flag = valNext
		return
	} else if it.flag == valNext && it.field != nil {
		it.newVal.Elem().Set(*it.field)
	}
}

func (p *processor) processUintptr(it *item) {
	it.newVal.SetUint(it.val.Uint())
}

func (p *processor) processUnsafePtr(it *item) {
	v := it.val.Pointer()
	it.newVal.SetPointer(unsafe.Pointer(v))
}

func (p *processor) processString(it *item) {
	it.newVal.SetString(it.val.String())
}

func (p *processor) processStruct(it *item) {

}

func (p *processor) processUint(it *item) {
	it.newVal.SetUint(it.val.Uint())
}
