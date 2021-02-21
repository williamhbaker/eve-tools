package sqlite

import (
	"math"
	"reflect"
	"sort"
	"testing"

	"github.com/wbaker85/eve-tools/pkg/models"
)

func TestEscapeString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			"basic string with no substitutions",
			"125mm Prototype Gauss Gun",
			"125mm Prototype Gauss Gun",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeString(tt.in)
			if got != tt.want {
				t.Errorf("got %s want %s", got, tt.want)
			}
		})
	}
}

func TestLoadData(t *testing.T) {
	sellStationID := 1234
	buyStationID := 4321

	sellData1 := map[int]*models.OrderItem{
		1: {
			ID:        1,
			Name:      "First Item",
			SellPrice: 10.1,
			BuyPrice:  5.5,
		},
		2: {
			ID:        2,
			Name:      "Second Item",
			SellPrice: 10.1,
			BuyPrice:  5.5,
		},
		3: {
			ID:        3,
			Name:      "Third Item",
			SellPrice: 10.1,
			BuyPrice:  5.5,
		},
	}

	buyData1 := map[int]*models.OrderItem{
		1: {
			ID:        1,
			Name:      "First Item",
			SellPrice: 9.1,
			BuyPrice:  6.5,
		},
		2: {
			ID:        2,
			Name:      "Second Item",
			SellPrice: 10.1,
			BuyPrice:  5.5,
		},
		3: {
			ID:        3,
			Name:      "Third Item",
			SellPrice: 11.1,
			BuyPrice:  4.5,
		},
	}

	want1 := []*models.MarginItem{
		{
			ID:        1,
			Name:      "First Item",
			SellPrice: 10.1,
			BuyPrice:  6.5,
			Margin:    (10.1 - 6.5) / 6.5 * 100,
		},
		{
			ID:        2,
			Name:      "Second Item",
			SellPrice: 10.1,
			BuyPrice:  5.5,
			Margin:    (10.1 - 5.5) / 5.5 * 100,
		},
		{
			ID:        3,
			Name:      "Third Item",
			SellPrice: 10.1,
			BuyPrice:  5.5,
			Margin:    (10.1 - 5.5) / 5.5 * 100,
		},
	}

	db, removeDB := newTestDB(t)
	defer removeDB()

	o := OrderModel{DB: db}

	o.LoadData(sellStationID, sellData1)
	o.LoadData(buyStationID, buyData1)

	got1 := o.GetAllMargins(sellStationID, buyStationID)

	assertMarginItems(t, got1, want1)

	// Now make sure that clearing the database actually works...

	o.init(sellStationID)
	o.init(buyStationID)

	got2 := o.GetAllMargins(sellStationID, buyStationID)
	want2 := []*models.MarginItem{}

	assertMarginItems(t, got2, want2)
}

func assertMarginItems(t *testing.T, m1, m2 []*models.MarginItem) {
	t.Helper()

	sortMarginItem(m1)
	sortMarginItem(m2)

	if !reflect.DeepEqual(m1, m2) {
		t.Errorf("got %#v want %#v", m1, m2)
	}
}

func sortMarginItem(m []*models.MarginItem) {
	for _, val := range m {
		val.Margin = math.Round(val.Margin*100) / 100
	}

	sort.Slice(m, func(i, j int) bool {
		return m[i].ID < m[j].ID
	})
}
