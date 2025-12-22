package general_models

import "github.com/Rawan-Temo/Baseet_Company_Registering.git/models"

type Office struct {
	models.NewGormModel
	Name string `gorm:"uniqueIndex:idx_user_name_active,where:deleted_at IS NULL;type:varchar(100);not null" json:"name"`
}
