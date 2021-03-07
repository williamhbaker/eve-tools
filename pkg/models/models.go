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
	ID           int
	Name         string
	SellPrice    float64
	BuyPrice     float64
	RecentOrders int
}

// MarginItem includes details about the margin from an item based on a buy station
// and a sell station
type MarginItem struct {
	ItemID       int
	Name         string
	SellPrice    float64
	BuyPrice     float64
	Margin       float64
	RecentOrders int
}

// ItemHistoryData includes information about an item in a region, with values
// derived from historical analysis.
type ItemHistoryData struct {
	RegionID    int
	ItemID      int
	NumDays     int
	OrdersAvg   int
	VolumeAvg   int
	YearMinSell float64
	YearMaxSell float64
	YearMinBuy  float64
	YearMaxBuy  float64
}

// ClientID is used by the esi login to register the specified client.
type ClientID struct {
	value string
}

// ClientSecret is also used by the esi login.
type ClientSecret struct {
	value string
}

// AuthToken is an ESI auth token
type AuthToken struct {
	AccessToken  string
	ExpiresIn    int
	RefreshToken string
}
