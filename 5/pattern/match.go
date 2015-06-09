// 2015-01-15 Adam Bryt

// Plik zawiera funkcjonalność związaną z dopasowywaniem wzorca.

package pattern

import (
	"fmt"
	"unicode/utf8"
)

// Match dopasowuje wzorzec pat w dowolnym miejscu stringu lin.
func Match(lin string, pat Pattern) bool {
	for i := range lin {
		ok, _ := amatch(lin, i, pat, 0)
		if ok {
			return true
		}
	}
	return false
}

// amatch dopasowuje wzorzec zaczynający się od pat[j] do stringu
// zaczynającego się od str[offset]. Jeśli wzorzec pasuje, to zwraca
// true i liczbę bajtów str pasujących do wzorca.
func amatch(str string, offset int, pat Pattern, j int) (bool, int) {
	i := offset
	for j < len(pat) {
		if pat[j] == closure {
			j += patsize(pat[j:]) // pomiń tag closure
			ii := i

			// dopasuj maksymalną ilość znaków do pat[j]
			for {
				ok, n := omatch(str, ii, pat, j)
				if ok {
					ii = ii + n
				} else {
					break
				}
			}

			// dopasuj pozostałą część stringu do pozostałej części
			// wzorca; jeśli się nie da do cofaj się w stringu
			for ii >= i {
				ok, n := amatch(str, ii, pat, j+patsize(pat[j:]))
				if ok {
					return true, n + ii - offset
				} else {
					// cofnij string o jeden znak
					_, n := utf8.DecodeLastRuneInString(str[:ii])
					ii = ii - n
				}
			}
			return false, 0
		} else {
			if ok, n := omatch(str, i, pat, j); ok {
				i += n
				j += patsize(pat[j:])
			} else {
				return false, 0
			}
		}
	}
	return true, i - offset
}

// omatch dopasowuje jeden element wzorca zaczynający się od pat[j] do
// stringu zaczynającego się od str[i]. Jeśli element wzorca pasuje
// to zwraca true i liczbę bajtów str pasujących do elementu wzorca.
func omatch(str string, i int, pat Pattern, j int) (bool, int) {
	tag := pat[j]
	switch tag {
	case bol:
		if i == 0 {
			return true, 0
		}
	case eol:
		if i < len(str) && str[i] == '\n' {
			return true, 0
		}
	case any:
		r, n := utf8.DecodeRuneInString(str[i:])
		if r != utf8.RuneError && r != '\n' {
			return true, n
		}
	case litchar:
		r1, n1 := utf8.DecodeRuneInString(str[i:])
		r2, _ := utf8.DecodeRuneInString(string(pat[j+1:]))
		if (r1 == r2) && (r1 != utf8.RuneError) {
			return true, n1
		}
	case ccl:
		r, n := utf8.DecodeRuneInString(str[i:])
		if locate(r, pat[j+1:]) {
			return true, n
		}
	case nccl:
		r, n := utf8.DecodeRuneInString(str[i:])
		if !locate(r, pat[j+1:]) && r != '\n' {
			return true, n
		}
	default:
		panic(fmt.Sprintf("omatch(): nie znany tag: %d", tag))
	}
	return false, 0
}

// locate sprawdza czy znak c znajduje się w klasie znaków zawartej w
// pat. Pierwszy bajt w pat zawiera liczbę znaków a pozostałe bajty
// zawierają znaki zakodowane w utf8.
func locate(c rune, pat Pattern) bool {
	num := int(pat[0])
	pat = pat[1:]
	for i := 0; i < num; i++ {
		r, n := utf8.DecodeRuneInString(string(pat))
		if r == c {
			return true
		}
		pat = pat[n:]
	}
	return false
}

// patsize zwraca rozmiar w bajtach pierwszego segmentu wzorca pat.
func patsize(pat Pattern) int {
	if len(pat) == 0 {
		return 0
	}

	tag := pat[0]

	switch tag {
	case bol, eol, any:
		return 1
	case litchar:
		_, n := utf8.DecodeRuneInString(string(pat[1:]))
		return 1 + n
	case ccl, nccl:
		nc := int(pat[1]) // liczba znaków w klasie
		p := pat[2:]      // początek znaków w klasie
		b := 0            // licznik bajtów
		for i := 0; i < nc; i++ {
			_, n := utf8.DecodeRuneInString(string(p))
			b += n
			p = p[n:]
		}
		return 2 + b
	case closure:
		return 1
	default:
		panic(fmt.Sprintf("patsize(): nie znany tag: %d", tag))
	}
}
