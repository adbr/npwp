// 2014-06-26 Adam Bryt

// Narzędzia Programistyczne w Pascalu, rozdział 2.1 "Przywracanie
// znaków tabulacji", program entab. Wersja rozszerzona - dodałem
// obsługę plików i opcje do ustawiania progów tabulacji.
//
// NAZWA
//
// entab - zamienia sekwencje spacji na znaki tabulacji
//
// SPOSÓB UŻYCIA
//
// entab [-hi] [-t tablist] [file ...]
//
// OPIS
//
// Program entab czyta tekst z podanych plików, zamienia sekwencje
// spacji na znaki tabulacji i drukuje wynikowy tekst na stdout. Jeśli
// nie podano plików to czyta dane z stdin. Domyślnie progi tabulacji są
// ustawione w co 8 kolumnie, ale można to zmienić opcją -t. Poprawnie
// obsługuje tekst zakodowany w UTF-8.
//
// Opcje:
//
//	-h
//	   Wyświetla krótki help
//
//	-i
//	   Przetwarza tylko początkowe spacje (indent only)
//
//	-t tablist
//	   Ustawia listę progów tabulatora. Jeśli tablist jest jedną
//	   liczbą całkowitą to ustawia odstęp między progami tabulatora
//	   (domyślnie 8).  Jeśli tablist jest listą liczb oddzielonych
//	   przecinkiem (a,b,c,...) to ustawia progi tabulatora w tych
//	   kolumnach.
//
// PRZYKŁADY
//
// Zamienia spacje na znaki tabulacji w pliku file; progi tabulatora są
// ustawione standardowo, w co 8 kolumnie:
//
//	$ entab file
//
// Zamienia spacje na znaki tabulacji w pliku file; progi tabulatora są
// ustawione w co 4 kolumnie:
//
//	$ entab -t 4 file
//
// Zamienia spacje na znaki tabulacji w pliku file; progi tabulatora są
// ustawione w kolumnach 3, 6, 12. Spacje poza ostatnim progiem
// tabulacji nie są zamieniane:
//
//	$ entab -t 3,6,12 file
//
// UWAGI
//
// Nie testowano przypadku gdy na wejściu są znaki tabulacji - może źle
// działać.
//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/adbr/npwp/tab"
)

const usageStr = "usage: entab [-hi] [-t tablist] [file ...]"
const helpStr = usageStr + `

Program entab zamienia sekwencje spacji na znaki tabulacji.

Opcje:
	-h
	   Wyświetla ten help

	-i
	   Przetwarza tylko początkowe spacje (indent only)

	-t tablist
	   Ustawia listę progów tabulatora. Jeśli tablist jest jedną
	   liczbą całkowitą to ustawia odstęp między progami tabulatora
	   (domyślnie 8).  Jeśli tablist jest listą liczb oddzielonych
	   przecinkiem (a,b,c,...) to ustawia progi tabulacji w tych
	   kolumnach.
`

func usage() {
	fmt.Fprintln(os.Stderr, usageStr)
	os.Exit(1)
}

func help() {
	fmt.Println(helpStr)
	os.Exit(0)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "entab: %s\n", err)
	os.Exit(1)
}

// Funkcja entab czyta tekst z r, zamienia sekwencje spacji na znaki
// tabulacji i zapisuje wynikowy tekst do w. Położenie progów tabulatora
// jest określone w t. Jeśli indentOnly ma wartość true to są zamieniane
// tylko początkowe spacje w wierszu (spacje tworzące wcięcie). W
// przypadku błędu zwraca error różny od nil.
func entab(w io.Writer, r io.Reader, t tab.Tabulator, indentOnly bool) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()
	br := bufio.NewReader(r)

	const (
		space   = ' '
		tabchar = '\t'
		newline = '\n'
	)

	process := true // wyłącznik przetwarzania spacji
	col := 0        // nr bieżącej kolumny
	nsp := 0        // licznik spacji
	for {
		r, _, err := br.ReadRune()
		if err == io.EOF {
			// opróżnij "bufor" spacji
			for nsp > 0 {
				_, err = bw.WriteRune(space)
				if err != nil {
					return err
				}
				nsp--
			}
			return nil
		}
		if err != nil {
			return err
		}

		if process && r != space && indentOnly {
			// wyłącz przetwarzanie spacji - koniec wcięcia
			process = false
		}

		// Jeśli przed progiem tabulacji są 2 lub więcej spacji
		// to zastąp je znakiem tab.
		if t.IsTab(col) && nsp >= 2 {
			_, err = bw.WriteRune(tabchar)
			if err != nil {
				return err
			}
			nsp = 0
		}

		// Jeśli przed progiem tabulacji jest jedna spacja i
		// znak na progu tabulacji nie jest spacją to nie
		// zastępuj jednej spacji nakiem tabulacji. Zachowanie
		// jak unexpand z OpenBSD.
		if t.IsTab(col) && nsp == 1 && r != space {
			_, err = bw.WriteRune(space)
			if err != nil {
				return err
			}
			nsp = 0
		}

		// Jeśli przed progiem tabulacji jest jedna spacja i
		// znak na progu tabulacji jest spacją to zastąp jedną
		// spację znakiem tabulacji. Zachowanie jak unexpand z
		// OpenBSD.
		if t.IsTab(col) && nsp == 1 && r == space {
			_, err = bw.WriteRune(tabchar)
			if err != nil {
				return err
			}
			nsp = 0
		}

		if r == space && process {
			nsp++
			col++
		} else {
			// opróżnij "bufor" spacji
			for nsp > 0 {
				_, err = bw.WriteRune(space)
				if err != nil {
					return err
				}
				nsp--
			}

			_, err = bw.WriteRune(r)
			if err != nil {
				return err
			}
			col++
		}

		if r == newline {
			col = 0
			nsp = 0
			process = true
		}
	}
}

func main() {
	var (
		helpFlag       bool
		tablistFlag    string
		indentOnlyFlag bool
	)

	flag.BoolVar(&helpFlag, "h", false, "wyświetl help")
	flag.BoolVar(&helpFlag, "help", false, "wyświetl help")
	flag.StringVar(&tablistFlag, "t", "8", "lista pozycji tabulatora")
	flag.BoolVar(&indentOnlyFlag, "i", false, "przetwarzaj tylko wcięcia")

	flag.Usage = usage
	flag.Parse()

	if helpFlag {
		help()
	}

	tablist, err := tab.ParseTablist(tablistFlag)
	if err != nil {
		fatal(err)
	}
	tabul, err := tab.NewTabulator(tablist)
	if err != nil {
		fatal(err)
	}

	if flag.NArg() == 0 {
		err := entab(os.Stdout, os.Stdin, tabul, indentOnlyFlag)
		if err != nil {
			fatal(err)
		}
	} else {
		for _, fname := range flag.Args() {
			file, err := os.Open(fname)
			if err != nil {
				fatal(err)
			}

			err = entab(os.Stdout, file, tabul, indentOnlyFlag)
			if err != nil {
				file.Close()
				fatal(err)
			}

			err = file.Close()
			if err != nil {
				fatal(err)
			}
		}
	}
}
