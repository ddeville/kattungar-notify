package client

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	APIKey    string `mapstructure:"api_key"`
	ServerURL string `mapstructure:"server_url"`
}

var C config

func init() {
	viper.AddConfigPath("/etc/kattungar-notify")
	viper.AddConfigPath("$XDG_CONFIG_HOME/kattungar-notify")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("KATTUNGAR_NOTIFY")
	viper.AutomaticEnv()
	// We need to bind environment variables for Unmarshal to use them...
	viper.BindEnv("api_key")
	viper.BindEnv("server_url")

	viper.SetDefault("server_url", "https://notify.home.kattungar.net")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("fatal error config file: %v", err)
		}
	}

	if err := viper.Unmarshal(&C); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	if C.APIKey == "" {
		log.Fatalln("Missing 'api_key' in config (or 'KATTUNGAR_NOTIFY_API_KEY' environment variable)")
	}
	if C.ServerURL == "" {
		log.Fatalln("Missing 'server_url' in config (or 'KATTUNGAR_NOTIFY_SERVER_URL' environment variable)")
	}
}
