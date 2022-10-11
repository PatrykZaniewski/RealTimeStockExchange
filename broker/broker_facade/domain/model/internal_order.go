package model

type InternalOrder struct {
	AssetName    string
	Quantity     int
	OrderType    string
	OrderSubtype string
	ClientId     string
	Id           string
}
