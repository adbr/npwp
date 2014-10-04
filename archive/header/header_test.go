// 2014-10-04 Adam Bryt

package header

import (
	"testing"
)

func TestNew(t *testing.T) {
	// header pliku istniejącego
	h, err := New("testfiles/a.txt")
	if err != nil {
		t.Error(err)
	}
	o := Header{ // oczekiwany header
		"-h-",
		"testfiles/a.txt",
		15,
	}
	if *h != o {
		t.Errorf("oczekiwano: %v jest %v", o, *h)
	}

	// header pliku nie istniejącego
	h, err = New("testfiles/notexists")
	if err == nil {
		t.Error("oczekiwano błedu: plik nie istnieje")
	}
}

func TestParse(t *testing.T) {
	type test struct {
		s string // string wejściowy
		e bool   // czy powinien wystąpić błąd
		h Header // oczekiwany header
	}
	tests := []test{
		{
			"-h- abc.txt 123\n",
			false,
			Header{
				"-h-",
				"abc.txt",
				123,
			},
		},
		{
			// różne odstępy między polami
			// działa, ale nie powinno wystąpić w archiwum
			"  -h-   a.txt\t123 \n",
			false,
			Header{
				"-h-",
				"a.txt",
				123,
			},
		},
		{
			// zerowa długość pliku
			"-h- a.txt 0\n",
			false,
			Header{
				"-h-",
				"a.txt",
				0,
			},
		},
		{
			// ujemna długość pliku
			"-h- a.txt -2\n",
			true,
			Header{},
		},
		{
			// zła ilość pól w nagłówku
			"-h- a.txt\n",
			true,
			Header{},
		},
		{
			// zły znacznik nagłówka (magic string)
			"h- a.txt 123\n",
			true,
			Header{},
		},
	}

	for i, tc := range tests {
		h, err := Parse(tc.s)
		if tc.e { // powinien wystąpić błąd
			if err == nil {
				t.Errorf("tc %d: powinien wystąpić błąd", i)
			}
		} else {
			if err != nil {
				t.Error(err)
				continue
			}
			if *h != tc.h {
				t.Errorf("tc %d: oczekiwano: %v jest: %v", i, tc.h, *h)
			}
		}
	}
}

func TestString(t *testing.T) {
	type test struct {
		h Header
		s string
	}
	tests := []test{
		{
			Header{}, // zero value header
			"  0\n",
		},
		{
			Header{
				"-h-",
				"a.txt",
				123,
			},
			"-h- a.txt 123\n",
		},
	}

	for i, tc := range tests {
		s := tc.h.String()
		if s != tc.s {
			t.Errorf("tc %d: oczekiwano: %q jest: %q", i, tc.s, s)
		}
	}
}
