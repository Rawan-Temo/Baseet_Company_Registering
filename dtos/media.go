package dtos

import "time"

// Media Type DTOs
type CreateMediaTypeRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"max=500"`
}

type UpdateMediaTypeRequest struct {
	Name        *string `json:"name" validate:"omitempty,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

type MediaTypeResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Image DTOs
type CreateImageRequest struct {
	Src         string `json:"src" validate:"required,max=100"`
	CompanyId   uint   `json:"company_id" validate:"required"`
	MediaTypeId int    `json:"media_type_id" validate:"required"`
}

type UpdateImageRequest struct {
	Src         *string `json:"src" validate:"omitempty,max=100"`
	CompanyId   *uint   `json:"company_id"`
	MediaTypeId *int    `json:"media_type_id"`
}

type ImageResponse struct {
	ID          uint      `json:"id"`
	Src         string    `json:"src"`
	CompanyId   uint      `json:"company_id"`
	MediaTypeId int       `json:"media_type_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Document DTOs
type CreateDocumentRequest struct {
	Src         string `json:"src" validate:"required,max=100"`
	CompanyId   uint   `json:"company_id" validate:"required"`
	MediaTypeId int    `json:"media_type_id" validate:"required"`
}

type UpdateDocumentRequest struct {
	Src         *string `json:"src" validate:"omitempty,max=100"`
	CompanyId   *uint   `json:"company_id"`
	MediaTypeId *int    `json:"media_type_id"`
}

type DocumentResponse struct {
	ID          uint      `json:"id"`
	Src         string    `json:"src"`
	CompanyId   uint      `json:"company_id"`
	MediaTypeId int       `json:"media_type_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
