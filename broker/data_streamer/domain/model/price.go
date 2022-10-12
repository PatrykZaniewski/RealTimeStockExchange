package model

type Price struct {
	AssetName string
	BuyPrice  float32
	SellPrice float32
}

type PriceMessage struct {
	Type string
	Price
}
