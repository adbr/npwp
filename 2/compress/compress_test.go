// 2014-08-10 Adam Bryt

package main

import (
	"bufio"
	"bytes"
	"testing"
)

type test struct {
	input  string // tekst wejściowy
	output string // oczekiwany tekst wynikowy
}

var tests = []test{
	// przypadki bez sekwencji
	{
		"",
		"",
	},
	{
		"a",
		"a",
	},
	{
		"ab",
		"ab",
	},
	{
		"abc",
		"abc",
	},
	{
		"abc de \n xy\n",
		"abc de \n xy\n",
	},
	{
		"aąeęsścć\n",
		"aąeęsścć\n",
	},
	// sekwencje poniżej progu 4 znaków
	{
		"aa  bb",
		"aa  bb",
	},
	{
		"aaa",
		"aaa",
	},
	// sekwencje powyżej progu 4 znaków
	{
		"aaaa",
		"~Da",
	},
	{
		"aaaaa      bbbbbbb\n\n\n\n\n\n\n\n",
		"~Ea~F ~Gb~H\n",
	},
	// tylda
	{
		"~",
		"~A~",
	},
	{
		"~~",
		"~B~",
	},
	{
		"~~~",
		"~C~",
	},
	{
		"~~~~",
		"~D~",
	},
	{
		"aaaaa~",
		"~Ea~A~",
	},
	{
		"aaa~~~bbb~~~\n\n\n",
		"aaa~C~bbb~C~\n\n\n",
	},
	{
		"aaabbbbb~cccc",
		"aaa~Eb~A~~Dc",
	},
}

func TestCompress(t *testing.T) {
	for i, tc := range tests {
		r := bytes.NewBufferString(tc.input)
		w := bytes.NewBuffer([]byte{})

		err := compress(w, r)
		if err != nil {
			t.Error(err)
		}

		out := w.String()
		if out != tc.output {
			t.Errorf("tc #%d: oczekiwano: %q, jest: %q", i, tc.output, out)
		}
	}
}

type testPutrep struct {
	n   int
	c   rune
	out string
}

var testsPutrep = []testPutrep{
	{0, 'a', ""},
	{1, 'a', "a"},
	{2, 'a', "aa"},
	{3, 'a', "aaa"}, // poniżej progu 4 znaków - nie koduje
	{4, 'a', "~Da"}, // przekroczenie progu 4 znaków
	{5, 'a', "~Ea"},
	{26, 'a', "~Za"}, // najwyższa wartość licznika znaków
	{27, 'a', "~Zaa"},
	{28, 'a', "~Zaaa"},
	{29, 'a', "~Zaaaa"},
	{30, 'a', "~Za~Da"},
	{52, 'a', "~Za~Za"},
	{53, 'a', "~Za~Zaa"},
	// znaki utf-8
	{3, 'ą', "ąąą"},
	{5, 'ą', "~Eą"},
	// kodowanie znaku ~
	{0, '~', ""},
	{1, '~', "~A~"}, // znak ~ jest kodowany także poniżej progu 4 znaków
	{2, '~', "~B~"},
	{3, '~', "~C~"},
	{4, '~', "~D~"},
	{5, '~', "~E~"},
	{26, '~', "~Z~"},
	{27, '~', "~Z~~A~"},
	{28, '~', "~Z~~B~"},
}

func TestPutrep(t *testing.T) {
	for i, tc := range testsPutrep {
		w := bytes.NewBuffer([]byte{})
		bw := bufio.NewWriter(w) // bufio.Writer jest wymagany przez putrep

		err := putrep(bw, tc.n, tc.c)
		if err != nil {
			t.Error(err)
		}
		bw.Flush()

		out := w.String()
		if out != tc.out {
			t.Errorf("tc #%d: oczekiwano: %q, jest: %q", i, tc.out, out)
		}
	}
}
