package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var EnvValues struct {
	Svc2Port string `mapstructure:"SVC2_PORT"`
}

var Postgresdb struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

func init() {
	loadConfig()
}

func loadConfig() {
	switch os.Getenv("ENVIRONMENT") {
	case "DOCKER":
		// nothing to do
	case "KUBERNETES":
		// nothing to do
	default:
		if err := godotenv.Load("config/.env"); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	Postgresdb.Host = os.Getenv("DB_HOST")
	Postgresdb.User = os.Getenv("DB_USER")
	Postgresdb.Password = os.Getenv("DB_PASSWORD")
	Postgresdb.Name = os.Getenv("DB_NAME")
	Postgresdb.Port = os.Getenv("DB_PORT")

	EnvValues.Svc2Port = os.Getenv("SVC2_PORT")

}
