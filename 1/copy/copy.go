// 2014-06-02 Adam Bryt

// Narzędzia Programistyczne w Pascalu,
// rozdział 1.1 "Kopiowanie plików",
// program copy.
//
// NAME
//
// copy - kopiuje wejście na wyjście
//
// SYNOPSIS
//
// usage: copy [ <file1 >file2 ]
//
// DESCRIPTION
//
// Program copy kopiuje standardowe wejscie na standardowe wyjście.
// Może być używany do kopiowania pliku na plik lub do wyświetlania
// pliku na terminal.
//
// EXAMPLES
//
// Kopiowanie standardowego wejścia na standardowe wyjście:
//
// 	$ copy
//
// Kopiowanie pliku file1 do pliku file2:
//
// 	$ copy <file1 >file2
//
package main

import (
	"io"
	"log"
	"os"
)

func main() {
	log.SetPrefix("copy: ")
	log.SetFlags(0)

	_, err := io.Copy(os.Stdout, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
}
