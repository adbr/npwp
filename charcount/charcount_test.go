// 2014-07-28 Adam Bryt

package main

import (
	"bytes"
	"testing"
)

// Typ test reprezentuje przypadek testowy.
type test struct {
	in     string // dane wejściowe
	n      int64  // liczba znaków
	isUtf8 bool   // czy dane są kodowane w UTF8
}

var tests = []test{
	{"", 0, false},
	{"", 0, true},
	{"a", 1, false},
	{"a", 1, true},
	{"abc", 3, false},
	{"abc\n", 4, false},
	{"ab\ncd\n", 6, false},
	{"ab\ncd", 5, false},
	{"ąąęę12\n", 7, true},
	{"ąąęę12\n", 11, false},
	{"\n", 1, false},
	{"\x12\xc3", 2, false},
	{"\x00\x00\x00", 3, false},
}

func TestCharcount(t *testing.T) {
	for i, tc := range tests {
		r := bytes.NewBufferString(tc.in)
		n, err := charcount(r, tc.isUtf8)
		if err != nil {
			t.Error(err)
		}
		if n != tc.n {
			t.Errorf("tc #%d: oczekiwano: %v, jest %v", i, tc.n, n)
		}
	}
}
