package models

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Role enum type
type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// User model
type User struct {
	gorm.Model
	UserName  string         `gorm:"column:username;type:varchar(100);not null" json:"username"`
	Password  string         `gorm:"type:varchar(100);not null" json:"password"`
	Email     string         `gorm:"uniqueIndex;type:varchar(100);not null" json:"email"`
	Role      Role           `gorm:"type:varchar(20);default:'user'" json:"role"`
	ExpiresAt time.Time      `gorm:"type:timestamp;default:null" json:"expires_at"`
	CompanyId *uint           `gorm:"column:company_id;type:integer" json:"company_id"`
	Company  *Company        `gorm:"foreignKey:CompanyId" json:"company"`
}

// BeforeCreate handles validation and password hashing
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Basic required fields
	if strings.TrimSpace(u.UserName) == "" {
		return errors.New("username is required")
	}
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("email is required")
	}
	if strings.TrimSpace(u.Password) == "" {
		return errors.New("password is required")
	}

	// Optional: company validation for non-admins
	if u.Role != RoleAdmin && u.CompanyId == nil {
		return errors.New("company_id is required for non-admin users")
	}

	// Default expiry
	if u.ExpiresAt.IsZero() {
		u.ExpiresAt = time.Now().Add(30 * 24 * time.Hour)
	}

	// Validate role
	if u.Role != RoleAdmin && u.Role != RoleUser {
		u.Role = RoleUser
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)

	return nil
}
