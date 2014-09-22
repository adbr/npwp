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
// Jeśli dla opcji -c (create) nie poda się nazw plików, to tworzone jest
// archiwum puste - poprzednia zawartość pliku jest usuwana.
//
// Archiwum jest sekwencją plików, z których każdy jest poprzedzony
// metryczką (header) postaci stringu zakończonego znakiem \n:
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
	"io"
	"io/ioutil"
	"os"
)

const (
	archhdr  = "-h-"     // początek nagłówka pliku w archiwum
	tempname = "archive" // prefix nazwy tymczasowego pliku archiwum
)

var (
	fnames []string // nazwy plików będących argumentami polecenia
	fstats []bool   // czy i-ty plik z fnames jest już w archiwum
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: archive -cdptux archname [file ...]\n")
	os.Exit(1)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "archive: %s\n", err)
	os.Exit(2)
}

// getFnames wstawia do fnames nazwy plików podane jako argumenty
// polecenia archive. Inicjuje fstats. Sprawdza czy nazwy plików się
// nie powtarzają - jeśli tak, to zwraca error.
func getFnames() error {
	fnames = os.Args[3:]
	fstats = make([]bool, len(fnames))
	for i := range fstats {
		fstats[i] = false
	}

	// sprawdzenie czy nazwy plików są unikalne
	for i := 0; i < len(fnames)-1; i++ {
		for j := i + 1; j < len(fnames); j++ {
			if fnames[i] == fnames[j] {
				return fmt.Errorf("plik się powtarza: %q", fnames[i])
			}
		}
	}

	return nil
}

// update uaktualnia lub dodaje pliki do archiwum. Działa na pliku
// tymczasowym, który na końcu jest kopiowany do archiwum aname.
// Jeśli podczas kopiowania archiwum tymczasowego wystąpi błąd, to
// plik tymczasowy nie jest usuwany i następuje wyjście z programu.
func update(aname, cmd string) error {
	tfile, err := ioutil.TempFile("", tempname)
	if err != nil {
		return err
	}
	tname := tfile.Name()
	defer func() {
		tfile.Close()
		os.Remove(tname)
	}()

	if cmd == "-u" {
		afile, err := os.Open(aname)
		if err != nil {
			return err
		}
		err = replace(afile, tfile, "-u")
		if err != nil {
			afile.Close()
			return err
		}
		afile.Close()
	}

	for i := 0; i < len(fstats); i++ {
		if fstats[i] == false {
			err := addfile(fnames[i], tfile)
			if err != nil {
				return err
			}
			fstats[i] = true
		}
	}

	err = tfile.Close()
	if err != nil {
		return err
	}

	err = fcopy(aname, tname)
	if err != nil {
		// tymczasowe archiwum powinno pozostać, ponieważ archiwum aname
		// mogło zostać uszkodzone przez błąd podczas fcopy
		fmt.Fprintf(os.Stderr, "archive: %s\n", err)
		fmt.Fprintf(os.Stderr, "archive: archiwum tymczasowe: %s\n", tname)
		os.Exit(2) // nie wykonuje defer
	}

	return nil
}

// addfile dodaje plik fname na koniec archiwum file.
func addfile(fname string, file *os.File) error {
	nf, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer nf.Close()

	hdr, err := makeHeader(fname)
	if err != nil {
		return err
	}

	_, err = file.WriteString(hdr)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, nf)
	if err != nil {
		return err
	}

	return nil
}

// makeHeader tworzy i zwraca nagłówek pliku fname.
func makeHeader(fname string) (string, error) {
	size, err := fsize(fname)
	if err != nil {
		return "", err
	}
	hdr := fmt.Sprintf("%s %s %d\n", archhdr, fname, size)
	return hdr, nil
}

// fsize zwraca rozmiar pliku w bajtach.
func fsize(fname string) (int64, error) {
	fi, err := os.Stat(fname)
	if err != nil {
		return 0, err
	}
	n := fi.Size()
	return n, nil
}

// fcopy kopiując zawartość pliku src do pliku dst.
func fcopy(dst, src string) error {
	fsrc, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fsrc.Close()

	fdst, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fdst.Close()

	_, err = io.Copy(fdst, fsrc)
	if err != nil {
		return err
	}

	return nil
}

func replace(af, tf *os.File, cmd string) error {
	return nil
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
	err := getFnames()
	if err != nil {
		fatal(err)
	}

	switch cmd {
	case "-c", "-u":
		err := update(aname, cmd)
		if err != nil {
			fatal(err)
		}
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
