package model

type RestConfig struct {
	RestStockConfig int `mapstructure:"stock"`
}

type RestStockConfig struct {
	OrderCollector string `mapstructure:"orderCollector"`
}
