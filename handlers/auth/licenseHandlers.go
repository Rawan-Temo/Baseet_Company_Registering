package handlers

import (
	"errors"
	"strings"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/dtos"
	auth_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/auth"
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
	allowedCols := []string{"company_id", "ID", "created_at", "updated_at", "deleted_at", "start_date", "expiration_date"}
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"total":   total,
		"results": len(licenses),
		"data":    licenses,
	})
}

func CreateLicense(c *fiber.Ctx) error {
	db := database.DB
	var req dtos.CreateLicenseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not parse json",
			"error":   err.Error(),
		})
	}

	license := auth_models.License{
		CompanyId:      req.CompanyId,
		StartDate:      req.StartDate,
		ExpirationDate: req.ExpirationDate,
		Image:          req.Image ,
	}

	if err := db.Create(&license).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not create license",
			"error":   err.Error(),
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

	var req dtos.UpdateLicenseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}

	sanitized := map[string]interface{}{}

	if req.StartDate != nil {
		sanitized["start_date"] = req.StartDate
	}
	if req.ExpirationDate != nil {
		sanitized["expiration_date"] = req.ExpirationDate
	}
	if req.Image != nil {
		sanitized["image"] = req.Image
	}

	if len(sanitized) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No valid fields"})
	}

	var license auth_models.License
	res := db.Model(&license).Where("id = ?", id).Updates(sanitized)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Update failed"})
	}
	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "License not found"})
	}
	if err := db.First(&license, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Fetch after update failed"})
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
		"message": "License updated successfully",
		"data":    response,
	})
}

func DeleteLicense(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	// Use RowsAffected to check if the record exists
	res := db.Delete(&auth_models.License{}, id)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete License",
		})
	}

	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "License not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "License deleted successfully",
	})
}





