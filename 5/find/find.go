// 2014-12-01 Adam Bryt

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"npwp/5/pattern"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: find PATTERN")
	os.Exit(1)
}

func find(w io.Writer, r io.Reader, pat pattern.Pattern) error {
	br := bufio.NewReader(r)
	for done := false; !done; {
		lin, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			if len(lin) > 0 {
				// wiersz nie zakończony znakiem '\n'
				done = true
			} else {
				return nil
			}
		}
		if pattern.Match(lin, pat) {
			io.WriteString(w, lin)
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	pat, err := pattern.Getpat(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	err = find(os.Stdout, os.Stdin, pat)
	if err != nil {
		log.Fatal(err)
	}
}
