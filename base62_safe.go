//go:build appengine

package base62

func byteSliceToString(src []byte) string {
	return string(src)
}

func stringToByteSlice(s string) []byte {
	return []byte(s)
}
