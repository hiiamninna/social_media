package controller

import (
	"social_media/library"
	"social_media/repository"

	"github.com/google/uuid"
)

type Controller struct {
	Image Image
	User  User
}

func NewController(repo repository.Repository, jwt library.JWT, bcryptSalt int, s3 library.S3) Controller {
	return Controller{
		Image: NewImageController(s3),
		User:  NewUserController(repo, jwt, bcryptSalt),
	}
}

func generateUUID() string {
	return uuid.NewString()
}
