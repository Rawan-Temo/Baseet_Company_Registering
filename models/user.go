package models

import (
	"errors"
	"fmt"
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
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserName  string         `gorm:"type:varchar(100);not null" json:"username"`
	Password  string         `gorm:"type:varchar(100);not null" json:"password"`
	Email     string         `gorm:"uniqueIndex;type:varchar(100);not null" json:"email"`
	Role      Role           `gorm:"type:varchar(20);default:'user'" json:"role"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

// BeforeCreate handles validation and password hashing
func (u *User) BeforeCreate() error {
	// 1️⃣ Basic validation
	if strings.TrimSpace(u.UserName) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("email is required")
	}
	fmt.Println(u)
	if strings.TrimSpace(u.Password) == "" {
		return errors.New("password is required")
	}

	// 2️⃣ Validate role (enum)
	validRoles := map[string]bool{
		"user":  true,
		"admin": true,
		"staff": true,
	}
	if _, ok := validRoles[string(u.Role)]; !ok {
		u.Role = RoleUser // default role
	}

	// 3️⃣ Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)

	return nil
}
