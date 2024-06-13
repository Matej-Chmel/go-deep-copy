package internal

import (
	r "reflect"
	"unsafe"

	gs "github.com/Matej-Chmel/go-generic-stack"
)

type Processor struct {
	stack gs.Stack[*Item]
}

func NewProcessor(v *r.Value) Processor {
	p := Processor{stack: gs.Stack[*Item]{}}
	p.push(None, v)
	return p
}

func NewPtr(val *r.Value) *r.Value {
	n := r.New(val.Type())
	return &n
}

func newVal(val *r.Value) *r.Value {
	n := NewPtr(val).Elem()
	return &n
}

func (p *Processor) finalize(it *Item) {
	if p.stack.Len() == 1 {
		it.Flag = Result
	} else if p.stack.Len() > 1 {
		p.stack.Pop()

		top, _ := p.stack.Top()

		if it.Flag2 == ValNext {
			top.Field2 = it.NewVal
		} else {
			top.Field = it.NewVal
		}
	}
}

func (p *Processor) push(flag int, val *r.Value) {
	p.stack.Push(NewItem(flag, nil, 0, val))
}

func (p *Processor) Result() *r.Value {
	top, err := p.stack.Pop()

	if err != nil {
		return nil
	}

	return top.NewVal
}

func (p *Processor) Run() {
	for p.stack.HasItems() {
		top, _ := p.stack.Top()

		if top.Flag == Result {
			break
		}

		p.processItem(top)
	}
}

func (p *Processor) processItem(it *Item) {
	if it.Flag == StructData {
		p.processStruct(it)
		return
	}

	kind := it.Val.Kind()

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

	it.NewVal = newVal(it.Val)

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
		it.NewVal = it.Val
	}

	p.finalize(it)
}

func (p *Processor) processArray(it *Item, isArray bool) {
	if it.Flag == None {
		var n r.Value

		if isArray {
			elem := r.ArrayOf(it.Val.Len(), it.Val.Type().Elem())
			n = r.New(elem).Elem()
		} else {
			elem := r.SliceOf(it.Val.Type().Elem())
			n = r.MakeSlice(elem, it.Val.Len(), it.Val.Cap())
		}

		it.NewVal = &n
		it.Flag = ValNext
	} else if it.Flag != ValNext {
		return
	}

	if it.Ix == it.Val.Len() {
		p.finalize(it)
		return
	}

	if it.Field == nil {
		field := it.Val.Index(it.Ix)
		p.push(None, &field)
	} else {
		it.NewVal.Index(it.Ix).Set(*it.Field)
		it.Field = nil
		it.Ix++
	}
}

func (p *Processor) processBool(it *Item) {
	it.NewVal.SetBool(it.Val.Bool())
}

func (p *Processor) processComplex(it *Item) {
	it.NewVal.SetComplex(it.Val.Complex())
}

func (p *Processor) processFloat(it *Item) {
	it.NewVal.SetFloat(it.Val.Float())
}

func (p *Processor) processInt(it *Item) {
	it.NewVal.SetInt(it.Val.Int())
}

func (p *Processor) processMap(it *Item) {
	if it.Flag == None {
		it.Keys = it.Val.MapKeys()
		it.NewVal = newVal(it.Val)
		it.Flag = KeyNext
	}

	if it.Ix == it.Val.Len() {
		p.finalize(it)
		return
	}

	if it.Flag == KeyNext {
		if it.Field == nil {
			p.push(None, &it.Keys[it.Ix])
		} else {
			it.Flag = ValNext
		}
	} else if it.Flag == ValNext {
		if it.Field2 == nil {
			v := it.Val.MapIndex(it.Keys[it.Ix])
			n := NewItem(None, nil, 0, &v)
			n.Flag2 = ValNext
		} else {
			it.NewVal.SetMapIndex(*it.Field, *it.Field2)
			it.Flag = KeyNext
			it.Ix++
		}
	}
}

func (p *Processor) processPtr(it *Item) {
	if it.Flag == None {
		if it.Val.Kind() == r.Struct {
			elem := it.Val.Elem()
			it.Val = &elem
			it.NewVal = newVal(it.Val)
			it.Flag = StructData
			return
		}

		elem := it.Val.Elem()
		it.NewVal = NewPtr(&elem)
		p.push(None, &elem)
		it.Flag = ValNext
	} else if it.Flag == ValNext && it.Field != nil {
		it.NewVal.Elem().Set(*it.Field)
		p.finalize(it)
	}
}

func (p *Processor) processString(it *Item) {
	it.NewVal.SetString(it.Val.String())
}

func (p *Processor) processStruct(it *Item) {
	if it.Flag != StructData {
		tmp := NewPtr(it.Val)
		tmp.Elem().Set(*it.Val)
		elem := tmp.Elem()
		it.Val = &elem
		it.NewVal = newVal(it.Val)
		it.Flag = StructData
	}

	if it.Ix == it.Val.NumField() {
		p.finalize(it)
		return
	}

	if it.Field == nil {
		if oldField := it.Val.Field(it.Ix); oldField.CanInterface() {
			p.push(None, &oldField)
		}
	} else {
		if newField := it.NewVal.Field(it.Ix); newField.CanSet() {
			newField.Set(*it.Field)
		}

		it.Field = nil
		it.Ix++
	}
}

func (p *Processor) processUint(it *Item) {
	it.NewVal.SetUint(it.Val.Uint())
}

func (p *Processor) processUintptr(it *Item) {
	it.NewVal.SetUint(it.Val.Uint())
}

func (p *Processor) processUnsafePtr(it *Item) {
	v := it.Val.Pointer()
	it.NewVal.SetPointer(unsafe.Pointer(v))
}
