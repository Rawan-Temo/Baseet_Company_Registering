package company_models

import (
	"errors"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/models"
	"gorm.io/gorm"
)

// -----------------------------

// People represents an individual associated with a company, such as a partner,
// authorized representative, board member, or stakeholder.
type People struct {
	models.NewGormModel
	CompanyID    uint    `gorm:"index;not null" json:"company_id"`
	Company      Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	FullName     string  `gorm:"type:varchar(255);not null" json:"full_name"`
	Email        string  `gorm:"type:varchar(255)" json:"email"`
	Phone        string  `gorm:"type:varchar(50)" json:"phone"`
	Address      string  `gorm:"type:varchar(500)" json:"address"`
	Role         string  `gorm:"type:varchar(50);not null" json:"role"` // "Partner", "AuthorizedRep", "BoardMember", "StakeHolder"
	ExtraDetails string  `gorm:"type:text" json:"extra_details"`
}

func (p *People) BeforeSave(tx *gorm.DB) error {
	if p.CompanyID == 0 {
		return errors.New("company_id is required")
	}
	if p.FullName == "" {
		return errors.New("full_name is required")
	}
	if p.Role == "" {
		return errors.New("role is required")
	}
	return nil
}
