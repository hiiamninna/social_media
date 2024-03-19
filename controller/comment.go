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

type Comment struct {
	repo repository.Repository
}

func NewCommentController(repo repository.Repository) Comment {
	return Comment{
		repo: repo,
	}
}

func (c Comment) Create(ctx *fiber.Ctx) (int, string, interface{}, interface{}, error) {

	raw := ctx.Request().Body()
	input := collections.CommentInput{}
	err := json.Unmarshal([]byte(raw), &input)
	if err != nil {
		return http.StatusBadRequest, "unmarshal input", nil, nil, err
	}

	input.UserID, _ = library.GetUserID(ctx)
	if input.UserID == "" {
		return http.StatusForbidden, "please check your credential", nil, nil, errors.New("not login")
	}

	// TODO : validation

	err = c.repo.Comment.Create(input)
	if err != nil {
		return http.StatusBadRequest, "failed add comment", nil, nil, errors.New("failed add comment")
	}

	return http.StatusOK, "successfully added a comment", nil, nil, nil
}
