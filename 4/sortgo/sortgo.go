// 2014-11-07 Adam Bryt

// Sortowanie wewnętrzne (w pamięci) z wykorzystaniem funkcji bibliotecznej
// collator.SortStrings
//
// NAZWA
//
// sortshell - sortuje wiersze tekstu
//
// SPOSÓB UŻYCIA
//
// sortshell [<file1] [>file2]
//
// OPIS
//
// Program sortshell czyta wiersze tekstu z stdin, sortuje je i drukuje
// na stdout. Wszystkie wiersze są trzymane w pamięci i sortowane metodą
// Shella. Wiersze są porównywane zgodnie z kolejnością znaków języka
// polskiego (collate).
//
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"code.google.com/p/go.text/collate"
	"code.google.com/p/go.text/language"
)

// func czyta i zwraca wszystkie wiersze z r.
func gtext(r io.Reader) ([]string, error) {
	br := bufio.NewReader(r) // dla ReadString
	var lines []string
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			if len(line) > 0 {
				lines = append(lines, line)
			}
			return lines, nil
		}
		if err != nil {
			return lines, err
		}
		lines = append(lines, line)
	}
}

// ptext drukuje wiersze z lines do w.
func ptext(w io.Writer, lines []string) error {
	for _, line := range lines {
		_, err := fmt.Fprint(w, line)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lines, err := gtext(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	c := collate.New(language.Polish)
	c.SortStrings(lines)

	err = ptext(os.Stdout, lines)
	if err != nil {
		log.Fatal(err)
	}
}
