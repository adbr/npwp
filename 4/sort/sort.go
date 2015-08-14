// 2014-10-25 Adam Bryt
//
// Uwagi:
// - Użycie collate do porównywania napisów bardzo spowalnia sortowanie;
//   porównanie stringów funkcją collate.CompareStrings() jest około 250 razy
//   wolniejsze niż porównanie stringów wbudowanym operatorem
//   (npwp/4/bench/cmpstring_test.go).
// - Użycie collate powoduje, że kompilacja trwa długo i wymaga dużo
//   pamięci RAM - ponad 500 MB. Poprawka: pakiet collate nie był
//   skompilowany i zainstalowany - po instalacji collage (go install)
//   kompilacja jest szybka.
// - Algorytm quicksort jest nie optymalny - np bardzo długo sortuje
//   plik /usr/share/dict/words - może z tego powodu, że plik jest już
//   posortowany, a to jest najgorszy przypadek dla tego algorytmu.
// - Wybieranie najmniejszego wiersza podczas łączenia plików tymczasowych
//   jest robione przez przeszukiwanie liniowe, zamiast przy użyciu stogu
//   jak w książce.

// Narzędzia Programistyczne w Pascalu,
// rozdział 4.5 "Sortowanie dużych plików",
// program sort.
//
// NAZWA
//
// sort - sortuje wiersze tekstu
//
// SPOSÓB UŻYCIA
//
// sort [<file1] [>file2]
//
// OPIS
//
// Program sort sortuje wiersze tekstu czytane ze standardowego wejścia
// i drukuje posortowany tekst na standardowe wyjście. Wiersze tekstu są
// sortowane zewnętrznie - fragmenty tekstu są sortowane metodą quick sort,
// zapisywane w plikach tymczasowych, które następnie są łączone. Pliki
// tymczasowe mają nazwy stemp# gdzie # jest liczbą całkowitą.
//
// PRZYKŁADY
//
// Sortowanie pliku file1 i zapisanie wyniku do pliku file2:
//
//	$sort <file1 >file2
//
// UWAGI
//
// Każdy wiersz w pliku powinien być zakończony znakiem \n - jeśli ostatni
// wiersz w pliku nie jest zakończony znakiem \n to mogą wystąpić błędy
// sortowania spowodowane tym, że podczas łączenia plików tymczasowych
// wiersz bez \n i wiersz następny zostaną połączone w jeden.
//
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

const (
	tname      = "stemp" // prefix nazwy plików tymczasowych
	maxlines   = 100     // liczba wierszy sortowanych na raz
	mergeorder = 5       // liczba łaczonych plików
)

// gtext czyta i zwraca maksymalnie maxlines wierszy z r. Gdy wystąpił
// koniec pliku zwraca wczytane wiersze i io.EOF.
func gtext(r *bufio.Reader, maxlines int) ([]string, error) {
	var lines []string
	for {
		if len(lines) >= maxlines {
			return lines, nil
		}
		line, err := r.ReadString('\n')
		if err == io.EOF {
			if len(line) > 0 {
				lines = append(lines, line)
			}
			return lines, err
		}
		if err != nil {
			return lines, err
		}
		lines = append(lines, line)
	}
}

// ptext drukuje wiersze z lines do pliku w.
func ptext(w io.Writer, lines []string) error {
	for _, line := range lines {
		_, err := fmt.Fprint(w, line)
		if err != nil {
			return err
		}
	}
	return nil
}

// gname zwraca nazwę pliku tymczasowego o numerze n.
func gname(n int) string {
	return fmt.Sprintf("%s%d", tname, n)
}

