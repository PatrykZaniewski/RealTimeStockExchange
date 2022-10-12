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
	OrderPrice   float32
	ClientId     string
	BrokerId     string
	Id           string
	Status       OrderStatusEnum
}

type OrderStatusMessage struct {
	Type string
	OrderStatus
}
