package godeepcopy_test

import (
	"fmt"
	"runtime"
	"testing"
	"unsafe"

	gd "github.com/Matej-Chmel/go-deep-copy"
)

type Example struct {
	A int
	B string
	C rune
}

type NestedExample struct {
	Example
	B    string
	Next *NestedExample
}

type NestedUnexported struct {
	Unexported
	C rune
}

type SliceExample struct {
	Bytes []uint16
	Ints  []int
}

// Wrapper around original test type
type tester struct {
	failed bool
	*testing.T
}

// Constructs new tester
func newTester(t *testing.T) *tester {
	return &tester{failed: false, T: t}
}

// Fail test with a formatted message containing the line number
// if it's available
func (t *tester) fail(skip int, format string, data ...any) {
	_, _, line, ok := runtime.Caller(skip)
	reason := fmt.Sprintf(format, data...)

	if ok {
		format = fmt.Sprintf("(line %d) %s", line, reason)
	}

	t.Errorf(format, reason)
	t.failed = true
}

type Unexported struct {
	a int
	B string
}

func boolStr(b bool) string {
	if b {
		return ""
	}

	return "NOT"
}

func check[T any](data T, t *tester) {
	checkImpl(func(d T) string {
		return fmt.Sprintf("%v", d)
	}, data, t)
}

// Create a deep copy of data and check whether
// copy and original live on different memory addresses
// and their string representations match
func checkImpl[T any](conv func(T) string, data T, t *tester) {
	if t.failed {
		return
	}

	dataCopy := gd.DeepCopy(data)

	if gd.IsCopyable(data) {
		copyAddr := unsafe.Pointer(&dataCopy)
		dataAddr := unsafe.Pointer(&data)

		if dataAddr == copyAddr {
			t.fail(3, "Same address %p", dataAddr)
			return
		}
	}

	actual := conv(dataCopy)
	expected := conv(data)

	if actual != expected {
		t.fail(3, "%s != %s", actual, expected)
	}
}

// Create pointer to data and check deep copy of that pointer
func checkPointer[T any](data T, t *tester) {
	checkImpl(func(d *T) string {
		return fmt.Sprintf("%v", *d)
	}, &data, t)
}

func TestArrays(ot *testing.T) {
	t := newTester(ot)
	check([]byte{65, 66, 67}, t)
	check([]float64{65.3, 66.4, 67.5}, t)
	check([]int{12, 34, 56}, t)
	check([]rune{'A', 'B', 'C'}, t)
	check([...]string{"hello", "world"}, t)
}

func TestBasicTypes(ot *testing.T) {
	t := newTester(ot)
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

func TestPointer(ot *testing.T) {
	t := newTester(ot)
	checkPointer(false, t)
	checkPointer(true, t)

	checkPointer(make(chan int), t)

	checkPointer(1+1i, t)
	checkPointer(2.4567+3.45678i, t)

	checkPointer(float64(12.3456), t)
	checkPointer(128.993456, t)

	checkPointer(int8(-128), t)
	checkPointer(int16(-32768), t)
	checkPointer(int32(-65536), t)
	checkPointer(int64(-128000), t)

	checkPointer("hello world", t)
	checkPointer("hi 123", t)

	checkPointer(uint8(255), t)
	checkPointer(uint16(65535), t)
	checkPointer(uint32(128000), t)
	checkPointer(uint64(10030030030303333333), t)

	checkPointer(uintptr(0x12345678), t)
	checkPointer(unsafe.Pointer(uintptr(0x45678902)), t)
}

func TestStructs(ot *testing.T) {
	t := newTester(ot)
	example := Example{12, "hello", '*'}
	nested := NestedExample{Example{34, "world", '%'}, "super", nil}
	slice := SliceExample{[]uint16{40, 20}, []int{1, 2, 3}}

	check(example, t)
	checkPointer(example, t)
	check(nested, t)
	checkPointer(nested, t)
	check(slice, t)
	checkPointer(slice, t)
}
