// 2014-08-23 Adam Bryt

// Narzędzia Programistyczne w Pascalu,
// rozdział 2.6 "Transliteracja znaków",
// program translit.
//
// NAZWA
//
// translit - transliteruje (zastępuje) znaki
//
// SPOSÓB UŻYCIA
//
// translit [^]from [to]
//
// OPIS
//
// Program translit czyta tekst z stdin, zastępuje znaki występujące
// w zbiorze znaków "from" odpowiednimi znakami ze zbioru "to" i
// drukuje przetworzony tekst na stdout. W najprostszym, typowym
// przypadku, znaki występujące w parametrze "from" są zamieniane
// na odpowiadające im znaki z parametru "to"; pozostałe znaki są
// kopiowane bez zmian.
//
// Parametry "from" i "to" mogą zawierać specyfikacje zakresu znaków
// w postaci c1-c2, co oznacza: wszystkie znaki z zakresu od znaku c1
// do c2 włącznie. Oba znaki c1 i c2 powinny być albo cyframi, albo
// dużymi literami, albo małymi literami.
//
// Jeśli parametr "to" nie występuje to znaki z parametru "from" są
// usuwane z tekstu.
//
// Jeśli parametr "to" jest krótszy niż "from", to wszystkie znaki
// z "from" nie mające swoich odpowiedników w "to", oraz znak
// odpowiadający ostatniemu znakowi z "to", są zastępowane ostatnim
// znakiem z "to", a ponadto sekwencje takich znaków są redukowane do
// pojedynczego znaku.
//
// Jeśli parametr "from" jest poprzedzony znakiem ^, to wtedy parametr
// "from" oznacza wszystkie znaki nie występujące w tym parametrze
// (negacja). W takim przypadku znaki te zostaną zredukowane do
// ostatniego znaku z parametru to, lub usunięte jeśli parametr to nie
// występuje.
//
// W parametrach mogą występować sekwencje kodujące znaki specjalne:
// @n oznacza znak nowego wiersza, @t oznacza znak tabulacji. Sam znak @
// należy zapisać jako @@.
//
// PRZYKŁADY
//
// Zamiana znaków abc na xyz (a na x, b na y, c na z):
//
//	translit "abc" "xyz"
//
// Zamiana małych liter na duże:
//
//	translit "a-z" "A-Z"
//
// Usunięcie znaków nie będących cyframi:
//
//	translit "^0-9"
//
// Zastąpienie sekwencji spacji, znaków tabulacji i znaków nowego
// wiersza, pojedynczym znakiem nowego wiersza (umieszcza każdy wyraz
// w nowym wierszu):
//
//	translit " @t@n" "@n"
//
// UWAGI
//
// Przetwarza poprawnie znaki UTF-8, ale zakresy znaków w parametrze
// "from" są rozwijane zgodnie z kodowaniem ASCII.
//
// Zakłada się, że w specyfikacji zakresu znaków, znak początkowy
// i końcowy należą do jednej grupy: cyfry (0-9), duże litery (A-Z)
// i małe litery (a-z). Jeśli znak początkowy i końcowy zakresu
// należą do różnych grup, to po rozwinięciu zbiór może zawierać
// znaki nie alfanumeryczne, ponieważ w zbiorze ASCII grupy znaków
// alfanumerycznych nie sąsiadują ze sobą (np zakres "9-B" jest
// rozwijany do "9:;<=>?@AB").
//
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	escape = '@' // poprzedza znak wyróżniony, np: @t oznacza znak tabulacji
	negate = '^' // oznacza negację zbioru znaków, np: ^abc
	dash   = '-' // oznacza zakres znaków, np: A-Z
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: translit string [string]")
	os.Exit(1)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "translit: %s\n", err)
	os.Exit(2)
}

// Funkcja index zwraca pierwszą pozycję (indeks) znaku r w stringu s,
// lub -1 gdy r nie występuje w r.
func index(s []rune, r rune) int {
	for i, c := range s {
		if c == r {
			return i
		}
	}
	return -1
}

// Funkcja xindex zwraca indeks znaku r w stringu s, z uwzględnieniem
// negacji zbioru znaków s. Jeśli allbut jest false (brak negacji) to
// zwraca indeks znaku r w s, lub -1 jeśli r nie występuje w s. Jeśli
// allbut jest true (negacja znaków), to jeśli znak r występuje w s
// zwraca -1, a jeśli nie występuje w s, to zwraca lastto+1 (w celu
// zredukowania lub usunięcia znaku).
func xindex(s []rune, r rune, allbut bool, lastto int) int {
	if !allbut { // bez negacji - zachowanie jak index()
		return index(s, r)
	}

	// negacja znaków
	if index(s, r) >= 0 {
		return -1
	} else {
		return lastto + 1
	}
}

