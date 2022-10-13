package model

type Price struct {
	AssetName string  `json:"assetName"`
	BuyPrice  float32 `json:"buyPrice"`
	SellPrice float32 `json:"sellPrice"`
}

type PriceMessage struct {
	Type string `json:"type"`
	Price
}
