package config

import (
	"boarding-week2/service_1/domain"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func connectToDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", Db.Host, Db.User, Db.Password, Db.Name, Db.Port)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Couldn't connect to DB, error: ", err)
	}

	// Create a connection string without specifying a database to connect to the PostgreSQL server
	serverDSN := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable",
		Db.Host,
		Db.User,
		Db.Password,
		Db.Port,
	)

	// Connect to the PostgreSQL server
	serverDB, err := gorm.Open(postgres.Open(serverDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Couldn't connect to PostgreSQL server. Error:", err)
	}

	// Checking if the database exists
	var exists bool
	checkDBQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", Db.Name)
	if err := serverDB.Raw(checkDBQuery).Scan(&exists).Error; err != nil {
		log.Fatal("Couldn't check if database exists. Error:", err)
	}

	// If the database does not exist, creating it
	if !exists {
		createDBQuery := fmt.Sprintf("CREATE DATABASE %s", Db.Name)
		if err := serverDB.Exec(createDBQuery).Error; err != nil {
			log.Fatal("Couldn't create database. Error:", err)
		}
		log.Println("Database created successfully:", Db.Name)
	} else {
		log.Println("Database", Db.Name, "already exists. Okay")
	}
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Couldn't connect to DB. Error:", err)
	}

}

func migrateTables() {
	if err := DB.AutoMigrate(&domain.User{}); err != nil {
		log.Fatal("Couldn't migrate tables. Error:", err)
	}
}
