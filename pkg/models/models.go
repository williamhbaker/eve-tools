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
	Issued       int64
}

// CharacterOrder is the data for an order from a character
type CharacterOrder struct {
	Name       string
	Duration   int
	IsBuy      bool
	Issued     string
	LocationID int
	MinVolume  int
	OrderID    int
	Price      float64
	Range      string
	RegionID   int
	State      string
	TypeID     int
	VolRemain  int
	VolTotal   int
}

// CharacterAsset is the data for an asset of a character
type CharacterAsset struct {
	Name         string
	TypeID       int
	LocationFlag string
	LocationID   int
	LocationType string
	Quantity     int
}

type ProfitItem struct {
	Name      string
	SoldQty   int
	BuyQty    int
	AvgSell   float64
	AvgBuy    float64
	SoldVal   float64
	BoughtVal float64
	Profit    float64
}

type DateSale struct {
	Date       string
	TotalSales float64
}
