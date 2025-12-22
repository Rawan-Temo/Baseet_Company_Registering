package general_models

import (
	"gorm.io/gorm"
)

type Office struct {
	gorm.Model
	Name string `gorm:"uniqueIndex:idx_user_name_active,where:deleted_at IS NULL;type:varchar(100);not null" json:"name"`
}
