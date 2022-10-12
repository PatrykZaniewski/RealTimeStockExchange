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

type PubSubBrokerConsumerConfig struct {
	BrokerOrdersTopicId string `mapstructure:"brokerOrdersTopicId"`
	BrokerOrdersSubId   string `mapstructure:"brokerOrdersSubId"`
}

type PubSubBrokerPublisherConfig struct {
	BrokerInternalOrdersStatusTopicId string `mapstructure:"brokerInternalOrdersStatusTopicId"`
}
