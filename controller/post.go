package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"social_media/collections"
	"social_media/library"
	"social_media/repository"
	"time"

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
		return http.StatusBadRequest, UNMARSHAL_INPUT, nil, nil, err
	}

	input.UserID = library.GetUserID(ctx)

	message, err := library.ValidateInput(input)
	if err != nil {
		return http.StatusBadRequest, message, nil, nil, err
	}

	err = c.repo.Post.Create(input)
	if err != nil {
		return http.StatusBadRequest, "failed add post", nil, nil, errors.New("failed add post")
	}

	return http.StatusOK, "successfully added a post", nil, nil, nil
}

func (c Post) List(ctx *fiber.Ctx) (int, string, interface{}, interface{}, error) {

	input := collections.PostInputParam{}
	if err := ctx.QueryParser(&input); err != nil {
		return http.StatusInternalServerError, "list post error", nil, nil, err
	}

	if input.Limit == 0 {
		input.Limit = 5
		input.Offset = 0
	}

	input.UserID = library.GetUserID(ctx)

	result := []collections.Post{}

	posts, ids, meta, err := c.repo.Post.List(input)
	if err != nil {
		return http.StatusInternalServerError, "post list error", nil, nil, err
	}

	comments, err := c.repo.Comment.List(ids)
	if err != nil {
		return http.StatusInternalServerError, "comment list error", nil, nil, err
	}

	creator, err := c.repo.User.List()
	if err != nil {
		return http.StatusInternalServerError, "creator list error", nil, nil, err
	}

	for _, p := range posts {
		result = append(result, collections.Post{
			PostID: p.PostID,
			PostData: struct {
				PostInHtml string    "json:\"postInHtml\""
				Tags       []string  "json:\"tags\""
				CreatedAt  time.Time "json:\"createdAt\""
			}{
				PostInHtml: p.PostData.PostInHtml,
				Tags:       p.PostData.Tags,
				CreatedAt:  p.PostData.CreatedAt,
			},
			Comments: getComments(p.PostID, comments),
			Creator:  getCreator(p.UserID, creator),
		})
	}

	return http.StatusOK, "ok", result, meta, nil

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
