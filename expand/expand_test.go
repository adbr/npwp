// 2014-08-20 Adam Bryt

package main

import (
	"bytes"
	"strings"
	"testing"
)

type test struct {
	in  string // dane wejściowe
	out string // oczekiwane dane wyjściowe
}

var tests = []test{
	{
		"",
		"",
	},
	{
		"abc def\nghi\n",
		"abc def\nghi\n",
	},
	{
		"~Eabc",
		"aaaaabc",
	},
	{
		"abc ~Fx",
		"abc xxxxxx",
	},
	{
		"abc~Exdef",
		"abcxxxxxdef",
	},
	{
		"abc ~Dx\n~Dydef~Ez\n",
		"abc xxxx\nyyyydefzzzzz\n",
	},
	{
		"~E~",
		"~~~~~",
	},
	{
		"~A~x~B~x~C~x~D~x~E~",
		"~x~~x~~~x~~~~x~~~~~",
	},

	// przypadki błędnych sekwencji kodujących

	// EOF po znaku ~
	{
		"abc~",
		"abc~",
	},

	// EOF po kodzie liczbę znaków
	{
		"abc~D",
		"abc~D",
	},

	// Błędny znak kodujący liczbę znaków
	{
		"abc~@x",
		"abc~@x",
	},
	{
		"~~~x",
		"~~~x",
	},

	// znaki UTF-8

	{
		"ąęśabcżźć",
		"ąęśabcżźć",
	},
	{
		"xx~Eąxx",
		"xxąąąąąxx",
	},
}

func TestExpand(t *testing.T) {
	for i, tc := range tests {
		w := new(bytes.Buffer)
		r := strings.NewReader(tc.in)

		err := expand(w, r)
		if err != nil {
			t.Error(err)
		}

		out := w.String()
		if out != tc.out {
			t.Errorf("tc #%d: oczekiwano: %q, jest: %q", i, tc.out, out)
		}
	}
}
