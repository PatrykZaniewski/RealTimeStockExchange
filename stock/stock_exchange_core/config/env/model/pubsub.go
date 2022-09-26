package model

type PubSubConfig struct {
	ListenerType   string               `mapstructure:"listenerType"`
	PublisherType  string               `mapstructure:"publisherType"`
	PublisherLocal bool                 `mapstructure:"publisherLocal"`
	Stock          PubSubStockConfig    `mapstructure:"stock"`
	Broker         []PubSubBrokerConfig `mapstructure:"broker"`
}

type PubSubStockConfig struct {
	ProjectId string                     `mapstructure:"projectId"`
	Publisher PubSubStockPublisherConfig `mapstructure:"publisher"`
	Consumer  PubSubStockConsumerConfig  `mapstructure:"consumer"`
}

type PubSubBrokerConfig struct {
	Id        string                      `mapstructure:"id"`
	ProjectId string                      `mapstructure:"projectId"`
	Publisher PubSubBrokerPublisherConfig `mapstructure:"publisher"`
}

type PubSubStockPublisherConfig struct {
	PricesTopicId string `mapstructure:"pricesTopicId"`
}

type PubSubStockConsumerConfig struct {
	InternalOrdersTopicId string `mapstructure:"internalOrdersTopicId"`
	InternalOrdersSubId   string `mapstructure:"internalOrdersSubId"`
}

type PubSubBrokerPublisherConfig struct {
	OrdersStatusTopicId string `mapstructure:"ordersStatusTopicId"`
}
