package dtos

import "time"

// User DTOs
type CreateUserRequest struct {
	FullName  string `json:"full_name" validate:"required,min=3,max=150"`
	UserName  string `json:"username" validate:"required,min=3,max=100"`
	Password  string `json:"password" validate:"required,min=6"`
	Email     string `json:"email" validate:"required,email"`
	Role      string `json:"role" validate:"required,oneof=admin user"`
	CompanyId *uint  `json:"company_id"`
}
type UserTokenClaim struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Exp  int64  `json:"exp"`
	Iat  int64  `json:"iat"`
}
type UpdateUserRequest struct {
	FullName  *string `json:"full_name" validate:"required,min=3,max=150"`
	Email  *string `json:"email" validate:"omitempty,email"`
	Password  *string `json:"password" validate:"omitempty,min=6"`
	Active *bool   `json:"active"`
}

type UserResponse struct {
	ID        uint       `json:"id"`
	FullName  string     `json:"full_name"`
	UserName  string     `json:"username"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	CompanyId *uint      `json:"company_id"`
	Active    bool       `json:"active"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// License DTOs
type CreateLicenseRequest struct {
	CompanyId      uint      `json:"company_id" validate:"required"`
	StartDate      time.Time `json:"start_date" validate:"required"`
	ExpirationDate time.Time `json:"expiration_date" validate:"required"`
	Image          *string    `json:"image"`
}

type UpdateLicenseRequest struct {
	StartDate      *time.Time `json:"start_date"`
	ExpirationDate *time.Time `json:"expiration_date"`
	Image          *string    `json:"image"`
}

type LicenseResponse struct {
	ID             uint      `json:"id"`
	CompanyId      uint      `json:"company_id"`
	StartDate      time.Time `json:"start_date"`
	ExpirationDate time.Time `json:"expiration_date"`
	Image          *string    `json:"image"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Login DTOs
type LoginRequest struct {
	UserName string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}
