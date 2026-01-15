package handlers

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/dtos"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/helpers"
	auth_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/auth"
	company_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/company"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllLicenses(c *fiber.Ctx) error {
	db := database.DB
	var total int64
	var licenses []auth_models.License
	queryArgs := c.Context().QueryArgs()
	queries := map[string][]string{}

	queryArgs.VisitAll(func(key, value []byte) {
		k := string(key)
		v := string(value)
		queries[k] = append(queries[k], v)
	})
	allowedCols := []string{"company_id", "image", "id", "created_at", "updated_at", "deleted_at", "start_date", "expiration_date"}
	queryBuilder := utils.NewQueryBuilder(db, queries, allowedCols)
	queryBuilder.Filter().Sort().LimitFields().Paginate()
	if err := queryBuilder.Apply().Find(&licenses).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch licenses",
			"details": err.Error(),
		})
	}
	// count total licenses
	if err := utils.NewQueryBuilder(db, queries, allowedCols).Filter().Apply().Model(&auth_models.License{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to count licenses",
			"details": err.Error(),
		})
	}

	// Convert licenses to response DTOs
	var licenseResponses []dtos.LicenseResponse
	for _, license := range licenses {
		licenseResponses = append(licenseResponses, dtos.LicenseResponse{
			ID:             license.ID,
			CompanyId:      license.CompanyId,
			StartDate:      license.StartDate,
			ExpirationDate: license.ExpirationDate,
			Image:          license.Image,
			CreatedAt:      license.CreatedAt,
			UpdatedAt:      license.UpdatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"total":   total,
		"results": len(licenseResponses),
		"data":    licenseResponses,
	})
}

func CreateLicense(c *fiber.Ctx) error {
	db := database.DB
	committed := false
	var req dtos.CreateLicenseRequest
	contentType := c.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		if err := helpers.ValidateMultiPartFormLicense(c, &req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid request format",
				"error":   err.Error(),
			})
		}
	} else {
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid request format",
				"error":   "Expected JSON or multipart form data",
			})
		}
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "validation error",
			"error":   err,
		})
	}
	// Now create the license
	license := auth_models.License{
		CompanyId:      req.CompanyId,
		StartDate:      req.StartDate,
		ExpirationDate: req.ExpirationDate,
		Image:          req.Image,
	}
	tx := db.Begin()
	imageConfig := utils.DefaultImageConfig()
	imageConfig.UploadDir = "./uploads/licenses/"
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()
	if err := tx.Create(&license).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": "License already exists for this company",
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not create license",
			"error":   err.Error(),
		})

	}
	var company company_models.Company
	if err := tx.First(&company, license.CompanyId).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not fetch created license",
			"error":   err.Error(),
		})
	}
	if license.ExpirationDate.After(company.License) {
		company.License = license.ExpirationDate
		if err := tx.Save(&company).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "fail",
				"message": "Could not fetch created license",
				"error":   err.Error(),
			})
		}
	}

	if uploadErr := utils.UploadImage(c, "image", imageConfig, *req.Image); uploadErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "fail",
			"error":  uploadErr.Error(),
		})
	}
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not create license",
			"error":   err.Error(),
		})
	}
	committed = true
	response := helpers.GetLicenseResponse(license)
	response.Company = company
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func GetLicenseByID(c *fiber.Ctx) error {
	db := database.DB
	var license auth_models.License
	id := c.Params("id")
	if err := db.First(&license, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "License not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch license",
		})
	}

	response := dtos.LicenseResponse{
		ID:             license.ID,
		CompanyId:      license.CompanyId,
		StartDate:      license.StartDate,
		ExpirationDate: license.ExpirationDate,
		Image:          license.Image,
		CreatedAt:      license.CreatedAt,
		UpdatedAt:      license.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func UpdateLicense(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	// Get current license
	var license auth_models.License
	if err := db.Preload("Company").First(&license, id).Error; err != nil {
		return handleNotFoundOrError(err)
	}
	var req dtos.UpdateLicenseRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "validation error",
			"error":   err,
		})
	}
	utils.UpdateStruct(&license, &req)

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Update license
	if err := tx.Save(&license).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not update license",
			"error":   err.Error(),
		})
	}

	if license.ExpirationDate.After(license.Company.License) {
		license.Company.License = license.ExpirationDate
		if err := tx.Save(&license.Company).Error; err != nil {

			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "fail",
				"message": "Could not update company license date",
				"error":   err.Error(),
			})
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Transaction commit failed",
			"error":   err.Error(),
		})
	}
	response := dtos.LicenseResponse{
		ID:             license.ID,
		CompanyId:      license.CompanyId,
		Company:        license.Company,
		StartDate:      license.StartDate,
		ExpirationDate: license.ExpirationDate,
		Image:          license.Image,
		CreatedAt:      license.CreatedAt,
		UpdatedAt:      license.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "License updated successfully",
		"data":    response,
	})
}

func DeleteLicense(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	// Get license first (outside transaction for file cleanup later)
	var license auth_models.License
	if err := db.Preload("Company").First(&license, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "License not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch license",
		})
	}

	// Store image path before deletion
	var imagePath string
	if license.Image != nil {
		imagePath = *license.Image
	}
	// Use a single SQL query to handle everything
	result := db.Transaction(func(tx *gorm.DB) error {
		// 1. Delete the license
		if err := tx.Delete(&license).Error; err != nil {
			return err
		}

		// 2. Check if this was the current license and update company if needed
		// We need to check this AFTER deletion
		licenseDate := license.ExpirationDate.Truncate(24 * time.Hour)
		companyLicenseDate := license.Company.License.Truncate(24 * time.Hour)

		if licenseDate.Equal(companyLicenseDate) {
			// Get the max expiration date from remaining licenses
			var maxDate time.Time
			err := tx.Model(&auth_models.License{}).
				Where("company_id = ?", license.CompanyId).
				Select("COALESCE(MAX(expiration_date), '0001-01-01')").
				Scan(&maxDate).Error

			if err != nil {
				return err
			}

			// Update company license
			updateValue := maxDate
			if maxDate.IsZero() || maxDate.Year() == 1 { // Check for zero time
				updateValue = time.Time{}
			}

			return tx.Model(&license.Company).
				Update("license", updateValue).Error
		}

		return nil
	})

	if result != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete license",
			"details": result.Error(),
		})
	}

	// Clean up image file after successful deletion
	if imagePath != "" && !strings.HasPrefix(imagePath, "http://") && !strings.HasPrefix(imagePath, "https://") {
		if err := utils.DeleteImage(imagePath, "./uploads/licenses/"); err != nil {
			log.Printf("Warning: Failed to delete license image: %v", err)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "License deleted successfully",
	})
}

// ========================== Helpers ==========================

func handleNotFoundOrError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "License not found")
	}
	return fiber.NewError(fiber.StatusInternalServerError, "Database error")
}
