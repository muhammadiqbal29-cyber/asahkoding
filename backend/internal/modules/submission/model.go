package submission

import (
	"time"
)

// Submission menyimpan riwayat kode yang dikirim peserta beserta nilainya
type Submission struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"not null" json:"user_id"`
	ProblemID     uint      `gorm:"not null" json:"problem_id"`
	Language      string    `gorm:"type:varchar(50);not null" json:"language"` // Contoh: "go"
	Code          string    `gorm:"type:text;not null" json:"code"`
	Status        string    `gorm:"type:varchar(50);not null" json:"status"`   // Accepted, Wrong Answer, Compile Error, TLE
	ExecutionTime int       `json:"execution_time_ms"` // Waktu eksekusi rata-rata
	CreatedAt     time.Time `json:"created_at"`
}
