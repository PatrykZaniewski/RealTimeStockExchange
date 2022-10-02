package model

type Order struct {
	AssetName string  `json:"assetName"`
	Quantity  float64 `json:"quantity"`
}
