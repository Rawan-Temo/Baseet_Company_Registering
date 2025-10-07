package models

import (
	"gorm.io/gorm"
)


type Company struct{
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Address     string `gorm:"type:varchar(200)" json:"address"`
	Description string `gorm:"type:varchar(500)" json:"description"`
	Email       string `gorm:"uniqueIndex;type:varchar(100)" json:"email"`
	PhoneNumber string `gorm:"type:varchar(15)" json:"phone_number"`
	TypeID      uint `gorm:"column:type_id;index" json:"type_id"`
	Type        Type `gorm:"foreignKey:TypeID" json:"type,omitempty"`
	TradingActivities []TradingActivity `gorm:"many2many:company_trading_activities;" json:"trading_activities"`
}