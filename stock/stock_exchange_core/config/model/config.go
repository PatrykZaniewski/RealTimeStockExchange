package model

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Rest   RestConfig   `mapstructure:"rest"`
	Debug  DebugConfig  `mapstructure:"debug"`
	PubSub PubSubConfig `mapstructure:"pubsub"`
}

var StockConfig Config

func ConfigSetup() {
	viper.SetConfigName("settings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error occured during env set. Err: %s", err)
	}

	err := viper.Unmarshal(&StockConfig)
	if err != nil {
		return
	}
	err = viper.Unmarshal(&StockConfig)
	if err != nil {
		return
	}
}
