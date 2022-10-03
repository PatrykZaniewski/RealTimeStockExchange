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
}

type PubSubBrokerPublisherConfig struct {
	BrokerInternalClientOrdersTopicId string `mapstructure:"brokerInternalClientOrdersTopicId"`
}
