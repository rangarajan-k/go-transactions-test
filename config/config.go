package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type MainConfig struct {
	Port     int      `json:"port" mapstructure:"port"`
	GinMode  string   `json:"gin_mode" mapstructure:"gin_mode"`
	DBConfig DBConfig `json:"db_config" mapstructure:"db_config"`
}

type DBConfig struct {
	Host         string `json:"host" mapstructure:"host"`
	Port         int    `json:"port" mapstructure:"port"`
	Username     string `json:"username" mapstructure:"username"`
	Password     string `json:"password" mapstructure:"password"`
	DatabaseName string `json:"database_name" mapstructure:"database_name"`
}

func LoadMainConfig(filepath string) *MainConfig {
	var mainConfig MainConfig
	viper.SetConfigFile(filepath)
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	err = viper.Unmarshal(&mainConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshalling config file: %s", err))
	}
	return &mainConfig
}
