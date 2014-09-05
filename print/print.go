// 2014-09-04 Adam Bryt
// todo
//	- dodać wyświetlanie daty w nagłówku

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
// PRZYKŁADY
//
// UWAGI
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

func skip(n int) {
	for i := 0; i < n; i++ {
		fmt.Printf("\n")
	}
}

func head(name string, page int) {
	const format = "2006-01-02 15:04"
	t := time.Now().Format(format)

	fmt.Printf("%s  %s  page %d\n", name, t, page)
}

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

		// wydrukuj nagłówek
		if lineno == 0 {
			pageno++
			skip(margin1)
			head(fname, pageno)
			skip(margin2)
			lineno = margin1 + margin2 + 1
		}

		// wydrukuj wiersz
		if line[len(line)-1] != '\n' {
			line += string('\n')
			fmt.Fprintf(os.Stderr, "dodanie newline\n") // debug
		}
		_, err = fmt.Print(line)
		if err != nil {
			return err
		}

		lineno++

		// dolny margines
		if lineno >= bottom {
			fmt.Fprintf(os.Stderr, "skip lineno: %d\n", lineno) //debug
			skip(pagelen-lineno)
			lineno = 0
		}
	}

	fmt.Fprintf(os.Stderr, "last lineno: %d\n", lineno) // debug

	// wypełnij stronę
	if lineno > 0 {
		skip(pagelen-lineno)
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
