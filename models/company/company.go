package company_models

import (
	general_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/general"
	"gorm.io/gorm"
)

// TODO 	add vaildation tags and dont casecade create for related models
// Company represents a registered company with its details and relationships.



type Company struct {
	gorm.Model
	Name              string            `gorm:"type:varchar(100);not null" json:"name"`
	TradeNames        string            `gorm:"type:varchar(200)" json:"trade_names"`
	LocalIdentifier   string            `gorm:"type:varchar(100)" json:"local_identifier"`  // رقم التعريف الممنوح / من هيئة عامة أخرى
	Address           string            `gorm:"type:varchar(200)" json:"address"`
	Description       string            `gorm:"type:varchar(500)" json:"description"`
	Email             string            `gorm:"uniqueIndex;type:varchar(100)" json:"email"`
	PhoneNumber       string            `gorm:"type:varchar(15)" json:"phone_number"`
	CompanyTypeID     uint              `gorm:"column:type_id;index" json:"type_id"`
	CompanyType       general_models.CompanyType       `gorm:"foreignKey:CompanyTypeID" json:"type,omitempty"`
	OfficeId            int               `gorm:"type:integer" json:"officeId"`
	Office            general_models.Office            `gorm:"foreignKey:OfficeId" json:"office,omitempty"`
	TradingActivities []TradingActivity `gorm:"many2many:company_trading_activities;" json:"trading_activities"`
	People   []People		  `gorm:"foreignKey:CompanyID" json:"people"`
	IsLicensed        bool              `gorm:"type:boolean;default:false" json:"is_licensed"`
	Duration          string            `gorm:"type:varchar(100)" json:"duration"`

}
