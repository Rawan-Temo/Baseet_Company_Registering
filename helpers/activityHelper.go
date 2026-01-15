package helpers

import (
	"errors"
	"strconv"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/dtos"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
)

func ValidateMultiPartFormActivity(c *fiber.Ctx, req *dtos.CreateCompanyActivityRequest) error {
	_, err := c.MultipartForm()
	if err == nil {
		companyIdStr := c.FormValue("company_id")
		tradingActivityIDStr := c.FormValue("trading_activity_id")
		imageFile, err := c.FormFile("image")
		if err != nil {
			return err
		}

		companyId, err := strconv.ParseUint(companyIdStr, 10, 64)
		if err != nil {
			return err
		}
		req.CompanyId = uint(companyId)

		tradingActivityID, err := strconv.ParseUint(tradingActivityIDStr, 10, 64)
		if err != nil {
			return err
		}
		req.TradingActivityID = uint(tradingActivityID)

		fileName := utils.GenerateFileName(imageFile)
		req.Image = fileName
	}

	// Validate required fields
	if req.CompanyId == 0 {
		return errors.New("company_id is required")
	}
	if req.TradingActivityID == 0 {
		return errors.New("trading_activity_id is required")
	}
	return nil
}
