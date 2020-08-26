package main

import (
	"errors"

	"github.com/spf13/viper"
)

type config struct {
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     uint
	DbName     string
}

func readConf() (*config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("$HOME/.fattingo")
	viper.AddConfigPath(".")

	// defaults
	viper.SetDefault("db.user", "root")
	viper.SetDefault("db.password", "s3cr3t")
	viper.SetDefault("db.host", "127.0.0.1")
	viper.SetDefault("db.port", "3306")
	viper.SetDefault("db.name", "fattingo")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("Configuration file not found in ./config.toml or $HOME/.fattingo/config.toml")
		} else {
			return nil, err
		}
	}

	cfg := &config{
		DbUser:     viper.GetString("db.user"),
		DbPassword: viper.GetString("db.password"),
		DbHost:     viper.GetString("db.host"),
		DbPort:     viper.GetUint("db.port"),
		DbName:     viper.GetString("db.name"),
	}

	return cfg, nil
}
