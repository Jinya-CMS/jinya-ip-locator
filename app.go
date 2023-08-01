package main

import (
	"compress/gzip"
	"encoding/json"
	"github.com/IncSW/geoip2"
	"github.com/go-co-op/gocron"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func downloadIpToLocationDb() {
	date := time.Now()
	res, err := http.Get("https://download.db-ip.com/free/dbip-city-lite-" + date.AddDate(0, -1, 0).Format("2006-01") + ".mmdb.gz")
	if err != nil {
		log.Println("Failed to download IP2L database")
		return
	}

	gzReader, err := gzip.NewReader(res.Body)
	defer gzReader.Close()
	if err != nil {
		log.Println("Failed to download IP2L database")
		return
	}

	data, err := io.ReadAll(gzReader)
	if err != nil {
		log.Println("Failed to extract IP2L database")
		return
	}

	err = os.MkdirAll("./data", 0775)
	if err != nil {
		log.Println("Failed to create IP2L database directory")
		return
	}
	err = os.WriteFile("./data/ip2l.mmdb", data, 0775)
	if err != nil {
		log.Println("Failed to write IP2L database")
	}
}

func main() {
	downloadIpToLocationDb()

	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(1).Day().Do(downloadIpToLocationDb)
	if err != nil {
		log.Println("Couldn't setup cron job please be aware, that the database will not be refreshed")
	}

	s.StartAsync()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !r.URL.Query().Has("ip") {
			http.NotFound(w, r)
			return
		}

		type data struct {
			City       string `json:"city"`
			Country    string `json:"country"`
			Region     string `json:"region"`
			ProvidedBy string `json:"providedBy"`
		}

		returnValue := data{}

		mmdb, err := geoip2.NewCityReaderFromFile("./data/ip2l.mmdb")
		if err != nil {
			returnValue = data{
				City:       "Unknown city",
				Country:    "Unknown country",
				Region:     "Unknown region",
				ProvidedBy: "https://db-ip.com",
			}
		}

		cityResult, err := mmdb.Lookup(net.ParseIP(r.URL.Query().Get("ip")))
		if err == nil {
			region := "Unknown region"
			if len(cityResult.Subdivisions) > 0 {
				region = cityResult.Subdivisions[0].Names["en"]
			}
			country := cityResult.Country.Names["en"]
			if country == "" {
				country = "Unknown country"
			}

			city := cityResult.City.Names["en"]
			if city == "" {
				city = "Unknow city"
			}

			returnValue = data{
				City:       city,
				Country:    country,
				Region:     region,
				ProvidedBy: "https://db-ip.com",
			}
		}

		marshal, err := json.Marshal(returnValue)

		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("X-DB-Provided-By", "https://db-ip.com")
		_, _ = w.Write(marshal)
	})

	log.Println("Listening on 0.0.0.0:1212...")
	err = http.ListenAndServe(":1212", nil)
	log.Panicln(err)
}
