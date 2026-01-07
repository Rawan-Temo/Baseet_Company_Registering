package helpers

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/dtos"
	auth_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/auth"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ValidateMultiPartForm(c *fiber.Ctx, req *dtos.CreateLicenseRequest) error {
	_, err := c.MultipartForm()
	if err == nil {
		companyIdStr := c.FormValue("company_id")

		if companyIdStr != "" {
			if companyId, parseErr := strconv.ParseUint(companyIdStr, 10, 32); parseErr == nil {
				req.CompanyId = uint(companyId)
			}
		}

		if startDateStr := c.FormValue("start_date"); startDateStr != "" {
			if startDate, parseErr := utils.ParseDate(startDateStr); parseErr == nil {
				req.StartDate = startDate
			}
		}

		if expirationDateStr := c.FormValue("expiration_date"); expirationDateStr != "" {
			if expirationDate, parseErr := utils.ParseDate(expirationDateStr); parseErr == nil {
				req.ExpirationDate = expirationDate
			}
		}

		// Handle image file upload
		file, err := c.FormFile("image")
		if err == nil {
			fileName := utils.GenerateFileName(file)
			req.Image = &fileName
			fmt.Println(*req.Image)

		}
	}

	// Validate required fields
	if req.CompanyId == 0 {
		return errors.New("company_id is required")
	}

	if req.StartDate.Equal((time.Time{})) {
		return errors.New("start Date is required")
	}

	if req.ExpirationDate.Equal((time.Time{})) {
		return errors.New("ExpirationDate is required")
	}
	return nil
}

func GetLicenseResponse(license auth_models.License) dtos.LicenseResponse {
	return dtos.LicenseResponse{
		ID:             license.ID,
		CompanyId:      license.CompanyId,
		StartDate:      license.StartDate,
		ExpirationDate: license.ExpirationDate,
		Image:          license.Image,
		CreatedAt:      license.CreatedAt,
		UpdatedAt:      license.UpdatedAt,
	}
}
func Rollback(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
	}
}
