package lib

import (
	"reflect"
	"testing"

	"github.com/wbaker85/eve-tools/pkg/models"
)

func TestMakeAggregates(t *testing.T) {
	testList1 := []*models.Transaction{}

	testList2 := []*models.Transaction{
		{
			Name:     "item 1",
			Quantity: 2,
			Value:    -10.5,
			Type:     "Buy",
			Tax:      0,
		},
		{
			Name:     "item 1",
			Quantity: 1,
			Value:    10.5,
			Type:     "Sell",
			Tax:      -1.5,
		},
		{
			Name:     "item 2",
			Quantity: 1,
			Value:    -10.5,
			Type:     "Buy",
			Tax:      0,
		},
		{
			Name:     "item 3",
			Quantity: 1,
			Value:    25.0,
			Type:     "Sell",
			Tax:      -2.5,
		},
	}

	tests := []struct {
		name         string
		transactions []*models.Transaction
		want         map[string]*Aggregate
	}{
		{
			"Empty list",
			testList1,
			make(map[string]*Aggregate),
		},
		{
			"List with stuff",
			testList2,
			map[string]*Aggregate{
				"item 1": {Name: "item 1", Bought: 2, Sold: 1, Tax: -1.5, Spent: -10.5, Earned: 10.5, Profit: -1.5},
				"item 2": {Name: "item 2", Bought: 1, Sold: 0, Tax: 0, Spent: -10.5, Earned: 0, Profit: -10.5},
				"item 3": {Name: "item 3", Bought: 0, Sold: 1, Tax: -2.5, Spent: 0, Earned: 25, Profit: 22.5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MakeAggregates(tt.transactions)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got != want")
			}
		})
	}

}
