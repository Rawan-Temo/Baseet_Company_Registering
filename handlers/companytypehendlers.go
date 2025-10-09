package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/models"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AllCompanyTypes(c *fiber.Ctx) error {
	db := database.DB
	var total int64
	var companyTypes []models.CompanyType
	queryArgs := c.Context().QueryArgs()
	fmt.Println(queryArgs)
	queries := map[string][]string{}

	queryArgs.VisitAll(func(key, value []byte) {
		k := string(key)
		v := string(value)
		queries[k] = append(queries[k], v)
	})
	allowedCols := []string{"name", "ID", "created_at", "updated_at", "deleted_at"}
	queryBuilder := utils.NewQueryBuilder(db, queries, allowedCols)
	queryBuilder.Filter().Sort().LimitFields().Paginate()
	if err := queryBuilder.Apply().Find(&companyTypes).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch companyTypes",
			"details": err.Error(),
		})
	}
	// count total companyTypes
	if err := utils.NewQueryBuilder(db, queries, allowedCols).Filter().Apply().Model(&models.CompanyType{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to count companyTypes",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"total":   total,
		"results": len(companyTypes),
		"data":    companyTypes,
	})
}
func CreateCompanyType(c *fiber.Ctx) error {
	db := database.DB
	var companyType models.CompanyType
	if err := c.BodyParser(&companyType); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not parse json",
			"error":   err.Error(),
		})
	}

	if err := db.Create(&companyType).Error; err != nil {

		if strings.Contains(err.Error(), "duplicate key value") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not create type",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   companyType,
	})
}
func GetCompanyType(c *fiber.Ctx) error {
	db := database.DB
	var companyType models.CompanyType
	id := c.Params("id")
	if err := db.First(&companyType, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Company not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch company",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   companyType,
	})
}

func UpdateCompanyType(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var input map[string]interface{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}

	allowed := map[string]bool{"name": true}
	sanitized := map[string]interface{}{}
	for k, v := range input {
		if allowed[strings.ToLower(k)] {
			sanitized[strings.ToLower(k)] = v
		}
	}
	if len(sanitized) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No valid fields"})
	}

	var companyType models.CompanyType
	res := db.Model(&companyType).Where("id = ?", id).Updates(sanitized)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Update failed"})
	}
	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Company not found"})
	}
	if err := db.First(&companyType, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Fetch after update failed"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Company updated successfully",
		"data":    companyType,
	})
}

// DELETE /company-types/:id
func DeleteCompanyType(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	// Use RowsAffected to check if the record exists
	res := db.Delete(&models.CompanyType{}, id)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete CompanyType",
		})
	}

	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "CompanyType not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "CompanyType deleted successfully",
	})
}
