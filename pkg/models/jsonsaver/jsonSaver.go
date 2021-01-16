package jsonsaver

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wbaker85/eve-tools/pkg/models"
)

// SaveAggregates is for saving a json file from a provided slice of data
func SaveAggregates(path string, data []*models.Aggregate) {
	var jStr strings.Builder
	json.NewEncoder(&jStr).Encode(data)
	fmt.Println(jStr.String())
}
