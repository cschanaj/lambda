package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
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

		// ignore default_off
		if len(r.Default_off) > 0 {
			continue
		}

		ignore_file := false

		// ignore file with wildcard target
		for _, target := range r.Targets {
			if strings.Contains(target.Host, "*") {
				ignore_file = true
				break
			}
		}

		if ignore_file {
			continue
		}


		// ignore if there is some exclusions
		if len(r.Exclusions) != 0 {
			continue
		}

		// ignore if there is non-trivial rules
		for _, rule := range r.Rules {
			if rule.From != "^http:" || rule.To != "https:" {
				ignore_file = true
			}
		}

		if ignore_file {
			continue
		}


		// ignore file with wildcard securecookie
		for _, sc := range r.SecureCookies {
			temp := 0
			host := sc.Host
			name := sc.Name

			if host == ".+" || host == ".*" || host == "." {
				temp++
			}

			if name == ".+" || name == ".*" || name == "." {
				temp++
			}

			if temp == 2 {
				ignore_file = true
			}
		}

		if ignore_file {
			continue
		}

		for _, target := range r.Targets {
			ok   := false
			host := target.Host

			for _, sc := range r.SecureCookies {
				cookiehost := sc.Host
				cookiename := sc.Name

				if matched, err := regexp.Match(cookiehost, []byte(host)); err != nil {
					log.Print(err)
					ignore_file = true
					break
				} else if matched {
					if cookiename == ".+" || cookiename == ".*" || cookiename == "." || cookiename == `^\w` {
						ok = true
						break
					}
				}
			}

			if !ok {
				ignore_file = true
				break
			}
		}

		if ignore_file {
			continue
		}

		// https://github.com/EFForg/https-everywhere/blob/master/utils/hsts-prune/index.js#L186
		// const target_regex = RegExp(`\n[ \t]*<target host=\\s*"${target}"\\s*/>\\s*?\n`);
		pattern := `\n[ \t]*<securecookie\s+host=\s*"[^"]*"\s+name=\s*"[^"]*"\s*/>`
		re := regexp.MustCompile(pattern)

		x := re.FindAll(xmlFile, -1)
		if len(x) != len(r.SecureCookies) {
			log.Printf("%s: Unspecified behaviour: %d != %d", file.Name(), len(x), len(r.SecureCookies))
			continue
		}

		for i := 1; i< len(x); i++ {
			xmlFile = bytes.Replace(xmlFile, x[i], []byte(""), 1)
		}

		for i := 0; i < 1; i++ {
			xmlFile = bytes.Replace(xmlFile, x[0], []byte("\n\t" + `<securecookie host=".+" name=".+" />`), 1)
		}

		if err := ioutil.WriteFile(filename, xmlFile, 0644); err != nil {
			log.Print(err)
		}
	}
}
