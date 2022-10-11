package model

type OrderStatusEnum string

const (
	CREATED   OrderStatusEnum = "CREATED"
	FULFILLED                 = "FULFILLED"
)

type OrderStatus struct {
	AssetName    string
	Quantity     int
	OrderType    string
	OrderSubtype string
	ClientId     string
	BrokerId     string
	Id           string
	Status       OrderStatusEnum
}
