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
