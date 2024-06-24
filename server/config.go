package server

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	ResponseTimeoutMs                     int
	ValidationExcludeRequestBody          bool
	ValidationExcludeRequestQueryParams   bool
	ValidationExcludeResponseBody         bool
	ValidationExcludeReadOnlyValidations  bool
	ValidationExcludeWriteOnlyValidations bool
	ValidationIncludeResponseStatus       bool
	AuthStore                             *AuthStoreConfig
	Authenticator                         *AuthConfig
}

func LoadConfig(path string) *Config {
	config := Config{}

	viper.SetDefault("ResponseTimeoutMs", 100)
	viper.SetDefault("ValidationExcludeRequestBody", false)
	viper.SetDefault("ValidationExcludeRequestQueryParams", false)
	viper.SetDefault("ValidationExcludeResponseBody", false)
	viper.SetDefault("ValidationExcludeReadOnlyValidations", false)
	viper.SetDefault("ValidationExcludeWriteOnlyValidations", false)
	viper.SetDefault("ValidationIncludeResponseStatus", false)

	viper.AutomaticEnv()
	if path != "" {
		viper.SetConfigFile(path)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("app")
		viper.SetConfigType("yaml")
	}

	envKeysMap := &map[string]interface{}{}
	if err := mapstructure.Decode(config, &envKeysMap); err != nil {
		log.Fatalf("Error decoding config: %v", err)
	}
	for k := range *envKeysMap {
		if bindErr := viper.BindEnv(k); bindErr != nil {
			log.Fatalf("Error decoding config: %v", bindErr)
		}
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error reading env file", err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Println("Error unmarshalling config", err)
	}

	return &config
}
