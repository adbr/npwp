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
// Kolejność wypisywania plików na standardowe wyjście (parametr -p) jest
// taka jak kolejność plików w archiwum, a nie jak kolejność nazw plików
// w poleceniu.
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
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"npwp/archive/header"
)

const (
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

// getfnames wstawia do fnames nazwy plików podane jako argumenty
// polecenia archive. Inicjuje fstats. Sprawdza czy nazwy plików się
// nie powtarzają - jeśli tak, to zwraca error.
func getfnames() error {
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
	t, err := ioutil.TempFile("", tempname)
	if err != nil {
		return err
	}
	defer func() {
		t.Close()
		os.Remove(t.Name())
	}()

	tw := bufio.NewWriter(t)
	defer tw.Flush()

	if cmd == "-u" {
		a, err := os.Open(aname)
		if err != nil {
			return err
		}
		ar := bufio.NewReader(a)
		err = replace(ar, tw, "-u")
		if err != nil {
			a.Close()
			return err
		}
		a.Close()
	}

	for i := 0; i < len(fstats); i++ {
		if fstats[i] == false {
			err := addfile(fnames[i], tw)
			if err != nil {
				return err
			}
			fstats[i] = true
		}
	}

	err = tw.Flush()
	if err != nil {
		return err
	}

	err = t.Close()
	if err != nil {
		return err
	}

	err = fcopy(aname, t.Name())
	if err != nil {
		// tymczasowe archiwum powinno pozostać, ponieważ archiwum aname
		// mogło zostać uszkodzone przez błąd podczas fcopy
		fmt.Fprintf(os.Stderr, "archive: %s\n", err)
		fmt.Fprintf(os.Stderr, "archive: archiwum tymczasowe: %s\n", t.Name())
		os.Exit(2) // nie wykonuje defer
	}

	return nil
}

// addfile dodaje plik fname na koniec archiwum w.
func addfile(fname string, w *bufio.Writer) error {
	defer w.Flush()

	nf, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer nf.Close()

	hdr, err := header.New(fname)
	if err != nil {
		return err
	}

	err = header.Write(w, hdr)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, nf)
	if err != nil {
		return err
	}

	return nil
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

// replace kopiuje archiwum ar do pliku tymczasowego tw zastępując lub
// usuwając pliki podane w wywołaniu archive.
func replace(ar *bufio.Reader, tw *bufio.Writer, cmd string) error {
	defer tw.Flush()

	for {
		hdr, err := header.Read(ar)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if filearg(hdr.Name) {
			if cmd == "-u" {
				err := addfile(hdr.Name, tw)
				if err != nil {
					return err
				}
			}
			err := fskip(ar, hdr.Size)
			if err != nil {
				return err
			}
		} else {
			err := header.Write(tw, hdr)
			if err != nil {
				return err
			}
			err = acopy(tw, ar, hdr.Size)
			if err != nil {
				return err
			}
		}
	}
}

// table drukuje wykaz zawartości archiwum aname.
func table(aname string) error {
	file, err := os.Open(aname)
	if err != nil {
		return err
	}
	defer file.Close()
	br := bufio.NewReader(file) // dla czytania wiersza funkcją ReadString

	for {
		hdr, err := header.Read(br)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if filearg(hdr.Name) {
			tprint(hdr)
		}

		err = fskip(br, hdr.Size)
		if err != nil {
			return err
		}
	}

	notfound()
	return nil
}

// filearg sprawdza czy plik name jest na liście parametrów. Jeśli
// lista parametrów jest pusta to zwraca zawsze true.
func filearg(name string) bool {
	if len(fnames) == 0 {
		return true
	}
	for i, f := range fnames {
		if name == f {
			fstats[i] = true
			return true
		}
	}
	return false
}

// tprint drukuje na stdout treść nagłówka hdr.
func tprint(hdr *header.Header) {
	fmt.Printf("%s %d\n", hdr.Name, hdr.Size)
}

// fskip czyta n bajtów z pliku r.
func fskip(r *bufio.Reader, n int64) error {
	for i := int64(0); i < n; i++ {
		_, err := r.ReadByte()
		if err != nil {
			return fmt.Errorf("błąd podczas fskip: %s", err)
		}
	}
	return nil
}

// notfound drukuje info o plikach występujących w liście parametrów
// ale nie znalezionych w archiwum.
func notfound() {
	for i, f := range fnames {
		if fstats[i] == false {
			fmt.Fprintf(os.Stderr, "%s: nie ma w archiwum\n", f)
		}
	}
}

// extract wydobywa pliki z archiwum.
func extract(aname, cmd string) error {
	f, err := os.Open(aname)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)

	for {
		hdr, err := header.Read(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if !filearg(hdr.Name) {
			err := fskip(r, hdr.Size)
			if err != nil {
				return err
			}
			continue
		}

		if cmd == "-p" {
			w := bufio.NewWriter(os.Stdout)
			err = acopy(w, r, hdr.Size)
			if err != nil {
				w.Flush()
				return err
			}
			w.Flush()
		} else { // "-x"
			ef, err := os.Create(hdr.Name)
			if err != nil {
				return err
			}
			w := bufio.NewWriter(ef)

			err = acopy(w, r, hdr.Size)
			if err != nil {
				w.Flush()
				ef.Close()
				return err
			}

			err = w.Flush()
			if err != nil {
				return err
			}
			err = ef.Close()
			if err != nil {
				return err
			}
		}
	}
	notfound()
	return nil
}

// acopy kopiuje n bajtów z src do dst.
func acopy(dst *bufio.Writer, src *bufio.Reader, n int64) error {
	for i := int64(0); i < n; i++ {
		c, err := src.ReadByte()
		if err != nil {
			return fmt.Errorf("błąd podczas acopy: %s", err)
		}

		err = dst.WriteByte(c)
		if err != nil {
			return fmt.Errorf("błąd podczas acopy: %s", err)
		}
	}
	return nil
}

// delete usuwa podane pliki z archiwum.
func delete(aname string) error {
	if len(fnames) == 0 {
		return errors.New("parametr -d wymaga podania nazw plików")
	}

	t, err := ioutil.TempFile("", tempname)
	if err != nil {
		return err
	}
	defer func() {
		t.Close()
		os.Remove(t.Name())
	}()

	tw := bufio.NewWriter(t)
	defer tw.Flush()

	a, err := os.Open(aname)
	if err != nil {
		return err
	}
	defer a.Close()

	ar := bufio.NewReader(a)

	err = replace(ar, tw, "-d")
	if err != nil {
		return err
	}

	notfound()

	err = a.Close()
	if err != nil {
		return err
	}

	err = tw.Flush()
	if err != nil {
		return err
	}
	err = t.Close()
	if err != nil {
		return err
	}

	err = fcopy(aname, t.Name())
	if err != nil {
		// tymczasowe archiwum powinno pozostać, ponieważ archiwum aname
		// mogło zostać uszkodzone przez błąd podczas fcopy
		fmt.Fprintf(os.Stderr, "archive: %s\n", err)
		fmt.Fprintf(os.Stderr, "archive: archiwum tymczasowe: %s\n", t.Name())
		os.Exit(2) // nie wykonuje defer
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		usage()
	}

	cmd := os.Args[1]
	aname := os.Args[2]
	err := getfnames()
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
		err := table(aname)
		if err != nil {
			fatal(err)
		}
	case "-x", "-p":
		err := extract(aname, cmd)
		if err != nil {
			fatal(err)
		}
	case "-d":
		err := delete(aname)
		if err != nil {
			fatal(err)
		}
	default:
		usage()
	}
}
