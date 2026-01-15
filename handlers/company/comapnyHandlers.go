package handlers

import (
	"fmt"
	"time"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/dtos"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/helpers"
	company_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/company"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
)

func AllCompanies(c *fiber.Ctx) error {
	db := database.DB
	var companies []company_models.Company
	queryArgs := c.Context().QueryArgs()
	queries := map[string][]string{}
	queryArgs.VisitAll(func(key, value []byte) {
		k := string(key)
		v := string(value)
		queries[k] = append(queries[k], v)
	})
	allowedCols := []string{"id", "company_category", "licesnse", "type", "office", "people", "trade_names", "authority_name", "authority_number", "name", "address", "ceo_name", "ceo_email", "ceo_phone", "ceo_address", "phone", "email", "created_at", "updated_at"}
	apiFeatures := utils.NewQueryBuilder(db, queries, allowedCols)

	apiFeatures.Filter().Sort().LimitFields().Paginate()

	if err := apiFeatures.Apply().Preload("Office").Find(&companies).Error; err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch companies",
			"details": err.Error(),
		})
	}
	// count total matching companies (filter only)
	var total int64
	utils.NewQueryBuilder(db, queries, allowedCols).Filter().Apply().Model(&company_models.Company{}).Count(&total)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"results": len(companies),
		"total":   total,
		"data":    companies,
	})
}
func CreateCompany(c *fiber.Ctx) error {

	db := database.DB
	var req dtos.CreateCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not parse json",
			"error":   err.Error(),
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "validation error",
			"error":   err,
		})
	}
	defaultLicense := time.Now().AddDate(0, 1, 0) // default to 30 days from now

	company := company_models.Company{
		Name:            req.Name,
		TradeNames:      req.TradeNames,
		AuthorityNumber: req.AuthorityNumber,
		LocalAddress:    req.LocalAddress,
		Description:     req.Description,
		Email:           req.Email,
		PhoneNumber:     req.PhoneNumber,
		CompanyCategory: company_models.CompanyCategory(req.CompanyCategory),
		OfficeId:        req.OfficeId,
		License:         defaultLicense,
		People:          req.People,
		CEOName:         req.CEOName,
		CEOPhone:        req.CEOPhone,
		CEOEmail:        req.CEOEmail,
		CEOAddress:      req.CEOAddress,
		Duration:        req.Duration,
	}

	if err := db.Create(&company).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create company",
			"error":   err.Error(),
		})
	}

	response := GetCompanyResponse(company)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Company created successfully",
		"data":    response,
	})
}

func SingleCompany(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")

	var company company_models.Company
	if err := db.Preload("Office").Preload("People").First(&company, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Company not found",
			"error":   err.Error(),
		})
	}

	response := GetCompanyResponse(company)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}
func UpdateCompany(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")

	var company company_models.Company
	if err := db.First(&company, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Company not found",
			"error":   err.Error(),
		})
	}

	// Parse dynamic JSON body into a DTO
	var req dtos.UpdateCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "validation error",
			"error":   err,
		})
	}
	utils.UpdateStruct(&company, &req)
	// // Update only provided fields
	// the utility functions basically does this part for us
	// if req.Name != nil {
	// 	company.Name = *req.Name
	// }

	if err := db.Save(&company).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not update company",
			"error":   err.Error(),
		})
	}
	response := GetCompanyResponse(company)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}
func DeleteCompany(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")

	var company company_models.Company
	if err := db.First(&company, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Company not found",
			"error":   err.Error(),
		})
	}

	if err := db.Delete(&company).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not delete company",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Company deleted successfully",
	})
}
func RegisterNewCompany(c *fiber.Ctx) error {
	db := database.DB
	var req dtos.RegisterCompanyRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not parse json",
			"error":   err.Error(),
		})
	}
	fmt.Println(req.Company)

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "validation error",
			"error":   err,
		})
	}

	// Create Company
	tx := db.Begin()
	err, company, user := helpers.RegisterCompanyAndUser(req.User, req.Company, tx)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not register company",
			"error":   err.Error(),
		})
	}
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not commit transaction",
			"error":   err.Error(),
		})
	}
	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"company": company,
		"user":    user,
	})
}
func GetCompanyResponse(company company_models.Company) dtos.CompanyResponse {
	return dtos.CompanyResponse{
		ID:              company.ID,
		Name:            company.Name,
		TradeNames:      company.TradeNames,
		AuthorityNumber: company.AuthorityNumber,
		LocalAddress:    company.LocalAddress,
		Description:     company.Description,
		Email:           company.Email,
		PhoneNumber:     company.PhoneNumber,
		CompanyCategory: string(company.CompanyCategory),
		OfficeId:        company.OfficeId,
		Office:          company.Office,
		License:         company.License,
		CEOName:         company.CEOName,
		CEOPhone:        company.CEOPhone,
		CEOEmail:        company.CEOEmail,
		CEOAddress:      company.CEOAddress,
		Duration:        company.Duration,
		People:          company.People,
		CreatedAt:       company.CreatedAt,
		UpdatedAt:       company.UpdatedAt,
	}
}
