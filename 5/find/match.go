// 2015-01-15 Adam Bryt

// Plik zawiera funkcjonalność związaną z dopasowywaniem wzorca.

package main

import (
	"fmt"
	"unicode/utf8"
)

// match dopasowuje wzorzec par w dowolnym miejscu stringu lin.
func match(lin string, pat pattern) bool {
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
func amatch(str string, offset int, pat pattern, j int) (bool, int) {
	i := offset
	for j < len(pat) {
		if pat[j] == CLOSURE {
			j += patsize(pat[j:]) // pomiń tag CLOSURE
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
func omatch(str string, i int, pat pattern, j int) (bool, int) {
	tag := pat[j]
	switch tag {
	case BOL:
		if i == 0 {
			return true, 0
		}
	case EOL:
		if i < len(str) && str[i] == '\n' {
			return true, 0
		}
	case ANY:
		r, n := utf8.DecodeRuneInString(str[i:])
		if r != utf8.RuneError && r != '\n' {
			return true, n
		}
	case LITCHAR:
		r1, n1 := utf8.DecodeRuneInString(str[i:])
		r2, _ := utf8.DecodeRuneInString(string(pat[j+1:]))
		if (r1 == r2) && (r1 != utf8.RuneError) {
			return true, n1
		}
	case CCL:
		r, n := utf8.DecodeRuneInString(str[i:])
		if locate(r, pat[j+1:]) {
			return true, n
		}
	case NCCL:
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
func locate(c rune, pat pattern) bool {
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
func patsize(pat pattern) int {
	if len(pat) == 0 {
		return 0
	}

	tag := pat[0]

	switch tag {
	case BOL, EOL, ANY:
		return 1
	case LITCHAR:
		_, n := utf8.DecodeRuneInString(string(pat[1:]))
		return 1 + n
	case CCL, NCCL:
		nc := int(pat[1]) // liczba znaków w klasie
		p := pat[2:]      // początek znaków w klasie
		b := 0            // licznik bajtów
		for i := 0; i < nc; i++ {
			_, n := utf8.DecodeRuneInString(string(p))
			b += n
			p = p[n:]
		}
		return 2 + b
	case CLOSURE:
		return 1
	default:
		panic(fmt.Sprintf("patsize(): nie znany tag: %d", tag))
	}
}
