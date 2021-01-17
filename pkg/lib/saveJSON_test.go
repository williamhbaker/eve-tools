package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestSaveJSON(t *testing.T) {
	data1 := []*Aggregate{}
	data2 := []*Aggregate{
		{
			Name: "hello",
		},
		{
			Name: "goodbye",
		},
	}

	tests := []struct {
		name string
		data []*Aggregate
	}{
		{
			"Empty list",
			data1,
		},
		{
			"Non-empty list",
			data2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var want strings.Builder
			json.NewEncoder(&want).Encode(tt.data)

			SaveJSON("./test.json", tt.data)

			writtenFile, _ := os.Open("./test.json")
			written, _ := ioutil.ReadAll(writtenFile)
			os.Remove("./test.json")

			if string(written) != want.String() {
				t.Errorf("got %s want %s", string(written), want.String())
			}
		})
	}
}
