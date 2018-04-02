// 2015-06-24 Adam Bryt

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"unicode"
	"unicode/utf8"
)

var errNotNumber = errors.New("not number")
var errMissingNumber = errors.New("missing number")

type syntaxError struct {
	line string // parsowany string
	pos  int    // miejsce (indeks) w line zawierające błąd
	err  error  // rodzaj błędu
}

func (e *syntaxError) Error() string {
	return fmt.Sprintf("syntax error in %q on position %d: %v",
		e.line, e.pos, e.err)
}

const usageText = "sposób użycia: edit [plik]"

func main() {
	h := flag.Bool("h", false, "display usage")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usageText)
	}
	flag.Parse()
	if *h {
		fmt.Println(usageText)
		os.Exit(0)
	}

	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		line := stdin.Text()
		w, err := getlist(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "edit: %s\n", err)
			continue
		}
		fmt.Print(lnums)
		fmt.Printf("reszta wiersza: %s\n", line[w:])
	}
	if err := stdin.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "edit: %s\n", err)
		os.Exit(1)
	}
}

// Typ Lnums zawiera informacje o numerach wierszy dla polecenia.
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
// '5+1', '5'. Jeśli nie ma numeru wiersza to zwraca błąd errNotNumber
// i num oraz width równe 0. Jeśli
//
// TODO: zgłaszanie błędów składniowych - np. po operatorze brakuje
// liczby
func getone(s string) (num, width int, err error) {
	i := 0 // indeks w stringu s

	// pierwszy operand
	num, w, err := getnum(s)
	if err != nil {
		return 0, 0, err
	}
	i += w

	// czy wystąpił operator?
	r, w := utf8.DecodeRuneInString(s[i:])
	switch r {
	case '+':
		i += w
		n, w, err := getnum(s[i:])
		if err == errNotNumber {
			return 0, 0, &syntaxError{
				line: s,
				pos:  i,
				err:  errMissingNumber}
		}
		if err != nil {
			return 0, 0, err
		}
		i += w
		num += n
	case '-':
		i += w
		n, w, err := getnum(s[i:])
		if err == errNotNumber {
			return 0, 0, &syntaxError{
				line: s,
				pos:  i,
				err:  errMissingNumber}
		}
		if err != nil {
			return 0, 0, err
		}
		i += w
		num -= n
	}

	return num, i, nil
}

// getnum parsuje numer wiersza znajdujący się na początku stringu s.
// Zwraca numer wiersza, jego długość w stringu s i błąd jeśli
// wystąpił. Numer wiersza może być liczbą całkowitą (jak w funkcji
// parseNumber), znakiem '.', znakiem '$' lub wzorcem. Pomija
// początkowe białe znaki. Używa zmiennej globalnej lnums (tylko do
// odczytu) w celu pobrania wartości dla '.' i '$'. Jeśli na początku
// stringu nie ma numeru wiersza to zwraca błąd errNotNumber oraz num
// i width równe 0.
//
// TODO: obsługa wzorca
func getnum(s string) (num, width int, err error) {
	i := 0 // indeks w stringu s

	w := skipSpace(s)
	i += w

	r, w := utf8.DecodeRuneInString(s[i:])
	switch r {
	case '.':
		i += w
		return lnums.curln, i, nil
	case '$':
		i += w
		return lnums.lastln, i, nil
	default:
		num, w, err = parseNumber(s[i:])
		if err != nil {
			return 0, 0, err
		}
		i += w
		return num, i, nil
	}
}

// parseNumber parsuje liczbę całkowitą znajdującą się na początku
// stringu s. Zwraca liczbę, jej długość w stringu i błąd jeśli
// wystąpił. Białe znaki występujące przed liczbą są pomijane; liczba
// może być poprzedzona znakiem + lub -; parsowanie liczby kończy się
// po napotkania znaku nie będącego cyfrą lub końca stringu. Nie
// sprawdza przepełnienia gdy liczba w stringu jest większa niż
// maksymalna wartość typu int. Gdy na początku stringu nie ma liczby
// to zwraca błąd errNotNumber oraz num i width równe 0.
func parseNumber(s string) (num, width int, err error) {
	i := 0 // indeks w stringu s

	// pomiń początkowe spacje
	w := skipSpace(s)
	i += w

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
		return 0, 0, errNotNumber
	}
	return n, i, nil
}

// skipSpace zwraca długość (w bajtach) początkowych białych znaków w
// stringu s. Białym znakiem jest znak spełniający warunek
// unicode.IsSpace(). Mając długość początkowych białych znaków w,
// można je pominąć przy użyciu wyrażenia s[w:].
func skipSpace(s string) int {
	i := 0
	for {
		r, w := utf8.DecodeRuneInString(s[i:])
		if !unicode.IsSpace(r) {
			break
		}
		i += w
	}
	return i
}
