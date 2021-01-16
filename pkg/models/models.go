package models

// Aggregate represents the aggregated data for a single item that you have
// done transactions for.
type Aggregate struct {
	Name   string
	Bought int
	Sold   int
	Tax    float64
	Spent  float64
	Earned float64
	Profit float64
}

// Transaction represents a single transaction from the jEveAssets history
type Transaction struct {
	Date     string
	Name     string
	Quantity int
	Price    float64
	Tax      float64
	Value    float64
	Owner    string
	Station  string
	Region   string
	Client   string
	Type     string
}
