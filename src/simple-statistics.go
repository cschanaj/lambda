/**
 * @file simpile-statistics.go
 *
 * Statistics on the structure of https-everywhere code base
 * https://github.com/EFForg/https-everywhere/issues/10378
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

	"golang.org/x/net/publicsuffix"

	"libse/rules"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s <path/https-everywhere/rules\n")
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// map of 'default_off'
	dmap := make(map[string]int)

	// map of 'platform'
	pmap := make(map[string]int)

	// map of 'rule'
	rmap := make(map[string]int)

	// map of 'securecookie'
	smap := make(map[string]int)

	// map of 'exclusion'
	emap := make(map[string]int)
	nlah := 0

	// map of 'target'
	tmap := make(map[string]int)

	// count left-wildcard, right-wildcard 'target'
	lwc_target := 0
	rwc_target := 0

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

		// count 'default_off'
		if len(r.Default_off) > 0 {
			if _, ok := dmap[r.Default_off]; ok {
				dmap[r.Default_off]++
			} else {
				dmap[r.Default_off] = 1
			}
		}

		// count 'platform'
		if len(r.Platform) > 0 {
			if _, ok := pmap[r.Platform]; ok {
				pmap[r.Platform]++
			} else {
				pmap[r.Platform] = 1
			}
		}

		// count 'securecookie'
		if len(r.SecureCookies) > 0 {
			for _, sc := range r.SecureCookies {
				key := sc.Host + sc.Name
				if _, ok := smap[key]; ok {
					smap[key]++
				} else {
					smap[key] = 1
				}
			}
		}

		// count 'rule'
		if len(r.Rules) > 0 {
			for _, rule := range r.Rules {
				key := rule.From + rule.To
				if _, ok := rmap[key]; ok {
					rmap[key]++
				} else {
					rmap[key] = 1
				}
			}
		}

		// count 'exclusion'
		if len(r.Exclusions) > 0 {
			for _, ex := range r.Exclusions {
				key := ex.Pattern

				if strings.Contains(key, "?!") {
					nlah++
				}

				if _, ok := emap[key]; ok {
					emap[key]++
				} else {
					emap[key] = 1
				}
			}
		}

		// count 'target'
		if len(r.Targets) > 0 {
			for _, target := range r.Targets {
				if strings.HasPrefix(target.Host, "*.") {
					lwc_target += 1
				}

				if strings.HasSuffix(target.Host, ".*") {
					rwc_target += 1
				}

				d, err := publicsuffix.EffectiveTLDPlusOne(target.Host)
				if err != nil {
					log.Println(err)
					break
				}

				if _, ok := tmap[d]; ok {
					tmap[d]++
				} else {
					tmap[d] = 1
				}
			}
		}
	}


	fmt.Printf("|  | %d | %d | %d | %d |\n", MyMapSum(pmap), len(pmap), MyMapSum(dmap), len(dmap))
	fmt.Printf("|  | %d | %d | %d |\n", pmap["mixedcontent"], pmap["cacert"], pmap["cacert mixedcontent"])
	fmt.Printf("|  | %d |\n", dmap["failed ruleset test"])
	fmt.Printf("|  | %d | %d | %d | %d |\n", MyMapSum(tmap), len(tmap), lwc_target, rwc_target)
	fmt.Printf("|  | %d | %d | %d |\n", MyMapSum(emap), len(emap), nlah)
	fmt.Printf("|  | %d | %d |\n", MyMapSum(smap), len(smap))
	fmt.Printf("|  | %d | %d |\n", MyMapSum(rmap), len(rmap))
}

func MyMapSum(pmap map[string]int) int {
	ret := 0
	for _, val := range pmap {
		ret += val
	}
	return ret
}
