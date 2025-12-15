package general_models

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Office struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;type:varchar(100);not null" json:"name"`
}

func (o *Office) BeforeCreate(tx *gorm.DB) error {
	if strings.TrimSpace(o.Name) == "" {
		return errors.New("name is required")
	}
	return nil
}