package main

import (
	"encoding/xml"
	"io/ioutil"
	"fmt"
	"sync"
	"log"
	"regexp"
	"net"
	"os"
	"runtime"
	"path/filepath"
	"strings"

	"libse/rules"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: dnsp [FILE.xml]...")
		fmt.Println("Remove 'target' which is gone from DNS")
		os.Exit(1)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	for _, file := range os.Args[1:] {
		_, fn := filepath.Split(file)

		if strings.HasSuffix(fn, ".xml") == false {
			os.Exit(1)
			continue
		}

		xmlFile, err := ioutil.ReadFile(file)
		if err != nil {
			os.Exit(1)
		}

		var r Rules.Ruleset
		xml.Unmarshal(xmlFile, &r)

		logger := log.New(os.Stdout, "", 0)

		hmap := make(map[string]bool)
		var mutex = &sync.Mutex{}
		var wg sync.WaitGroup

		for _, target := range r.Targets {

			go func(host string) {

				wg.Add(1)
				defer wg.Done()

				if isGoneFromDNS(host) {
					logger.Printf("%s\n", host)

					mutex.Lock()
					hmap[host] = true
					mutex.Unlock()
				}
			}(target.Host)
		}

		wg.Wait()

		for k, _ := range hmap {
			h := regexp.QuoteMeta(k)
			x := `\n[ \t]*<target\s+host\s*=\s*"` + h + `"\s*/>\s*?\n`

			re := regexp.MustCompile(x)
			xmlFile = re.ReplaceAll(xmlFile, []byte("\n"))
		}

		if err := ioutil.WriteFile(file, xmlFile, 0644); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
}


func isGoneFromDNS(host string) bool {
	_, err := net.LookupHost(host)
	return (err != nil)
}
