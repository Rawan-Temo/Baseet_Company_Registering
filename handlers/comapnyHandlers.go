package handlers

import (
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/models"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
)

func AllCompanies(c *fiber.Ctx) error {
	db := database.DB
	var companies []models.Company
	queryArgs := c.Context().QueryArgs()
	queries := map[string][]string{}
	queryArgs.VisitAll(func(key, value []byte) {
		k := string(key)
		v := string(value)
		queries[k] = append(queries[k], v)
	})
	allowedCols := []string{"id" , "name" , "address" , "phone" , "email" , "created_at" , "updated_at"}
	apiFeatures := utils.NewQueryBuilder(db ,queries ,allowedCols)

	apiFeatures.Filter().Sort().LimitFields().Paginate()


	if err := apiFeatures.Apply().Find(&companies).Error ; err!= nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "Failed to fetch companies",
			"details": err.Error(),
		})
	}
	// count total matching companies (filter only)
	var total int64
	utils.NewQueryBuilder(db, queries, allowedCols).Filter().Apply().Model(&models.Company{}).Count(&total)


	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"results": len(companies),
		"total": total,
		"data": companies,
	})
}
func CreateCompany(c *fiber.Ctx) error {

	db := database.DB
	var company models.Company
	if err := c.BodyParser(&company); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not parse json",
			"error":   err.Error(),
		})
	}
	if err := db.Create(&company).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create company",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Company created successfully",
		"data":    company,
	})
}

func SingleCompany(c *fiber.Ctx) error {
	return nil

}
func UpdateCompany(c *fiber.Ctx) error {
	return nil

}
func DeleteCompany(c *fiber.Ctx) error {
	return nil

}