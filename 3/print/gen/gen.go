// 2014-09-04 Adam Bryt
// Program generuje sekwencję NUM wierszy STRING; wiersze są numerowane.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gen NUM STRING\n")
	os.Exit(1)
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		usage()
	}
	str := os.Args[2]

	// bufor dla zwiększenia szybkości
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	for i := 1; i <= n; i++ {
		fmt.Fprintf(w, "%6d %s\n", i, str)
	}
}
