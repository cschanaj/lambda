package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"libse/rules"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: fdt <path/https-everywhere/rules\n")
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(os.Args[1	])
	if err != nil {
		log.Fatal(err)
	}

	// map of all 'target'
	tmap := make(map[string]int)

	for _, file := range files {
		// full path of file
		filename := filepath.Join(os.Args[1], file.Name())

		// extension must be ".xml"
		if strings.HasSuffix(filename, ".xml") == false {
			log.Printf("skipping %s", filename)
			continue
		}

		xmlFile, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Print(err)
			continue
		}

		// parse xml
		var r Rules.Ruleset
		xml.Unmarshal(xmlFile, &r)

		for _, target := range r.Targets {
			host := target.Host

			if val, ok := tmap[host]; ok {
				tmap[host] = val + 1
			} else {
				tmap[host] = 1
			}
		}
	}

	for k, v := range tmap {
		if v > 1 {
			fmt.Println(k)
		}
	}
}
