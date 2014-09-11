// 2014-09-08 Adam Bryt

// Narzędzia Programistyczne w Pascalu,
// rozdział 3.8 "Archiwum",
// program archive.
//
// NAZWA
//
// archive - obsługuje archiwum plików
//
// SPOSÓB UŻYCIA
//
// archive -cmd aname [file ...]
//
// OPIS
//
// Program archive służy do tworzenia i zarządzania archiwum plików
// aname, czyli jednym plikiem zawierającym wiele plików składowych.
// Umożliwia dodawanie, usuwanie, zastępowanie, wydobywanie elementów
// archiwum. Parametr -cmd określa jaką operację należy wykonać i
// może mieć jedną z wartości:
//
//	-c  utworzenie nowego archiwum zawierającego podane pliki
//	-d  usunięcie podanych plików z archiwum
//	-p  wypisanie podanych plików na standardowe wyjście
//	-t  wypisanie wykazu plików zawartych w archiwum
//	-u  uaktualnienie lub dodanie podanych plików
//	-x  wydobycie podanych plików z archiwum
//
// Jeśli nie poda się nazw plików to operacje dotyczą wszystkich plików
// w archiwum (nie dotyczy -d)
//
// Archiwum jest sekwencją plików, z których każdy jest poprzedzony
// metryczką (header) postaci:
//
//	-h- nazwa długość
//
// gdzie nazwa jest nazwą pliku, a długość jest długością pliku.
//
// PRZYKŁADY
//
// Zastąpienie dwóch plików i dodanie jednego nowego:
//
//	archive -u archfile old1 old2 new1
//
// Wypisanie skorowidza archiwum:
//
//	archive -t archfile
//
package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: archive -cdptux archname [file ...]\n")
	os.Exit(1)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "archive: %s\n", err)
	os.Exit(2)
}

// getFnames zwraca nazwy plików podane jako argumenty polecenia archive.
// Sprawdza czy nazwy plików się nie powtarzają.
func getFnames() ([]string, error) {
	fns := os.Args[3:]
	for i := 0; i < len(fns) - 1; i++ {
		for j := i + 1; j < len(fns); j++ {
			if fns[i] == fns[j] {
				return fns, fmt.Errorf("nazwa pliku się powtarza: %q", fns[i])
			}
		}
	}
	return fns, nil
}

func update(aname, cmd string) {
}

func table(aname string) {
}

func extract(aname, cmd string) {
}

func delete(aname string) {
}

func main() {
	if len(os.Args) < 3 {
		usage()
	}

	cmd := os.Args[1]
	aname := os.Args[2]
	fnames, err := getFnames()
	if err != nil {
		fatal(err)
	}

	_ = fnames

	switch cmd {
	case "-c", "-u":
		update(aname, cmd)
	case "-t":
		table(aname)
	case "-x", "-p":
		extract(aname, cmd)
	case "-d":
		delete(aname)
	default:
		usage()
	}
}
