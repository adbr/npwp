// 2014-09-04 Adam Bryt

// Narzędzia Programistyczne w Pascalu,
// rozdział 3.5 "Redaktor wydruków",
// program print.
//
// NAZWA
//
// print - drukuje pliki z nagłówkami
//
// SPOSÓB UŻYCIA
//
// print [file ...]
//
// OPIS
//
// Program print czyta tekst z podanych plików i drukuje na stdout,
// dodając nagłówek i dolny margines na każdej stronie. Nagłówek
// zawiera datę, nazwę pliku i numer strony. Każdy z podanych plików
// zaczyna się od nowej strony z nową numeracją stron od 1. Jeśli
// nie podano plików to drukuje dane ze standardowego wejścia. Program
// nie zmienia tekstu - nie łamie długich wierszy, nie zastępuje
// znaków tabulacji.
//
// PRZYKŁADY
//
// Drukowanie plików file1 file2:
//
//	print file1 file2
//
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: print [file ...]")
	os.Exit(1)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "print: %s\n", err)
	os.Exit(2)
}

// skip drukuje na stdout n pustych wierszy.
func skip(n int) {
	for i := 0; i < n; i++ {
		fmt.Printf("\n")
	}
}

// head drukuje na stdout nagłówek strony, zawierający: nazwę pliku
// name, aktualną datę i czas oraz numer strony page.
func head(name string, page int) {
	const format = "2006-01-02 15:04"
	t := time.Now().Format(format)

	fmt.Printf("%s  %s  page %d\n", name, t, page)
}

// print czyta tekst z pliku r, dzieli go na strony z nagłówkiem i
// dolnym marginesem i drukuje na stdout. W nagłówku jest umieszczana
// nazwa pliku przekazana w argumencie fname. Zwraca error jeśli
// wystąpił.
func print(r io.Reader, fname string) error {
	br := bufio.NewReader(r)

	const (
		margin1 = 1
		margin2 = 2
		bottom  = 64
		pagelen = 66
	)
	var (
		lineno int = 0
		pageno int = 0
	)

	pageno = 1
	skip(margin1)
	head(fname, pageno)
	skip(margin2)
	lineno = margin1 + margin2 + 1

	for {
		// wczytaj wiersz
		line, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF && len(line) == 0 {
			break
		}

		if lineno == 0 {
			// wydrukuj nagłówek
			pageno++
			skip(margin1)
			head(fname, pageno)
			skip(margin2)
			lineno = margin1 + margin2 + 1
		}

		// wydrukuj wiersz
		if line[len(line)-1] != '\n' {
			line += string('\n')
		}
		_, err = fmt.Print(line)
		if err != nil {
			return err
		}

		lineno++

		// dolny margines i nowa strona
		if lineno >= bottom {
			skip(pagelen - lineno)
			lineno = 0
		}
	}

	// wypełnij stronę
	if lineno > 0 {
		skip(pagelen - lineno)
	}
	return nil
}

func main() {
	if len(os.Args) == 1 {
		err := print(os.Stdin, "")
		if err != nil {
			fatal(err)
		}
	} else {
		for _, fname := range os.Args[1:] {
			file, err := os.Open(fname)
			if err != nil {
				fatal(err)
			}

			err = print(file, fname)
			if err != nil {
				fatal(err)
			}

			err = file.Close()
			if err != nil {
				fatal(err)
			}
		}
	}
}
