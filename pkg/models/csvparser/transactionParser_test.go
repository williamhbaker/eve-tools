package csvparser

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/wbaker85/eve-tools/pkg/models"
)

func TestParseTransactions(t *testing.T) {
	testFile1 := `
	Date,Name,Quantity,Price,Tax,Value,Owner,Station,Region,Transaction Price,Transaction Margin,Transaction Profit +,Transaction Profit %,Client,Type,For,Wallet Division
`

	testFile2 := `
Date,Name,Quantity,Price,Tax,Value,Owner,Station,Region,Transaction Price,Transaction Margin,Transaction Profit +,Transaction Profit %,Client,Type,For,Wallet Division
2021-01-15 07:31,Item 1,3,73410,-100,100,YourName Here,Jita IV - Moon 4 - Caldari Navy Assembly Plant,The Forge,100,100,28790,139%,Some Client,Sell,Personal,1
2021-01-15 07:30,Item 1,1,73410,,-100,YourName Here,Jita IV - Moon 4 - Caldari Navy Assembly Plant,The Forge,100,100,28790,139%,Some Client,Buy,Personal,1
2021-01-15 07:29,Item 1,1,73410,,-200,YourName Here,Jita IV - Moon 4 - Caldari Navy Assembly Plant,The Forge,100,100,28790,139%,Some Client,Buy,Personal,1
`

	tests := []struct {
		name string
		file io.Reader
		want []*models.Transaction
	}{
		{
			"Empty file",
			strings.NewReader(testFile1),
			[]*models.Transaction{},
		},
		{
			"Non-empty file",
			strings.NewReader(testFile2),
			[]*models.Transaction{
				{ID: 0, Date: "2021-01-15 07:31", Name: "Item 1", Quantity: 3, Price: 73410, Tax: -100, Value: 100, Owner: "YourName Here", Station: "Jita IV - Moon 4 - Caldari Navy Assembly Plant", Region: "The Forge", Client: "Some Client", Type: "Sell"},
				{ID: 0, Date: "2021-01-15 07:30", Name: "Item 1", Quantity: 1, Price: 73410, Tax: 0, Value: -100, Owner: "YourName Here", Station: "Jita IV - Moon 4 - Caldari Navy Assembly Plant", Region: "The Forge", Client: "Some Client", Type: "Buy"},
				{ID: 0, Date: "2021-01-15 07:29", Name: "Item 1", Quantity: 1, Price: 73410, Tax: 0, Value: -200, Owner: "YourName Here", Station: "Jita IV - Moon 4 - Caldari Navy Assembly Plant", Region: "The Forge", Client: "Some Client", Type: "Buy"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := TransactionParser{File: tt.file}
			got := parser.ParseTransactions()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got != want")
			}
		})
	}
}
