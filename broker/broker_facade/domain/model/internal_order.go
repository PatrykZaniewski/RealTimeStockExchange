package model

type InternalOrder struct {
	AssetName    string
	Quantity     int
	OrderType    string
	OrderSubtype string
	OrderPrice   float32
	ClientId     string
	Id           string
}
