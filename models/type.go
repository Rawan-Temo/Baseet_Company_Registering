package models

import "gorm.io/gorm"

type Type struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;type:varchar(100);not null" json:"name"`

}

