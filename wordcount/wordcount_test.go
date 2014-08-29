// 2014-07-30 Adam Bryt

package main

import (
	"bytes"
	"testing"
)

type test struct {
	in string // string wejściowy
	n  int64  // liczba słów
}

var tests = []test{
	{"", 0},
	{" \t  \n  ", 0}, // zero słów
	{"abc", 1},       // jedno słowo w różnych miejscach
	{"  abc", 1},
	{"  abc  \n  ", 1},
	{"abc 123", 2}, // dwa słowa w różnych miejscach
	{" \n\tabc 123", 2},
	{"abc  123   ", 2},
	{"a bb  cccc\td\nee", 5},
	{"ąą bb ", 2}, // polskie znaki (UTF-8)
	{"  abc  ", 1},
	{"aa\nbb\ncc\ndd", 4},
	{"aa\rbb\ncc", 3},
	{"123  !@#$% , {}[]", 4},
	{"programo-\nwanie", 2}, // słowo dzielone
}

func TestWordcount(t *testing.T) {
	for i, tc := range tests {
		b := bytes.NewBufferString(tc.in)
		n, err := wordcount(b)
		if err != nil {
			t.Error(err)
		}
		if n != tc.n {
			t.Errorf("tc #%v: oczekiwano %v, jest %v", i, tc.n, n)
		}
	}
}
