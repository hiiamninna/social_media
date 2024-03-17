package controller

import (
	"social_media/library"
	"social_media/repository"

	"github.com/google/uuid"
)

type Controller struct {
}

func NewController(repo repository.Repository, jwt library.JWT, bcryptSalt int, s3 library.S3) Controller {
	return Controller{}
}

func generateUUID() string {
	return uuid.NewString()
}
