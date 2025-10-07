package utils

import (
	"fmt"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/models"
	"gorm.io/gorm"
)

func CreateDefaultAdmin(db *gorm.DB)error {
	var user models.User
	if err := db.Where("role = ?" , models.RoleAdmin).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			admin := models.User{
				UserName: "admin",
				Email:    "admin@example.com",
				Password: "admin123",
				Role:     models.RoleAdmin,
				CompanyId: nil,
			}
			
			if err := db.Create(&admin).Error; err != nil {
				// Handle error
				return fmt.Errorf("failed to create default admin user: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check for existing admin user: %w", err)
		}
	}
	return nil
}
