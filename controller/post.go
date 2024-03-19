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

func (c Post) List(ctx *fiber.Ctx) (int, string, interface{}, interface{}, error) {

	raw := ctx.Request().Body()
	input := collections.PostInputParam{}
	err := json.Unmarshal([]byte(raw), &input)
	if err != nil {
		return http.StatusBadRequest, "unmarshal input", nil, nil, err
	}

	// TODO : query params

	input.UserID, _ = library.GetUserID(ctx)
	if input.UserID == "" {
		return http.StatusForbidden, "please check your credential", nil, nil, errors.New("not login")
	}

	result := []collections.PostList{}

	posts, err := c.repo.Post.List()
	if err != nil {
		return http.StatusInternalServerError, "post list error", nil, nil, err
	}

	comments, err := c.repo.Comment.List()
	if err != nil {
		return http.StatusInternalServerError, "comment list error", nil, nil, err
	}

	creator, err := c.repo.User.List()
	if err != nil {
		return http.StatusInternalServerError, "creator list error", nil, nil, err
	}

	for _, p := range posts {
		result = append(result, collections.PostList{
			Posts:    p,
			Comments: getComments(p.PostID, comments),
			Creator:  getCreator(p.UserID, creator),
		})
	}

	return http.StatusOK, "ok", result, nil, nil

}

func getComments(postID string, comments []collections.Comment) []collections.Comment {

	temp := []collections.Comment{}

	for _, c := range comments {
		if c.PostID == postID {
			temp = append(temp, c)
		}
	}

	return temp
}

func getCreator(userID string, creator []collections.UserAsFriend) collections.UserAsFriend {

	for _, c := range creator {
		if c.UserId == userID {
			return c
		}
	}

	return collections.UserAsFriend{}
}
