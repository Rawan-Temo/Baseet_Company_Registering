package helpers

import (
	"time"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/dtos"
	auth_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/auth"
	company_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/company"
	"gorm.io/gorm"
)

func RegisterCompanyAndUser(userReq dtos.CreateUserRequest, companyReq dtos.CreateCompanyRequest, tx *gorm.DB) (error, company_models.Company, auth_models.User) {

	//  company creation
	defaultLicense := time.Now().AddDate(0, 1, 0) // default to 30 days from now

	company := company_models.Company{
		Name:                      companyReq.Name,
		ForeignBranchName:         companyReq.ForeignBranchName,
		ForeignRegistrationNumber: companyReq.ForeignRegistrationNumber,
		TradeNames:                companyReq.TradeNames,
		AuthorityName:             companyReq.AuthorityName,
		AuthorityNumber:           companyReq.AuthorityNumber,
		LocalAddress:              companyReq.LocalAddress,
		Description:               companyReq.Description,
		Email:                     companyReq.Email,
		PhoneNumber:               companyReq.PhoneNumber,
		CompanyCategory:           company_models.CompanyCategory(companyReq.CompanyCategory),
		OfficeId:                  companyReq.OfficeId,
		License:                   defaultLicense,
		People:                    companyReq.People,
		CEOName:                   companyReq.CEOName,
		CEOPhone:                  companyReq.CEOPhone,
		CEOEmail:                  companyReq.CEOEmail,
		CEOAddress:                companyReq.CEOAddress,
		Duration:                  companyReq.Duration,
	}

	if err := tx.Create(&company).Error; err != nil {
		return err, company_models.Company{}, auth_models.User{}
	}
	companyId := &company.ID
	// User creation
	user := auth_models.User{
		FullName:  userReq.FullName,
		UserName:  userReq.UserName,
		Password:  userReq.Password,
		Email:     userReq.Email,
		Role:      auth_models.Role(userReq.Role),
		CompanyId: companyId,
	}

	if err := tx.Create(&user).Error; err != nil {
		return err, company_models.Company{}, auth_models.User{}
	}
	return nil, company, user
}
