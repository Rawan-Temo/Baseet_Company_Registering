package general_models

import (
	"errors"

	"gorm.io/gorm"
)

type CompanyType struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;type:varchar(100);not null" json:"name"`
}

func (ct *CompanyType) BeforeCreate(tx *gorm.DB) error {
	// check if name exists

	if ct.Name == "" {
		return errors.New("name is required")
	}
	ct.ID = 0
	return nil
}
