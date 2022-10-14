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
	BrokerInternalClientOrdersTopicId string `mapstructure:"brokerInternalClientOrdersTopicId"`
}

type PubSubBrokerConsumerConfig struct {
	ClientMockOrdersTopicId string `mapstructure:"clientMockOrdersTopicId"`
	ClientMockOrdersSubId   string `mapstructure:"clientMockOrdersSubId"`
}
