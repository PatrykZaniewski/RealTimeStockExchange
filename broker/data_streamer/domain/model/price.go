package model

type Price struct {
	AssetName string  `json:"assetName"`
	BuyPrice  float64 `json:"buyPrice"`
	SellPrice float64 `json:"sellPrice"`
}

type PriceMessage struct {
	Type string `json:"type"`
	Price
}
