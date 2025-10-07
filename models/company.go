package models




import (
	"gorm.io/gorm"
)


type Company struct{
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Address     string `gorm:"type:varchar(200);not null" json:"address"`
	Email       string `gorm:"uniqueIndex;type:varchar(100);not null" json:"email"`
	PhoneNumber string `gorm:"type:varchar(15);not null" json:"phone_number"`
}