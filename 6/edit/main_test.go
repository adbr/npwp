// 2015-07-21 Adam Bryt

package main

import (
	"fmt"
	"testing"
)

func TestSkipSpace(t *testing.T) {
	type test struct {
		s string // string wejściowy
		w int    // liczba początkowych białych znaków
	}
	var tests = []test{
		{"", 0},
		{" ", 1},
		{"  ", 2},
		{"\t", 1},
		{" \t\v", 3},
		{"aaa", 0},
		{" ąąą", 1},
		{"  ąąą", 2},
		{"   ąąą", 3},
		{" \t 123", 3},
		{" \t\v  123   ", 5},
	}

	for _, tc := range tests {
		w := skipSpace(tc.s)
		if w != tc.w {
			t.Errorf("skipSpace(%q) = %d, oczekiwano %d",
				tc.s, w, tc.w)
		}
	}
}

func TestParseNumber(t *testing.T) {
	type testCase struct {
		s     string // argument wejściowy
		num   int
		width int
		err   error
	}

	tests := []testCase{
		{"123", 123, 3, nil},
		{"5", 5, 1, nil},
		{"0", 0, 1, nil},
		{"+123", 123, 4, nil},
		{"-123", -123, 4, nil},
		{"000", 0, 3, nil},
		{"+00", 0, 3, nil},
		{"-00", 0, 3, nil},
		{"00123", 123, 5, nil},
		{" 123", 123, 4, nil},
		{"   +123", 123, 7, nil},
		{"   -000", 0, 7, nil},
		{"   \t\v123", 123, 8, nil},
		{"123xxx", 123, 3, nil},
		{"123.4", 123, 3, nil},
		{"123  ", 123, 3, nil},
		{"1234ąę", 1234, 4, nil},
		{"-123+", -123, 4, nil},
		{"   +123xxx", 123, 7, nil},
		{"", 0, 0, errNotNumber},
		{"abc123", 0, 0, errNotNumber},
		{"   abc", 0, 0, errNotNumber},
		{"   +abc", 0, 0, errNotNumber},
		{"   -abc", 0, 0, errNotNumber},
		{"+a", 0, 0, errNotNumber},
		{"ą", 0, 0, errNotNumber},
		{"++123", 0, 0, errNotNumber},
	}

	for _, tc := range tests {
		name := fmt.Sprintf("parseNumber(%q)", tc.s)
		check := func(t *testing.T) {
			n, w, err := parseNumber(tc.s)
			if n != tc.num || w != tc.width || err != tc.err {
				t.Errorf("(%d, %d, %#v), oczekiwano: (%d, %d, %#v)",
					n, w, err, tc.num, tc.width, tc.err)
			}
		}
		t.Run(name, check)
	}
}

func TestGetnum(t *testing.T) {
	type testCase struct {
		s     string
		num   int
		width int
		err   error
	}

	// ustawienie lnums dla testów
	lnums.curln = 22
	lnums.lastln = 55

	tests := []testCase{
		{"12print", 12, 2, nil},
		{"12", 12, 2, nil},
		{".print", 22, 1, nil},
		{".", 22, 1, nil},
		{"$print", 55, 1, nil},
		{"$", 55, 1, nil},
		{"0print", 0, 1, nil},
		{" \t12print", 12, 4, nil},
		{" .print", 22, 2, nil},
		{"  $", 55, 3, nil},
		{" ..", 22, 2, nil},
		{"$.", 55, 1, nil},
		{"print", 0, 0, errNotNumber},
		{"", 0, 0, errNotNumber},
		// TODO: obsługa wzorca
	}

	for _, tc := range tests {
		name := fmt.Sprintf("getnum(%q)", tc.s)
		check := func(t *testing.T) {
			n, w, err := getnum(tc.s)
			if n != tc.num || w != tc.width || err != tc.err {
				t.Errorf("wynik: (%d, %d, %#v), oczekiwano: (%d, %d, %#v)",
					n, w, err, tc.num, tc.width, tc.err)
			}
		}
		t.Run(name, check)
	}
}

