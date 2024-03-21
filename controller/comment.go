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
		return http.StatusBadRequest, UNMARSHAL_INPUT, nil, nil, err
	}

	input.UserID = library.GetUserID(ctx)

	message, err := library.ValidateInput(input)
	if err != nil {
		return http.StatusBadRequest, message, nil, nil, err
	}

	post, err := c.repo.Post.GetById(input.PostID)
	if err != nil {
		return http.StatusNotFound, "post not found", nil, nil, err
	}

	friend, err := c.repo.Friend.GetByUser(collections.FriendInput{
		UserID:   input.UserID,
		FriendID: post.UserID,
	})
	if err != nil || friend.Id == 0 {
		return http.StatusBadRequest, "can not comment, not your friends", nil, nil, errors.New("not friend")
	}

	err = c.repo.Comment.Create(input)
	if err != nil {
		return http.StatusInternalServerError, "failed add comment", nil, nil, errors.New("failed add comment")
	}

	return http.StatusOK, "successfully added a comment", nil, nil, nil
}
