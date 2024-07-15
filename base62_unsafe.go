//go:build !appengine && !go1.21

package base62

// reflect.{StringHeader,SliceHeader} are deprecated
// (see https://github.com/golang/go/issues/53003)

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
