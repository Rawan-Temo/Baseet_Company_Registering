package handlers

import (
	"time"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/dtos"
	auth_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/auth"
	company_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/company"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func AllUsers(c *fiber.Ctx) error {
	db := database.DB

	var users []auth_models.User
	queryArgs := c.Context().QueryArgs()
	queries := map[string][]string{}
	queryArgs.VisitAll(func(key, value []byte) {
		k := string(key)
		v := string(value)
		queries[k] = append(queries[k], v)
	})

	allowedCols := []string{"id", "username", "email", "age", "created_at", "updated_at", "deleted_at"}

	queryBuild := utils.NewQueryBuilder(db, queries, allowedCols)

	queryBuild.Paginate().Sort().Filter().LimitFields()
	if err := queryBuild.Apply().Find(&users).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "No users found",
			"error":   err.Error(),
		})

	}

	// count total matching users (filter only)
	var total int64
	utils.NewQueryBuilder(db, queries, allowedCols).Filter().Apply().Model(&auth_models.User{}).Count(&total)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total":   total,
		"status":  "success",
		"results": len(users),
		"data":    users,
	})

}

func CreateUser(c *fiber.Ctx) error {

	db := database.DB
	var req dtos.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not parse json",
			"error":   err.Error(),
		})
	}
	user := auth_models.User{
		UserName:  req.UserName,
		Password:  req.Password,
		Email:     req.Email,
		Role:      auth_models.Role(req.Role),
		CompanyId: req.CompanyId,
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not create user",
			"error":   err.Error(),
		})
	}

	response := dtos.UserResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		Role:      string(user.Role),
		CompanyId: user.CompanyId,
		Active:    user.Active,
		ExpiresAt: user.ExpiresAt,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})

}
func SingleUser(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	var user auth_models.User

	if err := db.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "User not found",
			"error":   err.Error(),
		})
	}

	response := dtos.UserResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		Role:      string(user.Role),
		CompanyId: user.CompanyId,
		Active:    user.Active,
		ExpiresAt: user.ExpiresAt,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}
func UpdateUser(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	var user auth_models.User

	if err := db.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "User not found",
			"error":   err.Error(),
		})
	}

	var req dtos.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not parse json",
			"error":   err.Error(),
		})
	}

	if req.Active != nil {
		user.Active = *req.Active
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Role != nil {
		user.Role = auth_models.Role(*req.Role)
	}

	if err := db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not update user",
			"error":   err.Error(),
		})
	}

	response := dtos.UserResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		Role:      string(user.Role),
		CompanyId: user.CompanyId,
		Active:    user.Active,
		ExpiresAt: user.ExpiresAt,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	var user auth_models.User

	if err := db.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "User not found",
			"error":   err.Error(),
		})
	}

	if err := db.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not delete user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User deleted",
	})
}

func Login(c *fiber.Ctx) error {
	
	db := database.DB
	// Parse input
	var req dtos.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Find user
	var user auth_models.User
	if err := db.Where("username = ?", req.UserName).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password")
	}

	// Validate password
	if !CheckPasswordHash(req.Password, user.Password) {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password")
	}

	// Handle expiration (only for non-admin users)
	if user.Role != auth_models.RoleAdmin && time.Now().After(*user.ExpiresAt) {
		// Try to fetch and deactivate company
		var company company_models.Company
		if err := db.First(&company, user.CompanyId).Error; err == nil {
			company.IsLicensed = false
			_ = db.Save(&company).Error
		}

		// Deactivate user
		user.Active = false
		_ = db.Save(&user).Error

		return fiber.NewError(fiber.StatusUnauthorized, "User account expired, contact admin")
	}

	// Check active status
	if !user.Active {
		return fiber.NewError(fiber.StatusUnauthorized, "User is inactive, contact admin")
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate token")
	}

	// Clean sensitive fields
	user.Password = ""

	userResponse := dtos.UserResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		Role:      string(user.Role),
		CompanyId: user.CompanyId,
		Active:    user.Active,
		ExpiresAt: user.ExpiresAt,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	loginResponse := dtos.LoginResponse{
		User:  userResponse,
		Token: token,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   loginResponse,
	})
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