// quicksort sortuje a metodą quicksort zgodnie z kolejnością (collate)
// określoną przez c.
func quicksort(a []string, c *collate.Collator) {
	if len(a) < 2 {
		return
	}
	xi := len(a) - 1 // indeks pivota
	x := a[xi]       // wartość pivota
	i, j := 0, len(a)-1
	// zamieniaj miejscami elementy większe i mniejsze od pivota
	for {
		if i == j {
			break
		}
		for {
			if i == j || c.CompareString(a[i], x) > 0 {
				break
			}
			i++
		}
		for {
			if j == i || c.CompareString(a[j], x) < 0 {
				break
			}
			j--
		}
		if i < j {
			a[i], a[j] = a[j], a[i]
		}
	}
	// wstaw pivot w miejsce spotkania i j
	a[i], a[xi] = a[xi], a[i]
	quicksort(a[:i], c)
	quicksort(a[i:], c)
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

// gopen otwiera do czytania i zwraca pliki tymczasowe o numerach od a do b.
func gopen(a int, b int) ([]*os.File, error) {
	n := b - a + 1
	files := make([]*os.File, n)
	for i := 0; i < n; i++ {
		f, err := os.Open(gname(a + i))
		if err != nil {
			return nil, err
		}
		files[i] = f
	}
	return files, nil
}

// gremove zamyka i usuwa pliki tymczasowe files.
func gremove(files []*os.File) error {
	for _, f := range files {
		err := f.Close()
		if err != nil {
			return err
		}
		err = os.Remove(f.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

// minLine zwraca indeks najmniejszego wiersza w lines lub -1 gdy lines
// jest pusta lub zawiera same stringi puste. Stringi puste ("" - zero value)
// są pomijane.
func minLine(lines []string, c *collate.Collator) int {
	min := -1
	for i := range lines {
		if lines[i] == "" {
			continue
		}
		if min == -1 {
			// pierwszy string != nil
			min = i
			continue
		}
		if c.CompareString(lines[i], lines[min]) < 0 {
			min = i
		}
	}
	return min
}

// merge łączy posortowane pliki files i zapisuje posortowane wiersze
// do pliku out.
func merge(out *os.File, files []*os.File, c *collate.Collator) error {
	nf := len(files)

	// utworzenie obiektów bufio.Reader dla plików - dla użycia ReadString
	bfiles := make([]*bufio.Reader, nf)
	for i, f := range files {
		bfiles[i] = bufio.NewReader(f)
	}

	// pobranie pierwszego wiersza z każdego pliku
	lines := make([]string, nf)
	for i, f := range bfiles {
		line, err := f.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if len(line) == 0 && err == io.EOF {
			return fmt.Errorf("merge: pusty plik tymczasowy: %s",
				files[i].Name())
		}
		lines[i] = line
	}

	// łączenie
	for nf > 0 {
		i := minLine(lines, c)
		if i < 0 {
			panic("merge: ujemny indeks najmniejszego wiersza")
		}

		_, err := out.WriteString(lines[i])
		if err != nil {
			return err
		}

		line, err := bfiles[i].ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if len(line) > 0 {
			lines[i] = line
		} else {
			lines[i] = ""
			nf--
		}
	}
	return nil
}

func main() {
	// czytanie porcji wierszy, sortowanie i zapisywanie do plików tmp
	c := collate.New(language.Polish)
	r := bufio.NewReader(os.Stdin)
	high := 0
	done := false
	for !done {
		lines, err := gtext(r, maxlines)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if err == io.EOF {
			if len(lines) == 0 {
				break
			} else {
				done = true
			}
		}

		quicksort(lines, c)

		high++
		file, err := os.Create(gname(high))
		if err != nil {
			log.Fatal(err)
		}
		err = ptext(file, lines)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}

	// łaczenie serii plików, aż zostanie tylko jeden plik
	for low := 1; low <= high; low += mergeorder {
		lim := min(low+mergeorder-1, high)
		files, err := gopen(low, lim)
		if err != nil {
			log.Fatal(err)
		}

		high++
		out, err := os.Create(gname(high))
		if err != nil {
			log.Fatal(err)
		}

		err = merge(out, files, c)
		if err != nil {
			log.Fatal(err)
		}

		err = out.Close()
		if err != nil {
			log.Fatal(err)
		}
		err = gremove(files)
		if err != nil {
			log.Fatal(err)
		}
	}

	// drukuj plik końcowy na stdout i usuń plik
	f, err := os.Open(gname(high))
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(os.Stdout, f)
	if err != nil {
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove(f.Name())
	if err != nil {
		log.Fatal(err)
	}
}
