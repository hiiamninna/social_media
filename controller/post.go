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

	// TODO : check again if tags:[null,"abc"], should be error or not
	for i := range input.Tags {
		if input.Tags[i] == "" {
			return http.StatusBadRequest, "input tags is empty", nil, nil, err
		}
	}

	err = c.repo.Post.Create(input)
	if err != nil {
		return http.StatusBadRequest, "failed add post", nil, nil, errors.New("failed add post")
	}

	return http.StatusOK, "successfully added a post", nil, nil, nil
}

func (c Post) List(ctx *fiber.Ctx) (int, string, interface{}, interface{}, error) {

	input := collections.PostInputParam{
		UserID: library.GetUserID(ctx),
		Tags:   []string{},
		Search: "",
		Limit:  5,
		Offset: 0,
	}
	if err := ctx.QueryParser(&input); err != nil {
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

	result := []collections.Post{}

	posts, ids, meta, err := c.repo.Post.List(input)
	if err != nil {
		return http.StatusInternalServerError, "post list error", nil, nil, err
	}

	comments, err := c.repo.Comment.List(ids)
	if err != nil {
		return http.StatusInternalServerError, "comment list error", nil, nil, err
	}

	for _, p := range posts {
		result = append(result, collections.Post{
			PostID: p.PostID,
			PostData: struct {
				PostInHtml   string    `json:"postInHtml"`
				Tags         []string  `json:"tags"`
				CreatedAt    time.Time `json:"-"`
				CreatedAtStr string    `json:"createdAt"`
			}{
				PostInHtml:   p.PostData.PostInHtml,
				Tags:         p.PostData.Tags,
				CreatedAtStr: p.PostData.CreatedAtStr,
			},
			Comments: comments[p.ID],
			Creator:  p.Creator,
		})
	}

	return http.StatusOK, "ok", result, meta, nil

}
