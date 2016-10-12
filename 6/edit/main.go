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
		i, err := getlist(line, 0)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Print(lnums)
		fmt.Printf("reszta wiersza: %s\n", line[i:])
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

// getlist parsuje numery wierszy zawarte w lin począwszy od znaku o
// indeksie i. Po sparsowaniu numerów wierszy ustawia pola zmiennej
// globalnej lnums. Zwraca indeks znaku występującego po sparsowanym
// fragmancie lin i błąd jeśli wystąpił.
func getlist(lin string, i int) (ii int, err error) {
	var num int
	lnums.nlines = 0

	num, ii, err = getone(lin, i)
	if err != nil {
		return i, err
	}
	if ii > i {
		// istnieje co najmniej jeden numer wiersza
		// todo: EOF ?
		lnums.line2 = num
		lnums.nlines++
	}
	for {
		r, w := utf8.DecodeRuneInString(lin[ii:])
		if (r != ',') && (r != ';') {
			break
		}
		ii += w
		if r == ';' {
			lnums.curln = num
		}

		num, ii, err = getone(lin, ii)
		if err != nil {
			return i, err
		}
		lnums.line1 = lnums.line2
		lnums.line2 = num
		lnums.nlines++
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
	return
}

// getone parsuje wyrażenie opisujące numer wiersza zawarte w lin,
// zaczynając parsowanie od znaku o indeksie i. Wyrażenie może
// zawierać operatory + i -. Zwraca obliczony numer wiersza i następną
// pozycję w stringu lin. Przykłady wyrażeń: '.+3', '$-5'.
func getone(lin string, i int) (num int, ii int, err error) {
	num, ii, err = getnum(lin, i)
	if err != nil {
		return num, ii, err
	}
	r, w := utf8.DecodeRuneInString(lin[ii:])
	if r == '+' || r == '-' {
		ii += w
		var n int
		n, ii, err = getnum(lin, ii)
		if err != nil {
			return n, ii, err
		}
		if r == '+' {
			num += n
		}
		if r == '-' {
			num -= n
		}
	}
	return num, ii, err
}

// getnum parsuje jeden element wyrażenia opisującego numer wiersza
// (liczba całkowita, . (kropka), $ lub wzorzec) zaczynając od znaku o
// indeksie i. Zwraca pobrany numer i następną pozycję w stringu lin.
// Używa (tylko do czytania) zmiennej globalnej lnums.
func getnum(lin string, i int) (num int, ii int, err error) {
	r, w := utf8.DecodeRuneInString(lin[i:])
	if r == '.' {
		num = lnums.curln
		ii = i + w
		return num, ii, nil
	}
	if r == '$' {
		num = lnums.lastln
		ii = i + w
		return num, ii, nil
	}

	// todo: obsługa wzorca

	num, ii = strToNum(lin, i)
	return num, ii, nil
}

// strToNum parsuje liczbę całkowitą zawartą w s od indeksu i, zwraca
// wartość liczby jako num i indeks ii wskazujący na pierwszy znak
// poza liczbą. Kończy parsowanie na znaku nie będącym cyfrą. Liczba
// może być poprzedzona sekwencją białych znaków i znakiem liczby. Gdy
// na pozycji i nie ma liczby, zwraca 0 i indeks ii równy początkowemu
// indeksowi i, czyli w przypadku błędu zwraca 0 i indeks nie jest
// przsuwany.
func strToNum(s string, i int) (num, ii int) {
	ii = i
	// pomiń białe znaki
	for {
		r, w := utf8.DecodeRuneInString(s[i:])
		if !unicode.IsSpace(r) {
			break
		}
		i += w
	}

	// parsuj znak liczby
	r, w := utf8.DecodeRuneInString(s[i:])
	sign := 1
	switch r {
	case '+':
		i += w
	case '-':
		sign = -1
		i += w
	}

	// parsuj liczbę
	isNum := false
	for {
		r, w := utf8.DecodeRuneInString(s[i:])
		if !unicode.IsDigit(r) {
			break
		}
		d := int(r - '0')
		num = num*10 + d
		i += w
		isNum = true
	}

	num *= sign
	if isNum {
		// zmień indeks tylko gdy została wykryta liczba
		ii = i
	}
	return
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
