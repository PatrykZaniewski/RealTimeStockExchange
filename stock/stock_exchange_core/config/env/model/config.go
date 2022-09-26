package model

type Config struct {
	Rest   RestConfig   `mapstructure:"rest"`
	Debug  DebugConfig  `mapstructure:"debug"`
	PubSub PubSubConfig `mapstructure:"pubsub"`
}
