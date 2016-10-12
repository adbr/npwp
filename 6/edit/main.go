// 2015-06-24 Adam Bryt

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"unicode"
	"unicode/utf8"
)

var ErrNotNumber = errors.New("koniec danych wejściowych")

func usage() {
	usageStr := "sposób użycia: edit [plik]"
	fmt.Fprintln(os.Stderr, usageStr)
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		w, err := getlist(line)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Print(lnums)
		fmt.Printf("reszta wiersza: %s\n", line[w:])
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// Typ Lnums zawiera informacje o numerach wierszy dla następnego
// polecenia.
type Lnums struct {
	line1  int // pierwszy numer wiersza
	line2  int // drugi numer wiersza
	nlines int // liczba podanych numerów wierszy
	curln  int // wiersz bieżący - wartość kropki
	lastln int // wiersz ostatni - wartość $
}

// Funkcja pomocnicza zwracająca wartość typu Lnums jako string.
// Implementuje interfejs fmt.Stringer.
func (l Lnums) String() string {
	s := ""
	s += fmt.Sprint("Lnums: {\n")
	s += fmt.Sprintf("\tline1:\t%d\n", l.line1)
	s += fmt.Sprintf("\tline2:\t%d\n", l.line2)
	s += fmt.Sprintf("\tnlines:\t%d\n", l.nlines)
	s += fmt.Sprintf("\tcurln:\t%d\n", l.curln)
	s += fmt.Sprintf("\tlastln:\t%d\n", l.lastln)
	s += fmt.Sprint("}\n")
	return s
}

// Zmienna globalna lnums zawiera informacje o numerach wierszy
// ostatniego polecenia.
var lnums Lnums

// getlist parsuje listę wyrażeń opisujących numery wierszy w stringu
// s i ustawia pola w zmiennej globalnej lnums. Zwraca długość listy
// wyrażeń w stringu s i błąd jeśli wystąpił. Elementy listy mogą myć
// oddzielone znakiem ',' lub ';'. Przykłady listy wyrażeń: '12,34',
// '12;34', '12,23,45', '.+1,$-2'.
func getlist(s string) (width int, err error) {
	i := 0 // indeks w stringu s
	lnums.nlines = 0

	num, w, err := getone(s[i:])
	if err != nil {
		return w, err
	}
	lnums.line2 = num
	lnums.nlines++
	i += w

	for {
		r, w := utf8.DecodeRuneInString(s[i:])
		if (r != ',') && (r != ';') {
			break
		}
		i += w
		if r == ';' {
			lnums.curln = num
		}

		num, w, err := getone(s[i:])
		if err != nil {
			return w, err
		}
		lnums.line1 = lnums.line2
		lnums.line2 = num
		lnums.nlines++
		i += w
	}
	
	if lnums.nlines > 2 {
		lnums.nlines = 2
	}
	if lnums.nlines == 1 {
		lnums.line1 = lnums.line2
	}
	if lnums.nlines == 0 {
		lnums.line1 = lnums.curln
		lnums.line2 = lnums.curln
	}
	return i, nil
}

// getone parsuje wyrażenie opisujące numer wiersza znajdujące się na
// początku stringu s. Zwraca obliczony numer wiersza, długość
// wyrażenia w stringu s i błąd jeśli wystąpił. Wyrażenie może
// zawierać operatory '+' i '-'. Przykłady wyrażeń: '.+3', '$-5',
// '5+1', '5'.
func getone(s string) (num, width int, err error) {
	i := 0 // indeks w stringu s

	// pierwszy operand (numer) musi wystąpić
	num, w, err := getnum(s)
	if err != nil {
		return num, w, err
	}
	i += w

	// czy wystąpił operator?
	r, w := utf8.DecodeRuneInString(s[i:])
	switch r {
	case '+':
		i += w
		n, w, err := getnum(s[i:])
		if err != nil {
			return n, w, err
		}
		i += w
		num += n
	case '-':
		i += w
		n, w, err := getnum(s[i:])
		if err != nil {
			return n, w, err
		}
		i += w
		num -= n
	}
	
	return num, i, nil
}

// getnum parsuje numer wiersza znajdujący się na początku stringu s.
// Zwraca numer wiersza, jego długość w stringu s i błąd jeśli
// wystąpił. Numer wiersza może być liczbą całkowitą, znakiem '.'
// (kropka), znakiem '$' (dolar) lub wzorcem. Używa zmiennej globalnej
// lnums (tylko do czytania).
func getnum(s string) (num, width int, err error) {
	r, w := utf8.DecodeRuneInString(s)
	switch r {
	case '.':
		num = lnums.curln
		return num, w, nil
	case '$':
		num = lnums.lastln
		return num, w, nil
	default:
		// TODO: obsługa wzorca
		num, w, err := parseNumber(s)
		return num, w, err
	}
}

// parseNumber parsuje liczbę całkowitą znajdującą się na początku
// stringu s. Zwraca liczbę, jej długość w stringu i błąd jeśli
// wystąpił. Białe znaki występujące przed liczbą są pomijane; liczba
// może być poprzedzona znakiem + lub -; parsowanie liczby kończy się
// po napotkania znaku nie będącego cyfrą lub końca stringu. Nie
// sprawdza przepełnienia gdy liczba w stringu jest większa niż
// maksymalna wartość typu int. Gdy na początku stringu nie ma liczby
// to zwraca błąd ErrNotNumber oraz num i width równe 0.
func parseNumber(s string) (num, width int, err error) {
	i := 0 // indeks w stringu s

	// pomiń początkowe spacje
	for {
		r, w := utf8.DecodeRuneInString(s[i:])
		if !unicode.IsSpace(r) {
			break
		}
		i += w
	}

	// parsuj znak liczby
	sign := 1
	r, w := utf8.DecodeRuneInString(s[i:])
	switch r {
	case '+':
		i += w
	case '-':
		i += w
		sign = -1
	}

	// parsuj liczbę całkowitą
	n := 0
	isnum := false
	for {
		r, w := utf8.DecodeRuneInString(s[i:])
		if !unicode.IsDigit(r) {
			break
		}
		d := int(r - '0')
		n = n*10 + d
		isnum = true
		i += w
	}
	n *= sign

	if !isnum {
		return 0, 0, ErrNotNumber
	}
	return n, i, nil
}
