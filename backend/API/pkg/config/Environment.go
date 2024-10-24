package config

import (
	"github.com/spf13/viper"
)

func GetEnv(key string) (string, error) {
	err := viper.ReadInConfig()

	if err != nil {
		return "", err
	}

	value := viper.GetString(key)

	return value, nil
}
