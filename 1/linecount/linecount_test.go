// 2014-07-30 Adam Bryt

package main

import (
	"bytes"
	"testing"
)

type test struct {
	in string
	n  int64
}

var tests = []test{
	{"", 0},
	{"abc", 0},
	{"\n", 1},
	{"abc\n", 1},
	{"abc\ndefg\n\n", 3},
	{"abc\ndef", 1},
	{"\n\n\n\n\n", 5},
	{"ąęś\nżźć\n", 2},
}

func TestLinecount(t *testing.T) {
	for i, tc := range tests {
		b := bytes.NewBufferString(tc.in)
		n, err := linecount(b)
		if err != nil {
			t.Error(err)
		}
		if n != tc.n {
			t.Errorf("tc #%v: oczekiwano: %v, jest %v", i, tc.n, n)
		}
	}
}
