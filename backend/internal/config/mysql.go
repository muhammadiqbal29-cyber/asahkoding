package config

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMySQL() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Println("WARNING: MYSQL_DSN is empty in .env. Skipping MySQL connection.")
		return
	}

	// Failover mechanism: try to connect, but don't panic if it fails
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("FAILOVER: Failed to connect to MySQL: %v\n", err)
		log.Println("WARNING: API will continue running, but database features will be disabled.")
		return
	}

	DB = db
	log.Println("SUCCESS: Connected to MySQL database via GORM")
}
