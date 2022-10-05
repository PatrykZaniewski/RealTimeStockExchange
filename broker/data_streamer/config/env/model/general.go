package model

type GeneralConfig struct {
	Rest GeneralRestConfig `mapstructure:"rest"`
}

type GeneralRestConfig struct {
	Port string `mapstructure:"port"`
}
