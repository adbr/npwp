// 2014-06-10 Adam Bryt

// Narzędzia Programistyczne w Pascalu, rozdział 1.5 "Usuwanie znaków
// tabulacji", program detab. Wersja rozszerzona - dodałem obsługę
// plików i opcje do ustawiania progów tabulacji.
//
// NAZWA
//
// detab - zamienia znaki tabulacji na spacje
//
// SPOSÓB UŻYCIA
//
// detab [-hi] [-t tablist] [file ...]
//
// OPIS
//
// Program detab czyta tekst z podanych plików, zamienia znaki tabulacji
// na odpowiednią liczbę spacji i drukuje wynikowy tekst na stdout.
// Jeśli nie podano plików to czyta dane z stdin. Domyślnie progi
// tabulacji są ustawione w co 8 kolumnie, ale można to zmienić opcją
// -t. Poprawnie obsługuje tekst zakodowany w UTF-8.
//
// Opcje:
//
//	-h
//	   Wyświetla krótki help
//
//	-i
//	   Przetwarza tylko początkowe znaki tabulacji (indent only)
//
//	-t tablist
//	   Ustawia listę progów tabulatora. Jeśli tablist jest jedną
//	   liczbą całkowitą to ustawia odstęp między progami tabulatora
//	   (domyślnie 8).  Jeśli tablist jest listą liczb oddzielonych
//	   przecinkiem (a,b,c,...) to ustawia progi tabulacji w tych
//	   kolumnach.
//
// PRZYKŁADY
//
// Zamienia znaki tabulacji w pliku file; progi tabulatora są ustawione
// standardowo, w co 8 kolumnie:
//
//	$ detab file
//
// Zamienia znaki tabulacji na spacje w pliku file; progi tabulacji są
// ustawione w co 4 kolumnie:
//
//	$ detab -t 4 file
//
// Zamienia znaki tabulacji na spacje w pliku file; progi tabulacji są
// ustawione w kolumnach 3, 6, 12. Znaki tabulacji poza ostatnim progiem
// są zamieniane na pojedynczą spacje:
//
//	$ detab -t 3,6,12 file
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

const usageStr = "usage: detab [-hi] [-t tablist] [file ...]"
const helpStr = usageStr + `

Program detab zamienia znaki tabulacji na spacje.

Opcje:
	-h 
	   Wyświetla ten help

	-i
	   Przetwarza tylko początkowe znaki tabulacji (indent only)

	-t tablist
	   Ustawia listę progów tabulatora. Jeśli tablist jest jedną
	   liczbą całkowitą to ustawia odstęp między progami tabulatora
	   (domyślnie 8).  Jeśli tablist jest listą liczb oddzielonych
	   przecinkiem (a,b,c,...) to ustawia progi tabulatora w tych
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
	fmt.Fprintf(os.Stderr, "detab: %s\n", err)
	os.Exit(1)
}

// Funkcja detab czyta tekst z r, zamienia znaki '\t' na odpowiednią
// liczbę spacji i zapisuje wynik do w. Położenie progów tabulacji jest
// określone w t. Jeśli argument indentOnly ma wartość true to są
// zamieniane tylko początkowe znaki tabulacji w wierszu (tylko znaki
// tab tworzące wcięcie). Argument pad zawiera znak jakim są zastępowane
// znaki tab - normalnie jest to spacja, a argument pad jest używany w
// testach dla ułatwienia liczenia znaków. W przypadku błędu zwraca
// error różny od nil.
func detab(w io.Writer, r io.Reader, t tab.Tabulator, indentOnly bool, pad rune) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()
	br := bufio.NewReader(r)

	process := true // wyłącznik przetwarzania znaków tab
	col := 0        // licznik kolumn
	const (
		tabchar = '\t'
		newline = '\n'
	)

	for {
		r, _, err := br.ReadRune()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if r != tabchar && indentOnly {
			process = false
		}

		switch {
		case r == tabchar && process:
			// Drukuj spacje aż do progu tabulacji
			for {
				_, err := bw.WriteRune(pad)
				if err != nil {
					return err
				}
				col++
				if t.IsTab(col) || t.BeyondLastTab(col) {
					break
				}
			}
		case r == newline:
			_, err := bw.WriteRune(r)
			if err != nil {
				return err
			}
			process = true
			col = 0
		default:
			_, err := bw.WriteRune(r)
			if err != nil {
				return err
			}
			col++
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

	const space = ' '
	if flag.NArg() == 0 {
		err := detab(os.Stdout, os.Stdin, tabul, indentOnlyFlag, space)
		if err != nil {
			fatal(err)
		}
	} else {
		for _, fname := range flag.Args() {
			file, err := os.Open(fname)
			if err != nil {
				fatal(err)
			}

			err = detab(os.Stdout, file, tabul, indentOnlyFlag, space)
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
