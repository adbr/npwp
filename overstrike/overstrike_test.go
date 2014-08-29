// 2014-08-08 Adam Bryt

package main

import (
	"bytes"
	"testing"
)

type test struct {
	input  string
	output string
}

var tests = []test{
	// różnie przypadki bez znaków \b
	//
	{ // brak danych na wejściu
		"",
		"",
	},
	{
		"a",
		" a",
	},
	{
		"abc",
		" abc",
	},
	{
		"ab cde efg",
		" ab cde efg",
	},
	{
		" abc  cde ",
		"  abc  cde ",
	},
	{
		"abc \n cde\nefg\n",
		" abc \n  cde\n efg\n",
	},
	// przypadki ze znakami \b
	//
	{
		"abc\b1234",
		" abc\n" +
		"+  1234",
	},
	{
		"abc\b\b1234",
		" abc\n" +
		"+ 1234",
	},
	{
		"abc\b\b\b1234",
		" abc\n" +
		"+1234",
	},
	{
		"abc\b\b\b\b\b1234",
		" abc\n" +
		"+1234",
	},
	{
		"\b\b\b123",
		"\n" + // czy powinna być spacja na początku?
		"+123",
	},
	{
		"abc\b_def\b_",
		" abc\n" +
		"+  _def\n" +
		"+     _",
	},
	{	"a\b_b\b_c\b_d\b_",
		" a\n" +
		"+_b\n" +
		"+ _c\n" +
		"+  _d\n" +
		"+   _",
	},
	{
		"asdf_",
		" asdf\n" +
		"+   _",
	},

}

func TestOverstrike(t *testing.T) {
	for i, tc := range tests {
		r := bytes.NewBufferString(tc.input)
		w := bytes.NewBuffer([]byte{})

		err := overstrike(w, r)
		if err != nil {
			t.Error(err)
		}

		out := w.String()
		if out != tc.output {
			t.Errorf("tc #%d: oczekiwano: %q, jest %q", i, tc.output, out)
		}
	}
}
