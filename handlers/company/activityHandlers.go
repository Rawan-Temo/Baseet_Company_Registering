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

func AllTradingActivity(c *fiber.Ctx) error {
	db := database.DB
	var total int64
	var tradingActivities []company_models.TradingActivity
	queryArgs := c.Context().QueryArgs()
	queries := map[string][]string{}

	queryArgs.VisitAll(func(key, value []byte) {
		k := string(key)
		v := string(value)
		queries[k] = append(queries[k], v)
	})

	allowedCols := []string{"name", "id", "created_at", "updated_at", "deleted_at"}
	queryBuilder := utils.NewQueryBuilder(db, queries, allowedCols)
	queryBuilder.Filter().Sort().LimitFields().Paginate()
	if err := queryBuilder.Apply().Find(&tradingActivities).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch trading activities",
			"details": err.Error(),
		})
	}
	// count total trading activities
	if err := utils.NewQueryBuilder(db, queries, allowedCols).Filter().Apply().Model(&company_models.TradingActivity{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to count trading activities",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"total":   total,
		"results": len(tradingActivities),
		"data":    tradingActivities,
	})
}

func CreateTradingActivity(c *fiber.Ctx) error {
	db := database.DB
	var req dtos.CreateTradingActivityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
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

	activity := company_models.TradingActivity{
		Name: req.Name,
	}

	if err := db.Create(&activity).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not create trading activity",
			"error":   err.Error(),
		})
	}

		response := GetTradingResponse(activity)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func GetTradingActivityByID(c *fiber.Ctx) error {
	db := database.DB
	var tradingActivity company_models.TradingActivity
	id := c.Params("id")
	if err := db.First(&tradingActivity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Trading activity not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch trading activity",
		})
	}
	response := GetTradingResponse(tradingActivity)

	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func UpdateTradingActivity(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var req dtos.UpdateTradingActivityRequest
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

	var tradingActivity company_models.TradingActivity
	res := db.Model(&tradingActivity).Where("id = ?", id).Updates(map[string]interface{}{"name": *req.Name})
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Update failed"})
	}
	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Trading activity not found"})
	}
	if err := db.First(&tradingActivity, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Fetch after update failed"})
	}

	response := GetTradingResponse(tradingActivity)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Trading activity updated successfully",
		"data":    response,
	})
}

func DeleteTradingActivity(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	// Use RowsAffected to check if the record exists
	res := db.Delete(&company_models.TradingActivity{}, id)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete Trading activity",
		})
	}

	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Trading activity not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Trading activity deleted successfully",
	})
}


func GetTradingResponse(activity company_models.TradingActivity) dtos.TradingActivityResponse{
	return dtos.TradingActivityResponse{
		ID: activity.ID,
		Name:              activity.Name,
		CreatedAt:         activity.CreatedAt,
		UpdatedAt:         activity.UpdatedAt,
	}
}