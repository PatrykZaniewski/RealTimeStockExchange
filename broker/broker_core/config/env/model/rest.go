package model

type RestConfig struct {
	RestStockConfig int `mapstructure:"stock"`
}

type RestStockConfig struct {
	OrderExecutor string `mapstructure:"orderExecutor"`
}
