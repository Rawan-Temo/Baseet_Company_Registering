package handlers

import (
	"time"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/dtos"
	auth_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/auth"
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

	allowedCols := []string{"id", "full_name", "username", "email", "age", "created_at", "updated_at", "deleted_at"}

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
	user := auth_models.User{
		FullName:  req.FullName,
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
		FullName:  user.FullName,
		Email:     user.Email,
		Role:      string(user.Role),
		CompanyId: user.CompanyId,
		Active:    user.Active,
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
		FullName:  user.FullName,
		UserName:  user.UserName,
		Email:     user.Email,
		Role:      string(user.Role),
		CompanyId: user.CompanyId,
		Active:    user.Active,
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
	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "validation error",
			"error":   err,
		})
	}

	utils.UpdateStruct(&user, &req)

	if err := db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not update user",
			"error":   err.Error(),
		})
	}

	response := dtos.UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		UserName:  user.UserName,
		Email:     user.Email,
		Role:      string(user.Role),
		CompanyId: user.CompanyId,
		Active:    user.Active,
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
	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "validation error",
			"error":   err,
		})
	}
	// Find user
	var user auth_models.User

	if err := db.Preload("Company").Where("username = ? And deleted_at IS NULL", req.UserName).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password")
	}

	// Validate password
	if !CheckPasswordHash(req.Password, user.Password) {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password")
	}

	// Handle expiration (only for non-admin users)
	if user.Role != auth_models.RoleAdmin && time.Now().After(user.Company.License) {
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

	userResponse := dtos.UserResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		FullName:  user.FullName,
		Email:     user.Email,
		Role:      string(user.Role),
		CompanyId: user.CompanyId,
		Active:    user.Active,
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

func GetUserFromToken(c *fiber.Ctx) error {
	database := database.DB
	currentUser := c.Locals("currentUser").(dtos.UserTokenClaim)

	var user auth_models.User
	database.Where("id = ?", currentUser.UserID).First(&user)
	var userResponse = dtos.UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		UserName:  user.UserName,
		Email:     user.Email,
		Role:      string(user.Role),
		CompanyId: user.CompanyId,
		Active:    user.Active,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   userResponse,
	})
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
