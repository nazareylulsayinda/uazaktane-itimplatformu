package models

import "gorm.io/gorm"

// User represents a user in the system (teacher or student)
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"uniqueIndex" json:"email"`
	Password string `json:"-"`
	Role     string `gorm:"default:'student'" json:"role"` // teacher or student
}
