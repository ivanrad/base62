package base62

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

var testCases = []struct {
	decoded, encoded string
}{
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
		i := i
		t.Run(fmt.Sprintf("CharToBitsTable[%d]", i), func(t *testing.T) {
			t.Parallel()
			want := byte(strings.IndexByte(encodingAlphabet, byte(i)))
			got := charToBitsTable[i]
			if got != want {
				t.Errorf("charToBitsTable[%d] = %02x; want = %02x", i, got, want)
			}
		})
	}
}

func TestEncodedLen(t *testing.T) {
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%q", tc.decoded), func(t *testing.T) {
			t.Parallel()
			n := EncodedLen(len(tc.decoded))
			if n < len(tc.encoded) {
				t.Errorf("EncodedLen(%q) = %v; want greater or equal than %v", tc.decoded,
					n, len(tc.encoded))
			}
		})
	}
}

func TestDecodedLen(t *testing.T) {
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%q", tc.encoded), func(t *testing.T) {
			t.Parallel()
			n := DecodedLen(len(tc.encoded))
			if n < len(tc.decoded) {
				t.Errorf("DecodedLen(%q) = %v; want greater or equal than %v", tc.encoded,
					n, len(tc.decoded))
			}
		})
	}
}

func TestEncodePairs(t *testing.T) {
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%q to %q", tc.decoded, tc.encoded), func(t *testing.T) {
			t.Parallel()
			encoded := EncodeToString([]byte(tc.decoded))
			if encoded != tc.encoded {
				t.Errorf("EncodeToString(%q) = %q; want: %q", tc.decoded, encoded,
					tc.encoded)
			}
		})
	}
}

func TestDecodePairs(t *testing.T) {
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%q to %q", tc.encoded, tc.decoded), func(t *testing.T) {
			t.Parallel()
			decoded, err := DecodeString(tc.encoded)
			if err != nil {
				t.Fatalf("Error while decoding(%q): %v", tc.encoded, err)
			}
			if string(decoded) != tc.decoded {
				t.Errorf("DecodeString(%q) = %q; wanted: %q", tc.encoded, string(decoded),
					tc.decoded)
			}
		})
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
		tc := tc
		t.Run(fmt.Sprintf("%q", tc.input), func(t *testing.T) {
			t.Parallel()
			_, err := DecodeString(tc.input)
			var inputErr InputError
			switch {
			case err == nil && tc.corruptIdx == -1:
				// input ok
			case errors.As(err, &inputErr):
				if int64(inputErr) != int64(tc.corruptIdx) {
					t.Errorf("Corrupted input(%q) at %v; want %v", tc.input,
						int64(inputErr), tc.corruptIdx)
				}
			case err != nil && tc.corruptIdx == -1:
				t.Errorf("Corrupted input(%q): %v", tc.input, err)
			default:
				t.Errorf("Unexpected error %v", err)
			}
		})
	}
}

func TestBigLen(t *testing.T) {
	encoded := EncodeToString(bigBuf)
	decoded, err := DecodeString(encoded)
	if err != nil {
		t.Fatalf("Error while decoding: %v", err)
	}
	if !bytes.Equal(bigBuf, decoded) {
		if len(decoded) != len(bigBuf) {
			t.Fatalf("Decoded length: %v; want: %v", len(decoded), len(bigBuf))
		}
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
	var inputErr InputError
	switch {
	case errors.As(err, &inputErr):
		if int64(inputErr) != -1 {
			t.Errorf("Expected input truncated InputError(-1); got %v",
				int64(inputErr))
		}
	default:
		t.Errorf("Expected InputError; got %v", err)
	}
}

func BenchmarkEncodePairs(b *testing.B) {
	for _, tc := range testCases {
		if len(tc.decoded) < 10 {
			continue
		}
		b.Run(fmt.Sprintf("input:%q", tc.decoded), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = EncodeToString([]byte(tc.decoded))
			}
		})
	}
}

func BenchmarkDecodePairs(b *testing.B) {
	for _, tc := range testCases {
		if len(tc.decoded) < 10 {
			continue
		}
		b.Run(fmt.Sprintf("input:%q", tc.encoded), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = DecodeString(tc.encoded)
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
