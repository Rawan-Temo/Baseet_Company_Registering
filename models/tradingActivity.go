package models

import "gorm.io/gorm"

type TradingActivity struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;type:varchar(100);not null" json:"name"`
	// add many 2 many relation with company
	Companies []Company `gorm:"many2many:company_trading_activities;" json:"companies"`
}