package main

import (
	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath("$XDG_CONFIG_HOME/kattungar-notify-admin")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("kattungar_notify")

	viper.ReadInConfig()
}

func getApiKey() string {
	return viper.GetString("api_key")
}
