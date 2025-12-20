package handlers

import (
	"fmt"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/dtos"
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
	allowedCols := []string{"id", "type_id" , "company_id", "type" , "office", "people", "trade_names", "authority_name", "authority_number", "name", "address" , "ceo_name" , "ceo_email" , "ceo_phone" , "ceo_address", "phone", "email", "created_at", "updated_at"}
	apiFeatures := utils.NewQueryBuilder(db, queries, allowedCols)

	apiFeatures.Filter().Sort().LimitFields().Paginate()

	if err := apiFeatures.Apply().Preload("CompanyType").Preload("Office").Find(&companies).Error; err != nil {

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
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not parse json",
			"error":   err.Error(),
		})
	}

	// people := []company_models.People{}
	// for _, personReq := range req.People{
	// 	people = append(people, company_models.People{
	// 		FullName: personReq.FullName,
	// 		Email: personReq.Email,
	// 		Phone: personReq.Phone,
	// 		Address: personReq.Address,
	// 		Role: personReq.Role,
	// 		ExtraDetails: personReq.ExtraDetails,
	// 	})
		
	// }
	company := company_models.Company{
		Name:            req.Name,
		TradeNames:      req.TradeNames,
		AuthorityNumber: req.AuthorityNumber,
		LocalAddress:    req.LocalAddress,
		Description:     req.Description,
		Email:           req.Email,
		PhoneNumber:     req.PhoneNumber,
		CompanyTypeID:   req.CompanyTypeID,
		OfficeId:        req.OfficeId,
		IsLicensed:      req.IsLicensed,
		People:req.People,
		CEOName:req.CEOName ,
		CEOPhone: req.CEOPhone,
		CEOEmail:req.CEOEmail ,
		CEOAddress: req.CEOAddress,
		Duration:        req.Duration,
	}

	fmt.Printf("people : %v ", company)
	if err := db.Create(&company).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create company",
			"error":   err.Error(),
		})
	}

	response := dtos.CompanyResponse{
		ID:              company.ID,
		Name:            company.Name,
		TradeNames:      company.TradeNames,
		AuthorityNumber: company.AuthorityNumber,
		LocalAddress:         company.LocalAddress,
		Description:     company.Description,
		Email:           company.Email,
		PhoneNumber:     company.PhoneNumber,
		CompanyTypeID:   company.CompanyTypeID,
		OfficeId:        company.OfficeId,
		IsLicensed:      company.IsLicensed,
		CEOName:company.CEOName ,
		CEOPhone: company.CEOPhone,
		CEOEmail:company.CEOEmail ,
		CEOAddress: company.CEOAddress,
		Duration:        company.Duration,
		CreatedAt:       company.CreatedAt,
		UpdatedAt:       company.UpdatedAt,
	}

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
	if err := db.Preload("CompanyType").Preload("Office").First(&company, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Company not found",
			"error":   err.Error(),
		})
	}

	response := dtos.CompanyResponse{
		ID:              company.ID,
		Name:            company.Name,
		TradeNames:      company.TradeNames,
		AuthorityNumber: company.AuthorityNumber,
		LocalAddress:    company.LocalAddress,
		Description:     company.Description,
		Email:           company.Email,
		PhoneNumber:     company.PhoneNumber,
		CompanyTypeID:   company.CompanyTypeID,
		OfficeId:        company.OfficeId,
		IsLicensed:      company.IsLicensed,
		Duration:        company.Duration,
		CreatedAt:       company.CreatedAt,
		UpdatedAt:       company.UpdatedAt,
	}

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

	// Update only provided fields
	if req.Name != nil {
		company.Name = *req.Name
	}
	if req.TradeNames != nil {
		company.TradeNames = *req.TradeNames
	}
	if req.AuthorityNumber != nil {
		company.AuthorityNumber = *req.AuthorityNumber
	}
	if req.LocalAddress != nil {
		company.LocalAddress = *req.LocalAddress
	}
	if req.Description != nil {
		company.Description = *req.Description
	}
	if req.Email != nil {
		company.Email = *req.Email
	}
	if req.PhoneNumber != nil {
		company.PhoneNumber = *req.PhoneNumber
	}
	if req.IsLicensed != nil {
		company.IsLicensed = *req.IsLicensed
	}
	if req.Duration != nil {
		company.Duration = *req.Duration
	}

	if err := db.Save(&company).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not update company",
			"error":   err.Error(),
		})
	}

	response := dtos.CompanyResponse{
		ID:              company.ID,
		Name:            company.Name,
		TradeNames:      company.TradeNames,
		AuthorityNumber: company.AuthorityNumber,
		LocalAddress:    company.LocalAddress,
		Description:     company.Description,
		Email:           company.Email,
		PhoneNumber:     company.PhoneNumber,
		CompanyTypeID:   company.CompanyTypeID,
		OfficeId:        company.OfficeId,
		IsLicensed:      company.IsLicensed,
		Duration:        company.Duration,
		CreatedAt:       company.CreatedAt,
		UpdatedAt:       company.UpdatedAt,
	}

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
