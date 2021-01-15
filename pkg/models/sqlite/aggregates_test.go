package sqlite

import (
	"testing"

	"github.com/wbaker85/eve-tools/pkg/models"
)

func TestAddMany(t *testing.T) {
	testModel := AggregateModel{}

	testItems := []*models.Aggregate{}
	testItems = append(testItems, &models.Aggregate{
		"test name",
		1,
		2,
		3,
		4,
		5,
		6,
	})

	testModel.AddMany(testItems)
}
