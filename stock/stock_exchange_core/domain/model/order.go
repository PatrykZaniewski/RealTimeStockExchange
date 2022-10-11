package model

type StockOrder struct {
	AssetName    string
	Quantity     int
	OrderType    string
	OrderSubtype string
	ClientId     string
	BrokerId     string
	Id           string
}
