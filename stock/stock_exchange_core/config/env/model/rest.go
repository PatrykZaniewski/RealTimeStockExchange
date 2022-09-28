package model

type RestConfig struct {
	Broker []RestBrokerConfig `mapstructure:"broker"`
}

type RestBrokerConfig struct {
	Id           string `mapstructure:"id"`
	ordersStatus string `mapstructure:"ordersStatus"`
}
