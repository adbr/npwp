// 2014-09-06 Adam Bryt

// Narzędzia Programistyczne w Pascalu,
// rozdział 3.7 "Dynamiczne tworzenie plików",
// program makecopy.
//
// NAZWA
//
// makecopy - kopiuje plik
//
// SPOSÓB UŻYCIA
//
// makecopy filesrc filedst
//
// OPIS
//
// Program makecopy kopiuje plik filesrc do pliku filedst. Jeśli plik
// docelowy nie istnieje to zostanie utworzony, a jeśli istnieje to
// jego zawartość zostanie obcięta do zera przed kopiowaniem.
//
// PRZYKŁADY
//
// Skopiowanie pliku file1 do pliku file2:
//
//	makecopy file1 file2
//
// UWAGI
//
// Skopiowanie pliku do samego siebie powoduje usunięcie zawartości pliku.
//
package main

import (
	"fmt"
	"io"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: makecopy filesrc filedst\n")
	os.Exit(1)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "makecopy: %s\n", err)
	os.Exit(2)
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	nsrc := os.Args[1]
	ndst := os.Args[2]

	fsrc, err := os.Open(nsrc)
	if err != nil {
		fatal(err)
	}
	defer fsrc.Close()

	fdst, err := os.Create(ndst)
	if err != nil {
		fatal(err)
	}
	defer fdst.Close()

	_, err = io.Copy(fdst, fsrc)
	if err != nil {
		fatal(err)
	}
}
