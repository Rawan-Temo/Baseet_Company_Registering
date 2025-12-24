package handlers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/dtos"
	auth_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/auth"
	company_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/company"
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
	allowedCols := []string{"company_id","image", "id", "created_at", "updated_at", "deleted_at", "start_date", "expiration_date"}
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
            if startDate, parseErr := utils.ParseDate(startDateStr); parseErr == nil {
                req.StartDate = startDate
            }
        }
        
        if expirationDateStr := c.FormValue("expiration_date"); expirationDateStr != "" {
            if expirationDate, parseErr := utils.ParseDate(expirationDateStr); parseErr == nil {
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
	var company company_models.Company
    if err := db.First(&company, license.CompanyId).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not fetch created license",
			"error":   err.Error(),
		})
	}

	if license.ExpirationDate.After(company.License) {

        company.License = license.ExpirationDate
          if err := db.Save(&company).Error; err != nil {
           // Log error but don't fail the entire operation
           log.Printf("Failed to update company license date: %v", err)
        }
    }
	




    response := dtos.LicenseResponse{
        ID:             license.ID,
        CompanyId:      license.CompanyId,
		Company: company,
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
//TODO fix this shit 
func UpdateLicense(c *fiber.Ctx) error {
    id := c.Params("id")
    db := database.DB

    // Get current license
    var license auth_models.License
    if err := db.First(&license, id).Error; err != nil {
        return handleNotFoundOrError(err)
    }

    var req dtos.UpdateLicenseRequest
    contentType := c.Get("Content-Type")

    // Check content type
    if strings.Contains(contentType, "multipart/form-data") {
         _, err := c.MultipartForm()
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "Failed to parse form data",
                "details": err.Error(),
            })
        }

        // Handle dates
        if startDateStr := c.FormValue("start_date"); startDateStr != "" {
			fmt.Println(startDateStr)
            startDate, err := utils.ParseDate(startDateStr)
            if err == nil {
                req.StartDate = &startDate
            }
        }

        if expirationDateStr := c.FormValue("expiration_date"); expirationDateStr != "" {
            expirationDate, err := utils.ParseDate(expirationDateStr)
            if err == nil {
                req.ExpirationDate = &expirationDate
            }
        }

        // Handle image - first check for file upload
        if _, err := c.FormFile("image"); err == nil {
            imageConfig := utils.DefaultImageConfig()
            imageConfig.UploadDir = "./uploads/licenses/"
            
            if uploadedPath, uploadErr := utils.UploadImage(c, "image", imageConfig); uploadErr == nil {
                req.Image = &uploadedPath
            }
        } else if imageStr := c.FormValue("image"); imageStr != "" {
            // Use direct URL if no file uploaded
            req.Image = &imageStr
        }
    } else {
        // Handle JSON
        if err := c.BodyParser(&req); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "Invalid JSON",
            })
        }
    }

    // Apply updates
    updates := map[string]interface{}{}
    oldImagePath := ""
	
    if req.StartDate != nil {
        updates["start_date"] = *req.StartDate
    }
    if req.ExpirationDate != nil {
        updates["expiration_date"] = *req.ExpirationDate
    }
    if req.Image != nil {
        // Save old image for cleanup
        if license.Image != nil {
            oldImagePath = *license.Image
        }
        updates["image"] = *req.Image
    }

    if len(updates) == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "No fields to update",
        })
    }

    // Update in database
    if err := db.Model(&license).Updates(updates).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Update failed",
        })
    }

    // Clean up old image
    if oldImagePath != "" && !isURL(oldImagePath) {
        utils.DeleteImage(oldImagePath, "./uploads/licenses/")
    }

    // Return updated license
    if err := db.Preload("Company").First(&license, id).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch updated license",
        })
    }


	if license.ExpirationDate.After(license.Company.License) {
         license.Company.License = license.ExpirationDate
         if err := db.Save(&license.Company).Error; err != nil {
          // Log error but don't fail the entire operation
            log.Printf("Failed to update company license date: %v", err)
         }
   }

    response := dtos.LicenseResponse{
        ID:             license.ID,
        CompanyId:      license.CompanyId,
		Company: license.Company,
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

    // Get license first (outside transaction for file cleanup later)
    var license auth_models.License
    if err := db.Preload("Company").First(&license, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "License not found",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch license",
        })
    }

    // Store image path before deletion
    var imagePath string
    if license.Image != nil {
        imagePath = *license.Image
    }
    // Use a single SQL query to handle everything
    result := db.Transaction(func(tx *gorm.DB) error {
        // 1. Delete the license
        if err := tx.Delete(&license).Error; err != nil {
            return err
        }

        // 2. Check if this was the current license and update company if needed
        // We need to check this AFTER deletion
        licenseDate := license.ExpirationDate.Truncate(24 * time.Hour)
        companyLicenseDate := license.Company.License.Truncate(24 * time.Hour)
        
        if licenseDate.Equal(companyLicenseDate) {
            // Get the max expiration date from remaining licenses
            var maxDate time.Time
            err := tx.Model(&auth_models.License{}).
                Where("company_id = ?", license.CompanyId).
                Select("COALESCE(MAX(expiration_date), '0001-01-01')").
                Scan(&maxDate).Error
            
            if err != nil {
                return err
            }

            // Update company license
            updateValue := maxDate
            if maxDate.IsZero() || maxDate.Year() == 1 { // Check for zero time
                updateValue = time.Time{}
            }
            
            return tx.Model(&license.Company).
                Update("license", updateValue).Error
        }

        return nil
    })

    if result != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to delete license",
            "details": result.Error(),
        })
    }

    // Clean up image file after successful deletion
    if imagePath != "" && !strings.HasPrefix(imagePath, "http://") && !strings.HasPrefix(imagePath, "https://") {
        if err := utils.DeleteImage(imagePath, "./uploads/licenses/"); err != nil {
            log.Printf("Warning: Failed to delete license image: %v", err)
        }
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "License deleted successfully",
    })
}


//========================== Helpers ==========================
func isURL(path string) bool {
    return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func handleNotFoundOrError(err error) error {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return fiber.NewError(fiber.StatusNotFound, "License not found")
    }
    return fiber.NewError(fiber.StatusInternalServerError, "Database error")
}