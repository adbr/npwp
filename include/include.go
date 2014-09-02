// 2014-09-02 Adam Bryt

// Narzędzia Programistyczne w Pascalu,
// rozdział 3.3 "Dołączanie plików"
// program include.
//
// NAZWA
//
// include - dołącza pliki wewnątrz plików
//
// SPOSÓB UŻYCIA
//
// include [<file1 >file2]
//
// OPIS
//
// Program include kopiuje stdin na stdout zastępując wiersze
// zaczynające się dyrektywą postaci '#include "file"' zawartością
// pliku "file". Dołączane pliki są przetwarzane rekurencyjnie, czyli
// jeśli zawierają dyrektywy #include, to one też są przetwarzane.
//
// PRZYKŁADY
//
// Jeśli plik a.txt zawiera:
//
//	aaaaaa
//	#include "b.txt"
//	aaaaaa
//
// a plik b.txt zawiera:
//
//	bbbb
//	bbbb
//
// to polecenie:
//
// 	include <a.txt
//
// wypisze na standardowe wyjście:
//
//	aaaaaa
//	bbbb
//	bbbb
//	aaaaaa
//
// UWAGI
//
// Program nie wykrywa sytuacji gdy plik dołącza samego siebie, nastąpi
// wtedy przepełnienie jakichś zasobów systemowych procesu.
//
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "include: %s\n", err)
	os.Exit(1)
}

// Funkcja include kopiuje dane z r do w zastępując wiersze postaci:
// #include "file", zawartością pliku "file". Jeśli dołączany plik
// zawiera dyrektywę #include to funkcja include jest wywoływana
// rekurencyjnie. Zwraca error jeśli wystąpił.
func include(w io.Writer, r io.Reader) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()
	br := bufio.NewReader(r)

	const incl = "#include"

	for {
		line, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF && len(line) == 0 {
			break
		}

		words := strings.Fields(line)
		if len(words) >= 2 && words[0] == incl {
			// dyrektywa #include
			fname := words[1]
			fname = strings.Trim(fname, `"`)

			file, err := os.Open(fname)
			if err != nil {
				return err
			}

			err = include(bw, file)
			if err != nil {
				return err
			}

			err = file.Close()
			if err != nil {
				return err
			}
		} else {
			// zwykły wiersz - skopiuj
			_, err := bw.WriteString(line)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	err := include(os.Stdout, os.Stdin)
	if err != nil {
		fatal(err)
	}
}
