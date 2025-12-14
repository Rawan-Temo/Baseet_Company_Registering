package general_models

import (
	"errors"

	"gorm.io/gorm"
)

type Office struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;type:varchar(100);not null" json:"name"`
}

func (o *Office) BeforeCreate(tx *gorm.DB) error {
	// check if name exists
	if o.Name == "" {
		return errors.New("name is required")
	}
	o.ID = 0
	return nil
}
