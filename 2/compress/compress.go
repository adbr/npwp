// 2014-08-09 Adam Bryt

// Narzędzia Programistyczne w Pascalu, rozdział 2.3 "Zagęszczanie
// tekstu", program compress.
//
// NAZWA
//
// compress - kompresuje tekst kodując serie znaków
//
// SPOSÓB UŻYCIA
//
// compress [<file >file2]
//
// OPIS
//
// Program compress czyta tekst z stdin, zastępuje sekwencje co
// najmniej czterech jednakowych znaków pewnym kodem i drukuje wynik na
// stdout. Dzięki temu tekst wyjściowy jest na ogół krótszy niż
// tekst wejściowy (skompresowany). Sekwencja znaków x jest kodowana
// jako ~nx, gdzie n jest zakodowaną liczbą znaków: A oznacza 1, B
// - 2, ..., Z - 26. Sekwencje dłuższe niż 26 są podzielone kilka
// krótszych. Sekwencje krótsze niż cztery znaki nie są kodowane.
// Sekwencje znaków ~ (tylda) zawsze są kodowane - nawet gdy są
// krótsze niż cztery znaki.
//
// PRZYKŁADY
//
// Jeśli na wejściu programu compress pojawi się tekst:
//
//	"aaabbbbb~cccc"
//
// to na wyjściu uzyskamy:
//
//	"aaa~Eb~A~~Dc"
//
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const warn = '~' // znacznik zakodowanej sekwencji znaków

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "compress: %s\n", err)
	os.Exit(1)
}

func compress(w io.Writer, r io.Reader) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()
	br := bufio.NewReader(r)

	// wczytaj pierwszy znak
	lastc, _, err := br.ReadRune()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	n := 1 // licznik znaków w sekwencji

	for {
		c, _, err := br.ReadRune()
		if err == io.EOF {
			if n > 1 || (lastc == warn) {
				err := putrep(bw, n, lastc)
				if err != nil {
					return err
				}
			} else {
				_, err := bw.WriteRune(lastc)
				if err != nil {
					return err
				}
			}
			return nil
		}
		if err != nil {
			return err
		}

		if c == lastc {
			n++
		} else if n > 1 || (lastc == warn) { // koniec sekwencji
			err := putrep(bw, n, lastc)
			if err != nil {
				return err
			}
			n = 1
		} else {
			_, err := bw.WriteRune(lastc)
			if err != nil {
				return err
			}
		}
		lastc = c
	}
}

// Funkcja putrep zapisuje do w zakodowaną, zwięzłą reprezentację
// sekwencji n znaków c. Zwraca błąd jeśli wystąpił.
func putrep(w *bufio.Writer, n int, c rune) error {
	const (
		thresh = 4 // koduje tylko sekwencje >= thresh
		maxrep = 'Z' - 'A' + 1
	)

	for {
		if c != warn && n < thresh {
			break
		}
		if c == warn && n <= 0 {
			break
		}

		// tylda
		_, err := w.WriteRune(warn)
		if err != nil {
			return err
		}

		// liczba znaków - zakodowana ('A' == 1, 'B' == 2, ...)
		cc := 'A' + min(n, maxrep) - 1
		_, err = w.WriteRune(rune(cc))
		if err != nil {
			return err
		}

		// znak tworzący sekwencje
		_, err = w.WriteRune(c)
		if err != nil {
			return err
		}

		n -= maxrep
	}

	// pozostało mniej niż thresh znaków - wypisz normalnie
	for ; n > 0; n-- {
		_, err := w.WriteRune(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func main() {
	err := compress(os.Stdout, os.Stdin)
	if err != nil {
		fatal(err)
	}
}
