// 2014-08-04 Adam Bryt

package main

import (
	"bytes"
	"testing"

	"npwp/tab"
)

type test struct {
	in      string // tekst wejściowy
	out     string // tekst wyjściowy (oczekiwany)
	tablist []int  // lista progów tabulatora
}

var tests = []test{
	{
		"",
		"",
		[]int{},
	},
	{
		"a",
		"a",
		[]int{},
	},
	{
		"aaa bbb ccc", // pojedyncza spacja przed progiem
		"aaa bbb\tccc",
		[]int{},
	},
	{
		" abc ",
		" abc ",
		[]int{},
	},
	{
		"abc bbb\nccc\n",
		"abc bbb\nccc\n",
		[]int{},
	},
	{
		"a\n a\n  a\n   a\n    a\n     a\n      a\n       a\n        a\n",
		"a\n a\n  a\n   a\n    a\n     a\n      a\n       a\n\ta\n",
		[]int{},
	},
	{
		"aaaaaaaaa",
		"aaaaaaaaa",
		[]int{},
	},
	{
		"aaaaaaa a",
		"aaaaaaa\ta",
		[]int{},
	},
	{
		"aaaaaa  a",
		"aaaaaa\ta",
		[]int{},
	},
	{
		"aaaaa   a",
		"aaaaa\ta",
		[]int{},
	},
	{
		"aaaa    a",
		"aaaa\ta",
		[]int{},
	},
	{
		"aaa     a",
		"aaa\ta",
		[]int{},
	},
	{
		"aa      a",
		"aa\ta",
		[]int{},
	},
	{
		"a       a",
		"a\ta",
		[]int{},
	},
	{
		"        a",
		"\ta",
		[]int{},
	},
	{
		"                    a",
		"\t\t    a",
		[]int{},
	},
	{
		"                    ",
		"\t\t    ",
		[]int{},
	},
	{
		"a\tbb", // znak tab na wejściu
		"a\tbb",
		[]int{},
	},

}

func TestEntab2(t *testing.T) {
	for i, tc := range tests {
		r := bytes.NewBufferString(tc.in)
		w := bytes.NewBuffer([]byte{})
		tt, err := tab.NewTabulator(tc.tablist)
		if err != nil {
			t.Error(err)
		}

		err = entab(w, r, tt)
		if err != nil {
			t.Error(err)
		}

		out := w.String()
		if out != tc.out {
			t.Errorf("tc #%d: oczekiwano: %q, jest %q", i, tc.out, out)
		}
	}
}
