package model

import "time"

type PendingOrder struct {
	AssetName    string    `json:"assetName"`
	Quantity     int       `json:"quantity"`
	OrderType    string    `json:"orderType"`
	OrderSubtype string    `json:"orderSubtype"`
	OrderPrice   float64   `json:"orderPrice"`
	ClientId     string    `json:"clientId"`
	Id           string    `json:"id"`
	RequestTime  time.Time `json:"requestTime"`
}
