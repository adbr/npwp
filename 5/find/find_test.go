// 2014-12-02 Adam Bryt

package main

import (
	"bytes"
	"testing"
)

// Testowanie czytania i zapisywania wierszy.
func TestFind(t *testing.T) {
	type test struct {
		in  string
		out string
	}

	var tests = []test{
		{
			"",
			"",
		},
		{
			// jeden wiersz nie zakończony znakiem '\n'
			"aaa",
			"aaa",
		},
		{
			// ostatni wiersz nie zakończony znakiem '\n'
			"aaa\nbbb",
			"aaa\nbbb",
		},
		{
			// jeden pusty wiersz
			"\n",
			"\n",
		},
		{
			// kilka normalnych wierszy
			"aaa\n   \n\tbbb\nccc\n",
			"aaa\n   \n\tbbb\nccc\n",
		},
	}

	pat := pattern("")
	for i, tc := range tests {
		w := new(bytes.Buffer)
		r := bytes.NewBufferString(tc.in)
		err := find(w, r, pat)
		if err != nil {
			t.Error(err)
		}
		s := w.String()
		if s != tc.out {
			t.Errorf("#%d: oczekiwano: %q, jest: %q", i, tc.out, s)
		}
	}
}
