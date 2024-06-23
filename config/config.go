package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type MainConfig struct {
	Port     string   `yaml:"port" mapStructure:"port"`
	DBConfig DBConfig `yaml:"db_config" mapStructure:"db_config"`
	GinMode  string   `yaml:"gin_mode" mapStructure:"gin_mode"`
}

type DBConfig struct {
	Host         string `yaml:"host" mapStructure:"host"`
	Port         int    `yaml:"port" mapStructure:"port"`
	Username     string `yaml:"username" mapStructure:"username"`
	Password     string `yaml:"password" mapStructure:"password"`
	DatabaseName string `yaml:"database_name" mapStructure:"database_name"`
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
