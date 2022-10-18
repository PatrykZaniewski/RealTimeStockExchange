package model

import "time"

type ArchiveOrder struct {
	Quantity      int       `json:"quantity"`
	Price         float64   `json:"Price"`
	ClientId      string    `json:"clientId"`
	BrokerId      string    `json:"brokerId"`
	OrderId       string    `json:"orderId"`
	RequestTime   time.Time `json:"requestTime"`
	ExecutionTime time.Time `json:"executionTime"`
}
