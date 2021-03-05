package main
import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)
func main() {
	fmt.Printf("Requesting Data from Wikidata\n")
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://query.wikidata.org/sparql?query=SELECT%20%3Fitem%20%3FitemLabel%20%3FICAO_code%20%3FIATA_code%20%3Fcoordinates%20WHERE%20%7B%0A%20%20SERVICE%20wikibase%3Alabel%20%7B%20bd%3AserviceParam%20wikibase%3Alanguage%20%22%5BAUTO_LANGUAGE%5D%2Cen%22.%20%7D%0A%20%20%3Fitem%20wdt%3AP239%20%3FICAO_code.%0A%20%20FILTER%28REGEX%28STR%28%3FICAO_code%29%2C%20%22%5E%5BA-Z0-9%5D%7B4%7D%24%22%29%29%0A%20%20OPTIONAL%20%7B%20%3Fitem%20wdt%3AP238%20%3FIATA_code.%20%7D%0A%20%20OPTIONAL%20%7B%20%3Fitem%20wdt%3AP625%20%3Fcoordinates.%20%7D%0A%20%20MINUS%20%7B%20%3Fitem%20wdt%3AP582%20%3Fend_time.%20%7D%0A%20%20minus%20%7B%20%3Fitem%20wdt%3AP576%20%3Fdissolved__abolished_or_demolished_date.%20%7D%0A%7D%0Aorder%20by%20%3FICAO_code", nil)
	req.Header.Add("User-Agent", "github.com/SBanczyk/airports")
	req.Header.Add("Accept", "text/csv")
	if err != nil {
		log.Panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	r := csv.NewReader(resp.Body)
	outputFile, err := os.Create("airports.db")
	if err != nil {
		log.Panic(err)
	}
	duplicateFile, err := os.Create("airports_duplicate.db")
	if err != nil {
		log.Panic(err)
	}
	defer duplicateFile.Close()
	defer outputFile.Close()
	firstLinePassed := false
	recordCounter := 0
	duplicate := duplicateCounter{duplicateFile: duplicateFile}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panic(err)
		}
		if !firstLinePassed {
			firstLinePassed = true
			continue
		}
		item := record[0]
		itemLabel := record[1]
		icao := record[2]
		iata := record[3]
		point := record[4]
		if len(iata) != 3 {
			iata = ""
		}
		lat, lon := calculatePosition(point)
		_, err1 := fmt.Fprintf(outputFile, "%v|%v|%v|%v|%v|%v\n", item, itemLabel, icao, iata, lat, lon)
		recordCounter++
		if err1 != nil {
			log.Panic(err1)
		}
		duplicate.checkIcao(icao)
	}
	fmt.Printf("Records: %v | Duplicates: %v\n", recordCounter, duplicate.duplicateCounter)
}

func calculatePosition(point string) (string, string) {
	if !strings.HasPrefix(point, "Point") {
		return "", ""
	}
	position := point[6:strings.Index(point, ")")]
	latlon := strings.Split(position, " ")
	return latlon[1], latlon[0]
}

type duplicateCounter struct {
	lastIcao         string
	duplicateCounter int
	duplicateFile    *os.File
	isDuplicate      bool
}

func (c *duplicateCounter) checkIcao(icao string) {
	if c.lastIcao == icao && !c.isDuplicate {
		c.duplicateCounter++
		c.isDuplicate = true
		_, err2 := fmt.Fprintf(c.duplicateFile, "%v\n", icao)
		if err2 != nil {
			log.Panic(err2)
		}
	} else if c.lastIcao == icao {
		c.isDuplicate = true
	} else {
		c.isDuplicate = false
	}
	c.lastIcao = icao
}
