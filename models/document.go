package models

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	Src string `gorm:"uniqueIndex;type:varchar(100);not null" json:"src"`
	CompanyId uint   `gorm:"column:company_id;type:integer;index" json:"company_id"`
	Company   Company `gorm:"foreignKey:CompanyId" json:"company,omitempty"`
}

