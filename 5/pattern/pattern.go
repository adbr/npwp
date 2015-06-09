// 2015-01-15 Adam Bryt

// Pakiet pattern zawiera funkcjonalność związaną z tworzeniem
// i dopasowywaniem wzorców.
package pattern

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

const (
	maxChars = 255 // maksymalna liczba znaków z klasie znaków (1 bajt)
)

// Stałe oznaczające znaki wyróżnione występujące we wzorcu źródłowym.
const (
	s_bol     = '%'
	s_eol     = '$'
	s_any     = '?'
	s_ccl     = '['
	s_cclend  = ']'
	s_negate  = '^'
	s_closure = '*'
	s_esc     = '@'
)

// Stałe oznaczające tagi elementów wzorca w postaci skompilowanej.
const (
	litchar byte = iota
	bol
	eol
	any
	ccl
	nccl
	closure
)

// Typ Pattern reprezentuje skompilowany wzorzec.
type Pattern string

// String zwraca stringową reprezentację skompilowanego wzorca.
func (p Pattern) String() string {
	var out []byte
	for {
		if len(p) == 0 {
			break
		}

		t := p[0]
		p = p[1:]

		switch t {
		case bol:
			out = append(out, "<BOL>"...)
		case eol:
			out = append(out, "<EOL>"...)
		case any:
			out = append(out, "<ANY>"...)
		case litchar:
			out = append(out, "<LITCHAR>"...)
			_, n := utf8.DecodeRuneInString(string(p))
			out = append(out, p[:n]...)
			p = p[n:]
		case ccl:
			out = append(out, "<CCL>"...)
			nc := p[0] // liczba znaków
			p = p[1:]
			for i := byte(0); i < nc; i++ {
				_, n := utf8.DecodeRuneInString(string(p))
				out = append(out, p[:n]...)
				p = p[n:]
			}
		case nccl:
			out = append(out, "<NCCL>"...)
			nc := p[0] // liczba znaków
			p = p[1:]
			for i := byte(0); i < nc; i++ {
				_, n := utf8.DecodeRuneInString(string(p))
				out = append(out, p[:n]...)
				p = p[n:]
			}
		case closure:
			out = append(out, "<CLOSURE>"...)
		default:
			panic(fmt.Sprintf("Pattern.String: nie znany tag: %d", t))
		}
	}
	return string(out)
}

// Makepat kompiluje wzorzec str do reprezentacji wewnętrznej Pattern.
func Makepat(str string) (Pattern, error) {
	var out []byte
	s := str[:]
	last := 0 // początek ostatnio dodanego wzorca

	for {
		if len(s) == 0 {
			break
		}
		r, n := utf8.DecodeRuneInString(s)

		switch {
		case r == s_bol && len(s) == len(str):
			last = len(out)
			out = append(out, bol)
			s = s[n:]
		case r == s_eol && len(s) == n:
			last = len(out)
			out = append(out, eol)
			s = s[n:]
		case r == s_any:
			last = len(out)
			out = append(out, any)
			s = s[n:]
		case r == s_ccl:
			last = len(out)
			var (
				chars []byte
				isneg bool
				err   error
			)
			chars, isneg, s, err = getccl(s)
			if err != nil {
				return "", err
			}
			if isneg {
				out = append(out, nccl)
			} else {
				out = append(out, ccl)
			}
			nr := utf8.RuneCount(chars)
			if nr > maxChars {
				return Pattern(out), fmt.Errorf("klasa zawiera więcej niż %d znaków: %d", maxChars, nr)
			}
			out = append(out, byte(nr))
			out = append(out, chars...)
		case r == s_closure && len(out) > 0:
			tag := out[last]
			if tag == bol ||
				tag == eol ||
				tag == closure {
				return Pattern(out), errors.New("'*' nie może być po BOL, EOL, CLOSURE")
			}
			out = stclose(out, last)
			s = s[n:]
		default:
			last = len(out)
			out = append(out, litchar)

			var c rune
			c, s = esc(s)
			out = appendUtf8(out, c)
		}
	}

	return Pattern(out), nil
}

// stclose dodaje znacznik closure do wzorca pat przed segmentem
// zaczynającym się od indeksu last.
func stclose(pat []byte, last int) []byte {
	pat = append(pat, 0)               // zwiększenie rozmiaru pat
	_ = copy(pat[last+1:], pat[last:]) // przesunięcie w prawo
	pat[last] = closure
	return pat
}

// esc zwraca pierwszy znak ze stringu s z uwzgędnieniem escape'owania.
// Zwraca string pomniejszony o sekwencję escapeową.
func esc(s string) (rune, string) {
	r, n := utf8.DecodeRuneInString(s)
	s = s[n:]

	if r != s_esc {
		return r, s
	}

	if len(s) == 0 {
		return r, s
	}

	r, n = utf8.DecodeRuneInString(s)
	s = s[n:]

	switch r {
	case 't':
		return '\t', s
	case 'n':
		return '\n', s
	default:
		return r, s
	}
}

// getccl zwraca w chars znaki tworzące klasę znaków między '['
// i ']', rozwijając zakresy znaków ASCII typu 'a-z' i rozwijając
// sekwencje escapeowe. isneg ma wartość true jeśli pierwszym znakiem
// po '[' jest '^', czyli klasa znaków jest zanegowana. s zostaje
// pomniejszone o skonsumowane znaki.
func getccl(ss string) (chars []byte, isneg bool, s string, err error) {
	s = ss
	r, n := utf8.DecodeRuneInString(s)
	if r != s_ccl {
		err = fmt.Errorf("zły wzorzec: oczekiwano: %q, jest: %q", s_ccl, r)
		return
	}
	s = s[n:]

	r, n = utf8.DecodeRuneInString(s)
	if r == s_negate {
		isneg = true
		s = s[n:]
	}

	chars, s, err = dodash(s, s_cclend)
	if err != nil {
		return
	}

	return
}

// dodash rozwija zakresy znaków ASCII typu 'a-z', rozwija sekwencje
// escapeowe aż do ogranicznika delim i zwraca wynik w chars. Zmniejsza
// s o przetworzone znaki.
func dodash(ss string, delim rune) (chars []byte, s string, err error) {
	s = ss
	for {
		if len(s) == 0 {
			err = fmt.Errorf("zły wzorzec: brak ogranicznika %q", delim)
			return
		}
		r, n := utf8.DecodeRuneInString(s)
		if r == delim {
			s = s[n:]
			return
		}
		if r == '-' {
			r1, _ := utf8.DecodeLastRune(chars)
			r2, n2 := utf8.DecodeRuneInString(s[n:])
			if isAlphanum(r1) && isAlphanum(r2) && r1 <= r2 {
				for i := r1 + 1; i <= r2; i++ {
					chars = appendUtf8(chars, i)
				}
				s = s[n+n2:]
			} else {
				chars = appendUtf8(chars, r)
				s = s[n:]
			}
		} else {
			var c rune
			c, s = esc(s)
			chars = appendUtf8(chars, c)
		}
	}
}

// appendUtf8 wstawia rune r na koniec b jako utf8.
func appendUtf8(b []byte, r rune) []byte {
	var a [utf8.UTFMax]byte
	n := utf8.EncodeRune(a[:], r)
	b = append(b, a[:n]...)
	return b
}

// isAlphanum sprawdza czy znak r jest alfanumerycznym znakiem ASCII.
func isAlphanum(r rune) bool {
	if 'a' <= r && r <= 'z' {
		return true
	}
	if 'A' <= r && r <= 'Z' {
		return true
	}
	if '0' <= r && r <= '9' {
		return true
	}
	return false
}
