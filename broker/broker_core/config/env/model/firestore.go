package model

type FirestoreConfig struct {
	StockConfig FirestoreBrokerConfig `mapstructure:"broker"`
}

type FirestoreBrokerConfig struct {
	ProjectId string `mapstructure:"projectId"`
}
