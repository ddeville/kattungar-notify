package main

import (
	"fmt"
	"os"
	"reflect"
)

type Config struct {
	StorePath          string `env:"KATTUNGAR_STORE_PATH"`
	ApnsKeyId          string `env:"KATTUNGAR_APNS_KEY_ID"`
	ApnsKeyPath        string `env:"KATTUNGAR_APNS_KEY_PATH"`
	GoogleCredsPath    string `env:"KATTUNGAR_GOOGLE_CREDS_PATH"`
	GoogleRefreshToken string `env:"KATTUNGAR_GOOGLE_REFRESH_TOKEN"`
}

func LoadConfig() (*Config, error) {
	cfg := Config{}

	v := reflect.ValueOf(&cfg).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		env := field.Tag.Get("env")
		val, hasVal := os.LookupEnv(env)
		if !hasVal {
			return nil, fmt.Errorf("missing environment variable %s", env)
		}

		v.Field(i).SetString(val)
	}

	return &cfg, nil
}
