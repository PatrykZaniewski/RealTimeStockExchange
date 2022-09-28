package model

type Config struct {
	Rest      RestConfig      `mapstructure:"rest"`
	General   GeneralConfig   `mapstructure:"general"`
	PubSub    PubSubConfig    `mapstructure:"pubsub"`
	Firestore FirestoreConfig `mapstructure:"firestore"`
}
