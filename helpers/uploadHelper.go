package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// func SaveAvatar(ctx *fiber.Ctx, field string, maxSize int64, destDir string) (string, error) {
// 	// Try to get the file from the form
// 	file, err := ctx.FormFile(field)
// 	if err != nil {
// 		// If no file was uploaded (which is okay), just return empty string
// 		if err == fiber.ErrUnprocessableEntity {
// 			return "", nil
// 		}
// 		// For other errors, return it
// 		return "", fmt.Errorf("failed to read file: %w", err)
// 	}

// 	// Optional: validate file size
// 	if file.Size > maxSize {
// 		return "", fmt.Errorf("file size exceeds limit of %d bytes", maxSize)
// 	}

// 	// Ensure destination directory exists
// 	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
// 		return "", fmt.Errorf("failed to create directory: %w", err)
// 	}

// 	// Validate file extension
// 	ext := strings.ToLower(filepath.Ext(file.Filename))
// 	allowedExtensions := map[string]bool{
// 		".jpg":  true,
// 		".jpeg": true,
// 		".png":  true,
// 		".gif":  true,
// 	}
// 	if !allowedExtensions[ext] {
// 		return "", fmt.Errorf("invalid file type: %s", ext)
// 	}

// 	// Save file with unique name
// 	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
// 	fullPath := filepath.Join(destDir, fileName)

// 	if err := ctx.SaveFile(file, fullPath); err != nil {
// 		return "", fmt.Errorf("failed to save file: %w", err)
// 	}

// 	return fileName, nil
// }

func SaveAvatar(ctx *fiber.Ctx, field string, maxSize int64, destDir string) (fileName string, fullPath string, err error) {
	file, err := ctx.FormFile(field)
	if err != nil {
		if err == fiber.ErrUnprocessableEntity {
			return "", "", nil
		}
		return "", "", fmt.Errorf("failed to read file: %w", err)
	}

	if file.Size > maxSize {
		return "", "", fmt.Errorf("file size exceeds limit of %d bytes", maxSize)
	}

	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return "", "", fmt.Errorf("failed to create directory: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	if !allowedExtensions[ext] {
		return "", "", fmt.Errorf("invalid file type: %s", ext)
	}

	fileName = fmt.Sprintf("%s%s", uuid.New().String(), ext)
	fullPath = filepath.Join(destDir, fileName)

	if err := ctx.SaveFile(file, fullPath); err != nil {
		return "", "", fmt.Errorf("failed to save file: %w", err)
	}

	return fileName, fullPath, nil
}
