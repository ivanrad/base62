package base62

import (
	"bytes"
	"strings"
	"testing"
)

// test suite based on base64_test.go from Go standard library
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
		switch err := err.(type) {
		case InputError:
			if int(err) != tc.corruptIdx {
				t.Errorf("Corrupted input(%q) at %v; want %v", tc.input, int(err),
					tc.corruptIdx)
			}
		default:
			t.Errorf("Unexpected error %v", err)
		}
	}
}

func TestBigLen(t *testing.T) {
	buflen := 1 << 20
	buf := make([]byte, buflen)
	for i := 0; i < buflen; i++ {
		buf[i] = encodingAlphabet[i%len(encodingAlphabet)]
	}
	encoded := EncodeToString(buf)
	decoded, err := DecodeString(encoded)
	if err != nil {
		t.Fatalf("Error while decoding: %v", err)
	}
	if len(decoded) != buflen {
		t.Fatalf("Decoded length: %v; want: %v", len(decoded), buflen)
	}
	if !bytes.Equal(buf, decoded) {
		for i := 0; i < buflen; i++ {
			if buf[i] != decoded[i] {
				t.Errorf("Decoding failed at position: %v", i)
			}
		}
	}
}
