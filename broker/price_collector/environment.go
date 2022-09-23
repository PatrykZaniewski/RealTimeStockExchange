package main

import (
	"log"

	"github.com/spf13/viper"
)

func EnvInit() {
	viper.SetConfigName("settings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error occured during env set. Err: %s", err)
	}
}
