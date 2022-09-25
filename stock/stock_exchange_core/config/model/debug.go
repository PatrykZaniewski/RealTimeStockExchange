package model

type DebugConfig struct {
	Rest DebugRestConfig `mapstructure:"rest"`
}

type DebugRestConfig struct {
	OrderStatusCollectorPort int `mapstructure:"orderStatusCollectorPort"`
}
