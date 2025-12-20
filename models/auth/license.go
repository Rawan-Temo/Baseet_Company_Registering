package auth_models

import (
	"errors"
	"time"

	company_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/company"
	"gorm.io/gorm"
)

type License struct {
	gorm.Model

	CompanyId uint `gorm:"not null;index" json:"company_id"`
	// الشركة المالكة للترخيص

	Company company_models.Company `gorm:"foreignKey:CompanyId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"company,omitempty"`
	// الشركة (قراءة فقط)

	StartDate time.Time `gorm:"type:date;not null" json:"start_date"`
	// تاريخ بداية الترخيص

	ExpirationDate time.Time `gorm:"type:date;not null" json:"expiration_date"`
	// تاريخ انتهاء الترخيص

	Image *string `gorm:"type:varchar(255)" json:"image"`
	// صورة الترخيص
}

func (l *License) BeforeSave(tx *gorm.DB) (err error) {

	// ---------- Required ----------
	if l.CompanyId == 0 {
		return errors.New("company_id is required")
	}

	if l.StartDate.IsZero() {
		return errors.New("start_date is required")
	}

	if l.ExpirationDate.IsZero() {
		return errors.New("expiration_date is required")
	}

	// ---------- Date logic ----------
	if !l.ExpirationDate.After(l.StartDate) {
		return errors.New("expiration_date must be after start_date")
	}

	// ---------- FK existence ----------
	if err := tx.First(&company_models.Company{}, l.CompanyId).Error; err != nil {
		return errors.New("invalid company_id")
	}

	// ---------- Safety: block nested creation ----------
	if l.Company.ID != 0 {
		return errors.New("nested company creation is not allowed")
	}

	return nil
}