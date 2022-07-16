package main

import (
	"encoding/json"
	"flag"
	"github.com/IncSW/geoip2"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = ""
	}
	dbFile := ""
	flag.StringVar(&dbFile, "file", pwd+"/data/GeoLite2-City.mmdb", "")
	flag.Parse()

	mmdb, err := geoip2.NewCityReaderFromFile(dbFile)
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !r.URL.Query().Has("ip") {
			http.NotFound(w, r)
			return
		}

		cityResult, err := mmdb.Lookup(net.ParseIP(r.URL.Query().Get("ip")))
		if err != nil {
			http.NotFound(w, r)
			return
		}

		region := ""
		if len(cityResult.Subdivisions) > 0 {
			region = cityResult.Subdivisions[0].Names["en"]
		}

		marshal, err := json.Marshal(struct {
			City    string `json:"city"`
			Country string `json:"country"`
			Region  string `json:"region"`
		}{
			City:    cityResult.City.Names["en"],
			Country: cityResult.Country.Names["en"],
			Region:  region,
		})

		if err != nil {
			http.NotFound(w, r)
			return
		}

		_, _ = w.Write(marshal)
	})

	log.Println("Listening on 0.0.0.0:1212...")
	err = http.ListenAndServe(":1212", nil)
	log.Panicln(err)
}
