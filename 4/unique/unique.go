// 2014-11-09 Adam Bryt

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
// Jeśli jest użyta flaga -kf lub -kl to stosowany jest inny algorytm:
// wiersze nie muszą być posortowane; cały tekst wczytywany jest do
// pamięci. Opcja -kf powoduje, że drukowany jest pierwszy wiersz
// z duplikatów, a opcja -kl powoduje drukowanie ostatniego wiersza
// z duplikatów. Ten algorytm może być nie efektywny w przypadku
// dużych plików - wymaga dużo pamięci i wykonuje dużo porównań
// wierszy.
//
// Opcje:
//
//	-n
//		przed wierszem jest drukowana liczba jego wystąpień
//	-d
//		drukuje tylko wiersze mające duplikaty
//	-u
//		drukuje tylko wiersze bez duplikatów (unikalne)
//	-kf
//		drukuje pierwszy wiersz z duplikatów (keep first)
//	-kl
//		drukuje ostatni wiersz z duplikatów (keep last)
//
// PRZYKŁADY
//
// Usuwanie duplikatów w pliku file1 i zapisanie wyniku do pliku file2:
//
//	cat file1 | sort | unique >file2
//
// Usuwanie duplikatów w pliku nie posortowanym. Dany jest plik
// zawierający następujące wiersze:
//
//	$ cat x.txt
//	xxxx
//	ccc
//	bbb
//	xxxx
//	ddd
//	xxxx
//	eee
//
// W wyniku polecenia:
//
//	$ cat x.txt | unique -kl
//
// zostaną zwrócone dane:
//
//	ccc
//	bbb
//	ddd
//	xxxx
//	eee
//
// czyli zostały usunięte początkowe wiersze (duplikaty) "xxxx".
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

	// keep first - nie wymaga posortowanych wierszy ale czyta całość do pamięci
	kfFlag = flag.Bool("kf", false, "drukuj tylko pierwszy wiersz z powtarzających się (keep first)")
	// keep last - nie wymaga posortowanych wierszy ale czyta całość do pamięci
	klFlag = flag.Bool("kl", false, "drukuj tylko ostatni wiersz z powtarzających się (keep last)")
)

// print drukuje wiersz line do w uwzględniając flagi programu.
// Parametr n zawiera liczbę wystąpień wiersza na wejściu.
func print(w io.Writer, line string, n int) error {
	if *dFlag && !*uFlag && n == 1 {
		return nil
	}
	if *uFlag && !*dFlag && n > 1 {
		return nil
	}

	const nFormat = "%4d %s"
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

// unique kopiuje wiersze tekstu z r do w pomijając duplikaty. Zakłada
// się, że wiersze w r są posortowane.
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
				// na pewno nie duplikat
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

// readLines czyta wszystkie wiersze tekstu z r i zwraca je jako
// []string.
func readLines(r io.Reader) ([]string, error) {
	br := bufio.NewReader(r)
	lines := []string{}
	for {
		line, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return lines, err
		}
		if err == io.EOF {
			if len(line) > 0 {
				lines = append(lines, line)
			}
			return lines, nil
		}
		lines = append(lines, line)
	}
}

// count zwraca liczbę wystąpień wiersza s w lines.
func count(s string, lines []string) int {
	c := 0
	for i := 0; i < len(lines); i++ {
		if s == lines[i] {
			c++
		}
	}
	return c
}

// uniqueInternal usuwa duplikaty wierszy z r i drukuje wynik do w.
// Tekst nie musi być posortowany i jest wczytywany w całości do
// pamięci. Flagi programu kf (keep first) i kl (keep last) powodują,
// że jest stosowana ta funkcja i decydują czy zachować pierwszy
// czy ostatni wiersz spośród duplikatów. Algorytm nie efektywny w
// przypadku dużych ilości tekstu z powodu dużej liczby porównań
// wierszy.
func uniqueInternal(w io.Writer, r io.Reader) error {
	lines, err := readLines(r)
	if err != nil {
		return err
	}

	switch {
	case *kfFlag:
		for i := 0; i < len(lines); i++ {
			n := count(lines[i], lines[:i+1])
			if n == 1 {
				nn := count(lines[i], lines)
				print(w, lines[i], nn)
				if err != nil {
					return err
				}
			}
		}
	case *klFlag:
		for i := 0; i < len(lines); i++ {
			n := count(lines[i], lines[i:])
			if n == 1 {
				nn := count(lines[i], lines)
				print(w, lines[i], nn)
				if err != nil {
					return err
				}
			}
		}
	default:
		panic("powinna być ustawiona flaga -kf lub -kl")
	}
	return nil
}

func main() {
	flag.Parse()
	if *kfFlag || *klFlag {
		err := uniqueInternal(os.Stdout, os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := unique(os.Stdout, os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
	}
}
