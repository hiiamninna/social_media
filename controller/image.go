package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"social_media/collections"
	"social_media/library"

	"github.com/gofiber/fiber/v2"
)

const (
	maxFileSize = 1024 * 1024 * 2
	minFileSize = 1024 * 10
	JPG         = ".jpg"
	JPEG        = ".jpeg"
)

var allowedFormats = map[string]bool{
	JPG:  true,
	JPEG: true,
}

type Image struct {
	S3 library.S3
}

func NewImageController(s3 library.S3) Image {
	return Image{
		S3: s3,
	}
}

func (c Image) ImageUpload(ctx *fiber.Ctx) (int, string, interface{}, interface{}, error) {

	var url string
	file, err := ctx.FormFile("file")
	if err != nil {
		return http.StatusBadRequest, "upload image fail", nil, nil, err
	}

	// Check file size
	if file.Size < minFileSize {
		return http.StatusBadRequest, "file size must be at least 10 KB", nil, nil, fmt.Errorf("file size is less 10 KB")
	}
	if file.Size > maxFileSize {
		return http.StatusBadRequest, "file size must be no more than 2 MB", nil, nil, fmt.Errorf("file size is more than 2 MB")
	}
	ext := filepath.Ext(file.Filename)
	if !allowedFormats[ext] {
		return http.StatusBadRequest, "file type must be JPG or JPEG", nil, nil, fmt.Errorf("file type not JPG or JPEG")
	}

	name := generateUUID()

	uploadedFile, err := file.Open()
	if err != nil {
		return http.StatusInternalServerError, "upload image fail", nil, nil, err
	}
	url, err = c.S3.UploadFile(uploadedFile, name+ext)
	if err != nil {
		return http.StatusInternalServerError, "upload image fail", nil, nil, err
	}

	return http.StatusOK, "File uploaded sucessfully", collections.FileUpload{ImageUrl: url}, nil, err
}
