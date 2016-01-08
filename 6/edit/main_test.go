// 2015-07-21 Adam Bryt

package main

import (
	"testing"
)

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
		lnums Lnums  // oczekiwana wartość zmiennej globalnej lnums
		ii    int    // indeks znaku po sparsowaniu
		isErr bool   // czy powinien wystąpić błąd
	}

	tests := []TestCase{
		{
			"12,34print",
			0,
			Lnums{
				line1:  12,
				line2:  34,
				nlines: 2,
				curln:  34, // ?
				lastln: 0,  // ?
			},
			5,
			false,
		},
	}

	for i, tc := range tests {
		lnums = Lnums{}
		ii, err := getlist(tc.lin, tc.i)

		if lnums != tc.lnums {
			t.Errorf("tc %d: getlist(): lnums: oczekiwano: %v, jest: %v", i, tc.lnums, lnums)
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
