package fattingo

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	serverAddr string
	dbType     string // mysql or sqlite3
	dbPath     string // Path to sqlite file
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     uint
	dbName     string
}

func ReadConf() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("$HOME/.fattingo")
	viper.AddConfigPath(".")

	// defaults
	viper.SetDefault("server.address", ":5000")
	viper.SetDefault("db.path", "./example-db.sqlite")
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

	return &Config{
		serverAddr: viper.GetString("server.address"),
		dbType:     viper.GetString("db.type"),
		dbPath:     viper.GetString("db.path"),
		dbUser:     viper.GetString("db.user"),
		dbPassword: viper.GetString("db.password"),
		dbHost:     viper.GetString("db.host"),
		dbPort:     viper.GetUint("db.port"),
		dbName:     viper.GetString("db.name"),
	}, nil
}
