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
	Id        string                     `mapstructure:"id"`
	ProjectId string                     `mapstructure:"projectId"`
	Consumer  PubSubBrokerConsumerConfig `mapstructure:"consumer"`
}

type PubSubStockPublisherConfig struct {
	InternalOrdersTopicId string `mapstructure:"internalOrdersTopicId"`
}

type PubSubStockConsumerConfig struct {
	BrokerMockOrdersTopicId string `mapstructure:"brokerMockOrdersTopicId"`
	BrokerMockOrdersSubId   string `mapstructure:"brokerMockOrdersSubId"`
}

type PubSubBrokerConsumerConfig struct {
	BrokerPendingOrdersTopicId string `mapstructure:"brokerPendingOrdersTopicId"`
	BrokerPendingOrdersSubId   string `mapstructure:"brokerPendingOrdersSubId"`
}
