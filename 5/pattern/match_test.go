// 2015-01-19 Adam Bryt

package pattern

import (
	"testing"
)

func TestPatsize(t *testing.T) {
	tests := []struct {
		in   string
		size int
	}{
		{"a", 2},      // LITCHAR
		{"ą", 3},      // LITCHAR
		{"", 0},       // ?
		{"%abc", 1},   // BOL
		{"$", 1},      // EOL
		{"?", 1},      // ANY
		{"[abc]", 5},  // CCL
		{"[ąbc]", 6},  // CCL z utf8
		{"[^ąęś]", 8}, // NCCL
		{"*", 2},      // LITCHAR bo * na początku wzorca
		{"a*", 1},     // CLOSURE - po kompilacji CLOSURE jest przed a
	}

	for i, test := range tests {
		pat, err := Makepat(test.in)
		if err != nil {
			t.Error(err)
		}
		size := patsize(pat)
		if size != test.size {
			t.Errorf("#%d: patsize: %d, oczekiwano: %d", i, size, test.size)
		}
	}
}

func TestOmatch(t *testing.T) {
	tests := []struct {
		str string
		i   int
		pat string // wzorzec w postaci źródłowej
		j   int
		ok  bool // czy pasuje
		n   int  // ile bajtów str pasuje
	}{
		// ANY
		{
			"abc", 0,
			"?bc", 0,
			true, 1,
		},
		{
			"ąbc", 0,
			"?", 0,
			true, 2,
		},
		{
			"aąb", 1,
			"x?x", 2, // <LITCHAR>x<ANY><LITCHAR>x
			true, 2,
		},
		{
			"", 0,
			"?", 0,
			false, 0,
		},
		{
			"a", 1,
			"x?x", 2,
			false, 0,
		},
		{
			"\n", 0,
			"?", 0,
			false, 0,
		},
		// BOL
		{
			"abc", 0,
			"%xyz", 0,
			true, 0,
		},
		{
			"abc", 1, // nie na początku stringu
			"%xyz", 0,
			false, 0,
		},
		// EOL
		{
			"abc\n", 3,
			"$", 0,
			true, 0,
		},
		{
			"abc", 1,
			"$", 0,
			false, 0,
		},
		// LITCHAR
		{
			"abc", 0,
			"axy", 0,
			true, 1,
		},
		{
			"ąbc", 0,
			"ąxy", 0,
			true, 2,
		},
		{
			"aąb", 1,
			"xxą", 4,
			true, 2,
		},
		{
			"abc", 0,
			"xbc", 0,
			false, 0,
		},
		// CCL
		{
			"abc", 0,
			"[xay]z", 0,
			true, 1,
		},
		{
			"abc", 1,
			"?[a-z]", 1,
			true, 1,
		},
		{
			"ąbc", 0,
			"[xząęśq0]", 0,
			true, 2,
		},
		{
			"abc", 2,
			"[abx]", 0,
			false, 0,
		},
		// NCCL
		{
			"abc", 0,
			"[^xyz]", 0,
			true, 1,
		},
		{
			"\nabc", 0, // nie bierzemy pod uwagę '\n'
			"[^xyz]", 0,
			false, 0,
		},
		{
			"a", 0,
			"[^0-9A-Z]", 0,
			true, 1,
		},
		{
			"a", 0,
			"[^ąęśżźć]", 0,
			true, 1,
		},
		{
			"ą", 0,
			"[^a-z]", 0,
			true, 2,
		},
		{
			"a", 0,
			"[^a-z]", 0,
			false, 0,
		},
	}

	for i, test := range tests {
		pat, err := Makepat(test.pat)
		if err != nil {
			t.Error(err)
		}
		ok, n := omatch(test.str, test.i, pat, test.j)
		if ok != test.ok {
			t.Errorf("#%d: (ok) oczekiwano: %v, jest: %v", i, test.ok, ok)
		}
		if n != test.n {
			t.Errorf("#%d: (n) oczekiwano: %v, jest: %v", i, test.n, n)
		}
	}
}

