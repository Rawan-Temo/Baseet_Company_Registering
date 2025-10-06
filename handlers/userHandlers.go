package handlers

import (
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/models"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func AllUsers(c *fiber.Ctx) error {
	db := database.DB

	var users []models.User
	queryArgs := c.Context().QueryArgs()
	queries := map[string][]string{}
	queryArgs.VisitAll(func(key, value []byte) {
		k := string(key)
		v := string(value)
		queries[k] = append(queries[k], v)
	})

	allowedCols := []string{"id", "user_name", "email", "age", "created_at", "updated_at", "deleted_at"}

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
	utils.NewQueryBuilder(db, queries, allowedCols).Filter().Apply().Model(&models.User{}).Count(&total)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total":   total,
		"status":  "success",
		"results": len(users),
		"data":    users,
	})

}

func CreateUser(c *fiber.Ctx) error {

	db := database.DB
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not parse json",
			"error":   err.Error(),
		})
	}

	// validation and password hashing
	if err := user.BeforeCreate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "validation failed",
			"error":   err.Error(),
		})
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": "could not create user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})

}
func SingleUser(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "User not found",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "User not found",
			"error":   err.Error(),
		})
	}

	var updateData models.User
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not parse json",
			"error":   err.Error(),
		})
	}

	// Only update allowed fields
	user.UserName = updateData.UserName
	user.Email = updateData.Email
	// If password is present, hash it
	if updateData.Password != "" {
		user.Password = updateData.Password
		if err := user.BeforeCreate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": "Password validation failed",
				"error":   err.Error(),
			})
		}
	}

	if err := db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not update user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	var user models.User

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
	type LoginInput struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	var input LoginInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not parse JSON",
			"error":   err.Error(),
		})
	}

	var user models.User

	if err := db.Where("user_name = ?", input.UserName).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid username or password",
		})
	}

	if !CheckPasswordHash(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid username or password",
		})
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not generate token",
			"error":   err.Error(),
		})
	}

	// Remove password before sending response
	user.Password = ""

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
