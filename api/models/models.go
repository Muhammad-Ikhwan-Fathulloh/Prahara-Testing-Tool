package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex" json:"username"`
	Password  string         `json:"-"`
	Role      string         `gorm:"default:tester" json:"role"` // admin, tester
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type TestScript struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Content     string         `json:"content"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Runs        []TestRun      `gorm:"foreignKey:TestScriptID" json:"runs"`
}

type TestRun struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TestScriptID uint           `json:"test_script_id"`
	Name         string         `json:"name"`
	TargetURL    string         `json:"target_url"`
	Method       string         `json:"method"`
	VUs          int            `json:"vus"`
	Duration     string         `json:"duration"`
	Status       string         `json:"status"` // pending, running, completed, failed
	InfluxBucket string         `json:"influx_bucket"`
	StartedAt    time.Time      `json:"started_at"`
	FinishedAt   *time.Time     `json:"finished_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type TestingURL struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	URL       string         `json:"url"`
	Category  string         `json:"category"` // Frontend, Backend, API, etc.
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
