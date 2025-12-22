package dtos

import "time"

// Office DTOs
type CreateOfficeRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type UpdateOfficeRequest struct {
	Name *string `json:"name" validate:"omitempty,min=1,max=100"`
}

type OfficeResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Company Type DTOs
type CreateCompanyTypeRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type UpdateCompanyTypeRequest struct {
	Name *string `json:"name" validate:"omitempty,min=1,max=100"`
}

type CompanyTypeResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
