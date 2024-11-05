package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	connectToDB()
}

func connectToDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", Postgresdb.Host, Postgresdb.User, Postgresdb.Password, Postgresdb.Name, Postgresdb.Port)
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Couldn't connect to DB. Error:", err)
	}

}
