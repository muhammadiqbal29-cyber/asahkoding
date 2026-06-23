package problem

import (
	"time"
)

type Problem struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Description string     `gorm:"type:text;not null" json:"description"`
	Difficulty  string     `gorm:"type:varchar(50);not null" json:"difficulty"`
	TimeLimitMs int        `gorm:"default:2000" json:"time_limit_ms"`
	MemoryLimit int        `gorm:"default:128" json:"memory_limit"`
	TestCases   []TestCase `gorm:"foreignKey:ProblemID" json:"test_cases,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TestCase struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProblemID      uint      `json:"problem_id"`
	Input          string    `gorm:"type:text;not null" json:"input"`
	ExpectedOutput string    `gorm:"type:text;not null" json:"expected_output"`
	IsHidden       bool      `gorm:"default:false" json:"is_hidden"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
