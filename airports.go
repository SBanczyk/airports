package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, f := os.Open("query.csv")
	if f != nil {
		log.Fatal(f)
	}
	r := csv.NewReader(file)
	flaga := false
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if !flaga {
			flaga = true
			continue
		}
		icao := record[0]
		point := record[1]
		iata := record[2]
		lat, lon := calculatePosition(point)
		fmt.Printf("ICAO:%s | IATA: %s | Position:%s %s \n", icao, iata, lat, lon)
	}
}
func calculatePosition(point string) (string, string) {
	if !strings.HasPrefix(point, "Point") {
		return "", ""
	}
	position := point[6:strings.Index(point, ")")]
	latlon := strings.Split(position, " ")
	return latlon[1], latlon[0]

}
