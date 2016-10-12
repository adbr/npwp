// 2015-07-21 Adam Bryt

package main

import (
	"fmt"
	"testing"
)

func TestParseNumber(t *testing.T) {
	type testCase struct {
		s     string // string wejściowy
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
		{"", 0, 0, ErrNotNumber},
		{"abc123", 0, 0, ErrNotNumber},
		{"   abc", 0, 0, ErrNotNumber},
		{"   +abc", 0, 0, ErrNotNumber},
		{"   -abc", 0, 0, ErrNotNumber},
		{"+a", 0, 0, ErrNotNumber},
		{"ą", 0, 0, ErrNotNumber},
		{"++123", 0, 0, ErrNotNumber},
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

func TestStrToNum(t *testing.T) {
	type Test struct {
		s   string // string wejściowy
		i   int    // indeks, początek liczby
		num int    // wartość liczby
		ii  int    // uaktualniony indeks
	}

	tests := []Test{
		// Przypadki poprawne
		{
			"123", 0, 123, 3,
		},
		{
			"a123", 1, 123, 4,
		},
		{
			"a123b4c", 1, 123, 4,
		},
		{
			"a0a", 1, 0, 2,
		},
		{
			"a+0a", 1, 0, 3,
		},
		{
			"a-0a", 1, 0, 3,
		},
		{
			"a000a", 1, 0, 4,
		},
		{
			"  123a", 0, 123, 5,
		},
		{
			"  +123a", 0, 123, 6,
		},
		{
			"a+123b", 1, 123, 5,
		},
		{
			"a-123b", 1, -123, 5,
		},
		{
			"a+123.45", 1, 123, 5,
		},
		// Przypadki błędne - indeks nie jest przesuwany
		{
			"a123", 0, 0, 0,
		},
		{
			"abc", 0, 0, 0,
		},
		{
			"ąęś", 0, 0, 0,
		},
		{
			"  abc", 0, 0, 0,
		},
		{
			"", 0, 0, 0,
		},
		{
			"   ", 0, 0, 0,
		},
		{
			// indeks poza stringiem
			"abc", 3, 0, 3,
		},
		{
			// sam znak, bez liczby
			"a+b", 1, 0, 1,
		},
		{
			"a+", 1, 0, 1,
		},
		{
			"a-b", 1, 0, 1,
		},
		{
			"a++b", 1, 0, 1,
		},
	}

	for i, test := range tests {
		num, ii := strToNum(test.s, test.i)
		if num != test.num || ii != test.ii {
			t.Errorf(": %d: strToNum(): jest: [%d, %d], oczekiwano: [%d, %d]",
				i, num, ii, test.num, test.ii)
		}
	}
}

func TestGetnum(t *testing.T) {
	type TestCase struct {
		lin   string
		i     int
		num   int
		ii    int
		isErr bool
	}

	tests := []TestCase{
		{
			"12,34print",
			0,
			12,
			2,
			false,
		},
		{
			" +12,34print",
			0,
			12,
			4,
			false,
		},
		{
			"-12,34print",
			0,
			-12,
			3,
			false,
		},
		{
			"-12,34print",
			4,
			34,
			6,
			false,
		},
		{
			".print",
			0,
			22, // wartość lnums.curln
			1,
			false,
		},
		{
			"$print",
			0,
			55, // wartość lnums.lastln
			1,
			false,
		},
	}

	// przykładowe wartości dla testów
	lnums.curln = 22
	lnums.lastln = 55

	for i, tc := range tests {
		num, ii, err := getnum(tc.lin, tc.i)
		if num != tc.num {
			t.Errorf("tc %d: getnum(): num == %d, oczekiwano: %d", i, num, tc.num)
		}
		if ii != tc.ii {
			t.Errorf("tc %d: getnum(): ii == %d, oczekiwano: %d", i, ii, tc.ii)
		}
		if err != nil && tc.isErr == false {
			t.Errorf("tc %d: wystąpił bład: %s", i, err)
		}
		if err == nil && tc.isErr == true {
			t.Errorf("tc %d: nie wystąpił oczekiwany błąd", i)
		}
	}
}

func TestGetone(t *testing.T) {
	type TestCase struct {
		lin   string
		i     int
		num   int
		ii    int
		isErr bool
	}

	tests := []TestCase{
		{
			".+3print",
			0,
			25, // lnums.curln + 3 (22+3)
			3,
			false,
		},
		{
			".-2print",
			0,
			20,
			3,
			false,
		},
		{
			"$-3print",
			0,
			52,
			3,
			false,
		},
		{
			"$+2print",
			0,
			57,
			3,
			false,
		},
		{
			"2+5print",
			0,
			7,
			3,
			false,
		},
	}

	// przykładowe wartości dla testów
	lnums.curln = 22
	lnums.lastln = 55

	for i, tc := range tests {
		num, ii, err := getone(tc.lin, tc.i)
		if num != tc.num {
			t.Errorf("tc %d: getone(): num == %d, oczekiwano: %d", i, num, tc.num)
		}
		if ii != tc.ii {
			t.Errorf("tc %d: getone(): ii == %d, oczekiwano: %d", i, ii, tc.ii)
		}
		if err != nil && tc.isErr == false {
			t.Errorf("tc %d: wystąpił błąd: %s", i, err)
		}
		if err == nil && tc.isErr == true {
			t.Errorf("tc %d: nie wystąpił oczekiwany błąd", i)
		}
	}
}

func TestGetlist(t *testing.T) {
	type TestCase struct {
		lin   string // string wejściowy
		i     int    // indeks początku parsowania
		ln0   Lnums  // wartość początkowa zmiennej globalnej lnums
		ln1   Lnums  // wartość oczkiwana zmiennej globalnej lnums
		ii    int    // indeks znaku po sparsowaniu
		isErr bool   // czy powinien wystąpić błąd
	}

	var tests = []TestCase{
		// kilka numerów wierszy
		{
			"12,34print",
			0,
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
				curln:  50,
				lastln: 123,
			},
			5,
			false,
		},
		{
			"12,34,56print",
			0,
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
			false,
		},
		{
			"12,34,56,789print",
			0,
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
			false,
		},
		// jeden numer wiersza
		{
			"12print",
			0,
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
			false,
		},
		// zero numerów wierszy
		{
			"print",
			0,
			Lnums{
				line1:  1,
				line2:  2,
				nlines: 2,
				curln:  50,
				lastln: 123,
			},
			Lnums{
				line1:  50,
				line2:  50,
				nlines: 0,
				curln:  50,
				lastln: 123,
			},
			0,
			false,
		},
		// kilka numerów wierszy oddzielonych średnikiem
		{
			"12;34;567print",
			0,
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
				curln:  34,
				lastln: 123,
			},
			9,
			false,
		},
		// kilka numerów, przecinek i średnik
		{
			"12;34,567print",
			0,
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
			false,
		},
		// wyrażenia: dwa numery wiersza
		{
			".+1,$-2print",
			0,
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
			false,
		},
		// wyrażenia: nietypowe: bez . i $
		{
			"1+2,5-1print",
			0,
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
			false,
		},
		// wyrażenia: jeden numer wiersza: .
		{
			".print",
			0,
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
			false,
		},
		// wyrażenia: jeden numer wiersza: $
		{
			"$print",
			0,
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
			false,
		},
		// wyrażenia: więcej niż dwa numery wierszy
		{
			".-2,.+2,$print",
			0,
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
			false,
		},
		// wyrażenia: więcej niż dwa numery, z średnikiem
		// kolejne wyrażenia zmieniają wartość '.' czyli curln
		{
			".-2;.+3;.+$print",
			0,
			Lnums{
				line1:  1,
				line2:  1,
				nlines: 1,
				curln:  5,
				lastln: 50,
			},
			Lnums{
				line1:  6,
				line2:  56,
				nlines: 2,
				curln:  6,
				lastln: 50,
			},
			11,
			false,
		},
		// wyrażenia: kilka operatorów w numerze (takie
		// wyrażenia nie działają, nie są poprawnie
		// obsługiwane)
		{
			".+2-3+$-2print",
			0,
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
			false,
		},
		// wyrażenia: brak pierwszego operandu
		// wyrażenia: brak drugiego operandu
		// wyrażenia: nie poprawny operator, różny od '+' i '-'
	}

	for i, tc := range tests {
		lnums = tc.ln0 // globalna zmienna lnums
		ii, err := getlist(tc.lin, tc.i)

		if lnums != tc.ln1 {
			t.Errorf("tc %d: getlist(): lnums: oczekiwano: %v, jest: %v", i, tc.ln1, lnums)
		}

		if ii != tc.ii {
			t.Errorf("tc %d: getlist(): indeks po sparsowaniu: oczekiwano: %d, jest %d", i, tc.ii, ii)
		}

		if err != nil && tc.isErr == false {
			t.Errorf("tc %d: getlist(): wystąpił błąd: %s", i, err)
		}

		if err == nil && tc.isErr == true {
			t.Errorf("tc %d: getlist(): nie wystąpił oczekiwany błąd", i)
		}
	}
}
