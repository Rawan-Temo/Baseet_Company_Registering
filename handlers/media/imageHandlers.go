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

func GetAllImages(c *fiber.Ctx) error {
	db := database.DB
	var total int64
	var images []media_models.Image
	queryArgs := c.Context().QueryArgs()
	queries := map[string][]string{}

	queryArgs.VisitAll(func(key, value []byte) {
		k := string(key)
		v := string(value)
		queries[k] = append(queries[k], v)
	})
	allowedCols := []string{"src", "company_id", "media_type_id", "ID", "created_at", "updated_at", "deleted_at"}
	queryBuilder := utils.NewQueryBuilder(db, queries, allowedCols)
	queryBuilder.Filter().Sort().LimitFields().Paginate()
	if err := queryBuilder.Apply().Find(&images).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch images",
			"details": err.Error(),
		})
	}
	// count total images
	if err := utils.NewQueryBuilder(db, queries, allowedCols).Filter().Apply().Model(&media_models.Image{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to count images",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"total":   total,
		"results": len(images),
		"data":    images,
	})
}

func CreateImage(c *fiber.Ctx) error {
	db := database.DB
	var req dtos.CreateImageRequest
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

	image := media_models.Image{
		Src:         req.Src,
		CompanyId:   req.CompanyId,
		MediaTypeId: req.MediaTypeId,
	}

	if err := db.Create(&image).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not create image",
			"error":   err.Error(),
		})
	}

	response := dtos.ImageResponse{
		ID:          image.ID,
		Src:         image.Src,
		CompanyId:   image.CompanyId,
		MediaTypeId: image.MediaTypeId,
		CreatedAt:   image.CreatedAt,
		UpdatedAt:   image.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func GetImageByID(c *fiber.Ctx) error {
	db := database.DB
	var image media_models.Image
	id := c.Params("id")
	if err := db.First(&image, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Image not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch image",
		})
	}

	response := dtos.ImageResponse{
		ID:          image.ID,
		Src:         image.Src,
		CompanyId:   image.CompanyId,
		MediaTypeId: image.MediaTypeId,
		CreatedAt:   image.CreatedAt,
		UpdatedAt:   image.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func UpdateImage(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var req dtos.UpdateImageRequest
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

	var image media_models.Image
	if err := db.First(&image, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Image not found"})
	}

	if req.Src != nil {
		image.Src = *req.Src
	}
	if req.CompanyId != nil {
		image.CompanyId = *req.CompanyId
	}
	if req.MediaTypeId != nil {
		image.MediaTypeId = *req.MediaTypeId
	}

	res := db.Save(&image)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Update failed"})
	}

	response := dtos.ImageResponse{
		ID:          image.ID,
		Src:         image.Src,
		CompanyId:   image.CompanyId,
		MediaTypeId: image.MediaTypeId,
		CreatedAt:   image.CreatedAt,
		UpdatedAt:   image.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Image updated successfully",
		"data":    response,
	})
}

func DeleteImage(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	// Use RowsAffected to check if the record exists
	res := db.Delete(&media_models.Image{}, id)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete Image",
		})
	}

	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Image not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Image deleted successfully",
	})
}