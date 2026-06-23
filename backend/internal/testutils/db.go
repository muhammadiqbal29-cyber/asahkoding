package testutils

import (
	"log"

	"github.com/glebarez/sqlite"
	"github.com/muhammad-iqbal/leetcode-backend/internal/config"
	"gorm.io/gorm"
)

// SetupMockDB meracik database virtual di RAM (In-Memory) agar test cepat dan bersih
func SetupMockDB() {
	// Menggunakan driver glebarez/sqlite murni Golang (tanpa butuh C Compiler/CGO)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open mock database: %v", err)
	}

	config.DB = db
}
