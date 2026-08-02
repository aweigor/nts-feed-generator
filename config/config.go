package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Nts    NtsAPIProperties
	Server ServerProperties
	Auth   AuthConfig
}

type AuthConfig struct {
	Secret string `mapsructure: "secret"`
}

type NtsAPIProperties struct {
	APIV2Url string `mapsructure: "api_v2_url"`
}

type ServerProperties struct {
	PublicURL string `mapsructure: "public_url"`
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return nil, err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// read .env
	viper.Set("auth.secret", os.Getenv("API_SECRET"))

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
