// 2014-11-28 Adam Bryt

// Narzędzia Programistyczne w Pascalu,
// rozdział 4.8 "Indeks permutacyjny",
// program unrotate.
//
// NAZWA
//
// unrotate - formatuje wiersze do indeksu KWIC
//
// SPOSÓB UŻYCIA
//
// unrotate
//
// OPIS
//
// Program unrotate czyta z stdin wiersze tekstu wygenerowane przez
// program kwic, formatuje je tak, żeby powstał indeks KWIC i drukuje
// do stdout.
//
// PRZYKŁADY
//
// Mając dany plik y.txt utworzony poleceniem kwik | sort:
//
//	cat y.txt
//	This is a test.$
//	a test.$This is
//	is a test.$This
//	test.$This is a
//
// polecenie unrotate <y.txt wypisze następujący tekst:
//
//
//                                          This is a test.
//                                 This is  a test.
//                                    This  is a test.
//                               This is a  test.
//
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	maxout = 80
	middle = 40
	fold   = "$"
)

func unrotateLine(line string) string {
	line = strings.TrimSuffix(line, "\n")
	var out = make([]byte, maxout)

	for i := range out {
		out[i] = ' '
	}

	f := strings.Index(line, fold)
	j := middle - 1
	for i := len(line) - 1; i > f; i-- {
		out[j] = line[i]
		j--
		if j < 0 {
			j = maxout - 1
		}
	}

	j = middle + 1
	for i := 0; i < f; i++ {
		out[j] = line[i]
		j++
		if j >= maxout {
			j = 0
		}
	}

	out = bytes.TrimRight(out, " ")
	return string(out)
}

func unrotate(w io.Writer, r io.Reader) error {
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			if len(line) > 0 {
				line = unrotateLine(line)
				_, err := fmt.Fprintln(w, line)
				if err != nil {
					return err
				}
			}
			return nil
		}

		line = unrotateLine(line)
		_, err = fmt.Fprintln(w, line)
		if err != nil {
			return err
		}
	}
}

func main() {
	err := unrotate(os.Stdout, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
}
