//go:build !appengine

package base62

import (
	"reflect"
	"unsafe"
)

func byteSliceToString(src []byte) string {
	return *(*string)(unsafe.Pointer(&src))
}

func stringToByteSlice(s string) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(
		(*reflect.StringHeader)(unsafe.Pointer(&s)).Data)), len(s))
}
