package models

// MakeAggregates takes a slice of transactions, aggregates them by item name,
// and returns a map of the result
func MakeAggregates(transactions []*Transaction) map[string]*Aggregate {
	output := make(map[string]*Aggregate)

	for _, trans := range transactions {
		var i *Aggregate

		if output[trans.Name] == nil {
			output[trans.Name] = &Aggregate{Name: trans.Name}
		}
		i = output[trans.Name]

		thisType := trans.Type
		if thisType == "Buy" {
			i.Bought += trans.Quantity
			i.Spent += trans.Value
		} else {
			i.Sold += trans.Quantity
			i.Earned += trans.Value
			i.Tax += trans.Tax
		}

		i.Profit = i.Earned + i.Spent + i.Tax
	}

	return output
}
