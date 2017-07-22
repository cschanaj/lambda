package Alexa

import (
	"archive/zip"
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/publicsuffix"
)

type Ranking map[string]int

func NewRanking() Ranking {
	// URL to Alexa top-1m.csv.zip
	url := "https://s3.amazonaws.com/alexa-static/top-1m.csv.zip"

	// perform the network download
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// write the zipped content to temporary file
	tmpf, err := ioutil.TempFile("", "top-1m.csv.zip")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpf.Name())

	if _, err := tmpf.Write(data); err != nil {
		log.Fatal(err)
	}

	if err := tmpf.Close(); err != nil {
		log.Fatal(err)
	}

	// unzip
	r, err := zip.OpenReader(tmpf.Name())
	if err != nil {
		log.Fatal(err)
	}

	// construct return value
	alexa := make(map[string]int)

	for _, f := range r.File {
		if f.Name != "top-1m.csv" {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}

		r := csv.NewReader(rc)
		for i := 1; ; i++ {
			record, err := r.Read()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal(err)
			}

			alexa[record[1]] = i
		}
		rc.Close()
	}
	return Ranking(alexa)
}

func (alexa Ranking) GetRanking(domain string) int {
	d, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		log.Println(err)
		return 0
	}

	if val, ok := alexa[d]; ok {
		return val
	}
	return 0
}
