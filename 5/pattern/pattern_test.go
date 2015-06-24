// 2015-01-20 Adam Bryt

package pattern

import (
	"bytes"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestMakepat(t *testing.T) {
	tests := []struct {
		in  string // wzorzec w postaci źródłowej
		out string // postać stringowa wzorca skompilowanego
	}{
		{
			"",
			"",
		},
		{
			"a",
			"<LITCHAR>a",
		},
		{
			"abc",
			"<LITCHAR>a<LITCHAR>b<LITCHAR>c",
		},
		{
			"%",
			"<BOL>",
		},
		{
			"%ab",
			"<BOL><LITCHAR>a<LITCHAR>b",
		},
		{
			// BOL nie na początku wzorca
			"a%b",
			"<LITCHAR>a<LITCHAR>%<LITCHAR>b",
		},
		{
			// dwa znaki BOL na początku - druki zwykły znak
			"%%a",
			"<BOL><LITCHAR>%<LITCHAR>a",
		},
		{
			"$",
			"<EOL>",
		},
		{
			"ab$",
			"<LITCHAR>a<LITCHAR>b<EOL>",
		},
		{
			// EOL nie na końcu wzorca
			"a$b",
			"<LITCHAR>a<LITCHAR>$<LITCHAR>b",
		},
		{
			// dwa znaki EOL na końcu - pierwszy zwykły
			"a$$",
			"<LITCHAR>a<LITCHAR>$<EOL>",
		},
		{
			"%$",
			"<BOL><EOL>",
		},
		{
			"a?b",
			"<LITCHAR>a<ANY><LITCHAR>b",
		},
		{
			"??",
			"<ANY><ANY>",
		},
		{
			// wielobajtowe znaki UTF8
			"ąĘ",
			"<LITCHAR>ą<LITCHAR>Ę",
		},
		{"%ą123??@?$$",
			"<BOL><LITCHAR>ą<LITCHAR>1<LITCHAR>2<LITCHAR>3<ANY><ANY><LITCHAR>?<LITCHAR>$<EOL>",
		},
		// klasa znaków
		{
			// prosta, typowa klasa znaków
			"[aąbcę]",
			"<CCL>aąbcę",
		},
		{
			// klasa znaków z escapeowaniem
			"a[a@b@@@-k@]ę]x",
			"<LITCHAR>a<CCL>ab@-k]ę<LITCHAR>x",
		},
		{
			// klasa znaków zawierająca '['
			"[[ąę[]",
			"<CCL>[ąę[",
		},
		{
			// klasa znaków z zakresem
			"a[aa-d0-5]b",
			"<LITCHAR>a<CCL>aabcd012345<LITCHAR>b",
		},
		// zanegowana klasa znaków
		{
			"[^abąę]",
			"<NCCL>abąę",
		},
		{
			"a[^a-e^ą]x",
			"<LITCHAR>a<NCCL>abcde^ą<LITCHAR>x",
		},
		// domknięcie (*)
		{
			"ab*c",
			"<LITCHAR>a<CLOSURE><LITCHAR>b<LITCHAR>c",
		},
		{
			// znak '*' na początku
			"*a",
			"<LITCHAR>*<LITCHAR>a",
		},
		{
			"a*[^a-d]*b[0-9]*?*$",
			"<CLOSURE><LITCHAR>a<CLOSURE><NCCL>abcd<LITCHAR>b<CLOSURE><CCL>0123456789<CLOSURE><ANY><EOL>",
		},
	}

	for i, test := range tests {
		pat, err := Makepat(test.in)
		if err != nil {
			t.Error(err)
		}
		s := pat.String()
		if s != test.out {
			t.Errorf("#%d: oczekiwano: %q, jest: %q", i, test.out, s)
		}
	}

	// Test przekroczenia maksymalnej liczby znaków w klasie
	// najpierw maksymalna dozwolona ilość
	in := strings.Repeat("a", maxChars)
	in = "[" + in + "]"
	_, err := Makepat(in)
	if err != nil {
		t.Error(err)
	}

	// przekroczenie maksymalnej dozwolonej wartości
	in = strings.Repeat("a", maxChars+1)
	in = "[" + in + "]"
	_, err = Makepat(in)
	if err == nil {
		t.Error("oczekiwano błędu przekroczenia maksymalnej liczby znaków")
	}
}

func TestEsc(t *testing.T) {
	tests := []struct {
		in string
		r  rune
		s  string
	}{
		{
			"abc",
			'a',
			"bc",
		},
		{
			"a",
			'a',
			"",
		},
		{
			"@abc",
			'a',
			"bc",
		},
		{
			"@@abc",
			'@',
			"abc",
		},
		{
			// znak @ na końcu stringu
			"@",
			'@',
			"",
		},
		{
			"@tabc",
			'\t',
			"abc",
		},
		{
			"@nabc",
			'\n',
			"abc",
		},
		{
			"@???",
			'?',
			"??",
		},
		{
			"@[bc",
			'[',
			"bc",
		},
		{
			"",
			utf8.RuneError,
			"",
		},
	}

	for i, test := range tests {
		r, s := Esc(test.in)
		if r != test.r {
			t.Errorf("#%d oczekiwano: %q, jest: %q", i, test.r, r)
		}
		if s != test.s {
			t.Errorf("#%d oczekiwano: %q, jest: %q", i, test.s, s)
		}
	}
}

func TestIsAlphanum(t *testing.T) {
	var tests = []struct {
		r rune
		b bool
	}{
		{'a', true},
		{'k', true},
		{'z', true},
		{'A', true},
		{'M', true},
		{'Z', true},
		{'0', true},
		{'1', true},
		{'9', true},
		{'+', false},
		{' ', false},
		{'ą', false}, // tylko znaki ASCII
		{'Ż', false},
	}

	for i, test := range tests {
		b := isAlphanum(test.r)
		if b != test.b {
			t.Errorf("#%d: isAlphanum(%q) == %v, oczekiwano %v", i, test.r, b, test.b)
		}
	}
}

func TestDodash(t *testing.T) {
	tests := []struct {
		ss    string
		delim rune
		chars []byte
		s     string
		isErr bool
	}{
		// najprostszy, typowy przypadek
		{
			"abc1ąę]de",
			']',
			[]byte("abc1ąę"),
			"de",
			false,
		},
		// zmieniony ogranicznik
		{
			"abc]dxe",
			'x',
			[]byte("abc]d"),
			"e",
			false,
		},
		// sekwencje escapeowe
		{
			"a@b@@c@t@]de]fgh",
			']',
			[]byte("ab@c\t]de"),
			"fgh",
			false,
		},
		// zakres znaków
		{
			"ab-fgąę]hi",
			']',
			[]byte("abcdefgąę"),
			"hi",
			false,
		},
		// kilka zakresów znaków
		{
			"a-zA-Z0-9]",
			']',
			[]byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"),
			"",
			false,
		},
		// znak '-' na początku
		{
			"-klmn]",
			']',
			[]byte("-klmn"),
			"",
			false,
		},
		// znak '-' na końcu
		{
			"abc-]def",
			']',
			[]byte("abc-"),
			"def",
			false,
		},
		// znak '-' obok nie alphanum
		{
			"a.-bc--dą-ż]",
			']',
			[]byte("a.-bc--dą-ż"),
			"",
			false,
		},
		// znak '-' escapeowany
		{
			"a@-z]",
			']',
			[]byte("a-z"),
			"",
			false,
		},
		// brak ogranicznika - błąd
		{
			"abc",
			']',
			[]byte("abc"),
			"",
			true,
		},
		// string pusty: brak ogranicznika - błąd
		{
			"",
			']',
			[]byte(""),
			"",
			true,
		},
		// brak znaków przed ogranicznikiem
		{
			"]xyz",
			']',
			[]byte(""),
			"xyz",
			false,
		},
		// zakres znaków 'a-b', a > b
		{
			"az-ab]x",
			']',
			[]byte("az-ab"),
			"x",
			false,
		},
		// zakres znaków 'a-b', a==b
		{
			"ab-bc]x",
			']',
			[]byte("abc"),
			"x",
			false,
		},
		// znaki się powtarzają
		{
			"aabbb]xx",
			']',
			[]byte("aabbb"),
			"xx",
			false,
		},
	}

	for i, test := range tests {
		chars, s, err := dodash(test.ss, test.delim)
		if !bytes.Equal(chars, test.chars) {
			t.Errorf("#%d (chars): oczekiwano: %q, jest: %q", i, test.chars, chars)
		}
		if s != test.s {
			t.Errorf("#%d (s): oczekiwano: %q, jest: %q", i, test.s, s)
		}
		isErr := err != nil
		if isErr != test.isErr {
			t.Errorf("#%d (isErr): oczekiwano: %v, jest: %v", i, test.isErr, isErr)
		}
	}
}

func TestGetccl(t *testing.T) {
	tests := []struct {
		ss    string // argument getccl
		chars []byte
		isneg bool
		s     string
		iserr bool
	}{
		// typowy przypadek
		{
			"[abcąęść]xyz",
			[]byte("abcąęść"),
			false,
			"xyz",
			false,
		},
		// zanegowana klasa znaków
		{
			"[^abcąęś]xyz",
			[]byte("abcąęś"),
			true,
			"xyz",
			false,
		},
		// zakres znaków w klasie
		{
			"[ab-gąę]xyz",
			[]byte("abcdefgąę"),
			false,
			"xyz",
			false,
		},
		// escapeowanie znaków
		{
			"[@a@-@d@@@]x]xyz",
			[]byte("a-d@]x"),
			false,
			"xyz",
			false,
		},
		// znak '^' nie na początku klasy
		{
			"[a^bc]x",
			[]byte("a^bc"),
			false,
			"x",
			false,
		},
		// znak '[' wewnątrz klasy
		{
			"[a[b]x",
			[]byte("a[b"),
			false,
			"x",
			false,
		},
		// znak ']' wewnątrz klasy
		{
			"[ab]c]x",
			[]byte("ab"),
			false,
			"c]x",
			false,
		},
		{
			"[ab@]c]x",
			[]byte("ab]c"),
			false,
			"x",
			false,
		},
		// znak '-' na początku
		{
			"[-c-e]x",
			[]byte("-cde"),
			false,
			"x",
			false,
		},
		// znak '-' na końcu
		{
			"[ab-]x",
			[]byte("ab-"),
			false,
			"x",
			false,
		},
		// błąd: brak początkowego znaku '['
		{
			"abc]x",
			[]byte(""),
			false,
			"abc]x",
			true,
		},
		// błąd: brak końcowego znaku ']'
		{
			"[abcd",
			[]byte("abcd"),
			false,
			"",
			true,
		},
	}

	for i, test := range tests {
		chars, isneg, s, err := getccl(test.ss)
		if !bytes.Equal(chars, test.chars) {
			t.Errorf("#%d (chars): oczekiwano: %q, jest: %q", i, test.chars, chars)
		}
		if isneg != test.isneg {
			t.Errorf("#%d (isneg): oczekiwano: %v, jest: %v", i, test.isneg, isneg)
		}
		if s != test.s {
			t.Errorf("#%d (s): oczekiwano: %q, jest: %q", i, test.s, s)
		}
		iserr := err != nil
		if iserr != test.iserr {
			t.Errorf("#%d (iserr): oczekiwano: %v, jest: %v", i, test.iserr, iserr)
		}
	}
}
