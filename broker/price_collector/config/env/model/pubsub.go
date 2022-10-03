package model

type PubSubConfig struct {
	ListenerType   string             `mapstructure:"listenerType"`
	PublisherType  string             `mapstructure:"publisherType"`
	PublisherLocal bool               `mapstructure:"publisherLocal"`
	Stock          PubSubStockConfig  `mapstructure:"stock"`
	Broker         PubSubBrokerConfig `mapstructure:"broker"`
}

type PubSubStockConfig struct {
	ProjectId string                    `mapstructure:"projectId"`
	Consumer  PubSubStockConsumerConfig `mapstructure:"consumer"`
}

type PubSubBrokerConfig struct {
	ProjectId string                      `mapstructure:"projectId"`
	Publisher PubSubBrokerPublisherConfig `mapstructure:"publisher"`
}

type PubSubStockConsumerConfig struct {
	CorePricesTopicId string `mapstructure:"corePricesTopicId"`
	CorePricesSubId   string `mapstructure:"corePricesSubId"`
}

type PubSubBrokerPublisherConfig struct {
	BrokerInternalPricesTopicId string `mapstructure:"brokerInternalPricesTopicId"`
}
