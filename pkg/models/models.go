package models

// Transaction represents a single transaction from the jEveAssets history
type Transaction struct {
	ID       int
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
