package lib

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

// SaveJSON is for saving a json file from a provided slice of data and the desired
// path
func SaveJSON(path string, data interface{}) {
	var jStr strings.Builder
	json.NewEncoder(&jStr).Encode(data)

	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString(jStr.String())
}
