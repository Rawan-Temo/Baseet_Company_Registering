package auth_models

import (
	"errors"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/models"
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
	models.NewGormModel
	UserName  string     `gorm:"type:varchar(50);uniqueIndex:idx_user_name_active,where:deleted_at IS NULL;column:username;type:varchar(100);not null" json:"username"`
	FullName  string     `gorm:"type:varchar(150);uniqueIndex:idx_user_name_active,where:deleted_at IS NULL;not null" json:"full_name"`
	Password  string     `gorm:"type:varchar(100);not null" json:"-"`
	Email     string     `gorm:"type:varchar(100)" json:"email"`
	Role      Role       `gorm:"type:varchar(20);default:'user'" json:"role"`
	CompanyId *uint      `gorm:"column:company_id;index" json:"company_id"`
	Company   *company_models.Company   `gorm:"foreignKey:CompanyId" json:"company,omitempty"`
	Active    bool       `gorm:"type:boolean;default:true" json:"active"`
}

// BeforeSave handles validation and password hashing
func (u *User) BeforeSave(tx *gorm.DB) error {	
	// Optional: company validation for non-admins
	if u.Role != RoleAdmin && u.CompanyId == nil {
		return errors.New("company_id is required for non-admin users")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)

	return nil
}

