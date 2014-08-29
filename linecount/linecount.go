// 2014-06-08 Adam Bryt

// Narzędzia Programistyczne w Pascalu, rozdział 1.3 "Zliczanie
// wierszy", program linecount. Poza funkcjonalnością z książki dodałem
// obsługę plików.
//
// NAZWA
//
// linecount - zlicza wiersze
//
// SPOSÓB UŻYCIA
//
// linecount [-h] [file ...]
//
// OPIS
//
// Program linecount zlicza wiersze w podanych plikach i drukuje na
// stdout liczbę wierszy. Jeśli nie podano plików to zlicza wiersze na
// standardowym wejściu.
//
// Opcja -h drukuje krótki help.
//
// PRZYKŁADY
//
// Zliczanie wierszy w pliku:
//
//	$ linecount file
//
// Zliczanie wierszy na stdin:
//
//	$ linecount
//
// UWAGI
//
// Program zlicza znaki nowego wiersza '\n' - zakłada się, że każdy
// wiersz jest zakończony znakiem '\n'. Jeśli ostatni wiersz w pliku nie
// jest zakończony znakiem '\n' to nie jest liczony.
//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

const usageStr = "usage: linecount [-h] [file ...]"
const helpStr = usageStr + `

Program linecount zlicza wiersze w podanych plikach lub na stdin i
drukuje na stdout liczbę wierszy.

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
	fmt.Fprintf(os.Stderr, "linecount: %s\n", err)
	os.Exit(1)
}

// Funkcja linecount zwraca liczbę wierszy w obiekcie r oraz błąd jeśli
// wystąpił. Wiersze są zakończone znakiem '\n'.
func linecount(r io.Reader) (int64, error) {
	br := bufio.NewReader(r)
	var n int64
	for {
		b, err := br.ReadByte()
		if err == io.EOF {
			return n, nil
		}
		if err != nil {
			return n, err
		}
		if b == '\n' {
			n++
		}
	}
}

func main() {
	var helpFlag bool
	flag.BoolVar(&helpFlag, "h", false, "wyświetla help")
	flag.BoolVar(&helpFlag, "help", false, "wyświetla help")
	flag.Usage = usage
	flag.Parse()

	if helpFlag {
		help()
	}

	var num int64 // licznik wierszy

	if flag.NArg() == 0 {
		n, err := linecount(os.Stdin)
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

			n, err := linecount(file)
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
