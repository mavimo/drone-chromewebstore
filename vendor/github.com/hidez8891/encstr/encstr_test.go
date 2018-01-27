package encstr

import (
	"reflect"
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
)

type StrTest struct {
	raw []byte
	str string
	enc encoding.Encoding
}

var testdata = []StrTest{
	{
		str: "test str",
		raw: []byte{0x74, 0x65, 0x73, 0x74, 0x20, 0x73, 0x74, 0x72},
		enc: unicode.UTF8,
	},
	{
		str: "test str",
		raw: []byte{0x74, 0x65, 0x73, 0x74, 0x20, 0x73, 0x74, 0x72},
		enc: japanese.ShiftJIS,
	},
	{
		str: "!てすと!",
		raw: []byte{0x21, 0xe3, 0x81, 0xa6, 0xe3, 0x81, 0x99, 0xe3, 0x81, 0xa8, 0x21},
		enc: unicode.UTF8,
	},
	{
		str: "!てすと!",
		raw: []byte{0x21, 0x82, 0xc4, 0x82, 0xb7, 0x82, 0xc6, 0x21},
		enc: japanese.ShiftJIS,
	},
}

func TestStr(t *testing.T) {
	for _, d := range testdata {
		es := NewString2(d.raw, d.enc)
		str := es.Str()
		if str != d.str {
			t.Fatalf("EncString.Str(): got %q, want %q", str, d.str)
		}
	}
}

func TestSet(t *testing.T) {
	for _, d := range testdata {
		es := NewString2([]byte{}, d.enc)
		es.Set(d.str)
		if reflect.DeepEqual(es.data, d.raw) != true {
			t.Fatalf("EncString.Set(): got %v, want %v", es.data, d.raw)
		}
	}
}
