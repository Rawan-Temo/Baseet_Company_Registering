package handlers

import (
	"errors"
	"log"
	"strings"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/dtos"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/helpers"
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
	allowedCols := []string{"company_id", "trading_activity_id", "image", "ID", "created_at", "updated_at", "deleted_at"}
	queryBuilder := utils.NewQueryBuilder(db, queries, allowedCols)
	queryBuilder.Filter().Sort().LimitFields().Paginate()
	if err := queryBuilder.Apply().Preload("Company").Preload("TradingActivity").Find(&companyActivities).Error; err != nil {
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
	contentType := c.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		if err := helpers.ValidateMultiPartFormActivity(c, &req); err != nil {
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
	companyActivity := company_models.CompanyActivity{
		CompanyID:         req.CompanyId,
		TradingActivityID: req.TradingActivityID,
		Image:             req.Image,
	}
	tx := db.Begin()
	imageConfig := utils.DefaultImageConfig()
	imageConfig.UploadDir = "./uploads/tradingActivities/"
	if err := tx.First(&company, companyActivity.CompanyID).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "invalid company id",
		})
	}

	if company.CompanyCategory == company_models.CompanyCategoryRepresentationOffice {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "representation offices cannot legally have trading activities",
		})
	}

	if err := tx.Create(&companyActivity).Error; err != nil {
		tx.Rollback()
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
	if uploadErr := utils.UploadImage(c, "image", imageConfig, req.Image); uploadErr != nil {
		tx.Rollback()
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
	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "validation error",
			"error":   err,
		})
	}

	var companyActivity company_models.CompanyActivity
	if err := db.Preload("Company").First(&companyActivity, id).Error; err != nil {
		return handleNotFoundOrError(err)
	}
	companyActivity.TradingActivityID = req.TradingActivityID
	if err := db.Save(&companyActivity).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not update company activity",
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Company activity updated successfully",
		"data":    response,
	})
}

func DeleteCompanyActivity(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "ID required"})
	}

	db := database.DB
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Get record
	var activity company_models.CompanyActivity
	if err := tx.First(&activity, "id = ?", id).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "Not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": "Fetch failed", "error": err.Error()})
	}

	// Delete image first
	imagePath := activity.Image
	if imagePath != "" {
		if err := utils.DeleteImage(imagePath, "./uploads/tradingActivities/"); err != nil {
			tx.Rollback()
			// Try to restore image? Or log for manual cleanup
			log.Printf("Failed to delete image, transaction rolled back: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "fail",
				"message": "Failed to delete associated image",
			})
		}
	}

	// Delete DB record
	result := tx.Delete(&company_models.CompanyActivity{}, "id = ?", id)
	if result.Error != nil {
		tx.Rollback()
		// At this point image is deleted but DB not - need recovery logic
		// You could implement a compensation: try to restore image
		log.Printf("CRITICAL: DB delete failed after image deletion. Image %s orphaned", imagePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to delete record",
			"error":   result.Error.Error(),
		})
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "Not found"})
	}

	if err := tx.Commit().Error; err != nil {
		// DB commit failed but image already deleted
		log.Printf("CRITICAL: Commit failed after image deletion. Image %s orphaned", imagePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to commit deletion",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Company activity deleted successfully",
	})
}

// func CreateManyCompnayAcitivities(c *fiber.Ctx) error {

// }
func handleNotFoundOrError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "License not found")
	}
	return fiber.NewError(fiber.StatusInternalServerError, "Database error")
}