func TestAmatch(t *testing.T) {
	tests := []struct {
		str string
		i   int
		pat string // postać źródłowa wzorca
		j   int
		ok  bool
		n   int
	}{
		// LITCHAR proste dopasowanie znaków
		{
			"abc", 0,
			"abc", 0,
			true, 3,
		},
		{
			"abcąęśćxyz", 3,
			"ąęś", 0,
			true, 6,
		},
		{
			"abcd", 0,
			"xab", 0,
			false, 0,
		},
		// ANY
		{
			"abcde", 0,
			"a???", 0,
			true, 4,
		},
		{
			"axyzw", 0,
			"a???", 0,
			true, 4,
		},
		{
			"ąęść", 0,
			"??ś", 0,
			true, 6,
		},
		{
			"ab", 0,
			"????", 0,
			false, 0,
		},
		{
			"ab\n", 0,
			"???", 0,
			false, 0, // '?' nie pasuje do '\n'
		},
		// BOL
		{
			"abc", 0,
			"%ab", 0,
			true, 2,
		},
		{
			"abc", 1,
			"%bc", 0,
			false, 0,
		},
		{
			"a%bc", 0,
			"a%bc", 0, // '%' nie na początku pat jest zwykłym znakiem
			true, 4,
		},
		// EOL
		{
			"abc\n", 2,
			"c$", 0,
			true, 1,
		},
		{
			"abcd", 2,
			"c$", 0,
			false, 0,
		},
		{
			"abc", 2, // nie ma znaku '\n' na końcu
			"c$", 0,
			false, 0,
		},
		{
			"abc$$\n", 2,
			"c$$$", 0, // znak '$' nie na końcu jest zwykłym znakiem
			true, 3,
		},
		// CCL
		{
			"abcd", 0,
			"[a-z][a-z][cd][cd]", 0,
			true, 4,
		},
		{
			"klcd", 0,
			"[a-z][a-z][cd][cd]", 0,
			true, 4,
		},
		{
			"a-zd", 0,
			"a[a@-z]z", 0,
			true, 3,
		},
		{
			"a[0-9xxx", 0,
			"a@[0-9", 0,
			true, 5,
		},
		{
			"aA9xxx", 0,
			"[a-z][A-Z][0-9]", 0,
			true, 3,
		},
		{
			"ąę", 0,
			"[ąę][ąę]", 0,
			true, 4,
		},
		// NCCL
		{
			"xyz", 0,
			"[^abc][^a-w][xzy]", 0,
			true, 3,
		},
		{
			"xyz", 0,
			"[^xyz]y", 0,
			false, 0,
		},
		{
			"x", 0,
			"[xy^z]", 0, // '^' nie na początku klasy jest zwykłym znakiem
			true, 1,
		},
		{
			"ąą", 0,
			"[^ęś][^ęś]", 0,
			true, 4,
		},
		// CLOSURE
		{
			"aaaab", 0,
			"a*b", 0,
			true, 5,
		},
		{
			"bbbbbb", 0,
			"b*b", 0,
			true, 6,
		},
		{
			"aaa123", 0,
			"b*a*a[0-9]*", 0,
			true, 6,
		},
		// złożony wzorzec, ale bez domknięcia
		{
			"ab7ąx\n", 0,
			"%?b[0-9][^a-z]x$", 0,
			true, 6,
		},
		{
			"ab7ąx", 0, // jak poprzedni ale bez '\n'
			"%?b[0-9][^a-z]x$", 0,
			false, 0,
		},
	}

	for i, test := range tests {
		pat, err := Makepat(test.pat)
		if err != nil {
			t.Error(err)
		}
		ok, n := amatch(test.str, test.i, pat, test.j)
		if ok != test.ok {
			t.Errorf("#%d: (ok) oczekiwano: %v, jest: %v", i, test.ok, ok)
		}
		if n != test.n {
			t.Errorf("#%d: (n) oczekiwano: %v, jest: %v", i, test.n, n)
		}
	}
}

func TestMatch(t *testing.T) {
	tests := []struct {
		str string
		pat string
		ok  bool
	}{
		{
			"abc",
			"a?",
			true,
		},
		{
			"ala ma kota\n",
			"ma?",
			true,
		},
		{
			"ala ma kota\n",
			"?ta$",
			true,
		},
		{
			"abc 2015-01-15 \n",
			" [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9] ",
			true,
		},
		{
			"abc 1900-01-01 \n",
			" [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9] ",
			true,
		},
		{
			"abc 215-01-15 \n",
			" [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9] ",
			false,
		},
		{
			"abc 215-01-15 \n",
			"%[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9] ",
			false,
		},
		{
			"abc 215-01-15 \n",
			"[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]$",
			false,
		},
		{
			"abc 123 def\n",
			"[^a-z]",
			true,
		},
		{
			"abc ąęś xyz",
			"[ęś] ",
			true,
		},
		{
			"aaa 2015-01-24 bbb\n",
			"%a* [0-9]*-[0-9]*-[0-9]* b*$",
			true,
		},
		{
			"aaa 2015-01-24 bbb\n",
			"%a* [0-9-]*-[0-9-]*-[0-9-]* b*$",
			true,
		},
	}

	for i, test := range tests {
		pat, err := Makepat(test.pat)
		if err != nil {
			t.Error(err)
		}
		ok := Match(test.str, pat)
		if ok != test.ok {
			t.Errorf("#%d: Match() oczekiwano: %v, jest %v", i, test.ok, ok)
		}
	}
}
