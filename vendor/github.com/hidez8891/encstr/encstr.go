package encstr

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

func NewString(s string) *String {
	return &String{
		data: []byte(s),
		enc:  unicode.UTF8,
	}
}

func NewString2(b []byte, enc encoding.Encoding) *String {
	return &String{
		data: b,
		enc:  enc,
	}
}

type String struct {
	data []byte
	enc  encoding.Encoding
}

func (s *String) Str() string {
	b, _ := s.enc.NewDecoder().Bytes(s.data)
	return string(b)
}

func (s *String) Set(str string) {
	b, _ := s.enc.NewEncoder().Bytes([]byte(str))
	s.data = b
}

func (s *String) Check() bool {
	_, err := s.enc.NewDecoder().Bytes(s.data)
	return err == nil
}

func (s *String) Raw() []byte {
	return s.data
}

func (s *String) Encoding() encoding.Encoding {
	return s.enc
}

func (s *String) Convert(enc encoding.Encoding) {
	dec := s.enc
	b1, _ := dec.NewDecoder().Bytes(s.data)
	b2, _ := enc.NewEncoder().Bytes(b1)

	s.data = b2
	s.enc = enc
}

func (s *String) ForceConvert(enc encoding.Encoding) {
	s.enc = enc
}
