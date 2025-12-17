package dtos

import "time"

// Company DTOs
type CreateCompanyRequest struct {
	Name            string `json:"name" validate:"required,max=100"`
	TradeNames      string `json:"trade_names" validate:"max=200"`
	AuthorityNumber string `json:"authority_number" validate:"max=100"`
	LocalAddress    string `json:"local_address" validate:"required,max=200"`
	Description     string `json:"description" validate:"max=500"`
	Email           string `json:"email" validate:"omitempty,email"`
	PhoneNumber     string `json:"phone_number" validate:"max=15"`
	CompanyTypeID   uint   `json:"type_id" validate:"required"`
	OfficeId        uint   `json:"office_id" validate:"required"`
	IsLicensed      bool   `json:"is_licensed"`
	Duration        string `json:"duration" validate:"max=100"`
}

type UpdateCompanyRequest struct {
	Name            *string `json:"name" validate:"omitempty,max=100"`
	TradeNames      *string `json:"trade_names" validate:"omitempty,max=200"`
	AuthorityNumber *string `json:"authority_number" validate:"omitempty,max=100"`
	LocalAddress    *string `json:"local_address" validate:"omitempty,max=200"`
	Description     *string `json:"description" validate:"omitempty,max=500"`
	Email           *string `json:"email" validate:"omitempty,email"`
	PhoneNumber     *string `json:"phone_number" validate:"omitempty,max=15"`
	IsLicensed      *bool   `json:"is_licensed"`
	Duration        *string `json:"duration" validate:"omitempty,max=100"`
}

type CompanyResponse struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	TradeNames      string    `json:"trade_names"`
	AuthorityNumber string    `json:"authority_number"`
	LocalAddress    string    `json:"local_address"`
	Description     string    `json:"description"`
	Email           string    `json:"email"`
	PhoneNumber     string    `json:"phone_number"`
	CompanyTypeID   uint      `json:"type_id"`
	OfficeId        uint      `json:"office_id"`
	IsLicensed      bool      `json:"is_licensed"`
	Duration        string    `json:"duration"`
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
type CreatePersonRequest struct {
	CompanyID uint   `json:"company_id" validate:"required"`
	FullName  string `json:"full_name" validate:"required,max=255"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"max=50"`
	Address   string `json:"address" validate:"max=500"`
	Role      string `json:"role" validate:"required,max=50"`
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
