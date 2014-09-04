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
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: print [file ...]")
	os.Exit(1)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "print: %s\n", err)
	os.Exit(2)
}

func skip(w io.Writer, n int) error {
	for i := 0; i < n; i++ {
		_, err := fmt.Fprintf(w, "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func head(w io.Writer, name string, page int) error {
	_, err := fmt.Fprintf(w, "%s page %d\n", name, page)
	if err != nil {
		return err
	}
	return nil
}

func print(w io.Writer, r io.Reader, fname string) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()
	br := bufio.NewReader(r)

	const (
		margin1 = 1
		margin2 = 2
		bottom = 64
		pagelen = 66
	)
	var (
		lineno int
		pageno int
	)

	for {
		// nagłówek
		if lineno == 0 {
			pageno++
			err := skip(bw, margin1)
			if err != nil {
				return err
			}
			err = head(bw, fname, pageno)
			if err != nil {
				return err
			}
			err = skip(bw, margin2)
			if err != nil {
				return err
			}
			lineno = margin1 + margin2 + 1
		}

		// wczytaj wiersz
		line, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF && len(line) == 0 {
			break
		}

		// wydrukuj wiersz
		if line[len(line) - 1] != '\n' {
			line += string('\n')
		}
		_, err = bw.WriteString(line)
		if err != nil {
			return err
		}

		lineno++

		// dolny margines
		if lineno >= bottom {
			err := skip(bw, pagelen - lineno)
			if err != nil {
				return err
			}
			lineno = 0
		}
	}

	if lineno > 0 {
		err := skip(bw, pagelen - lineno)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if len(os.Args) == 1 {
		err := print(os.Stdout, os.Stdin, "")
		if err != nil {
			fatal(err)
		}
	} else {
		for _, fname := range os.Args[1:] {
			file, err := os.Open(fname)
			if err != nil {
				fatal(err)
			}

			err = print(os.Stdout, file, fname)
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
