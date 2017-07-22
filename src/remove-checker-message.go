/**
 * @file remove-checker-message.go 
 * 
 * Remove outdated comment from https-everywhere-checker
 */

package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"bytes"
	"regexp"

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

		// ruleset is still default_off
		if len(r.Default_off) > 0 {
			continue
		}

		start := bytes.Index(xml_ctx, []byte("<!--"))
		end   := bytes.Index(xml_ctx, []byte("-->"))

		if start == -1 || end == -1 || start > end {
			continue
		}

		regex1 := `Disabled by https-everywhere-checker because:\n((Fetch error|Non-2xx HTTP code):[^\n]*\n?)*`
		regex2 := `<!--\s*-->`
		regex3 := `--><ruleset`

		re1 := regexp.MustCompile(regex1)
		re2 := regexp.MustCompile(regex2)
		re3 := regexp.MustCompile(regex3)

		xml_ctx = re1.ReplaceAll(xml_ctx, []byte(""))
		xml_ctx = re2.ReplaceAll(xml_ctx, []byte("") )
		xml_ctx = re3.ReplaceAll(xml_ctx, []byte("-->\n<ruleset") )
		xml_ctx = bytes.TrimLeft(xml_ctx, "\t\r\n ")

		err = ioutil.WriteFile(fp, xml_ctx, 0644)
		if err != nil {
			log.Print(err)
		}
	}
}
