package model

type Config struct {
	Rest     RestConfig    `mapstructure:"rest"`
	General  GeneralConfig `mapstructure:"general"`
	PubSub   PubSubConfig  `mapstructure:"pubsub"`
	Identity string        `mapstructure:"identity"`
}
