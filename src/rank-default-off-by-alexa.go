/**
 * @file rank-default-off-by-alexa.go 
 * 
 * Rank default_off ruleset by Alexa top-1m
 */

package main

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"libse/alexa"
	"libse/rules"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: RDO <https-everywhere/rules>")
		fmt.Println("Rank default_off ruleset by Alexa top-1m")
		os.Exit(1)
	}

	// read file from https-everywhere/rules
	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	alexa_r := Alexa.NewRanking()

	for _, file := range files {
		fn := file.Name()
		fp := filepath.Join(os.Args[1], fn)

		// assert file extension is '.xml'
		if strings.HasSuffix(fn, ".xml") == false {
			log.Printf("Skipping %s", fn)
		}

		// read file into memory
		xml_ctx, err := ioutil.ReadFile(fp)
		if err != nil {
			log.Println(err)
			continue
		}

		var r Rules.Ruleset
		xml.Unmarshal(xml_ctx, &r)

		// work on default_off only
		if len(r.Default_off) == 0 {
			continue
		}

		min_rank := 0

		for _, target := range r.Targets {
			tmp := alexa_r.GetRanking(target.Host)

			if min_rank == 0 {
				min_rank = tmp
			} else if tmp < min_rank{
				min_rank = tmp
			}
		}

		if min_rank != 0 {
			h := sha256.New()
			h.Write(xml_ctx)

			fmt.Printf("%x %s %d\n", h.Sum(nil), fn, (min_rank/ 10) * 10)
		}
	}
}
