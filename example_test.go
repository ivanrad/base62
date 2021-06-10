package base62_test

import (
	"fmt"

	"github.com/ivanrad/base62"
)

func ExampleEncodeToString() {
	msg := []byte("hello, world")
	encoded := base62.EncodeToString(msg)
	fmt.Printf("%q\n", encoded)
	// Output:
	// "aGVsbG8WEDu3uTYyA"
}

func ExampleDecodeString() {
	decoded, err := base62.DecodeString("aGVsbG8WEDu3uTYyA")
	if err != nil {
		fmt.Printf("error occurred: %v", err)
		return
	}
	fmt.Printf("%q", decoded)
	// Output:
	// "hello, world"
}