func TestGetone(t *testing.T) {
	type testCase struct {
		s     string // argument wejściowy
		num   int
		width int
		err   error
	}

	// przykładowe wartości dla testów
	lnums.curln = 22
	lnums.lastln = 55

	tests := []testCase{
		{".+3print", 25, 3, nil},
		{".-2print", 20, 3, nil},
		{"$-3print", 52, 3, nil},
		{"$+2print", 57, 3, nil},
		{"2+5print", 7, 3, nil},
		{"  .+3print", 25, 5, nil},
		{"  $-3print", 52, 5, nil},
		{"  2+5print", 7, 5, nil},
		{"print", 0, 0, errNotNumber},
		{"+print", 0, 0, errNotNumber},
		{" -print", 0, 0, errNotNumber},
		{"2+print", 0, 0, &syntaxError{
			line: "2+print",
			pos:  2,
			err:  errMissingNumber,
		}},
		{"2-print", 0, 0, &syntaxError{
			line: "2-print",
			pos:  2,
			err:  errMissingNumber,
		}},
		{" .+print", 0, 0, &syntaxError{
			line: " .+print",
			pos:  3,
			err:  errMissingNumber,
		}},
		{" .-print", 0, 0, &syntaxError{
			line: " .-print",
			pos:  3,
			err:  errMissingNumber,
		}},
		// TODO: spacje dookoła operatorów?
	}

	for _, tc := range tests {
		name := fmt.Sprintf("getone(%q)", tc.s)
		check := func(t *testing.T) {
			n, w, err := getone(tc.s)
			if n != tc.num || w != tc.width {
				t.Errorf("wynik: (%d, %d), oczekiwano: (%d, %d)",
					n, w, tc.num, tc.width)
			}
			if e0, ok := tc.err.(*syntaxError); ok {
				e1, ok := err.(*syntaxError)
				if !ok || *e1 != *e0 {
					t.Errorf("error: %#v, oczekiwano: %#v",
						*e1, *e0)
				}
			} else {
				if err != tc.err {
					t.Errorf("error: %#v, oczekiwano: %#v",
						err, tc.err)
				}
			}
		}
		t.Run(name, check)
	}
}

