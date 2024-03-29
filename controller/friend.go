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

type Friend struct {
	repo repository.Repository
}

func NewFriendController(repo repository.Repository) Friend {
	return Friend{
		repo: repo,
	}
}

func (c Friend) List(ctx *fiber.Ctx) (int, string, interface{}, interface{}, error) {

	// set default value first
	// TODO : analize for indexing params
	input := collections.FriendInputParam{
		UserID:        library.GetUserID(ctx),
		Search:        "",
		OnlyFriendStr: "false",
		OnlyFriend:    false,
		OrderBy:       "desc",
		SortBy:        "createdAt",
		Limit:         5,
		Offset:        0,
	}
	err := ctx.QueryParser(&input)
	if err != nil {
		return http.StatusBadRequest, "unmarshal input", nil, nil, err
	}

	//TODO : temporary solution if limit= OR offset=
	//clue=>everything from url is a string
	maps := ctx.Queries()
	if val, ok := maps["limit"]; ok {
		if val == "" {
			return http.StatusBadRequest, "limit an empty value", nil, nil, errors.New("limit can not empty")
		}
		if val == "0" {
			input.Limit = 5
		}
	}

	if val, ok := maps["offset"]; ok {
		if val == "" {
			return http.StatusBadRequest, "offset an empty value", nil, nil, errors.New("offset can not empty")
		}
	}

	message, err := library.ValidateInput(input)
	if err != nil {
		return http.StatusBadRequest, message, nil, nil, err
	}

	if input.OnlyFriendStr == "true" {
		input.OnlyFriend = true
	}

	friends, err := c.repo.Friend.List(input)
	if err != nil {
		return http.StatusInternalServerError, "failed get list of friends", nil, nil, err
	}

	totalRow, err := c.repo.Friend.CountList(input)
	if err != nil {
		return http.StatusInternalServerError, "failed get list of friends", nil, nil, err
	}

	return http.StatusOK, "ok", friends, collections.Meta{Limit: input.Limit, Offset: input.Offset, Total: totalRow}, nil
}

func (c Friend) Create(ctx *fiber.Ctx) (int, string, interface{}, interface{}, error) {

	raw := ctx.Request().Body()
	input := collections.FriendInput{}
	err := json.Unmarshal([]byte(raw), &input)
	if err != nil {
		return http.StatusBadRequest, "unmarshal input", nil, nil, err
	}

	message, err := library.ValidateInput(input)
	if err != nil {
		return http.StatusBadRequest, message, nil, nil, err
	}

	input.UserID = library.GetUserID(ctx)

	// REFACTOR : not in test case, register not return id, can not test it
	// if input.FriendID == input.UserID {
	// 	return http.StatusBadRequest, "can not add your own self", nil, nil, errors.New("can not add your own self")
	// }

	newFriend, _ := c.repo.User.GetByID(input.FriendID)
	if newFriend.Id == "" {
		return http.StatusNotFound, "user not found", nil, nil, errors.New("new friend not found")
	}

	friend, _ := c.repo.Friend.GetByUser(input)
	if friend.Id != 0 {
		return http.StatusBadRequest, "already be friend", nil, nil, errors.New("already be friend")
	}

	err = c.repo.Friend.Create(input)
	if err != nil {
		return http.StatusBadRequest, "failed to be friend", nil, nil, err
	}

	return http.StatusOK, "successfully added as a friend", nil, nil, nil
}

func (c Friend) Delete(ctx *fiber.Ctx) (int, string, interface{}, interface{}, error) {

	raw := ctx.Request().Body()
	input := collections.FriendInput{}
	err := json.Unmarshal([]byte(raw), &input)
	if err != nil {
		return http.StatusBadRequest, "unmarshal input", nil, nil, err
	}

	message, err := library.ValidateInput(input)
	if err != nil {
		return http.StatusBadRequest, message, nil, nil, err
	}

	input.UserID = library.GetUserID(ctx)

	newFriend, _ := c.repo.User.GetByID(input.FriendID)
	if newFriend.Id == "" {
		return http.StatusNotFound, "friend not found", nil, nil, errors.New("friend not found")
	}

	friend, _ := c.repo.Friend.GetByUser(input)
	if friend.Id == 0 {
		return http.StatusBadRequest, "not your friend", nil, nil, errors.New("not your friend")
	}

	err = c.repo.Friend.Delete(friend.Id, input.UserID, input.FriendID)
	if err != nil {
		return http.StatusInternalServerError, "failed unfriend", nil, nil, errors.New("failed unfriend")
	}

	return http.StatusOK, "successfully unfriend", nil, nil, nil
}
