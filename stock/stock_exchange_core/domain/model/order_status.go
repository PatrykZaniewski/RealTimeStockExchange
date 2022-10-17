package model

type OrderStatusEnum string

const (
	CREATED   OrderStatusEnum = "CREATED"
	FULFILLED                 = "FULFILLED"
)

type OrderStatus struct {
	AssetName      string          `json:"assetName"`
	Quantity       int             `json:"quantity"`
	OrderType      string          `json:"orderType"`
	OrderSubtype   string          `json:"orderSubtype"`
	OrderPrice     float64         `json:"orderPrice"`
	ExecutionPrice float64         `json:"executionPrice"`
	ClientId       string          `json:"clientId"`
	BrokerId       string          `json:"brokerId"`
	Id             string          `json:"id"`
	Status         OrderStatusEnum `json:"status"`
}