func TestGetlist(t *testing.T) {
	type testCase struct {
		s     string // argument wejściowy
		ln0   Lnums  // wartość początkowa zmiennej globalnej lnums
		ln1   Lnums  // wartość oczkiwana zmiennej globalnej lnums
		width int    // dłogość sparsowanego stringu
		err   error  // czy i jaki włąd powinien wystąpić
	}

	tests := []testCase{
		// kilka numerów wierszy
		{
			"12,34print",
			Lnums{
				line1:  1,
				line2:  2,
				nlines: 0,
				curln:  50,
				lastln: 123,
			},
			Lnums{
				line1:  12,
				line2:  34,
				nlines: 2,
				curln:  50, // TODO: czy curln nie zmienione?
				lastln: 123,
			},
			5,
			nil,
		},
		{
			"12,34,56print",
			Lnums{
				line1:  0,
				line2:  0,
				nlines: 0,
				curln:  50,
				lastln: 123,
			},
			Lnums{
				line1:  34,
				line2:  56,
				nlines: 2,
				curln:  50,
				lastln: 123,
			},
			8,
			nil,
		},
		{
			"12,34,56,789print",
			Lnums{
				line1:  1,
				line2:  1,
				nlines: 1,
				curln:  50,
				lastln: 123,
			},
			Lnums{
				line1:  56,
				line2:  789,
				nlines: 2,
				curln:  50,
				lastln: 123,
			},
			12,
			nil,
		},
		// jeden numer wiersza
		{
			"12print",
			Lnums{
				line1:  0,
				line2:  0,
				nlines: 0,
				curln:  50,
				lastln: 123,
			},
			Lnums{
				line1:  12,
				line2:  12,
				nlines: 1,
				curln:  50,
				lastln: 123,
			},
			2,
			nil,
		},
		// zero numerów wierszy
		// TODO: błąd w teście - poprawić
		//{
		//	"print",
		//	Lnums{
		//		line1:  1,
		//		line2:  2,
		//		nlines: 2,
		//		curln:  50,
		//		lastln: 123,
		//	},
		//	Lnums{
		//		line1:  50,
		//		line2:  50,
		//		nlines: 0,
		//		curln:  50,
		//		lastln: 123,
		//	},
		//	0,
		//	nil,
		//},
		// kilka numerów wierszy oddzielonych średnikiem
		// TODO: błąd w teście - poprawić
		//{
		//	"12;34;567print",
		//	Lnums{
		//		line1:  1,
		//		line2:  1,
		//		nlines: 1,
		//		curln:  50,
		//		lastln: 123,
		//	},
		//	Lnums{
		//		line1:  34,
		//		line2:  567,
		//		nlines: 2,
		//		curln:  34,
		//		lastln: 123,
		//	},
		//	9,
		//	nil,
		//},
		// kilka numerów, przecinek i średnik
		{
			"12;34,567print",
			Lnums{
				line1:  1,
				line2:  1,
				nlines: 1,
				curln:  50,
				lastln: 123,
			},
			Lnums{
				line1:  34,
				line2:  567,
				nlines: 2,
				curln:  12,
				lastln: 123,
			},
			9,
			nil,
		},
		// wyrażenia: dwa numery wiersza
		{
			".+1,$-2print",
			Lnums{
				line1:  1,
				line2:  1,
				nlines: 1,
				curln:  5,
				lastln: 50,
			},
			Lnums{
				line1:  6,
				line2:  48,
				nlines: 2,
				curln:  5,
				lastln: 50,
			},
			7,
			nil,
		},
		// wyrażenia: nietypowe: bez . i $
		{
			"1+2,5-1print",
			Lnums{
				line1:  1,
				line2:  1,
				nlines: 1,
				curln:  5,
				lastln: 50,
			},
			Lnums{
				line1:  3,
				line2:  4,
				nlines: 2,
				curln:  5,
				lastln: 50,
			},
			7,
			nil,
		},
		// wyrażenia: jeden numer wiersza: .
		{
			".print",
			Lnums{
				line1:  1,
				line2:  1,
				nlines: 1,
				curln:  5,
				lastln: 50,
			},
			Lnums{
				line1:  5,
				line2:  5,
				nlines: 1,
				curln:  5,
				lastln: 50,
			},
			1,
			nil,
		},
		// wyrażenia: jeden numer wiersza: $
		{
			"$print",
			Lnums{
				line1:  1,
				line2:  1,
				nlines: 1,
				curln:  5,
				lastln: 50,
			},
			Lnums{
				line1:  50,
				line2:  50,
				nlines: 1,
				curln:  5,
				lastln: 50,
			},
			1,
			nil,
		},
		// wyrażenia: więcej niż dwa numery wierszy
		{
			".-2,.+2,$print",
			Lnums{
				line1:  1,
				line2:  1,
				nlines: 1,
				curln:  5,
				lastln: 50,
			},
			Lnums{
				line1:  7,
				line2:  50,
				nlines: 2,
				curln:  5,
				lastln: 50,
			},
			9,
			nil,
		},
		// wyrażenia: więcej niż dwa numery, z średnikiem
		// kolejne wyrażenia zmieniają wartość '.' czyli curln
		// TODO: błąd w teście - poprawić
		//{
		//	".-2;.+3;.+$print",
		//	Lnums{
		//		line1:  1,
		//		line2:  1,
		//		nlines: 1,
		//		curln:  5,
		//		lastln: 50,
		//	},
		//	Lnums{
		//		line1:  6,
		//		line2:  56,
		//		nlines: 2,
		//		curln:  6,
		//		lastln: 50,
		//	},
		//	11,
		//	nil,
		//},
		// wyrażenia: kilka operatorów w numerze (takie
		// wyrażenia nie działają, nie są poprawnie
		// obsługiwane)
		// TODO: czy powinien być zgłoszony błąd?
		{
			".+2-3+$-2print",
			Lnums{
				line1:  1,
				line2:  1,
				nlines: 1,
				curln:  5,
				lastln: 50,
			},
			Lnums{
				line1:  7,
				line2:  7,
				nlines: 1,
				curln:  5,
				lastln: 50,
			},
			3, // pozostały string "-3+$-2print"
			nil,
		},
		// wyrażenia: brak pierwszego operandu
		// wyrażenia: brak drugiego operandu
		// wyrażenia: nie poprawny operator, różny od '+' i '-'
	}

	for _, tc := range tests {
		name := fmt.Sprintf("getlist(%q)", tc.s)
		check := func(t *testing.T) {
			lnums = tc.ln0 // globalna zmienna lnums
			w, err := getlist(tc.s)
			if lnums != tc.ln1 {
				t.Errorf("lnums: %v, oczekiwano: %v", lnums, tc.ln1)
			}
			if w != tc.width || err != tc.err {
				t.Errorf("wynik: (%d, %#v), oczekiwano: (%d, %#v)",
					w, err, tc.width, tc.err)
			}
		}
		t.Run(name, check)
	}
}
