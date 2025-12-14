package media_models

import (
	company_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/company"
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	Src string `gorm:"uniqueIndex;type:varchar(100);not null" json:"src"`
	CompanyId uint   `gorm:"column:company_id;type:integer;index" json:"company_id"`
	Company   company_models.Company `gorm:"foreignKey:CompanyId" json:"company,omitempty"`
}

