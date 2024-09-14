package config

import (
	"log"

	"github.com/spf13/viper"
)

var EnvValues struct {
	Svc2Port string `mapstructure:"SVC2_PORT"`
}

var Postgresdb struct {
	DbHost     string `mapstructure:"DB_HOST"`
	DbUser     string `mapstructure:"DB_USER"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbName     string `mapstructure:"DB_NAME"`
	DbPort     string `mapstructure:"DB_PORT"`
}

func init() {
	loadConfig()
}

func loadConfig() {
	viper.AddConfigPath("config")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("error occured while reading env variables, error:", err)
	}

	err = viper.Unmarshal(&EnvValues)
	if err != nil {
		log.Fatalln("error occured while writing env values onto variables, error:", err)
	}

	err = viper.Unmarshal(&Postgresdb)
	if err != nil {
		log.Fatalln("error occured while writing env values onto variables, error:", err)
	}
}
