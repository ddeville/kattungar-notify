package client

import (
	"github.com/spf13/viper"
	"log"
)

func init() {
	viper.AddConfigPath("$XDG_CONFIG_HOME/kattungar-notify")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("kattungar_notify")

	viper.ReadInConfig()
}

func GetApiKey() string {
	apiKey := viper.GetString("api_key")
	if apiKey == "" {
		log.Fatalln("Missing 'api_key' in config (or 'KATTUNGAR_NOTIFY_API_KEY' environment variable)")
	}
	return apiKey
}

func GetServerUrl() string {
	url := viper.GetString("server_url")
	if url != "" {
		return url
	} else {
		return "https://notify.home.kattungar.net"
	}
}
