package handlers

import (
	"errors"
	"strings"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/dtos"
	company_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/company"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllCompanyActivities(c *fiber.Ctx) error {
	db := database.DB
	var total int64
	var companyActivities []company_models.CompanyActivity
	queryArgs := c.Context().QueryArgs()
	queries := map[string][]string{}

	queryArgs.VisitAll(func(key, value []byte) {
		k := string(key)
		v := string(value)
		queries[k] = append(queries[k], v)
	})
	allowedCols := []string{"company_id", "trading_activity_id", "ID", "created_at", "updated_at", "deleted_at"}
	queryBuilder := utils.NewQueryBuilder(db, queries, allowedCols)
	queryBuilder.Filter().Sort().LimitFields().Paginate()
	if err := queryBuilder.Apply().Find(&companyActivities).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch company activities",
			"details": err.Error(),
		})
	}
	// count total company activities
	if err := utils.NewQueryBuilder(db, queries, allowedCols).Filter().Apply().Model(&company_models.CompanyActivity{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to count company activities",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"total":   total,
		"results": len(companyActivities),
		"data":    companyActivities,
	})
}

func CreateCompanyActivity(c *fiber.Ctx) error {
	db := database.DB
	
	var company company_models.Company
	var req dtos.CreateCompanyActivityRequest
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
	companyActivity := company_models.CompanyActivity{
		CompanyID:         req.CompanyId,
		TradingActivityID: req.TradingActivityID,
		Image:             req.Image,
	}
	if err := db.First(&company, companyActivity.CompanyID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "invalid company id",
		})
	}

	if company.CompanyType.Name == "Free"{
		c.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": "Free type company cannot leggaly have any activities in syria",
		})
		return nil
	}

	if err := db.Create(&companyActivity).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not create company activity",
			"error":   err.Error(),
		})
	}

	response := dtos.CompanyActivityResponse{
		ID:                companyActivity.ID,
		CompanyId:         companyActivity.CompanyID,
		TradingActivityID: companyActivity.TradingActivityID,
		Image:             companyActivity.Image,
		CreatedAt:         companyActivity.CreatedAt,
		UpdatedAt:         companyActivity.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func GetCompanyActivityByID(c *fiber.Ctx) error {
	db := database.DB
	var companyActivity company_models.CompanyActivity
	id := c.Params("id")
	if err := db.First(&companyActivity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Company activity not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch company activity",
		})
	}

	response := dtos.CompanyActivityResponse{
		ID:                companyActivity.ID,
		CompanyId:         companyActivity.CompanyID,
		TradingActivityID: companyActivity.TradingActivityID,
		Image:             companyActivity.Image,
		CreatedAt:         companyActivity.CreatedAt,
		UpdatedAt:         companyActivity.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func UpdateCompanyActivity(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var req dtos.UpdateCompanyActivityRequest
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
	if req.Image == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No valid fields"})
	}

	var companyActivity company_models.CompanyActivity
	res := db.Model(&companyActivity).Where("id = ?", id).Updates(map[string]interface{}{"image": *req.Image})
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Update failed"})
	}
	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Company activity not found"})
	}
	if err := db.First(&companyActivity, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Fetch after update failed"})
	}

	response := dtos.CompanyActivityResponse{
		ID:                companyActivity.ID,
		CompanyId:         companyActivity.CompanyID,
		TradingActivityID: companyActivity.TradingActivityID,
		Image:             companyActivity.Image,
		CreatedAt:         companyActivity.CreatedAt,
		UpdatedAt:         companyActivity.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Company activity updated successfully",
		"data":    response,
	})
}

func DeleteCompanyActivity(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	// Use RowsAffected to check if the record exists
	res := db.Delete(&company_models.CompanyActivity{}, id)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete Company activity",
		})
	}

	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Company activity not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Company activity deleted successfully",
	})
}





