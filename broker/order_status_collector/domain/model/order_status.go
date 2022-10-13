package model

type OrderStatusEnum string

const (
	CREATED   OrderStatusEnum = "CREATED"
	FULFILLED                 = "FULFILLED"
)

type OrderStatus struct {
	AssetName    string          `json:"assetName"`
	Quantity     int             `json:"quantity"`
	OrderType    string          `json:"orderType"`
	OrderSubtype string          `json:"orderSubtype"`
	OrderPrice   float32         `json:"orderPrice"`
	ClientId     string          `json:"clientId"`
	BrokerId     string          `json:"brokerId"`
	Id           string          `json:"id"`
	Status       OrderStatusEnum `json:"status"`
}
