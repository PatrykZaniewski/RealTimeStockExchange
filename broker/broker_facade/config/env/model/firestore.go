package model

type FirestoreConfig struct {
	StockConfig FirestoreStockConfig `mapstructure:"stock"`
}

type FirestoreStockConfig struct {
	ProjectId string `mapstructure:"projectId"`
}
