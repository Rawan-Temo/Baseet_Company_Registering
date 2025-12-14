package handlers

import (
	"errors"
	"strings"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	general_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/general"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AllOffices(c *fiber.Ctx) error {
	db := database.DB
	var total int64
	var offices []general_models.Office
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
	if err := queryBuilder.Apply().Find(&offices).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch offices",
			"details": err.Error(),
		})
	}
	// count total offices
	if err := utils.NewQueryBuilder(db, queries, allowedCols).Filter().Apply().Model(&general_models.Office{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to count offices",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"total":   total,
		"results": len(offices),
		"data":    offices,
	})
}

func CreateOffice(c *fiber.Ctx) error {
	db := database.DB
	var office general_models.Office
	if err := c.BodyParser(&office); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not parse json",
			"error":   err.Error(),
		})
	}
	if err := db.Create(&office).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not create office",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   office,
	})
}

func GetOffice(c *fiber.Ctx) error {
	db := database.DB
	var office general_models.Office
	id := c.Params("id")
	if err := db.First(&office, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Office not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch office",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   office,
	})
}

func UpdateOffice(c *fiber.Ctx) error {
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

	var office general_models.Office
	res := db.Model(&office).Where("id = ?", id).Updates(sanitized)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Update failed"})
	}
	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Office not found"})
	}
	if err := db.First(&office, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Fetch after update failed"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Office updated successfully",
		"data":    office,
	})
}

func DeleteOffice(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	// Use RowsAffected to check if the record exists
	res := db.Delete(&general_models.Office{}, id)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete Office",
		})
	}

	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Office not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Office deleted successfully",
	})
}
