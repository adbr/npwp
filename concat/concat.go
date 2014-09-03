// 2014-09-03 Adam Bryt

// Narzędzia Programistyczne w Pascalu,
// rozdział 3.4 "Konkatenacja plików"
// program concat.
//
// NAZWA
//
// concat - konkatenuje (łączy) pliki
//
// SPOSÓB UŻYCIA
//
// concat [file ...]
//
// OPIS
//
// Program concat czyta kolejno dane z plików podanych jako parametry i
// wysyła je na stdout. Jeśli nie podano parametrów to czyta z stdin.
//
// PRZYKŁADY
//
// Połączenie plików file i file2 i zapisanie danych do file3:
//
//	concat file file2 >file3
//
// Wyświetlenie na terminalu zawartości pliku file.txt
//
//	concat file.txt
//
// UWAGI
//
// Polecenie:
//
//	./concat file >file
//
// usuwa zawartość pliku file.
//
package main

import (
	"fmt"
	"io"
	"os"
)

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "concat: %s\n", err)
	os.Exit(1)
}

func main() {
	if len(os.Args) == 1 {
		// nie podano argumentów - kopiuj stdin
		_, err := io.Copy(os.Stdout, os.Stdin)
		if err != nil {
			fatal(err)
		}
		return
	}

	for _, fname := range os.Args[1:] {
		file, err := os.Open(fname)
		if err != nil {
			fatal(err)
		}
		_, err = io.Copy(os.Stdout, file)
		if err != nil {
			fatal(err)
		}
	}
}
