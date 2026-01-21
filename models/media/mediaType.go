package media_models

import "github.com/Rawan-Temo/Baseet_Company_Registering.git/models"

type MediaType struct {
	models.NewGormModel
	Name        string `gorm:"type:varchar(100);unique;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
}
