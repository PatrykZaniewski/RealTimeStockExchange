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
	BrokerPendingOrdersTopicId            string `mapstructure:"brokerPendingOrdersTopicId"`
	brokerInternalCoreOrdersStatusTopicId string `mapstructure:"brokerInternalCoreOrdersStatusTopicId"`
}

type PubSubBrokerConsumerConfig struct {
	BrokerInternalClientOrdersTopicId string `mapstructure:"brokerInternalClientOrdersTopicId"`
	BrokerInternalClientOrdersSubId   string `mapstructure:"brokerInternalClientOrdersSubId"`
	BrokerInternalOrdersStatusTopicId string `mapstructure:"brokerInternalOrdersStatusTopicId"`
	BrokerInternalOrdersStatusSubId   string `mapstructure:"brokerInternalOrdersStatusSubId"`
	BrokerInternalPricesTopicId       string `mapstructure:"brokerInternalPricesTopicId"`
	BrokerInternalPricesSubId         string `mapstructure:"brokerInternalPricesSubId"`
}
