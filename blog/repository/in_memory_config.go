package repository

import (
	"github.com/rshlin/go-blog-api-assesment/blog/model"
	"github.com/spf13/viper"
	"log"
)

type InMemoryConfig struct {
	Posts []model.Post `json:"posts"`
}

func LoadInMemoryConfig(path string) *InMemoryConfig {
	cfg := InMemoryConfig{}
	if path != "" {
		viper.SetConfigFile(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("Error unmarshalling config", err)
	}

	return &cfg
}
