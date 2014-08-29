// 2014-08-19 Adam Bryt

// Narzędzia Programistyczne w Pascalu, rozdział 2.4 "Rozgęszczanie tekstu",
// program expand.
//
// NAZWA
//
// expand - dekompresuje (rozgęszcza) tekst
//
// SPOSÓB UŻYCIA
//
// expand [<file >file2]
//
// OPIS
//
// Program expand dekompresuje tekst skompresowany programem compress.
// Czyta tekst z stdin i zastępuje trójki znaków typu ~Nx sekwencją
// znaków x. Liczba znaków w sekwencji (N) jest zakodowana w
// pojedynczym znaku: A oznacza 1, B - 2, ..., Z - 26. Tekst wynikowy
// zapisuje do stdout.
//
// PRZYKŁADY
//
// Jeśli na wejściu wystąpi tekst:
//
//	"abc~Ed"
//
// to na wyjściu pojawi się:
//
//	"abcddddd"
//
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	warn = '~' // oznacza początek zakodowanej sekwencji znaków
)

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "expand: %s\n", err)
	os.Exit(1)
}

// Funkcja expand czyta tekst z r, rozgęszcza go (dekompresuje) i zapisuje
// do w. Zakłada się, że tekst został skompresowany programem compress,
// czyli sekwencje powtarzających się znaków są kodowane w postaci
// ~Nx (N to liczba powtórzeń znaku x zakodowana literą A do Z, gdzie
// A to 1, B - 2, ..., Z - 26). Zwraca błąd gdy wystąpił. W przypadku
// błędnych lub nie pełnych kodów sekwencji, wypisuje wypisuje wczytane
// znaki bez ich interpretacji.
func expand(w io.Writer, r io.Reader) error {
	// wrappery buffio dla czytania i pisania runów
	bw := bufio.NewWriter(w)
	defer bw.Flush()
	br := bufio.NewReader(r)

	for {
		c, _, err := br.ReadRune()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// zwykły znak
		if c != warn {
			_, err := bw.WriteRune(c)
			if err != nil {
				return err
			}
			continue
		}

		// zakodowana sekwencja

		// wczytaj znak-licznik
		cn, _, err := br.ReadRune()
		if err == io.EOF {
			_, err := bw.WriteRune(c)
			if err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			return err
		}

		if !walidCounter(cn) {
			_, err := bw.WriteRune(c)
			if err != nil {
				return err
			}
			_, err = bw.WriteRune(cn)
			if err != nil {
				return err
			}
			continue
		}

		// wczytaj znak sekwencji
		cc, _, err := br.ReadRune()
		if err == io.EOF {
			_, err := bw.WriteRune(c)
			if err != nil {
				return err
			}
			_, err = bw.WriteRune(cn)
			if err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			return err
		}

		// wypisz sekwencję
		n := cn - 'A' + 1
		for ; n > 0; n-- {
			_, err := bw.WriteRune(cc)
			if err != nil {
				return err
			}
		}
	}
}

// Funkcja walidCounter sprawdza czy znak c jest poprawnym znakiem
// reprezentującym krotność sekwencji.
func walidCounter(c rune) bool {
	if 'A' <= c && c <= 'Z' {
		return true
	} else {
		return false
	}
}

func main() {
	err := expand(os.Stdout, os.Stdin)
	if err != nil {
		fatal(err)
	}
}
