// 2014-08-29 Adam Bryt

// Narzędzia Programistyczne w Pascalu,
// rozdział 3.1 "Porównywanie plików",
// program compare.
//
// NAZWA
//
// compare - porównuje dwa pliki
//
// SPOSÓB UŻYCIA
//
// compare file1 file2
//
// OPIS
//
// Program compare porównuje wiersz po wierszu pliki o nazwach podanych
// jako parametry. Jeśli para wierszy jest różna to drukuje na
// standardowym wyjściu numer wiersza i oba różniące się wiersze;
// wiersze są numerowane od 0. Jeśli pliki są identyczne to nic nie
// wypisuje. W przypadku gdy któryś z plików jest krótszy, drukuje
// komunikat o osiągnięciu końca krótszego pliku.
//
// PRZYKŁADY
//
// Przykładowy komunikat w przypadku różnic w plikach:
//
//	compare a.txt b.txt
//
//	1:
//	abc def ghi
//	def ghi
//	4:
//	abc def ghi
//	xxx def ghi
//
// pliki różnią się w wierszach 1 i 4.
//
// UWAGI
//
// Niewielka różnica powodująca rozsynchronizowanie plików powoduje
// wygenerowanie dużych raportów.
//
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: compare file1 file2")
	os.Exit(1)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "compare: %s\n", err)
	os.Exit(2)
}

func diffmsg(n int, line1, line2 string) {
	// Dodaj na końcu wiersza znak '\n' jeśli nie ma - przypadek
	// gdy ostatni wiersz w pliku nie jest zakończony znakiem '\n'
	if len(line1) > 0 && line1[len(line1)-1] != '\n' {
		line1 += string('\n')
	}
	if len(line2) > 0 && line2[len(line2)-1] != '\n' {
		line2 += string('\n')
	}

	fmt.Printf("%d:\n", n)
	fmt.Print(line1)
	fmt.Print(line2)
}

func message(m string) {
	fmt.Fprintf(os.Stderr, "compare: %s\n", m)
}

// Funkcja compare porównuje pliki f1 i f2.
func compare(f1, f2 *os.File) error {
	br1 := bufio.NewReader(f1)
	br2 := bufio.NewReader(f2)

	var (
		lnum  int = 0
		line1 string
		line2 string
		err1  error
		err2  error
	)

	for {
		line1, err1 = br1.ReadString('\n')
		if err1 != nil && err1 != io.EOF {
			return err1
		}
		line2, err2 = br2.ReadString('\n')
		if err2 != nil && err2 != io.EOF {
			return err2
		}

		// Porównuj wiersze tylko gdy nie było EOF lub
		// gdy był EOF ale wczytano fragment wiersza
		if (err1 != io.EOF || (err1 == io.EOF && len(line1) > 0)) &&
			(err2 != io.EOF || (err2 == io.EOF && len(line2) > 0)) {
			if line1 != line2 {
				diffmsg(lnum, line1, line2)
			}
		}

		if err1 == io.EOF || err2 == io.EOF {
			break
		}
		lnum++
	}

	if err1 == io.EOF && err2 != io.EOF {
		message("EOF na pliku 1: " + f1.Name())
	} else if err2 == io.EOF && err1 != io.EOF {
		message("EOF na pliku 2: " + f2.Name())
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	fname1 := os.Args[1]
	fname2 := os.Args[2]

	file1, err := os.Open(fname1)
	if err != nil {
		fatal(err)
	}
	defer file1.Close()

	file2, err := os.Open(fname2)
	if err != nil {
		fatal(err)
	}
	defer file2.Close()

	err = compare(file1, file2)
	if err != nil {
		fatal(err)
	}
}
