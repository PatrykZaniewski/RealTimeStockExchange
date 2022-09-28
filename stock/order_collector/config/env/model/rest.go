package model

type RestConfig struct {
	Port int `mapstructure:"port"`
}

type RestStockConfig struct {
	StockExchangeCore string `mapstructure:"stockExchangeCore"`
}
