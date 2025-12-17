package media_models

import "gorm.io/gorm"








type MediaType struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);unique;not null"`
	Description string `gorm:"type:text"`
}