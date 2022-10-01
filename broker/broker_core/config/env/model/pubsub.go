package model

type PubSubConfig struct {
	ListenerType   string             `mapstructure:"listenerType"`
	PublisherType  string             `mapstructure:"publisherType"`
	PublisherLocal bool               `mapstructure:"publisherLocal"`
	Broker         PubSubBrokerConfig `mapstructure:"broker"`
}

type PubSubBrokerConfig struct {
	ProjectId string                      `mapstructure:"projectId"`
	Publisher PubSubBrokerPublisherConfig `mapstructure:"publisher"`
	Consumer  PubSubBrokerConsumerConfig  `mapstructure:"consumer"`
}

type PubSubBrokerPublisherConfig struct {
	BrokerPendingOrdersTopicId string `mapstructure:"brokerPendingOrdersTopicId"`
}

type PubSubBrokerConsumerConfig struct {
	BrokerInternalOrdersTopicId string `mapstructure:"brokerInternalOrdersTopicId"`
	BrokerInternalOrdersSubId   string `mapstructure:"brokerInternalOrdersSubId"`
}
