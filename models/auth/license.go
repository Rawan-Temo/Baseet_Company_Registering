package auth_models

import (
	"time"

	company_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/company"
	"gorm.io/gorm"
)

type License struct {
	gorm.Model
	CompanyId uint    `gorm:"type:integer;not null;index" json:"company_id"`
	Company   company_models.Company `gorm:"foreignKey:CompanyId" json:"company,omitempty"`
	StartDate time.Time `gorm:"type:date;not null" json:"start_date"`
	ExpirationDate   time.Time `gorm:"type:date;not null" json:"expiration_date"`
}

