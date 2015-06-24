// 2015-05-12 Adam Bryt

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"unicode/utf8"

	"npwp/5/pattern"
)

// Znak, który w tekście zastępującym oznacza dopasowany fragment tekstu.
const s_ditto = '&'

// Wartość znaku 'ditto' w skompilowanym tekście zastępującym.
// Wartość różna od kodów zwykłych znaków.
const ditto = 0

var usageStr = "sposób użycia: change <pattern> [<substitution>]"

var helpStr = `Program zamienia wzorce w tekście.
sposób użycia: change <pattern> [<substitution>]

Program change czyta wiersze tekstu ze standardowego wejścia,
zamienia wszystkie nie nakładające się fragmenty pasujące do
<pattern> na <substitution> i drukuje je na standardowe wyjście.
Reguły budowy wzorców <pattern> są takie same jak w programie
'find'.  Jeśli argument <substitution> nie istnieje to tekst
pasujący do <pattern> jest usuwany. Jeśli argument <substitution>
zawiera znaki '&' to te znaki są zastępowane fragmentem pasującym
do <pattern>.  Żeby znak '&' pozbawić specjalnego znaczenia,
należy zamiast niego użyć sekwencji '@&'.
`

func usage() {
	fmt.Fprintln(os.Stderr, usageStr)
	os.Exit(1)
}

func help() {
	fmt.Fprintln(os.Stdout, helpStr)
	os.Exit(0)
}

// getsub zwraca tekst zastępujący utworzony z argumentu s.
// Zastępuje znak 'ditto' przez kod będący wartością stałej ditto.
// Rozwija cytowania (escape'owania).
func getsub(s string) (string, error) {
	var new []byte
	for {
		if len(s) == 0 {
			break
		}
		r, n := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			err := errors.New("getsub: błąd dekodowania znaku utf8")
			return string(new), err
		}
		if r == s_ditto {
			new = append(new, ditto)
			s = s[n:]
		} else {
			r, s = pattern.Esc(s)
			var a [utf8.UTFMax]byte
			n := utf8.EncodeRune(a[:], r)
			new = append(new, a[:n]...)
		}
	}
	return string(new), nil
}

// subpart zwraca string zastępujący dopasowany fragment match.
// Argument sub może zawierać zaczniki 'ditto', które zostaną
// zastąpione przez dopasowany fragment match.
func subpart(match string, sub string) (string, error) {
	var new []byte
	for {
		if len(sub) == 0 {
			break
		}
		r, n := utf8.DecodeRuneInString(sub)
		if r == utf8.RuneError {
			err := errors.New("subpart: błąd dekodowania znaku utf8")
			return string(new), err
		}
		if r == ditto {
			new = append(new, match...)
		} else {
			new = append(new, sub[:n]...)
		}
		sub = sub[n:]
	}
	return string(new), nil
}

func subline(line string, pat pattern.Pattern, sub string) (string, error) {
	var new []byte
	var i int = 0 // indeks początku dopasowania w line
	for {
		if i >= len(line) {
			break
		}
		ok, n := pattern.Amatch(line, i, pat, 0)
		if ok {
			s, err := subpart(line[i:i+n], sub)
			if err != nil {
				return string(new), err
			}
			new = append(new, s...)
			i += n
		} else {
			r, rn := utf8.DecodeRuneInString(line[i:])
			if r == utf8.RuneError {
				err := errors.New("subline: błąd dekodowania znaku utf8")
				return string(new), err
			}
			new = append(new, line[i:i+rn]...)
			i += rn
		}
	}

	return string(new), nil
}

func change(w io.Writer, r io.Reader, pat pattern.Pattern, sub string) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		new, err := subline(line, pat, sub)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(w, new)
		if err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func main() {
	helpFlag := flag.Bool("h", false, "wyświetla help")
	flag.BoolVar(helpFlag, "help", false, "wyświetla help")
	flag.Usage = usage
	flag.Parse()

	if *helpFlag {
		help()
	}

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "zła liczba argumentów")
		usage()
	}
	pat, err := pattern.Makepat(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	sub := ""
	if flag.NArg() >= 2 {
		sub, err = getsub(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
	}

	err = change(os.Stdout, os.Stdin, pat, sub)
	if err != nil {
		log.Fatal(err)
	}
}
