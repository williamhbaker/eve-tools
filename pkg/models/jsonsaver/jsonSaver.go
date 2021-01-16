package jsonsaver

import (
	"encoding/json"
	"os"
	"strings"
)

// JSONSaver handles everything for saving to a JSON file.
type JSONSaver struct {
	file *os.File
}

// Save is for saving a json file from a provided slice of data
func (j *JSONSaver) Save(data interface{}) {
	var jStr strings.Builder
	json.NewEncoder(&jStr).Encode(data)

	j.file.WriteString(jStr.String())

	j.file.Close()
}
