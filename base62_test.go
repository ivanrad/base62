package base62

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

// Test suite based on base64_test.go from Go standard library.
type testPair struct {
	decoded, encoded string
}

var pairs = []testPair{
	{"", ""},
	{"f", "ZC"},
	{"fo", "ZmP"},
	{"foo", "Zm83B"},
	{"foob", "Zm83sC"},
	{"fooba", "Zm83sTB"},
	{"foobar", "Zm83sTC5A"},
	{"foobare", "Zm83sTC5MF"},
	{"foobared", "Zm83sTC5MrE"},
	{"\x53\xfe\x92", "U98kC"},
	{"hello", "aGVsbGP"},
	{"simple", "c2ltcGxl"},
	{"abhor", "YWJob3C"},
	{"abhorrently", "YWJob3JyZW50bHJ"},
	{"yolk", "eW82ND"},
	{"Yorkshireism", "WW85Nbm0NLkytLm2B"},
	{"commonplaceness", "Y282tre3ODYwsbK3Mrm5B"},
	{"\xff", "9H"},
	{"\xff\xff", "999B"},
	{"\xff\xff\xff", "9999P"},
	{"\xff\xff\xff\xff", "999999D"},
	{"\xff\xff\xff\xff\xff", "9999999f"},
	{"\xff\xff\xff\xff\xff\xff", "999999999H"},
}

var bigBuf []byte

func TestCharToBitsTable(t *testing.T) {
	for i := 0; i < len(charToBitsTable); i++ {
		want := byte(strings.IndexByte(encodingAlphabet, byte(i)))
		got := charToBitsTable[i]
		if got != want {
			t.Errorf("charToBitsTable[%d] = %02x; want = %02x", i, got, want)
		}
	}
}

func TestEncodedLen(t *testing.T) {
	for _, p := range pairs {
		n := EncodedLen(len(p.decoded))
		if n < len(p.encoded) {
			t.Errorf("EncodedLen(%q) = %v; want greater or equal than %v", p.decoded,
				n, len(p.encoded))
		}
	}
}

func TestDecodedLen(t *testing.T) {
	for _, p := range pairs {
		n := DecodedLen(len(p.encoded))
		if n < len(p.decoded) {
			t.Errorf("DecodedLen(%q) = %v; want greater or equal than %v", p.encoded,
				n, len(p.decoded))
		}
	}
}

func TestEncodePairs(t *testing.T) {
	for _, p := range pairs {
		encoded := EncodeToString([]byte(p.decoded))
		if encoded != p.encoded {
			t.Errorf("EncodeToString(%q) = %q; want: %q", p.decoded, encoded,
				p.encoded)
		}
	}
}

func TestDecodePairs(t *testing.T) {
	for _, p := range pairs {
		decoded, err := DecodeString(p.encoded)
		if err != nil {
			t.Fatalf("Error while decoding(%q): %v", p.encoded, err)
		}
		if string(decoded) != p.decoded {
			t.Errorf("DecodeString(%q) = %q; wanted: %q", p.encoded, string(decoded),
				p.decoded)
		}
	}
}

func TestCorruptInput(t *testing.T) {
	testCases := []struct {
		input      string
		corruptIdx int
	}{
		{"", -1},
		{"AA", -1},
		{"AA!", 2},
		{"AAA", -1},
		{"AAAA", -1},
		{"AA=A", 2},
		{"foobar", -1},
		{"foob-r", 4},
		{"xbar", -1},
		{" bar", 0},
		{"    ", 0},
	}
	for _, tc := range testCases {
		_, err := DecodeString(tc.input)
		if tc.corruptIdx == -1 {
			if err != nil {
				t.Errorf("Corrupted input(%q): %v", tc.input, err)
			}
			continue
		}
		var inputErr InputError
		switch {
		case errors.As(err, &inputErr):
			if int64(inputErr) != int64(tc.corruptIdx) {
				t.Errorf("Corrupted input(%q) at %v; want %v", tc.input,
					int64(inputErr), tc.corruptIdx)
			}
		default:
			t.Errorf("Unexpected error %v", err)
		}
	}
}

func TestBigLen(t *testing.T) {
	encoded := EncodeToString(bigBuf)
	decoded, err := DecodeString(encoded)
	if err != nil {
		t.Fatalf("Error while decoding: %v", err)
	}
	if len(decoded) != len(bigBuf) {
		t.Fatalf("Decoded length: %v; want: %v", len(decoded), len(bigBuf))
	}
	if !bytes.Equal(bigBuf, decoded) {
		for i := 0; i < len(bigBuf); i++ {
			if bigBuf[i] != decoded[i] {
				t.Errorf("Decoding failed at position: %v", i)
			}
		}
	}
}

func TestInputTruncated(t *testing.T) {
	buf, err := DecodeString("a")
	if len(buf) != 0 {
		t.Errorf("Decoded bytes: want 0; got %v", len(buf))
	}
	if err == nil {
		t.Errorf("Expected InputError; got nil")
	}
	var inputErr InputError
	switch {
	case errors.As(err, &inputErr):
		if int64(inputErr) != -1 {
			t.Errorf("Expected input truncated InputError(-1); got %v",
				int64(inputErr))
		}
	default:
		t.Errorf("Unexpected error %v", err)
	}
}

func BenchmarkEncodePairs(b *testing.B) {
	for _, p := range pairs {
		if len(p.decoded) < 10 {
			continue
		}
		b.Run(fmt.Sprintf("input:%q", p.decoded), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = EncodeToString([]byte(p.decoded))
			}
		})
	}
}

func BenchmarkDecodePairs(b *testing.B) {
	for _, p := range pairs {
		if len(p.decoded) < 10 {
			continue
		}
		b.Run(fmt.Sprintf("input:%q", p.encoded), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = DecodeString(p.encoded)
			}
		})
	}
}

func BenchmarkEncodeBigLen(b *testing.B) {
	encoded := make([]byte, EncodedLen(len(bigBuf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Encode(encoded, bigBuf)
	}
}

func BenchmarkDecodeBigLen(b *testing.B) {
	encoded := make([]byte, EncodedLen(len(bigBuf)))
	n := Encode(encoded, bigBuf)
	encoded = encoded[:n]
	decoded := make([]byte, len(bigBuf))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Decode(decoded, encoded)
	}
}

func init() {
	bigBuf = make([]byte, 1<<20)
	for i := 0; i < len(bigBuf); i++ {
		bigBuf[i] = byte(i)
	}
}
