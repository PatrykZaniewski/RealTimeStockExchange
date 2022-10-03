package model

type PubSubConfig struct {
	ListenerType   string             `mapstructure:"listenerType"`
	PublisherType  string             `mapstructure:"publisherType"`
	PublisherLocal bool               `mapstructure:"publisherLocal"`
	Broker         PubSubBrokerConfig `mapstructure:"broker"`
}

type PubSubBrokerConfig struct {
	ProjectId string                     `mapstructure:"projectId"`
	Consumer  PubSubBrokerConsumerConfig `mapstructure:"consumer"`
}

type PubSubBrokerConsumerConfig struct {
	BrokerInternalPricesTopicId string `mapstructure:"brokerInternalPricesTopicId"`
	BrokerInternalPricesSubId   string `mapstructure:"brokerInternalPricesSubId"`
}
