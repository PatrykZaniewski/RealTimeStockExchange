package model

type StockOrder struct {
	AssetName    string
	Quantity     int
	OrderType    string
	OrderSubtype string
	OrderPrice   float32
	ClientId     string
	BrokerId     string
	Id           string
}
