// 2014-08-04 Adam Bryt

// Narzędzia Programistyczne w Pascalu, rozdział 2.2 "Nakładanie
// wierszy", program overstrike.
//
// NAZWA
//
// overstrike - zastępuje znaki cofania przez nakładane wiersze
//
// SPOSÓB UŻYCIA
//
// overstrike [<file >file2]
//
// OPIS
//
// Program overstrike czyta tekst z stdin, zamienia znaki cofania ('\b')
// na sekwencję nakładanych wierszy i drukuje tekst wynikowy na stdout.
// Zakłada się, że urządzenie drukujące interpretuje pierwszy znak w
// wierszu jako znak sterujący: spacja powoduje przejście do nowego
// wiersza przed wydrukiem, a znak '+' powoduje wydruk wiersza bez
// przechodzenia do nowego wiersza, czyli nadrukowanie wiersza na
// wierszu poprzednim.
//
// PRZYKŁADY
//
// Jeśli plik file zawiera tekst:
//
// 	abcd\b123
//
// to polecenie:
//
// 	overstrike file
//
// powoduje wygenerowanie na stdout tekstu:
//
// 	 abcd
// 	+   123
//
// czyli znak '1' zostanie nadrukowany na 'd'.
//
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "overstrike: %s\n", err)
	os.Exit(1)
}

func overstrike(w io.Writer, r io.Reader) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()
	br := bufio.NewReader(r)

	const (
		newline   = '\n'
		backspace = '\b'
		space     = ' '
	)

	// znaki sterujące drukarką - pierwszy znak w wierszu
	const (
		skip   = ' ' // przejdź do nowego wiersza przed wydrukiem
		noskip = '+' // nie przechodź do nowego wiersza - nadrukuj
	)

	var (
		c   rune  // wczytany znak
		err error //
		col int   // numer kolumny
		nbs int   // licznik znaków \b
	)

	for {
		// zero lub więcej znaków backspace
		for {
			c, _, err = br.ReadRune()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			if c == backspace {
				nbs++
			} else {
				break
			}
		}

		if nbs > 0 { // wystąpiła nakładka
			_, err := bw.WriteRune(newline)
			if err != nil {
				return err
			}
			_, err = bw.WriteRune(noskip)
			if err != nil {
				return err
			}
			for i := 0; i < col-nbs; i++ {
				_, err := bw.WriteRune(space)
				if err != nil {
					return err
				}
			}
			col = col - nbs
			if col < 0 {
				col = 0
			}
			nbs = 0
		} else if col == 0 {
			_, err := bw.WriteRune(skip)
			if err != nil {
				return err
			}
		}

		// normalny znak
		_, err := bw.WriteRune(c)
		if err != nil {
			return err
		}

		if c == newline {
			col = 0
		} else {
			col++
		}
	}
}

func main() {
	err := overstrike(os.Stdout, os.Stdin)
	if err != nil {
		fatal(err)
	}
}
