package handlers

import (
	"errors"
	"strings"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/dtos"
	general_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/general"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AllCompanyTypes(c *fiber.Ctx) error {
	db := database.DB
	var total int64
	var companyTypes []general_models.CompanyType
	queryArgs := c.Context().QueryArgs()
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
	if err := utils.NewQueryBuilder(db, queries, allowedCols).Filter().Apply().Model(&general_models.CompanyType{}).Count(&total).Error; err != nil {
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
	var req dtos.CreateCompanyTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not parse json",
			"error":   err.Error(),
		})
	}

	if err:= utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "validation error",
			"error":   err,
		})
	}
	companyType := general_models.CompanyType{
		Name: req.Name,
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

	response := dtos.CompanyTypeResponse{
		ID:        companyType.ID,
		Name:      companyType.Name,
		CreatedAt: companyType.CreatedAt,
		UpdatedAt: companyType.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}
func GetCompanyType(c *fiber.Ctx) error {
	db := database.DB
	var companyType general_models.CompanyType
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

	response := dtos.CompanyTypeResponse{
		ID:        companyType.ID,
		Name:      companyType.Name,
		CreatedAt: companyType.CreatedAt,
		UpdatedAt: companyType.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func UpdateCompanyType(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var req dtos.UpdateCompanyTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}
	if err:= utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "validation error",
			"error":   err,
		})
	}

	if req.Name == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No valid fields"})
	}

	var companyType general_models.CompanyType
	res := db.Model(&companyType).Where("id = ?", id).Updates(map[string]interface{}{"name": *req.Name})
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Update failed"})
	}
	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Company not found"})
	}
	if err := db.First(&companyType, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Fetch after update failed"})
	}

	response := dtos.CompanyTypeResponse{
		ID:        companyType.ID,
		Name:      companyType.Name,
		CreatedAt: companyType.CreatedAt,
		UpdatedAt: companyType.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Company updated successfully",
		"data":    response,
	})
}

// DELETE /company-types/:id
func DeleteCompanyType(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	// Use RowsAffected to check if the record exists
	res := db.Delete(&general_models.CompanyType{}, id)
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
