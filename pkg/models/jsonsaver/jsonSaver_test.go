package jsonsaver

import (
	"os"
	"testing"

	"github.com/wbaker85/eve-tools/pkg/models"
)

func TestSave(t *testing.T) {
	data := []*models.Aggregate{
		{
			Name: "hello",
		},
		{
			Name: "goodbye",
		},
	}

	file, _ := os.Create("./test.json")

	j := JSONSaver{file: file}
	j.Save(data)
	os.Remove("./test.json")
}
