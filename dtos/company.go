package dtos

import (
	"time"

	company_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/company"
	general_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/general"
)

// Company DTOs
type CreateCompanyRequest struct {
	Name            string `json:"name" validate:"required,max=100"`
	ForeignBranchName string `json:"foreign_branch_name" validate:"required,max=100"`
	ForeignRegistrationNumber string `json:"foreign_registration_number" validate:"required,max=100"`
	TradeNames      string `json:"trade_names" validate:"max=500"`
	AuthorityName string `json:"authority_name" validate:"max=100"`
	AuthorityNumber string `json:"authority_number" validate:"max=100"`
	LocalAddress    string `json:"local_address" validate:"required,max=200"`
	ForeignAddress    string `json:"foreign_address" validate:"required,max=200"`
	Description     string `json:"description" validate:"max=600"`
	Email           string `json:"email" validate:"omitempty,email"`
	PhoneNumber     string `json:"phone_number" validate:"max=20"`
	CEOName string `json:"ceo_name" validate:"max=100"`
	CEOPhone string `json:"ceo_phone" validate:"max=20"`
	CEOEmail string `json:"ceo_email" validate:"omitempty,email"`
	CEOAddress string `json:"ceo_address" validate:"max=200"`
	CompanyTypeID   uint   `json:"type_id" validate:"required"`
	OfficeId        uint   `json:"office_id" validate:"required"`
	IsLicensed      bool   `json:"is_licensed"`
	People 		[]company_models.People `json:"people"`
	Duration        string `json:"duration" validate:"max=100"`
}

type UpdateCompanyRequest struct {
	// the pointer means that the field is optional mr nullable
	Name            *string `json:"name" validate:"omitempty,max=100"`
	ForeignBranchName *string `json:"foreign_branch_name" validate:"required,max=100"`
	ForeignRegistrationNumber *string `json:"foreign_registration_number" validate:"required,max=100"`
	TradeNames      *string `json:"trade_names" validate:"omitempty,max=200"`
	AuthorityName *string `json:"authority_name" validate:"max=100"`
	AuthorityNumber *string `json:"authority_number" validate:"max=100"`
	LocalAddress    *string `json:"local_address" validate:"omitempty,max=200"`
	ForeignAddress    *string `json:"foreign_address" validate:"required,max=200"`
	Description     *string `json:"description" validate:"omitempty,max=500"`
	Email           *string `json:"email" validate:"omitempty,email"`
	PhoneNumber     *string `json:"phone_number" validate:"omitempty,max=15"`
	IsLicensed      *bool   `json:"is_licensed"`
	CEOName *string `json:"ceo_name" validate:"max=100"`
	CEOPhone *string `json:"ceo_phone" validate:"max=20"`
	CEOEmail *string `json:"ceo_email" validate:"omitempty,email"`
	CEOAddress *string `json:"ceo_address" validate:"max=200"`
	CompanyTypeID   *uint   `json:"type_id" validate:"required"`
	OfficeId        *uint   `json:"office_id" validate:"required"`
	People 		*[]company_models.People `json:"people"`
	Duration        *string `json:"duration" validate:"omitempty,max=100"`
}

type CompanyResponse struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	TradeNames      string    `json:"trade_names"`
	AuthorityName   string    `json:"authority_name"`
	AuthorityNumber string    `json:"authority_number"`
	LocalAddress    string    `json:"local_address"`
	Description     string    `json:"description"`
	Email           string    `json:"email"`
	PhoneNumber     string    `json:"phone_number"`
	CompanyTypeID   uint      `json:"type_id"`
	CompanyType      general_models.CompanyType `json:"company_type"`
	OfficeId        uint      `json:"office_id"`
	Office      general_models.Office `json:"office"`
	IsLicensed      bool      `json:"is_licensed"`
	CEOName         string    `json:"ceo_name"`
	CEOPhone        string    `json:"ceo_phone"`
	CEOEmail        string    `json:"ceo_email"`
	CEOAddress      string    `json:"ceo_address"`
	Duration        string    `json:"duration"`
	People         []company_models.People `json:"people"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Trading Activity DTOs
type CreateTradingActivityRequest struct {
	Name string `json:"name" validate:"required,max=100"`
}

type UpdateTradingActivityRequest struct {
	Name *string `json:"name" validate:"omitempty,max=100"`
}

type TradingActivityResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Company Activity DTOs
type CreateCompanyActivityRequest struct {
	CompanyId         uint   `json:"company_id" validate:"required"`
	TradingActivityID uint   `json:"trading_activity_id" validate:"required"`
	Image             string `json:"image"`
}

type UpdateCompanyActivityRequest struct {
	Image *string `json:"image"`
}

type CompanyActivityResponse struct {
	ID                uint      `json:"id"`
	CompanyId         uint      `json:"company_id"`
	TradingActivityID uint      `json:"trading_activity_id"`
	Image             string    `json:"image"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// People DTOs
type CreatePersonRequest struct  {
	CompanyID uint   `json:"company_id" validate:"required"`
	FullName  string `json:"full_name" validate:"required,max=255"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"max=50"`
	Address   string `json:"address" validate:"max=500"`
	Role      string `json:"role" validate:"required,max=50"`
	ExtraDetails string `json:"extra_details" validate:"max=1000"`
}

type UpdatePersonRequest struct {
	FullName *string `json:"full_name" validate:"omitempty,max=255"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Phone    *string `json:"phone" validate:"omitempty,max=50"`
	Address  *string `json:"address" validate:"omitempty,max=500"`
	Role     *string `json:"role" validate:"omitempty,max=50"`
}

type PersonResponse struct {
	ID        uint      `json:"id"`
	CompanyID uint      `json:"company_id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