// Funkcja isAlphaNum zwraca true jeśli c jest literą lub cyfrą ASCII.
func isAlphaNum(c rune) bool {
	if ('0' <= c && c <= '9') ||
		('A' <= c && c <= 'Z') ||
		('a' <= c && c <= 'z') {
		return true
	} else {
		return false
	}
}

// Funkcja esc zwraca znak specjalny odpowiadający sekwencji @<znak>
// zaczynającej się od indeksu i w stringu s.
func esc(s []rune, i int) rune {
	if s[i] != escape {
		return s[i]
	}

	if i == len(s)-1 {
		// @ na końcu jest zwykłym znakiem
		return s[i]
	}

	switch s[i+1] {
	case 'n':
		return '\n'
	case 't':
		return '\t'
	default:
		return s[i+1]
	}
}

// Funkcja dodash rozwija skrótowy zapis przedziału znaków i rozwija
// sekwencje znaków specjalnych.
func dodash(s []rune) []rune {
	new := []rune{}
	for i := 0; i < len(s); i++ {
		if s[i] == escape {
			c := esc(s, i)
			new = append(new, c)
			i++
		} else if s[i] != dash {
			new = append(new, s[i])
		} else if (i == 0) || (i == len(s)-1) {
			new = append(new, s[i])
		} else if isAlphaNum(s[i-1]) &&
			isAlphaNum(s[i+1]) &&
			(s[i-1] <= s[i+1]) {
			for k := s[i-1] + 1; k <= s[i+1]; k++ {
				new = append(new, k)
			}
			i++
		} else {
			new = append(new, s[i])
		}
	}
	return new
}

// Funkcja translit czyta znaki z w, zastępuje znaki występujące
// w fromstr znakami z tostr i zapisuje do w. Zwraca błąd jeśli
// wystąpił.
func translit(w io.Writer, r io.Reader, fromstr string, tostr string) error {
	// wrapery umożliwiające czytanie i pisanie runów
	bw := bufio.NewWriter(w)
	defer bw.Flush()
	br := bufio.NewReader(r)

	// konwersja stringów do tablicy runów - dla dostępu przez indeks
	fromrune := []rune(fromstr)
	torune := []rune(tostr)

	// ustawienie flagi allbut (negacja znaków)
	allbut := false
	if len(fromrune) > 0 && fromrune[0] == negate {
		allbut = true
		fromrune = fromrune[1:]
	}

	// rozwinięcie zakresów i znaków wyróżnionych
	fromrune = dodash(fromrune)
	torune = dodash(torune)

	lastto := len(torune) - 1
	squash := (len(fromrune) > len(torune)) || allbut

	for {
		c, _, err := br.ReadRune()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return nil
		}

		i := xindex(fromrune, c, allbut, lastto)

		// zastępowanie ostatnim znakiem i redukcja
		if i >= lastto && lastto >= 0 && squash {
			cc := torune[lastto]
			_, err = bw.WriteRune(cc)
			if err != nil {
				return err
			}

			// zredukuj sekwencję
			for {
				c, _, err = br.ReadRune()
				if err == io.EOF {
					return nil
				}
				if err != nil {
					return err
				}

				i = xindex(fromrune, c, allbut, lastto)
				if i < lastto {
					break
				}
			}
		}

		if i >= 0 {
			if lastto >= 0 { // zastępowanie znaków
				cc := torune[i]
				_, err = bw.WriteRune(cc)
				if err != nil {
					return err
				}
			} else {
				// usuwanie znaków
			}
		} else if i == -1 { // kopiowanie znaków bez zmiany
			_, err := bw.WriteRune(c)
			if err != nil {
				return err
			}
		}
	}
}

func main() {
	var (
		from string
		to   string
	)
	if len(os.Args) == 3 {
		from = os.Args[1]
		to = os.Args[2]
	} else if len(os.Args) == 2 {
		from = os.Args[1]
		to = ""
	} else {
		usage()
	}

	err := translit(os.Stdout, os.Stdin, from, to)
	if err != nil {
		fatal(err)
	}
}
