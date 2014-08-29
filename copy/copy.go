// 2014-06-02 Adam Bryt

// Narzędzia Programistyczne w Pascalu, rozdział 1.1 "Kopiowanie
// plików", program copy.
//
// NAME
//
// copy - kopiuje wejście na wyjście
//
// SYNOPSIS
//
// usage: copy [-h] [ <file1 >file2 ]
//
// DESCRIPTION
//
// Program copy kopiuje standardowe wejscie na standardowe wyjście.
// Może być używany do kopiowania pliku na plik lub do wyświetlania
// pliku na terminal.
//
// Opcja -h wyświetla krótki help.
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
	"flag"
	"fmt"
	"io"
	"os"
)

const usageStr = "usage: copy [-h] [ <file1 >file2 ]"

const helpStr = usageStr + `
Program copy kopiuje stdin na stdout
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

func main() {
	var helpFlag bool
	flag.BoolVar(&helpFlag, "h", false, "wyświetl help")
	flag.BoolVar(&helpFlag, "help", false, "wyświetl help")
	flag.Usage = usage
	flag.Parse()

	if helpFlag {
		help()
	}

	_, err := io.Copy(os.Stdout, os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "copy: %v\n", err)
		os.Exit(1)
	}
}
