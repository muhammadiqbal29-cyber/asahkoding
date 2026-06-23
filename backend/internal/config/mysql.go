package config

import (
	"log"
	"os"
	"time"

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

	// Menambahkan mekanisme Retry agar tangguh terhadap keterlambatan MySQL (Docker-in-Docker issue)
	var db *gorm.DB
	var err error
	for i := 1; i <= 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Menunggu MySQL siap (percobaan %d/10): %v\n", i, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Printf("FAILOVER: Gagal terhubung ke MySQL setelah beberapa kali percobaan: %v\n", err)
		log.Println("WARNING: API will continue running, but database features will be disabled.")
		return
	}

	DB = db
	log.Println("SUCCESS: Connected to MySQL database via GORM")
}
