// 2014-10-24 Adam Bryt

// Narzędzia Programistyczne w Pascalu,
// rozdział 4.4 "Sortowanie szybkie",
// program sortquick.
//
// NAZWA
//
// sortquick - sortuje wiersze tekstu
//
// SPOSÓB UŻYCIA
//
// sortquick [<file1] [>file2]
//
// OPIS
//
// Program sortquick czyta wiersze tekstu z stdin, sortuje je i drukuje
// na stdout. Wszystkie wiersze są trzymane w pamięci i sortowane metodą
// Quicksort. Wiersze są porównywane zgodnie z kolejnością znaków języka
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

// quicksort sortuje a metodą quicksort zgodnie z kolejnością (collate)
// określoną przez c.
func quicksort(a []string, c *collate.Collator) {
	if len(a) < 2 {
		return
	}
	xi := len(a) - 1 // indeks pivota
	x := a[xi]       // wartość pivota
	i, j := 0, len(a)-1
	// zamieniaj miejscami elementy większe i mniejsze od pivota
	for {
		if i == j {
			break
		}
		for {
			if i == j || c.CompareString(a[i], x) > 0 {
				break
			}
			i++
		}
		for {
			if j == i || c.CompareString(a[j], x) < 0 {
				break
			}
			j--
		}
		if i < j {
			a[i], a[j] = a[j], a[i]
		}
	}
	// wstaw pivot w miejsce spotkania i j
	a[i], a[xi] = a[xi], a[i]
	quicksort(a[:i], c)
	quicksort(a[i:], c)
}

func main() {
	lines, err := gtext(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	c := collate.New(language.Polish)
	quicksort(lines, c)

	err = ptext(os.Stdout, lines)
	if err != nil {
		log.Fatal(err)
	}
}
