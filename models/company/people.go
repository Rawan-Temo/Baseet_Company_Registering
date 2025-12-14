package company_models

import (
	"gorm.io/gorm"
)

// -----------------------------

// People represents an individual associated with a company, such as a partner,
// authorized representative, board member, or stakeholder.
type People struct {
    gorm.Model
    CompanyID uint   `gorm:"index;not null" json:"company_id"`
    Company  Company `gorm:"foreignKey:CompanyID" json:"company"`
    FullName  string `gorm:"type:varchar(255);not null" json:"full_name"`
    Email     string `gorm:"type:varchar(255);uniqueIndex" json:"email"`
    Phone     string `gorm:"type:varchar(50)" json:"phone"`
    Address   string `gorm:"type:varchar(500)" json:"address"`
    Role      string `gorm:"type:varchar(50);not null" json:"role"` // "Partner", "AuthorizedRep", "BoardMember", "StakeHolder"
}

func (p *People) BeforeCreate(tx *gorm.DB) (err error) {
    var existing People
    if err := tx.Where("email = ?", p.Email).First(&existing).Error; err == nil {
        // Person exists, cancel creation and set ID to existing
        p.ID = existing.ID
        // Stop creating a new row
        return gorm.ErrRegistered
    }
    return nil
}