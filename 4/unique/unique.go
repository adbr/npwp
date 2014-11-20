// 2014-11-09 Adam Bryt
// TODO: usuwanie duplikatów w plikach nie posortowanych

// Narzędzia Programistyczne w Pascalu,
// rozdział 4.7 "Rozdzielanie funkcji - unique",
// program unique.
//
// NAZWA
//
// unique - usuwa przyległe duplikaty wierszy
//
// SPOSÓB UŻYCIA
//
// unique [opcje] [<file1] [>file2]
//
// OPIS
//
// Program unique czyta wiersze tekstu ze standardowego wejścia, usuwa
// jednakowe, sąsiadujące wiersze, zostawiając tylko pierwszy z
// nich i drukuje tekst wynikowy na standardowe wyjście. Program jest
// najbardziej przydatny dla posortowanego tekstu, w którym jednakowe
// wiersze znajdują się obok siebie.
//
// Opcje:
//
//	-n
//		przed wierszem jest drukowana liczba jego wystąpień
//	-d
//		drukuje tylko wiersze mające duplikaty
//	-u
//		drukuje tylko wiersze bez duplikatów (unikalne)
//
// PRZYKŁADY
//
// Usuwanie duplikatów w pliku file1 i zapisanie wyniku do pliku file2:
//
//	cat file1 | sort | unique >file2
//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	nFlag = flag.Bool("n", false, "drukuj liczbę wystąpień wiersza")
	dFlag = flag.Bool("d", false, "drukuj tylko wiersze mające duplikaty")
	uFlag = flag.Bool("u", false, "drukuj tylko wiersze unikalne")
)

// print drukuje wiersz line do w uwzględniając flagi programu. Parametr n
// zawiera liczbę wystąpień wiersza na wejściu.
func print(w io.Writer, line string, n int) error {
	if *dFlag && !*uFlag && n == 1 {
		return nil
	}
	if *uFlag && !*dFlag && n > 1 {
		return nil
	}

	const nFormat = "%4d %s" // format wydruku: liczba wystąpień, wiersz
	if *nFlag {
		_, err := fmt.Fprintf(w, nFormat, n, line)
		if err != nil {
			return err
		}
	} else {
		_, err := fmt.Fprint(w, line)
		if err != nil {
			return err
		}
	}
	return nil
}

// unique kopiuje wiersze tekstu z r do w pomijając duplikaty.
func unique(w io.Writer, r io.Reader) error {
	br := bufio.NewReader(r)
	var last string
	var cnt int
	for {
		line, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			if last != "" {
				err := print(w, last, cnt)
				if err != nil {
					return err
				}
			}
			if len(line) > 0 {
				err := print(w, last, 1)
				if err != nil {
					return err
				}
			}
			return nil
		}

		if line != last {
			if last != "" {
				err := print(w, last, cnt)
				if err != nil {
					return err
				}
			}
			last = line
			cnt = 1
		} else {
			cnt++
		}
	}
}

func main() {
	flag.Parse()
	err := unique(os.Stdout, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
}
