//go:build !appengine && go1.21

package base62

import (
	"unsafe"
)

func byteSliceToString(src []byte) string {
	return unsafe.String(unsafe.SliceData(src), len(src))
}

func stringToByteSlice(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
