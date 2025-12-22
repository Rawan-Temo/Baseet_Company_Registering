package handlers

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/dtos"
	auth_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/auth"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllLicenses(c *fiber.Ctx) error {
	db := database.DB
	var total int64
	var licenses []auth_models.License
	queryArgs := c.Context().QueryArgs()
	queries := map[string][]string{}

	queryArgs.VisitAll(func(key, value []byte) {
		k := string(key)
		v := string(value)
		queries[k] = append(queries[k], v)
	})
	allowedCols := []string{"company_id", "ID", "created_at", "updated_at", "deleted_at", "start_date", "expiration_date"}
	queryBuilder := utils.NewQueryBuilder(db, queries, allowedCols)
	queryBuilder.Filter().Sort().LimitFields().Paginate()
	if err := queryBuilder.Apply().Find(&licenses).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch licenses",
			"details": err.Error(),
		})
	}
	// count total licenses
	if err := utils.NewQueryBuilder(db, queries, allowedCols).Filter().Apply().Model(&auth_models.License{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to count licenses",
			"details": err.Error(),
		})
	}

	// Convert licenses to response DTOs
	var licenseResponses []dtos.LicenseResponse
	for _, license := range licenses {
		licenseResponses = append(licenseResponses, dtos.LicenseResponse{
			ID:             license.ID,
			CompanyId:      license.CompanyId,
			StartDate:      license.StartDate,
			ExpirationDate: license.ExpirationDate,
			Image:          license.Image,
			CreatedAt:      license.CreatedAt,
			UpdatedAt:      license.UpdatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"total":   total,
		"results": len(licenseResponses),
		"data":    licenseResponses,
	})
}

func CreateLicense(c *fiber.Ctx) error {
    db := database.DB
    var req dtos.CreateLicenseRequest
    
    // First, check if this is multipart form data
    _, err := c.MultipartForm()
    if err == nil {
        // Handle multipart form data
        companyIdStr := c.FormValue("company_id")
        if companyIdStr != "" {
            if companyId, parseErr := strconv.ParseUint(companyIdStr, 10, 32); parseErr == nil {
                req.CompanyId = uint(companyId)
            }
        }
        
        if startDateStr := c.FormValue("start_date"); startDateStr != "" {
            if startDate, parseErr := time.Parse("2006-01-02", startDateStr); parseErr == nil {
                req.StartDate = startDate
            }
        }
        
        if expirationDateStr := c.FormValue("expiration_date"); expirationDateStr != "" {
            if expirationDate, parseErr := time.Parse("2006-01-02", expirationDateStr); parseErr == nil {
                req.ExpirationDate = expirationDate
            }
        }
        
        // Handle image file upload
        _, err := c.FormFile("image")
        if err == nil {
            imageConfig := utils.DefaultImageConfig()
            imageConfig.UploadDir = "./uploads/licenses/"
            
            // Create directory if it doesn't exist
            if _, err := os.Stat(imageConfig.UploadDir); os.IsNotExist(err) {
                os.MkdirAll(imageConfig.UploadDir, 0755)
            }

            if uploadedPath, uploadErr := utils.UploadImage(c, "image", imageConfig); uploadErr != nil {
                return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                    "status":  "fail",
                    "message": "Image upload failed",
                    "error":   uploadErr.Error(),
                })
            } else {
                // Initialize Image pointer if nil
                if req.Image == nil {
                    imageStr := uploadedPath
                    req.Image = &imageStr
                } else {
                    *req.Image = uploadedPath
                }
            }
        } else if imageStr := c.FormValue("image"); imageStr != "" {
            // Use direct image URL if provided
            req.Image = &imageStr
        }
        
        // Validate required fields
        if req.CompanyId == 0 {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "fail",
                "message": "Company ID is required",
            })
        }
        
        if req.StartDate.Equal((time.Time{})) {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "fail",
                "message": "Start date is required",
            })
        }

        if req.ExpirationDate.Equal((time.Time{})) {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "fail",
                "message": "Expiration date is required",
            })
        }
        
    } else {
        // Try to parse as JSON
        if err := c.BodyParser(&req); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "fail",
                "message": "Invalid request format",
                "error":   "Expected JSON or multipart form data",
            })
        }
    }
    
    // Now create the license
    license := auth_models.License{
        CompanyId:      req.CompanyId,
        StartDate:      req.StartDate,
        ExpirationDate: req.ExpirationDate,
        Image:          req.Image,
    }
    
    if err := db.Create(&license).Error; err != nil {
        if strings.Contains(err.Error(), "duplicate key value") {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "fail",
                "message": "License already exists for this company",
                "error":   err.Error(),
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "fail",
            "message": "Could not create license",
            "error":   err.Error(),
        })
    }
    
    response := dtos.LicenseResponse{
        ID:             license.ID,
        CompanyId:      license.CompanyId,
        StartDate:      license.StartDate,
        ExpirationDate: license.ExpirationDate,
        Image:          license.Image,
        CreatedAt:      license.CreatedAt,
        UpdatedAt:      license.UpdatedAt,
    }
    
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "status": "success",
        "data":   response,
    })
}


