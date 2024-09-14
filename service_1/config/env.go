package config

import (
	//viper

	"log"

	viper "github.com/spf13/viper"
)

var Db struct {
	Host     string `mapstructure:"DB_HOST"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	Port     string `mapstructure:"DB_PORT"`
}

var EnvValues struct {
	Port     string `mapstructure:"SVC1_PORT"`
	Svc2Port string `mapstructure:"SVC2_PORT"`
	Svc2Url  string `mapstructure:"SVC2_URL"`
	RedisAddress string `mapstructure:"REDIS_ADDRESS"`
}

func init() {
	getEnvValues()
	connectToDB()
	migrateTables()
}

func getEnvValues() {
	viper.AddConfigPath("config")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("error occured while reading env variables, error:", err)
	}

	err = viper.Unmarshal(&Db)
	if err != nil {
		log.Fatalln("error occured while writing env values onto variables, error:", err)
	}

	err = viper.Unmarshal(&EnvValues)
	if err != nil {
		log.Fatalln("error occured while writing env values onto variables, error:", err)
	}

}
