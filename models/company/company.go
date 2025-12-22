package company_models

import (
	"errors"
	"strings"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/models"
	general_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/general"
	"gorm.io/gorm"
)

// TODO 	add vaildation tags and dont casecade create for related models
// Company represents a registered company with its details and relationships.



type Company struct {
	models.NewGormModel

	Name            string `gorm:"uniqueIndex:idx_user_name_active,where:deleted_at IS NULL;type:varchar(100);not null" json:"name" validate:"required"`
	// الاسم الرسمي للشركة

	ForeignBranchName string `gorm:"uniqueIndex:idx_user_name_active,where:deleted_at IS NULL;type:varchar(200);not null" json:"foreign_branch_name"`
	// اسم الفرع الأجنبي إن وجد

	ForeignRegistrationNumber string `gorm:"uniqueIndex:idx_user_name_active,where:deleted_at IS NULL;type:varchar(50);not null" json:"foreign_registration_number"`
	// رقم التسجيل

	TradeNames      string `gorm:"type:varchar(200)" json:"trade_names"`
	// الاسم (الأسماء) التجارية إن وجد

	AuthorityName   string `gorm:"type:varchar(100);not null" json:"authority_name" validate:"required"`
	// اسم الممنوح من هيئة عامة أخرى


	AuthorityNumber string `gorm:"type:varchar(100)" json:"authority_number"`
	// رقم التعريف الممنوح من هيئة عامة أخرى

	LocalAddress         string `gorm:"type:varchar(200);not null" json:"local_address" validate:"required"`
	// عنوان العمل الرئيسي

	ForeignAddress         string `gorm:"type:varchar(200);not null" json:"foreign_address" `
	// عنوان العمل الاوروبي


	Description     string `gorm:"type:varchar(500)" json:"description"`
	// وصف النشاط

	Email           string `gorm:"uniqueIndex:idx_user_name_active,where:deleted_at IS NULL;type:varchar(100)" json:"email" validate:"omitempty,email"`
	// البريد الإلكتروني

	PhoneNumber     string `gorm:"uniqueIndex:idx_user_name_active,where:deleted_at IS NULL;type:varchar(20)" json:"phone_number"`
	// رقم الهاتف

	CompanyTypeID   uint `gorm:"column:type_id;not null;index" json:"type_id" validate:"required"`
	// نوع الشركة (FK فقط – لا إنشاء تلقائي)

	CompanyType     general_models.CompanyType `gorm:"foreignKey:CompanyTypeID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"type,omitempty"`
	// نوع الشركة (للقراءة فقط)

	OfficeId        uint `gorm:"not null;index" json:"office_id" validate:"required"`
	// المكتب / الجهة المسجلة

	Office          general_models.Office `gorm:"foreignKey:OfficeId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"office,omitempty"`
	// المكتب (قراءة فقط)

	CompanyActivity []CompanyActivity `gorm:"foreignKey:CompanyID" json:"activities"`
	// الأنشطة التجارية (ربط فقط – لا إنشاء)

	People          []People `gorm:"foreignKey:CompanyID" json:"people"`
	// الشركاء، الممثلون، أعضاء مجلس الإدارة

	IsLicensed      bool   `gorm:"default:false" json:"is_licensed"`
	// هل النشاط مرخص

	Duration        string `gorm:"type:varchar(100)" json:"duration"`
	// مدة الشركة


	// ==============================
	// Default ceo for the company not in the poeple entity but embedded here
	// ==============================:
	CEOName     string `gorm:"type:varchar(100)" json:"ceo_name"`
	// اسم المدير العام
	CEOPhone    string `gorm:"type:varchar(20)" json:"ceo_phone"`
	// هاتف المدير العام
	CEOEmail    string `gorm:"type:varchar(100)" json:"ceo_email"`
	// ايميل المدير العام
	CEOAddress  string `gorm:"type:varchar(200)" json:"ceo_address"`
	// عنوان المدير العام
}


func (c *Company) BeforeSave(tx *gorm.DB) (err error) {

	// ---------- Basic normalization ----------
	c.Name = strings.TrimSpace(c.Name)
	c.LocalAddress = strings.TrimSpace(c.LocalAddress)
	c.Duration = strings.TrimSpace(c.Duration)

	if c.Email != "" {
		c.Email = strings.ToLower(strings.TrimSpace(c.Email))
	}

	// ---------- Required fields ----------
	if c.Name == "" {
		return errors.New("company name is required")
	}

	if c.LocalAddress == "" {
		return errors.New("company local address is required")
	}

	if c.CompanyTypeID == 0 {
		return errors.New("company type is required")
	}

	if c.OfficeId == 0 {
		return errors.New("office is required")
	}

	// ---------- FK existence checks (NO cascade) ----------

	// Check company type exists
	if err := tx.First(&general_models.CompanyType{}, c.CompanyTypeID).Error; err != nil {
		return errors.New("invalid company type")
	}

	// Check office exists
	if err := tx.First(&general_models.Office{}, c.OfficeId).Error; err != nil {
		return errors.New("invalid office")
	}

    if c.CompanyType.ID != 0 {
    	return errors.New("nested company_type creation is not allowed")
    }
    if c.Office.ID != 0 {
    	return errors.New("nested office creation is not allowed")
    }
	if len(c.CompanyActivity) > 0 {
		for _, a := range c.CompanyActivity {
			if a.TradingActivityID == 0 {
				return errors.New("trading_activity_id is required")
			}
		}
	}

	return nil
}
