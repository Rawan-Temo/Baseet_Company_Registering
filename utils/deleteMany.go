package utils

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Ids struct {
	IDs []uint `json:"ids"`
}

func DeleteMany(db *gorm.DB, model interface{}) fiber.Handler {
	return func (c *fiber.Ctx)error{
		var ids Ids
		if err:= c.BodyParser(&ids); err !=nil{
			return c.Status(400).JSON(fiber.Map{
				"status":  "fail",
				"message": "could not parse ids",
				"error":   err.Error(),
			})
		}

		if err := db.Where("id IN ?", ids.IDs).Delete(&model).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "fail",
				"message": "Could not delete records",
				"error":   err.Error(),
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Records deleted successfully",	
		})
		}
		
	}

