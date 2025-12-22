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

func GetAllDocuments(c *fiber.Ctx) error {
	db := database.DB
	var total int64
	var documents []media_models.Document
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
	if err := queryBuilder.Apply().Find(&documents).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch documents",
			"details": err.Error(),
		})
	}
	// count total documents
	if err := utils.NewQueryBuilder(db, queries, allowedCols).Filter().Apply().Model(&media_models.Document{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to count documents",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"total":   total,
		"results": len(documents),
		"data":    documents,
	})
}

func CreateDocument(c *fiber.Ctx) error {
	db := database.DB
	var req dtos.CreateDocumentRequest
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

	document := media_models.Document{
		Src:         req.Src,
		CompanyId:   req.CompanyId,
		MediaTypeId: req.MediaTypeId,
	}

	if err := db.Create(&document).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not create document",
			"error":   err.Error(),
		})
	}

	response := dtos.DocumentResponse{
		ID:          document.ID,
		Src:         document.Src,
		CompanyId:   document.CompanyId,
		MediaTypeId: document.MediaTypeId,
		CreatedAt:   document.CreatedAt,
		UpdatedAt:   document.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func GetDocumentByID(c *fiber.Ctx) error {
	db := database.DB
	var document media_models.Document
	id := c.Params("id")
	if err := db.First(&document, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Document not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch document",
		})
	}

	response := dtos.DocumentResponse{
		ID:          document.ID,
		Src:         document.Src,
		CompanyId:   document.CompanyId,
		MediaTypeId: document.MediaTypeId,
		CreatedAt:   document.CreatedAt,
		UpdatedAt:   document.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func UpdateDocument(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var req dtos.UpdateDocumentRequest
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

	var document media_models.Document
	if err := db.First(&document, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
	}

	if req.Src != nil {
		document.Src = *req.Src
	}
	if req.CompanyId != nil {
		document.CompanyId = *req.CompanyId
	}
	if req.MediaTypeId != nil {
		document.MediaTypeId = *req.MediaTypeId
	}

	res := db.Save(&document)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Update failed"})
	}

	response := dtos.DocumentResponse{
		ID:          document.ID,
		Src:         document.Src,
		CompanyId:   document.CompanyId,
		MediaTypeId: document.MediaTypeId,
		CreatedAt:   document.CreatedAt,
		UpdatedAt:   document.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Document updated successfully",
		"data":    response,
	})
}

func DeleteDocument(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	// Use RowsAffected to check if the record exists
	res := db.Delete(&media_models.Document{}, id)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete Document",
		})
	}

	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Document not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Document deleted successfully",
	})
}