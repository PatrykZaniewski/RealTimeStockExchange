package env

import (
	configModel "broker/price_collector/config/env/model"
	"github.com/spf13/viper"
	"log"
)

var AppConfig configModel.Config

func ConfigSetup() {
	viper.SetConfigName("settings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error occured during env set. Err: %s", err)
	}

	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		return
	}
}
