package utils

import (
	"fmt"

	auth_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/auth"
	"gorm.io/gorm"
)

func CreateDefaultAdmin(db *gorm.DB)error {
	var user auth_models.User
	if err := db.Where("role = ?" , auth_models.RoleAdmin).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			admin := auth_models.User{
				UserName: "admin",
				Email:    "admin@example.com",
				Password: "admin123",
				Role:     auth_models.RoleAdmin,
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
