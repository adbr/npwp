// 2014-08-23 Adam Bryt

package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestTranslit(t *testing.T) {
	type test struct {
		from string // lista znaków zamienianych
		to   string // lista znaków na które są zamieniane znaki from
		in   string // tekst wejściowy
		out  string // oczekiwany tekst wyjściowy
	}

	tests := []test{
		{
			"",
			"",
			"",
			"",
		},
		{
			"",
			"",
			"abc\ndef",
			"abc\ndef",
		},
		{
			"x",
			"y",
			"abxxcdxfx",
			"abyycdyfy",
		},
		{
			"xyz",
			"XYZ",
			"axyzbxcydz",
			"aXYZbXcYdZ",
		},
		{
			"ab",
			"12",
			"aabbccaadd",
			"1122cc11dd",
		},
		{ // runy wielobajtowe
			"abłąq",
			"xxyęć",
			" qweaabbxłłłxąąąx",
			" ćwexxxxxyyyxęęęx",
		},
		{
			"aabc", // powtórzenie znaku we from
			"xy12",
			"aaabbbcccddd",
			"xxx111222ddd",
		},
		{
			"a\n", // newline we from
			"1x",
			"abc\ndef",
			"1bcxdef",
		},
		// "from" jest krótsze niż "to" - zachowanie bez zmian
		{
			"ab",
			"12345",
			"aaabbbcccaaa",
			"111222ccc111",
		},
		// "from" jest dłuższe niż "to":
		// - znak z from odpowiadający ostatniemu znakowi z to jest redukowany
		// - znaki z from nie mające odpowiednika w to są zastępowane ostatnim
		//   znakiem z to, i dodatkowo sekwencje takich znaków są redukowane
		//   do jednego znaku
		{
			"abcd",
			"xy",
			"aaabbbcccdddeee",
			"xxxyeee",
		},
		// "to" jest pusty - usuń znaki z "from"
		{
			"ab",
			"",
			"aaabbbcccdddaaa",
			"cccddd",
		},
		{
			"a",
			"",
			"qqqaaawww",
			"qqqwww",
		},
		// negacja znaków z "from",
		// - tylko redukcja do ostatniego znaku z "to", lub usuwanie znaków
		{
			"^ab",
			"x",
			"12aqqqbbb3334",
			"xaxbbbx",
		},
		{
			"^a",
			"",
			"123qweaaazxc",
			"aaa",
		},
		{
			"^a", // len(from) < len(to) a mimo to redukcja bo występuje negacja
			"xyz",
			"123aaa455556",
			"zaaaz",
		},
		// zakres znaków - dokładniejsze testy są w TestDodash
		{
			"a-z",
			"A-Z",
			"abcdef",
			"ABCDEF",
		},
		// znaki wyróżnione
		// przykładowe zastosowanie: zastąpienie dowolnych sekwencji znaków
		// spacji, nowego weirsza i tabulacji pojedynczym znakiem nowego
		// wiersza - czyli umieszczenie każdego słowa w osobnym wierszu
		{
			" @t@n",
			"@n",
			"a   b\n \t\nc\t\t \nd",
			"a\nb\nc\nd",
		},
		// negacja i zakres znaków
		{
			"^x-z",
			"n",
			"abc32+-.,xxaaybbzzz",
			"nxxnynzzz",
		},
		{
			// negacja i zakres: usunięcie znaków różnych od a-z
			"^a-z",
			"",
			"abcABCxx123,./",
			"abcxx",
		},
	}

	for i, tc := range tests {
		w := new(bytes.Buffer)
		r := strings.NewReader(tc.in)

		err := translit(w, r, tc.from, tc.to)
		if err != nil {
			t.Error(err)
		}

		out := w.String()
		if out != tc.out {
			t.Errorf("tc #%d: oczekiwano: %q, jest: %q", i, tc.out, out)
		}
	}
}

func TestIsAlphaNum(t *testing.T) {
	type test struct {
		c rune
		b bool
	}

	tests := []test{
		{'0', true},
		{'1', true},
		{'9', true},
		{'A', true},
		{'B', true},
		{'Z', true},
		{'a', true},
		{'b', true},
		{'z', true},
		{'!', false},
		{'=', false},
		{'\n', false},
		{'\t', false},
		{'ą', false},
		{'ż', false},
		{'Ą', false},
		{'Ż', false},
	}

	for i, tc := range tests {
		v := isAlphaNum(tc.c)
		if v != tc.b {
			t.Errorf("tc #%d: oczekiwano: %v, jest: %v", i, tc.b, v)
		}
	}
}

func TestEsc(t *testing.T) {
	type test struct {
		s []rune
		i int
		c rune // oczekiwany znak
	}

	tests := []test{
		{
			[]rune("abc"),
			1,
			'b',
		},
		{
			[]rune("aąćę"),
			2,
			'ć',
		},
		{ // znak @ na końcu
			[]rune("ab@"),
			2,
			'@',
		},
		{ // @t
			[]rune("@tabc"),
			0,
			'\t',
		},
		{ // @n
			[]rune("a@n"),
			1,
			'\n',
		},
		{ // znak @ przed zwykłym znakiem
			[]rune("a@bc"),
			1,
			'b',
		},
		{
			[]rune("a@ąbc"),
			1,
			'ą',
		},
	}

	for i, tc := range tests {
		v := esc(tc.s, tc.i)
		if v != tc.c {
			t.Errorf("tc #%d: oczekiwano: %c, jest: %c", i, tc.c, v)
		}
	}
}

func TestDodash(t *testing.T) {
	type test struct {
		in  string
		out string
	}

	tests := []test{
		{
			"abc",
			"abc",
		},
		{ // rozwinięcie @t@n
			"ab@t@nc",
			"ab\t\nc",
		},
		{ // nie rozwija @ na końcu
			"ąabc@",
			"ąabc@",
		},
		{ // @ przed zwykłym znakiem
			"ab@c@@d",
			"abc@d",
		},
		{ // - na początku i na końcu
			"-abc@nxyz-",
			"-abc\nxyz-",
		},
		{ // zakres cyfr
			"x0-9y",
			"x0123456789y",
		},
		{ // zakres dużych liter
			"xA-Z",
			"xABCDEFGHIJKLMNOPQRSTUVWXYZ",
		},
		{ // zakres małych liter
			"a-z",
			"abcdefghijklmnopqrstuvwxyz",
		},
		{ // zakres - oba krańce takie same
			"ab-bc",
			"abc",
		},
		{ // błędny zakres - nie alphanum
			"ab$-Z",
			"ab$-Z",
		},
		{
			"ab-!z",
			"ab-!z",
		},
		{ // błędny zakres - pierwszy znak większy niż drugi
			"ae-cx",
			"ae-cx",
		},
		// zakres od początku cyfr do końca małych liter - dodatkowe znaki
		// wynikają z zestawu znaków ASCII
		{
			"0-z",
			"0123456789:;<=>?@" +
				"ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`" +
				"abcdefghijklmnopqrstuvwxyz",
		},
	}

	for i, tc := range tests {
		s := []rune(tc.in)
		out := dodash(s)
		outs := string(out)
		if outs != tc.out {
			t.Errorf("tc: #%d: oczekiwano: %q, jest: %q", i, tc.out, outs)
		}
	}
}
