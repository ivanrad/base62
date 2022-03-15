//go:build go1.18
// +build go1.18

package base62

import (
	"bytes"
	"testing"
)

func FuzzEncodeDecode(f *testing.F) {
	for _, seed := range [][]byte{
		{},
		{0},
		{0xff, 0xfe},
		{0, 1, 2, 3, 4, 5},
		[]byte("\x53\xfe\x92\xfe\xff\x00\xab"),
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		t.Helper()
		encoded := EncodeToString(b)
		decoded, err := DecodeString(encoded)
		if err != nil {
			t.Fatalf("Error while decoding: %v", err)
		}
		if len(decoded) != len(b) {
			t.Fatalf("Decoded length: %v; want: %v", len(decoded), len(b))
		}
		if !bytes.Equal(b, decoded) {
			for i := 0; i < len(b); i++ {
				if b[i] != decoded[i] {
					t.Errorf("Decoding failed at position: %v", i)
				}
			}
		}
	})
}
