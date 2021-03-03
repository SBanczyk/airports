package main
import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"net/http"
)
func main() {
	resp, err := http.Get("http://localhost:8080/dir/query.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	r := csv.NewReader(resp.Body)
	file1, f1 := os.Create("airports.db")
	if f1 != nil {
		log.Fatal(f1)
	}
	defer file1.Close()
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
		item := record[0]
		itemLabel := record[1]
		icao := record[2]
		iata := record[3]
		point := record[4]
		if len(iata) !=3 {
			iata=""
		}
		lat, lon := calculatePosition(point)
		_, err1 := fmt.Fprintf(file1, "%v|%v|%v|%v|%v|%v\n", item, itemLabel, icao, iata, lat, lon)
		if err1 != nil {
			log.Fatal(err1)
		}
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