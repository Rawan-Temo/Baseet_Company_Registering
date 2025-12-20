package general_models

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Office struct {
	gorm.Model
	Name string `gorm:"uniqueIndex:idx_user_name_active,where:deleted_at IS NULL;type:varchar(100);not null" json:"name"`
}

func (o *Office) BeforeSave(tx *gorm.DB) error {
	if strings.TrimSpace(o.Name) == "" {
		return errors.New("name is required")
	}
	return nil
}