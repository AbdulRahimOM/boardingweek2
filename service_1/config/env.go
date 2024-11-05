package config

import (
	//viper

	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Db struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

var EnvValues struct {
	Svc1Port           string
	Svc2Port           string
	Svc2Url            string
	RedisAddress       string
	CacheDurationInSec int
}

func init() {
	getEnvValues()
	connectToDB()
	migrateTables()
}

func getEnvValues() {
	switch os.Getenv("ENVIRONMENT") {
	case "DOCKER": //nothing to do
	case "KUBERNETES": //nothing to do
	default:
		if err := godotenv.Load("config/.env"); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	Db.Host = os.Getenv("DB_HOST")
	Db.User = os.Getenv("DB_USER")
	Db.Password = os.Getenv("DB_PASSWORD")
	Db.Name = os.Getenv("DB_NAME")
	Db.Port = os.Getenv("DB_PORT")

	EnvValues.Svc1Port = os.Getenv("SVC1_PORT")
	EnvValues.Svc2Port = os.Getenv("SVC2_PORT")
	EnvValues.Svc2Url = os.Getenv("SVC2_URL")
	EnvValues.RedisAddress = os.Getenv("REDIS_ADDRESS")

	cacheDurationInSecStr := os.Getenv("CACHE_DURATION_IN_SEC")
	if cacheDurationInSecStr == "" {
		EnvValues.CacheDurationInSec = 15
	} else {
		var err error
		EnvValues.CacheDurationInSec, err = strconv.Atoi(cacheDurationInSecStr)
		if err != nil {
			log.Fatalf("Error converting CACHE_DURATION_IN_SEC to int: %v", err)
		}
	}

	fmt.Println("Db: ", Db)
	fmt.Println("EnvValues: ", EnvValues)

}
