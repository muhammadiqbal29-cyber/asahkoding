package user

import (
	"time"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// User merepresentasikan tabel users di MySQL
type User struct {
	ID           uint      `gorm:"primaryKey"`
	Username     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"` // Menyimpan password yang sudah dienkripsi
	Role         Role      `gorm:"type:varchar(10);default:'user'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
