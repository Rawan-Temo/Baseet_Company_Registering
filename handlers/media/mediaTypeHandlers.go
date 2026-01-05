package handlers

import (
	"errors"
	"strings"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/dtos"
	media_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/media"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllMediaTypes(c *fiber.Ctx) error {
	db := database.DB
	var total int64
	var mediaTypes []media_models.MediaType
	queryArgs := c.Context().QueryArgs()
	queries := map[string][]string{}

	queryArgs.VisitAll(func(key, value []byte) {
		k := string(key)
		v := string(value)
		queries[k] = append(queries[k], v)
	})
	allowedCols := []string{"name", "description", "ID", "created_at", "updated_at", "deleted_at"}
	queryBuilder := utils.NewQueryBuilder(db, queries, allowedCols)
	queryBuilder.Filter().Sort().LimitFields().Paginate()
	if err := queryBuilder.Apply().Find(&mediaTypes).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch media types",
			"details": err.Error(),
		})
	}
	// count total media types
	if err := utils.NewQueryBuilder(db, queries, allowedCols).Filter().Apply().Model(&media_models.MediaType{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to count media types",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"total":   total,
		"results": len(mediaTypes),
		"data":    mediaTypes,
	})
}

func CreateMediaType(c *fiber.Ctx) error {
	db := database.DB
	var req dtos.CreateMediaTypeRequest
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

	mediaType := media_models.MediaType{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := db.Create(&mediaType).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not create media type",
			"error":   err.Error(),
		})
	}

	response := dtos.MediaTypeResponse{
		ID:          mediaType.ID,
		Name:        mediaType.Name,
		Description: mediaType.Description,
		CreatedAt:   mediaType.CreatedAt,
		UpdatedAt:   mediaType.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func GetMediaTypeByID(c *fiber.Ctx) error {
	db := database.DB
	var mediaType media_models.MediaType
	id := c.Params("id")
	if err := db.First(&mediaType, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Media type not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch media type",
		})
	}

	response := dtos.MediaTypeResponse{
		ID:          mediaType.ID,
		Name:        mediaType.Name,
		Description: mediaType.Description,
		CreatedAt:   mediaType.CreatedAt,
		UpdatedAt:   mediaType.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func UpdateMediaType(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var req dtos.UpdateMediaTypeRequest
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

	var mediaType media_models.MediaType
	if err := db.First(&mediaType, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Media type not found"})
	}

	if req.Name != nil {
		mediaType.Name = *req.Name
	}
	if req.Description != nil {
		mediaType.Description = *req.Description
	}

	res := db.Save(&mediaType)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Update failed"})
	}

	response := dtos.MediaTypeResponse{
		ID:          mediaType.ID,
		Name:        mediaType.Name,
		Description: mediaType.Description,
		CreatedAt:   mediaType.CreatedAt,
		UpdatedAt:   mediaType.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Media type updated successfully",
		"data":    response,
	})
}

func DeleteMediaType(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	// Use RowsAffected to check if the record exists
	res := db.Delete(&media_models.MediaType{}, id)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete Media type",
		})
	}

	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Media type not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Media type deleted successfully",
	})
}