package utils

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/anurag-rajawat/rest-api/pkg/types"
)

var db *gorm.DB

// Init initializes env variables and database
func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		log.Fatal("DB host is not specified")
	}

	dbUser, ok := os.LookupEnv("DB_USER")
	if !ok {
		log.Fatal("DB User is not specified")
	}

	dbPasswd, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		log.Fatal("DB Password is not specified")
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		log.Fatal("DB Name is not specified")
	}

	dbPort, ok := os.LookupEnv("DB_PORT")
	if !ok {
		log.Warn("DB Port is not specified. Defaults to 5432")
		dbPort = "5432"
	}
	db, err = ConnectToDb(dbHost, dbUser, dbPasswd, dbName, dbPort)

	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Successfully connected to DB")

	err = db.AutoMigrate(&types.User{})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func GetDb() *gorm.DB {
	return db
}
