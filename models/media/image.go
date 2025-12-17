package media_models

import (
	company_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/company"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	Src string `gorm:"uniqueIndex;type:varchar(100);not null" json:"src"`
	CompanyId uint    `gorm:"column:company_id;type:integer;index" json:"company_id"`
	Company   company_models.Company `gorm:"foreignKey:CompanyId" json:"company,omitempty"`
		MediaTypeId  int `gorm:"column:media_type_id;type:integer;index" json:"media_type_id"`
	MediaType	MediaType `gorm:"foreignKey:MediaTypeId" json:"media_type,omitempty"`
}

