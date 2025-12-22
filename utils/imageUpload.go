package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ImageUploadConfig holds configuration for image uploads
type ImageUploadConfig struct {
	UploadDir    string
	MaxFileSize  int64  // in bytes
	AllowedTypes []string
}

// DefaultImageConfig returns default configuration for image uploads
func DefaultImageConfig() ImageUploadConfig {
	return ImageUploadConfig{
		UploadDir:    "./uploads/",
		MaxFileSize:  5 * 1024 * 1024, // 5MB
		AllowedTypes: []string{"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"},
	}
}

// UploadImage handles single image file upload
func UploadImage(c *fiber.Ctx, fieldName string, config ImageUploadConfig) (string, error) {
	
	file, err := c.FormFile(fieldName)
	fmt.Print(file)
	if err != nil {
		return "", fmt.Errorf("failed to get file: %v", err)
	}

	// Validate file size
	if file.Size > config.MaxFileSize {
		return "", fmt.Errorf("file size exceeds maximum allowed size of %d bytes", config.MaxFileSize)
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	isAllowed := false

	isAllowed = slices.Contains(config.AllowedTypes , contentType)
	
	if !isAllowed {
		return "", fmt.Errorf("file type %s is not allowed. Allowed types: %v", contentType, config.AllowedTypes)
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(config.UploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d%s", timestamp, ext)
	filePath := filepath.Join(config.UploadDir, filename)

	// Save file
	fmt.Print(filePath)

	if err := c.SaveFile(file, filePath); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	// Return relative path for database storage
	return filename, nil
}

// UploadMultipleImages handles multiple image file uploads
func UploadMultipleImages(c *fiber.Ctx, fieldName string, config ImageUploadConfig) ([]string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, fmt.Errorf("failed to parse multipart form: %v", err)
	}

	files := form.File[fieldName]
	if len(files) == 0 {
		return nil, fmt.Errorf("no files found for field %s", fieldName)
	}

	var uploadedFiles []string

	for _, file := range files {
		// Validate file size
		if file.Size > config.MaxFileSize {
			return nil, fmt.Errorf("file %s size exceeds maximum allowed size of %d bytes", file.Filename, config.MaxFileSize)
		}

		// Validate file type
		contentType := file.Header.Get("Content-Type")
		isAllowed := false
		for _, allowedType := range config.AllowedTypes {
			if contentType == allowedType {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			return nil, fmt.Errorf("file type %s is not allowed for file %s. Allowed types: %v", contentType, file.Filename, config.AllowedTypes)
		}

		// Create upload directory if it doesn't exist
		if err := os.MkdirAll(config.UploadDir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create upload directory: %v", err)
		}

		// Generate unique filename
		ext := filepath.Ext(file.Filename)
		timestamp := time.Now().UnixNano()
		filename := fmt.Sprintf("%d%s", timestamp, ext)
		filePath := filepath.Join(config.UploadDir, filename)

		// Open uploaded file
		src, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open uploaded file %s: %v", file.Filename, err)
		}
		defer src.Close()

		// Create destination file
		dst, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create destination file for %s: %v", file.Filename, err)
		}
		defer dst.Close()

		// Copy file content
		if _, err := io.Copy(dst, src); err != nil {
			return nil, fmt.Errorf("failed to save file %s: %v", file.Filename, err)
		}

		uploadedFiles = append(uploadedFiles, filename)
	}

	return uploadedFiles, nil
}

// DeleteImage removes an image file from the filesystem
func DeleteImage(imagePath string, uploadDir string) error {
	if imagePath == "" {
		return nil // No file to delete
	}

	fullPath := filepath.Join(uploadDir, imagePath)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete image file %s: %v", fullPath, err)
	}

	return nil
}

// ValidateImageFile validates an image file without saving it
func ValidateImageFile(file *multipart.FileHeader, config ImageUploadConfig) error {
	// Validate file size
	if file.Size > config.MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", config.MaxFileSize)
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	isAllowed := false
	for _, allowedType := range config.AllowedTypes {
		if contentType == allowedType {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return fmt.Errorf("file type %s is not allowed. Allowed types: %v", contentType, config.AllowedTypes)
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	isValidExt := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			isValidExt = true
			break
		}
	}
	if !isValidExt {
		return fmt.Errorf("file extension %s is not allowed. Allowed extensions: %v", ext, allowedExts)
	}

	return nil
}