package model

type StockOrder struct {
	AssetName    string  `json:"assetName"`
	Quantity     int     `json:"quantity"`
	OrderType    string  `json:"orderType"`
	OrderSubtype string  `json:"orderSubtype"`
	OrderPrice   float64 `json:"orderPrice"`
	ClientId     string  `json:"clientId"`
	BrokerId     string  `json:"brokerId"`
	Id           string  `json:"id"`
}
