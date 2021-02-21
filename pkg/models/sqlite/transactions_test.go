package sqlite

import (
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wbaker85/eve-tools/pkg/models"
)

func TestLoadAndGetData(t *testing.T) {
	tests := []struct {
		name string
		list []*models.Transaction
		want []*models.Transaction
	}{
		{
			"Empty list",
			[]*models.Transaction{},
			[]*models.Transaction{},
		},
		{
			"List with one thing",
			[]*models.Transaction{
				{
					Name: "hello",
				},
			},
			[]*models.Transaction{
				{
					ID:   1,
					Name: "hello",
				},
			},
		},
		{
			"List with more than one thing",
			[]*models.Transaction{
				{
					Name: "hello",
				},
				{
					Name: "goodbye",
				},
			},
			[]*models.Transaction{
				{
					ID:   1,
					Name: "hello",
				},
				{
					ID:   2,
					Name: "goodbye",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, closer := newTestDB(t)
			defer closer()

			m := TransactionModel{DB: db}
			m.LoadData(tt.list)
			got := m.GetAll()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got != want")
			}
		})
	}
}
