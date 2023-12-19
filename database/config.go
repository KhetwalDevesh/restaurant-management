package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "restaurant"
)

var (
	db  *gorm.DB
	err error
)

func ConfigDB() {
	fmt.Println("inside ConfigDb")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		panic("Failed to connect to the database")
	}
	fmt.Println("Successfully connected to db", db)

	//err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic("Failed to auto-migrate the Todo model")
	}
}

// GetDB returns an initialized instance of gorm.DB
func GetDB() *gorm.DB {
	if db == nil {
		ConfigDB()
	}
	return db
}
