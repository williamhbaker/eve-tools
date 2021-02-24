package lib

import (
	"encoding/csv"
	"log"
	"os"
)

// SaveCSV saves a slice of records to the given path.
func SaveCSV(path string, records [][]string) {
	file, _ := os.Create(path)

	w := csv.NewWriter(file)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
