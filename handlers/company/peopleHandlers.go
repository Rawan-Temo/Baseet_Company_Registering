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

func GetAllPeople(c *fiber.Ctx) error {
	db := database.DB
	var total int64
	var people []company_models.People
	queryArgs := c.Context().QueryArgs()
	queries := map[string][]string{}

	queryArgs.VisitAll(func(key, value []byte) {
		k := string(key)
		v := string(value)
		queries[k] = append(queries[k], v)
	})
	allowedCols := []string{"company_id", "full_name", "email", "phone", "role", "ID", "created_at", "updated_at", "deleted_at"}
	queryBuilder := utils.NewQueryBuilder(db, queries, allowedCols)
	queryBuilder.Filter().Sort().LimitFields().Paginate()
	if err := queryBuilder.Apply().Find(&people).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch people",
			"details": err.Error(),
		})
	}
	// count total people
	if err := utils.NewQueryBuilder(db, queries, allowedCols).Filter().Apply().Model(&company_models.People{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to count people",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"total":   total,
		"results": len(people),
		"data":    people,
	})
}

func CreatePerson(c *fiber.Ctx) error {
	db := database.DB
	var req dtos.CreatePersonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not parse json",
			"error":   err.Error(),
		})
	}

	person := company_models.People{
		CompanyID: req.CompanyID,
		FullName:  req.FullName,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		Role:      req.Role,
	}

	if err := db.Create(&person).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not create person",
			"error":   err.Error(),
		})
	}

	response := dtos.PersonResponse{
		ID:        person.ID,
		CompanyID: person.CompanyID,
		FullName:  person.FullName,
		Email:     person.Email,
		Phone:     person.Phone,
		Address:   person.Address,
		Role:      person.Role,
		CreatedAt: person.CreatedAt,
		UpdatedAt: person.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func GetPersonByID(c *fiber.Ctx) error {
	db := database.DB
	var person company_models.People
	id := c.Params("id")
	if err := db.First(&person, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Person not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch person",
		})
	}

	response := dtos.PersonResponse{
		ID:        person.ID,
		CompanyID: person.CompanyID,
		FullName:  person.FullName,
		Email:     person.Email,
		Phone:     person.Phone,
		Address:   person.Address,
		Role:      person.Role,
		CreatedAt: person.CreatedAt,
		UpdatedAt: person.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func UpdatePerson(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var req dtos.UpdatePersonRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}

	var person company_models.People
	if err := db.First(&person, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Person not found"})
	}

	if req.FullName != nil {
		person.FullName = *req.FullName
	}
	if req.Email != nil {
		person.Email = *req.Email
	}
	if req.Phone != nil {
		person.Phone = *req.Phone
	}
	if req.Address != nil {
		person.Address = *req.Address
	}
	if req.Role != nil {
		person.Role = *req.Role
	}

	res := db.Save(&person)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Update failed"})
	}

	response := dtos.PersonResponse{
		ID:        person.ID,
		CompanyID: person.CompanyID,
		FullName:  person.FullName,
		Email:     person.Email,
		Phone:     person.Phone,
		Address:   person.Address,
		Role:      person.Role,
		CreatedAt: person.CreatedAt,
		UpdatedAt: person.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Person updated successfully",
		"data":    response,
	})
}

func DeletePerson(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	// Use RowsAffected to check if the record exists
	res := db.Delete(&company_models.People{}, id)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete Person",
		})
	}

	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Person not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Person deleted successfully",
	})
}