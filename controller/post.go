package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"social_media/collections"
	"social_media/library"
	"social_media/repository"

	"github.com/gofiber/fiber/v2"
)

type Post struct {
	repo repository.Repository
}

func NewPostController(repo repository.Repository) Post {
	return Post{
		repo: repo,
	}
}

func (c Post) Create(ctx *fiber.Ctx) (int, string, interface{}, interface{}, error) {

	raw := ctx.Request().Body()
	input := collections.PostInput{}
	err := json.Unmarshal([]byte(raw), &input)
	if err != nil {
		return http.StatusBadRequest, "unmarshal input", nil, nil, err
	}

	input.UserID, _ = library.GetUserID(ctx)
	if input.UserID == "" {
		return http.StatusForbidden, "please check your credential", nil, nil, errors.New("not login")
	}

	// TODO : validation

	err = c.repo.Post.Create(input)
	if err != nil {
		return http.StatusBadRequest, "failed add post", nil, nil, errors.New("failed add post")
	}

	return http.StatusOK, "successfully added a post", nil, nil, nil
}
