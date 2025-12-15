package auth_models

import (
	"errors"
	"strings"
	"time"

	company_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/company"
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
	UserName  string     `gorm:"uniqueIndex;column:username;type:varchar(100);not null" json:"username"`
	Password  string     `gorm:"type:varchar(100);not null" json:"-"`
	Email     string     `gorm:"type:varchar(100)" json:"email"`
	Role      Role       `gorm:"type:varchar(20);default:'user'" json:"role"`
	ExpiresAt *time.Time `gorm:"type:timestamp;"`
	CompanyId *uint      `gorm:"column:company_id;index" json:"company_id"`
	Company   *company_models.Company   `gorm:"foreignKey:CompanyId" json:"company,omitempty"`
	Active    bool       `gorm:"type:boolean;default:true" json:"active"`
}

// BeforeCreate handles validation and password hashing
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Basic required fields
	u.UserName = strings.TrimSpace(u.UserName)
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	if u.UserName == "" {
		return errors.New("username is required")
	}
	if strings.TrimSpace(u.Password) == "" {
		return errors.New("password is required")
	}

	// Optional: company validation for non-admins
	if u.Role != RoleAdmin && u.CompanyId == nil {
		return errors.New("company_id is required for non-admin users")
	}

	// Default expiry

	if u.Role != RoleAdmin && u.ExpiresAt == nil {
		expires := time.Now().Add(24 * time.Hour)
		u.ExpiresAt = &expires
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

