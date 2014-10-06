// 2014-06-9 Adam Bryt

// Narzędzia Programistyczne w Pascalu, rozdział 1.4 "Zliczanie słów",
// program wordcount. Dodałem obsługę plików.
//
// NAZWA
//
// wordcount - zlicza słowa
//
// SPOSÓB UŻYCIA
//
// wordcount [-h] [file ...]
//
// OPIS
//
// Program wordcount zlicza słowa w podanych plikach i drukuje liczbę
// słów na stdout.
//
// PRZYKŁADY
//
// Zliczanie słów w pliku:
//
//	$ wordcount file
//
// Zliczanie słów na stdin:
//
//	$ wordcount
//
// UWAGI
//
// Granice słów stanowią znaki unicode.IsSpace().
//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode"
)

const usageStr = "usage: wordcount [-h] [file ...]"
const helpStr = usageStr + `

Program wordcount zlicza słowa w podanych plikach lub na stdin i drukuje
na stdout liczbę słów.

Opcje:
	-h  wyświetla ten help
`

func usage() {
	fmt.Fprintln(os.Stderr, usageStr)
	os.Exit(1)
}

func help() {
	fmt.Println(helpStr)
	os.Exit(0)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "wordcount: %s\n", err)
	os.Exit(1)
}

// Funkcja wordcount zwraca liczbę słów w r lub error jeśli wystąpił.
func wordcount(r io.Reader) (int64, error) {
	br := bufio.NewReader(r)
	n := int64(0)
	inword := false
	for {
		r, _, err := br.ReadRune()
		if err == io.EOF {
			return n, nil
		}
		if err != nil {
			return n, err
		}

		if unicode.IsSpace(r) {
			inword = false
		} else if !inword {
			inword = true
			n++
		}
	}
}

func main() {
	helpFlag := false
	flag.BoolVar(&helpFlag, "h", false, "wyświetla help")
	flag.BoolVar(&helpFlag, "help", false, "wyświetla help")
	flag.Usage = usage
	flag.Parse()

	if helpFlag {
		help()
	}

	num := int64(0) // licznik słów

	if flag.NArg() == 0 {
		n, err := wordcount(os.Stdin)
		if err != nil {
			fatal(err)
		}
		num = n
	} else {
		for _, fname := range flag.Args() {
			file, err := os.Open(fname)
			if err != nil {
				fatal(err)
			}

			n, err := wordcount(file)
			if err != nil {
				file.Close()
				fatal(err)
			}

			err = file.Close()
			if err != nil {
				fatal(err)
			}

			num += n
		}
	}
	fmt.Println(num)
}
