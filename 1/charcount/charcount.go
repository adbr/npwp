// 2014-06-05 Adam Bryt

// Narzędzia Programistyczne w Pascalu, rozdział 1.2 "Zliczanie znaków",
// program charcount. Poza funkcjonalnością z książki dodałem obsługę
// plików i zliczanie znaków UTF-8.
//
// NAME
//
// charcount - zlicza znaki
//
// SYNOPSIS
//
// charcount [-hu] [file ...]
//
// DESCRIPTION
//
// Program charcount zlicza znaki w podanych plikach i drukuje na stdout
// liczbę znaków. Jeśli nie podano plików to zlicza znaki na
// standardowym wejściu. Domyślnie są zliczane bajty, opcja -u powoduje
// zliczanie znaków UTF-8.
//
// Opcje:
//	-h  wyświetla krótki help
//	-u  zlicza znaki UTF-8 (domyślnie są zliczane bajty)
//
// EXAMPLES
//
// Zliczanie znaków w pliku file:
//
//	$ charcount file
//
// Zliczanie znaków na stdin:
//
//	$ charcount
//
// Zliczanie znaków UTF-8 w pliku file:
//
//	$ charcount -u file
//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

const usageStr = "usage: charcount [-hu] [file ...]"
const helpStr = usageStr + `

Program charcount zlicza bajty czytane z stdin lub z podanych plików i
drukuje na stdout liczbę bajtów. Jeśli podano opcję -u to zakłada się,
że dane są tekstem Unicode UTF-8 i zamiast bajtów są zliczane znaki
(code points).

Opcje:
	-h  wyświetla ten help
	-u  zlicza znaki UTF-8 (code points)
`

func usage() {
	fmt.Fprintln(os.Stderr, usageStr)
	os.Exit(1)
}

func help() {
	fmt.Println(helpStr)
	os.Exit(0)
}

func warn(format string, a ...interface{}) {
	format = "charcount: warning: " + format
	fmt.Fprintf(os.Stderr, format, a...)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "charcount: %s\n", err)
	os.Exit(1)
}

// Funkcja charcount czyta z obiektu r wszystkie znaki (bajty lub runy -
// zależnie od parametru isUtf8) i zwraca liczbę odczytanych znaków oraz
// błąd jeśli wystąpił.
func charcount(r io.Reader, isUtf8 bool) (int64, error) {
	br := bufio.NewReader(r)
	var n int64 = 0

	if isUtf8 {
		for {
			r, _, err := br.ReadRune()
			if err == io.EOF {
				return n, nil
			}
			if err != nil {
				return n, err
			}
			if r == utf8.RuneError {
				warn("invalid rune nr %d\n", n)
			}
			n++
		}
	} else {
		for {
			_, err := br.ReadByte()
			if err == io.EOF {
				return n, nil
			}
			if err != nil {
				return n, err
			}
			n++
		}
	}
}

func main() {
	helpFlag := false
	utf8Flag := false

	flag.BoolVar(&helpFlag, "h", false, "wyświetla help")
	flag.BoolVar(&helpFlag, "help", false, "wyświetla help")
	flag.BoolVar(&utf8Flag, "u", false, "zlicza znaki UTF-8")
	flag.Usage = usage
	flag.Parse()

	if helpFlag {
		help()
	}

	var num int64 = 0 // licznik znaków

	if flag.NArg() == 0 {
		n, err := charcount(os.Stdin, utf8Flag)
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

			n, err := charcount(file, utf8Flag)
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
