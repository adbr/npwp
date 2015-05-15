// 2015-05-12 Adam Bryt

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func change(line, pat, sub string) (error, string) {
	return nil, "abc" + line
}

func main() {
	pat := "par"
	sub := "sub"
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		err, new := change(line, pat, sub)
		if err != nil {
			log.Fatal(err)
		}
		_, err = fmt.Println(new)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
