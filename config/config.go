package config

import (
	"fmt"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-transactions-test/docs"
	"net/url"
)

type MainConfig struct {
	Port          int           `json:"port" mapstructure:"port"`
	GinMode       string        `json:"gin_mode" mapstructure:"gin_mode"`
	DBConfig      DBConfig      `json:"db_config" mapstructure:"db_config"`
	SwaggerConfig SwaggerConfig `json:"swagger_config" mapstructure:"swagger_config"`
}

type DBConfig struct {
	Host         string `json:"host" mapstructure:"host"`
	Port         int    `json:"port" mapstructure:"port"`
	Username     string `json:"username" mapstructure:"username"`
	Password     string `json:"password" mapstructure:"password"`
	DatabaseName string `json:"database_name" mapstructure:"database_name"`
}
type SwaggerConfig struct {
	Version     string `json:"version" mapstructure:"version"`
	Host        string `json:"host" mapstructure:"host"`
	BasePath    string `json:"base_path" mapstructure:"base_path"`
	Title       string `json:"title" mapstructure:"title"`
	Description string `json:"description" mapstructure:"description"`
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

func NewGinSwaggerConfig(config SwaggerConfig) *ginSwagger.Config {
	docs.SwaggerInfo.Title = config.Title
	docs.SwaggerInfo.Description = config.Description
	docs.SwaggerInfo.Version = config.Version
	docs.SwaggerInfo.BasePath = config.BasePath

	url, _ := url.Parse(config.Host)
	docs.SwaggerInfo.Host = url.Host

	return &ginSwagger.Config{
		URL: "doc.json",
	}
}
