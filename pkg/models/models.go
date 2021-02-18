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

// OrderItem is the collected information about an item based on the active buy
// and sell orders at a given station.
type OrderItem struct {
	ID        int
	Name      string
	SellPrice float64
	BuyPrice  float64
}
