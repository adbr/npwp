// 2014-11-23 Adam Bryt

// Narzędzia Programistyczne w Pascalu,
// rozdział 4.8 "Indeks permutacyjny",
// program kwic.
//
// NAZWA
//
// kwic - generuje wiersze do indeksu KWIC
//
// SPOSÓB UŻYCIA
//
// kwic [<file1] [>file2]
//
// OPIS
//
// Program kwic czyta wiersze tekstu z stdin i dla każdego wiersza
// drukuje na stdout permutacje wiersza. Permutacje są tworzone tak,
// że dla każdego słowa w wierszu tworzy się wiersz przesunięty
// cyklicznie tak, że to słowo znajduje się na początku wiersza, a
// oryginalny znak nowego wiersza jest zastąpiony przez znak '$'.
// Program jest używany łącznie z programami sort i unrotate w celu
// utworzenia indeksu KWIC.
//
// PRZYKŁADY
//
// Wynik działania programu dla przykładowego wiersza:
//
//	kwik
//	This is a test.
//	This is a test.$
//	is a test.$This
//	a test.$This is
//	test.$This is a
//
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

const (
	fold = "$" // oznacza koniec wiersza w kolejnych permutacjach
)

// rotate zapisuje wiersz line do w zaczynając od indeksu i. Początek
// wiersza line przed i jest zapisywany na końcu. Logiczny koniec
// wiersza jest zaznaczany znakiem fold.
func rotate(w io.Writer, line string, i int) error {
	// usuń końcowy znak newline jeśli jest
	line = strings.TrimSuffix(line, "\n")

	// drukuj koniec wiersza - od i do końca
	_, err := fmt.Fprint(w, line[i:])
	if err != nil {
		return err
	}

	// drukuj znak oznaczający logiczny koniec wiersza
	_, err = fmt.Fprint(w, fold)
	if err != nil {
		return err
	}

	// drukuj początek wiersza - od początku do i
	_, err = fmt.Fprint(w, line[:i])
	if err != nil {
		return err
	}

	// drukuj znak nowego wiersza (fizyczny)
	_, err = fmt.Fprint(w, "\n")
	if err != nil {
		return err
	}

	return nil
}

func isAlphanum(c rune) bool {
	if unicode.IsLetter(c) || unicode.IsNumber(c) {
		return true
	}
	return false
}

// permutate tworzy permutacje wiersza tekstu line i zapisuje je do w.
func permutate(w io.Writer, line string) error {
	inword := false
	for i, c := range line {
		if isAlphanum(c) && !inword {
			err := rotate(w, line, i)
			if err != nil {
				return err
			}
			inword = true
		}
		if !isAlphanum(c) {
			inword = false
		}
	}
	return nil
}

func kwic(w io.Writer, r io.Reader) error {
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			if len(line) > 0 {
				err := permutate(w, line)
				if err != nil {
					return err
				}
			}
			return nil
		}

		err = permutate(w, line)
		if err != nil {
			return err
		}
	}
}

func main() {
	err := kwic(os.Stdout, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
}
