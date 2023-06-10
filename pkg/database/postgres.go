package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/anurag-rajawat/rest-api/pkg/types"
)

// ConnectToDb creates a connection to the postgres
func ConnectToDb(dbHost, dbUser, dbPasswd, dbName, dbPort string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		dbHost, dbUser, dbPasswd, dbName, dbPort)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Initialize initializes database connection
func Initialize(db *gorm.DB) error {
	err := db.AutoMigrate(&types.User{})
	if err != nil {
		return err
	}
	return nil
}
