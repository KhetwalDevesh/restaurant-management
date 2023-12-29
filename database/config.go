package config

import (
	"fmt"
	"github.com/KhetwalDevesh/restaurant-management/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbname   = os.Getenv("DB_NAME")
)

var (
	db  *gorm.DB
	err error
)

func ConfigDB() {
	fmt.Println("inside ConfigDb")
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		panic("Failed to connect to the database")
	}
	fmt.Println("Successfully connected to db", db)

	err = db.AutoMigrate(&models.User{}, &models.OrderItem{}, &models.Order{}, &models.Food{}, &models.Menu{}, &models.Table{}, &models.Invoice{})
	if err != nil {
		panic("Failed to auto-migrate the model")
	}
}

// GetDB returns an initialized instance of gorm.DB
func GetDB() *gorm.DB {
	if db == nil {
		ConfigDB()
	}
	return db
}
