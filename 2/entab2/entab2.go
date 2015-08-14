// 2014-08-04 Adam Bryt

// Program entab - inny, prostszy algorytm.
//
// Wersja programu entab z algorytmem zamiany spacji na tabulacje według
// pomysłu z książki - struktura programu powinna odpowiadać strukturze
// danych w strumieniu wejściowym - zamiast stosowania stanów, flag i
// przełączników.
//
// Z tym podejściem powstał prostszy program, i jest łatwiejszy do
// zrozumienia.
//
// UWAGI
//
// Nie testowano przypadku gdy na wejściu są znaki tab - może źle
// działać.
//
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/adbr/npwp/tab"
)

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "entab2: %s\n", err)
	os.Exit(1)
}

func entab(w io.Writer, r io.Reader, t tab.Tabulator) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()
	br := bufio.NewReader(r)

	const (
		space   = ' '
		tabchar = '\t'
		newline = '\n'
	)

	col := 0   // nr bieżącej kolumny
	nsp := 0   // licznik spacji
	var c rune // aktualny znak
	var err error

	for {
		// przetwarzaj sekwencję 0 lub więcej spacji
		for {
			c, _, err = br.ReadRune()
			if err == io.EOF {
				// opróżnij licznik spacji
				for ; nsp > 0; nsp-- {
					_, err := bw.WriteRune(space)
					if err != nil {
						return err
					}
				}
				return nil
			}
			if err != nil {
				return err
			}
			if t.IsTab(col) && nsp > 0 {
				_, err := bw.WriteRune(tabchar)
				if err != nil {
					return err
				}
				nsp = 0
			}
			if c != space {
				// opróżnij licznik spacji
				for ; nsp > 0; nsp-- {
					_, err := bw.WriteRune(space)
					if err != nil {
						return err
					}
				}
				break
			}
			col++
			nsp++
		}

		// tutaj znak jest różny od spacji
		_, err := bw.WriteRune(c)
		if err != nil {
			return err
		}

		if c == newline {
			col = 0
		} else {
			col++
		}
	}
}

func main() {
	const tabsize = 8
	t, err := tab.NewTabulator([]int{tabsize})
	if err != nil {
		fatal(err)
	}

	err = entab(os.Stdout, os.Stdin, t)
	if err != nil {
		fatal(err)
	}
}
