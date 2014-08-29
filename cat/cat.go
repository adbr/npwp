// 2014-06-03 Adam Bryt

// Program cat jest modyfikacją i rozszerzeniem programu copy (Narzędzia
// Programistyczne w Pascalu, rozdział 1.1 "Kopiowanie plików").
//
// NAME
//
// cat - konkatenuje pliki i drukuje na stdout
//
// SYNOPSIS
//
// cat [-h] [file ...]
//
// DESCRIPTION
//
// Program cat konkatenuje podane pliki i wysyła dane na stdout.  Jeśli
// nie podano plików to kopiuje stdin na stdout. Jeśli plik nie istnieje
// lub nie może zostać otwarty to jest zgłaszany błąd i program kończy
// działanie. Może być używany np do kopiowanie plików, do wyświetlania
// zawartości plików na terminal lub do wysyłania zawartości plików do
// innego polecenia.
//
// Opcja -h wyświetla krótki help.
//
// EXAMPLES
//
// Kopiowanie plików file1 i file2 do pliku file3:
//
// 	$ cat file1 file2 >file3
//
// Wyświetlanie na terminalu zawartości pliku file:
//
// 	$ cat file
//
// Kopiowanie danych z stdin do pliku file:
//
// 	$ cat >file
//
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const usageStr = "usage: cat [-h] [file ...]"
const helpStr = usageStr + `

Program cat konkatenuje podane pliki i wysyła dane na stdout. Jeśli nie
podano plików to kopiuje stdin na stdout.

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

func fatal(e error) {
	fmt.Fprintf(os.Stderr, "cat: %v\n", e)
	os.Exit(1)
}

func main() {
	var helpFlag bool
	flag.BoolVar(&helpFlag, "h", false, "wyświetla help")
	flag.BoolVar(&helpFlag, "help", false, "Wyświetla help")
	flag.Usage = usage
	flag.Parse()

	if helpFlag {
		help()
	}

	if flag.NArg() == 0 {
		_, err := io.Copy(os.Stdout, os.Stdin)
		if err != nil {
			fatal(err)
		}
	} else {
		for _, fname := range flag.Args() {
			file, err := os.Open(fname)
			if err != nil {
				fatal(err)
			}

			_, err = io.Copy(os.Stdout, file)
			if err != nil {
				file.Close()
				fatal(err)
			}

			err = file.Close()
			if err != nil {
				fatal(err)
			}
		}
	}
}