func GetLicenseByID(c *fiber.Ctx) error {
	db := database.DB
	var license auth_models.License
	id := c.Params("id")
	if err := db.First(&license, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "License not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch license",
		})
	}

	response := dtos.LicenseResponse{
		ID:             license.ID,
		CompanyId:      license.CompanyId,
		StartDate:      license.StartDate,
		ExpirationDate: license.ExpirationDate,
		Image:          license.Image,
		CreatedAt:      license.CreatedAt,
		UpdatedAt:      license.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func UpdateLicense(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	// First, get the current license to check for existing image
	var currentLicense auth_models.License
	if err := db.First(&currentLicense, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "License not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch license"})
	}

	var req dtos.UpdateLicenseRequest

	// Try to parse JSON first
	if err := c.BodyParser(&req); err != nil {
		// If JSON parsing fails, try multipart form (for file uploads)
		startDate, _ := time.Parse("2006-01-02", c.FormValue("start_date"))
		if startDate != (time.Time{}) {
			req.StartDate = &startDate
		}
		expirationDate, _ := time.Parse("2006-01-02", c.FormValue("expiration_date"))
		if expirationDate != (time.Time{}) {
			req.ExpirationDate = &expirationDate
		}
		if image := c.FormValue("image"); image != "" {
			req.Image = &image
		}

		// Handle image upload if file is provided
		if _, err := c.FormFile("image_file"); err == nil {
			imageConfig := utils.DefaultImageConfig()
			imageConfig.UploadDir = "./uploads/licenses/"
			if uploadedPath, uploadErr := utils.UploadImage(c, "image_file", imageConfig); uploadErr != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "fail",
					"message": "Image upload failed",
					"error":   uploadErr.Error(),
				})
			} else {
				req.Image = &uploadedPath
				// Delete old image if it exists
				if *currentLicense.Image != "" && *req.Image != *currentLicense.Image {
					os.Remove(*currentLicense.Image)
				}
			}
		}
	}

	sanitized := map[string]interface{}{}

	if req.StartDate != nil {
		sanitized["start_date"] = req.StartDate
	}
	if req.ExpirationDate != nil {
		sanitized["expiration_date"] = req.ExpirationDate
	}
	if req.Image != nil {
		sanitized["image"] = req.Image
		// Delete old image if a new image is being set and it's different
		if *currentLicense.Image != "" && *req.Image != *currentLicense.Image {
			os.Remove(*currentLicense.Image)
		}
	}

	if len(sanitized) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No valid fields"})
	}

	var license auth_models.License
	res := db.Model(&license).Where("id = ?", id).Updates(sanitized)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Update failed"})
	}
	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "License not found"})
	}
	if err := db.First(&license, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Fetch after update failed"})
	}

	response := dtos.LicenseResponse{
		ID:             license.ID,
		CompanyId:      license.CompanyId,
		StartDate:      license.StartDate,
		ExpirationDate: license.ExpirationDate,
		Image:          license.Image,
		CreatedAt:      license.CreatedAt,
		UpdatedAt:      license.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "License updated successfully",
		"data":    response,
	})
}

func DeleteLicense(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	// First, get the license to check for image
	var license auth_models.License
	if err := db.First(&license, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "License not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch license",
		})
	}

	// Delete the associated image file if it exists
	if *license.Image != "" {
		os.Remove(*license.Image)
	}

	// Use RowsAffected to check if the record exists
	res := db.Delete(&auth_models.License{}, id)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete License",
		})
	}

	if res.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "License not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "License deleted successfully",
	})
}





