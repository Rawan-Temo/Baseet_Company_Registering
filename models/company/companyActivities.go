package company_models

import (
	"errors"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/models"
	"gorm.io/gorm"
)




type CompanyActivity struct {
	models.NewGormModel


	CompanyID uint `gorm:"primaryKey" json:"company_id"`
	// معرف الشركة

	TradingActivityID uint `gorm:"primaryKey" json:"trading_activity_id"`
	// معرف النشاط التجاري

	Company Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// الشركة

	TradingActivity TradingActivity `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	// النشاط التجاري

	Image string `gorm:"type:varchar(255)" json:"image"`
	// صورة مرتبطة بالنشاط داخل الشركة
}

func (ca *CompanyActivity) BeforeSave(tx *gorm.DB) error {
	if ca.CompanyID == 0 {
		return errors.New("company_id is required")
	}
	if ca.TradingActivityID == 0 {
		return errors.New("trading_activity_id is required")
	}
	return nil
}