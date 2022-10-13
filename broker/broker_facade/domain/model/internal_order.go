package model

type InternalOrder struct {
	AssetName    string  `json:"assetName"`
	Quantity     int     `json:"quantity"`
	OrderType    string  `json:"orderType"`
	OrderSubtype string  `json:"orderSubtype"`
	OrderPrice   float32 `json:"orderPrice"`
	ClientId     string  `json:"clientId"`
	Id           string  `json:"id"`
}
