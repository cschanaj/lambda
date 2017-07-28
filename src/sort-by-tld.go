package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"

	"libse/tldsort"
)

func main() {
	// input source to read from
	file := os.Stdin

	if len(os.Args) > 1 {
		lsource, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		} else {
			file = lsource
		}
	}


	// domain list to be sorted
	domains := make([]string, 0)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if tmp := scanner.Text(); len(tmp) > 0 {
			domains = append(domains, tmp)
		}
	}

	sort.Sort(Tldsort.Order(domains))

	for _, domain := range domains {
		fmt.Println(domain)
	}
}
