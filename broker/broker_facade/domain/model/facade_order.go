package model

type FacadeOrder struct {
	AssetName    string  `json:"assetName"`
	Quantity     int     `json:"quantity"`
	OrderType    string  `json:"orderType"`
	OrderSubtype string  `json:"orderSubtype"`
	OrderPrice   float64 `json:"orderPrice"`
	Id           string  `json:"id"`
}
