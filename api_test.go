package godeepcopy_test

import (
	"runtime"
	"testing"
	"unsafe"

	ats "github.com/Matej-Chmel/go-any-to-string"
	gd "github.com/Matej-Chmel/go-deep-copy"
)

func check[T any](data T, t *testing.T) {
	dataCopy := gd.DeepCopy(data)
	actual := ats.AnyToString(dataCopy)
	expected := ats.AnyToString(data)

	if actual == expected {
		return
	}

	_, _, line, ok := runtime.Caller(1)

	if ok {
		t.Errorf("(line %d) %s != %s", line, actual, expected)
	} else {
		t.Errorf("%s != %s", actual, expected)
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
