// 2015-06-24 Adam Bryt

package main

import (
	"flag"
	"fmt"
	"os"
)

var usageStr = "sposób użycia: edit [plik]"

func usage() {
	fmt.Fprintln(os.Stderr, usageStr)
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()
}
