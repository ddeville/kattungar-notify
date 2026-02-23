package main

import (
	"log"

	"github.com/spf13/viper"
)

const defaultPort = 5000
const defaultAppleTeamID = "Q8B696Y8U4"
const defaultAppleAppID = "com.ddeville.kattungar-notify"

type config struct {
	// The port that the application should be served on.
	Port int `mapstructure:"port"`
	// The team ID for the Apple application that will be receing notifications.
	AppleTeamID string `mapstructure:"apple_team_id"`
	// The app ID for the Apple application that will be receing notifications.
	AppleAppID string `mapstructure:"apple_app_id"`
	// The path that should be used to write the sqlite3 store to disk.
	StorePath string `mapstructure:"store_path"`
	// The path that should be used to retrieve the JSON file containing the API keys.
	ServerAPIKeysPath string `mapstructure:"server_api_keys_path"`
	// The ID of the APNS key being used to send notifications.
	ApnsKeyID string `mapstructure:"apns_key_id"`
	// The path to the APNS key being used to send notifications.
	ApnsKeyPath string `mapstructure:"apns_key_path"`
	// The path to the config file containing credentials information for the Google Calendar account.
	GoogleCredsPath string `mapstructure:"google_creds_path"`
	// The OAuth2 refresh token that should be used with Google Calendar.
	GoogleRefreshToken string `mapstructure:"google_refresh_token"`
	// The ID of the Google Calendar instance that should be watched.
	GoogleCalendarID string `mapstructure:"google_calendar_id"`
}

var C config

func init() {
	viper.AddConfigPath("/etc/kattungar-notify-server")
	viper.AddConfigPath("$XDG_CONFIG_HOME/kattungar-notify-server")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("KATTUNGAR")
	viper.AutomaticEnv()
	// We need to bind environment variables for Unmarshal to use them...
	viper.BindEnv("port")
	viper.BindEnv("apple_team_id")
	viper.BindEnv("apple_app_id")
	viper.BindEnv("store_path")
	viper.BindEnv("server_api_keys_path")
	viper.BindEnv("apns_key_id")
	viper.BindEnv("apns_key_path")
	viper.BindEnv("google_creds_path")
	viper.BindEnv("google_refresh_token")
	viper.BindEnv("google_calendar_id")

	viper.SetDefault("port", defaultPort)
	viper.SetDefault("apple_team_id", defaultAppleTeamID)
	viper.SetDefault("apple_app_id", defaultAppleAppID)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("fatal error config file: %v", err)
		}
	}

	if err := viper.Unmarshal(&C); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	if C.StorePath == "" {
		log.Fatalln("Missing 'store_path' in config (or 'KATTUNGAR_STORE_PATH' environment variable)")
	}
	if C.ServerAPIKeysPath == "" {
		log.Fatalln("Missing 'server_api_keys_path' in config (or 'KATTUNGAR_SERVER_API_KEYS_PATH' environment variable)")
	}
	if C.ApnsKeyID == "" {
		log.Fatalln("Missing 'apns_key_id' in config (or 'KATTUNGAR_APNS_KEY_ID' environment variable)")
	}
	if C.ApnsKeyPath == "" {
		log.Fatalln("Missing 'apns_key_path' in config (or 'KATTUNGAR_APNS_KEY_PATH' environment variable)")
	}
	if C.GoogleCredsPath == "" {
		log.Fatalln("Missing 'google_creds_path' in config (or 'KATTUNGAR_GOOGLE_CREDS_PATH' environment variable)")
	}
	if C.GoogleRefreshToken == "" {
		log.Fatalln("Missing 'google_refresh_token' in config (or 'KATTUNGAR_GOOGLE_REFRESH_TOKEN' environment variable)")
	}
	if C.GoogleCalendarID == "" {
		log.Fatalln("Missing 'google_calendar_id' in config (or 'KATTUNGAR_GOOGLE_CALENDAR_ID' environment variable)")
	}
}
