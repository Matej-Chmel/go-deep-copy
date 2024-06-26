package godeepcopy_test

import (
	"fmt"
	"runtime"
	"testing"
	"unsafe"

	ats "github.com/Matej-Chmel/go-any-to-string"
	gd "github.com/Matej-Chmel/go-deep-copy"
)

func checkImpl[T any](data T, t *testing.T) {
	dataCopy := gd.DeepCopy(data)

	if gd.IsCopyable(data) {
		dataAddr := unsafe.Pointer(&data)
		copyAddr := unsafe.Pointer(&dataCopy)

		if dataAddr == copyAddr {
			throw(t, 3, "Same address %p", dataAddr)
		}
	}

	actual := ats.AnyToString(dataCopy)
	expected := ats.AnyToString(data)

	if actual == expected {
		return
	}

	throw(t, 3, "%s != %s", actual, expected)
}

func check[T any](data T, t *testing.T) {
	checkImpl(data, t)
}

func checkPtr[T any](data T, t *testing.T) {
	checkImpl(&data, t)
}

func throw(t *testing.T, skip int, format string, data ...any) {
	_, _, line, ok := runtime.Caller(skip)
	reason := fmt.Sprintf(format, data...)

	if ok {
		t.Errorf("(line %d) %s", line, reason)
	} else {
		t.Errorf("%s", reason)
	}
}

func TestArrays(t *testing.T) {
	check([]byte{65, 66, 67}, t)
	check([]float64{65.3, 66.4, 67.5}, t)
	check([]int{12, 34, 56}, t)
	check([]rune{'A', 'B', 'C'}, t)
	check([...]string{"hello", "world"}, t)
}

func TestBasicTypes(t *testing.T) {
	check(false, t)
	check(true, t)

	check(make(chan int), t)

	check(1+1i, t)
	check(2.4567+3.45678i, t)

	check(float64(12.3456), t)
	check(128.993456, t)

	check(int8(-128), t)
	check(int16(-32768), t)
	check(int32(-65536), t)
	check(int64(-128000), t)

	check("hello world", t)
	check("hi 123", t)

	check(uint8(255), t)
	check(uint16(65535), t)
	check(uint32(128000), t)
	check(uint64(10030030030303333333), t)

	check(uintptr(0x12345678), t)
	check(unsafe.Pointer(uintptr(0x45678902)), t)
}

func TestPtr(t *testing.T) {
	checkPtr(false, t)
	checkPtr(true, t)

	checkPtr(make(chan int), t)

	checkPtr(1+1i, t)
	checkPtr(2.4567+3.45678i, t)

	checkPtr(float64(12.3456), t)
	checkPtr(128.993456, t)

	checkPtr(int8(-128), t)
	checkPtr(int16(-32768), t)
	checkPtr(int32(-65536), t)
	checkPtr(int64(-128000), t)

	checkPtr("hello world", t)
	checkPtr("hi 123", t)

	checkPtr(uint8(255), t)
	checkPtr(uint16(65535), t)
	checkPtr(uint32(128000), t)
	checkPtr(uint64(10030030030303333333), t)

	checkPtr(uintptr(0x12345678), t)
	checkPtr(unsafe.Pointer(uintptr(0x45678902)), t)
}

type Example struct {
	A int
	B string
	C rune
}

type NestedExample struct {
	Example
	B string
	C rune
}

type SliceExample struct {
	Bytes []uint16
	Ints  []int
}

func TestStruct(t *testing.T) {
	a := Example{12, "hello", '*'}
	b := NestedExample{Example{34, "world", '%'}, "super", 'X'}
	c := SliceExample{[]uint16{40, 20}, []int{1, 2, 3}}

	check(a, t)
	checkPtr(a, t)
	check(b, t)
	checkPtr(b, t)
	check(c, t)
	checkPtr(c, t)
}

type Unexported struct {
	a int
	B string
}

func TestExport(t *testing.T) {
	e := Example{A: 1, B: "hello", C: 'a'}
	u := Unexported{a: 1, B: "hello"}

	eX := gd.IsFullyExported(e)
	eXP := gd.IsFullyExported(&e)
	uX := gd.IsFullyExported(u)
	uXP := gd.IsFullyExported(&u)

	if !eX {
		throw(t, 1, "Example should be fully exported")
	}

	if !eXP {
		throw(t, 1, "&Example should be fully exported")
	}

	if uX {
		throw(t, 1, "Unexported should NOT be fully exported")
	}

	if uXP {
		throw(t, 1, "&Unexported should NOT be fully exported")
	}
}
