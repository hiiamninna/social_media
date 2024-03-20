package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"social_media/collections"
	"social_media/library"
	"social_media/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	repo       repository.Repository
	jwt        library.JWT
	bcryptSalt int
}

func NewUserController(repo repository.Repository, jwt library.JWT, bcryptSalt int) User {
	return User{
		repo:       repo,
		jwt:        jwt,
		bcryptSalt: bcryptSalt,
	}
}

func (c User) Register(ctx *fiber.Ctx) (int, string, interface{}, interface{}, error) {

	raw := ctx.Request().Body()
	input := collections.UserRegisterInput{}
	err := json.Unmarshal([]byte(raw), &input)
	if err != nil {
		return http.StatusBadRequest, UNMARSHAL_INPUT, nil, nil, err
	}

	message, err := library.ValidateInput(input)
	if err != nil {
		return http.StatusBadRequest, message, nil, nil, err
	}

	if input.CredType == "phone" {
		isPhone := library.IsPhone(input.CredValue)
		if !isPhone {
			return http.StatusBadRequest, "not valid phone number", nil, nil, err
		}

		exist, _ := c.repo.User.GetByPhone(input.CredValue)
		if exist.Id != "" {
			return http.StatusConflict, "phone number registered", nil, nil, err
		}

		input.Phone = input.CredValue
	} else {
		isEmail := library.IsEmail(input.CredValue)
		if !isEmail {
			return http.StatusBadRequest, "not valid email", nil, nil, err
		}

		exist, _ := c.repo.User.GetByEmail(input.CredValue)
		if exist.Id != "" {
			return http.StatusConflict, "email registered", nil, nil, err
		}

		input.Email = input.CredValue
	}

	generated, err := bcrypt.GenerateFromPassword([]byte(input.Password), c.bcryptSalt)
	if err != nil {
		return http.StatusInternalServerError, "failed generate", nil, nil, err
	}

	input.Password = string(generated)

	id, err := c.repo.User.Create(input)
	if err != nil {
		return http.StatusInternalServerError, "User registered failed", nil, nil, err
	}

	token, err := c.jwt.CreateToken(strconv.Itoa(id))
	if err != nil {
		return http.StatusInternalServerError, "User registered failed", nil, nil, err
	}

	// TODO : make it more simple
	if input.CredType == "phone" {
		return http.StatusCreated, "User registered successfully", collections.UserRegisterByPhone{
			Phone:       input.Phone,
			Name:        input.Name,
			AccessToken: token,
		}, nil, err
	} else {
		return http.StatusCreated, "User registered successfully", collections.UserRegisterByEmail{
			Email:       input.Email,
			Name:        input.Name,
			AccessToken: token,
		}, nil, err
	}
}

func (c User) Login(ctx *fiber.Ctx) (int, string, interface{}, interface{}, error) {
	raw := ctx.Request().Body()
	input := collections.UserLoginInput{}
	err := json.Unmarshal([]byte(raw), &input)
	if err != nil {
		return http.StatusBadRequest, UNMARSHAL_INPUT, nil, nil, err
	}

	message, err := library.ValidateInput(input)
	if err != nil {
		return http.StatusBadRequest, message, nil, nil, err
	}

	user := collections.User{}
	if input.CredType == "phone" {

		isPhone := library.IsPhone(input.CredValue)
		if !isPhone {
			return http.StatusBadRequest, "not valid phone number", nil, nil, err
		}

		user, err = c.repo.User.GetByPhone(input.CredValue)
		if err != nil {
			return http.StatusNotFound, "User not found", nil, nil, err
		}
	} else {

		isEmail := library.IsEmail(input.CredValue)
		if !isEmail {
			return http.StatusBadRequest, "not valid email", nil, nil, err
		}

		user, err = c.repo.User.GetByEmail(input.CredValue)
		if err != nil {
			return http.StatusNotFound, "User not found", nil, nil, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return http.StatusBadRequest, "Check your password again", nil, nil, err
	}

	token, err := c.jwt.CreateToken(user.Id)
	if err != nil {
		return http.StatusInternalServerError, "User login failed", nil, nil, err
	}

	resp := collections.UserLogin{
		Email:       user.Email,
		Phone:       user.Phone,
		Name:        user.Name,
		AccessToken: token,
	}

	return http.StatusOK, "User logged successfully", resp, nil, err
}

func (c User) UpdateProfile(ctx *fiber.Ctx) (int, string, interface{}, interface{}, error) {
	raw := ctx.Request().Body()
	input := collections.UserUpdateInput{}
	err := json.Unmarshal([]byte(raw), &input)
	if err != nil {
		return http.StatusBadRequest, UNMARSHAL_INPUT, nil, nil, err
	}

	input.UserID = library.GetUserID(ctx)

	message, err := library.ValidateInput(input)
	if err != nil {
		return http.StatusBadRequest, message, nil, nil, err
	}

	isImage := library.IsImageUrl(input.ImageUrl)
	if !isImage {
		return http.StatusBadRequest, "not image url", nil, nil, errors.New("not image url")
	}

	err = c.repo.User.UpdateProfile(input)
	if err != nil {
		return http.StatusInternalServerError, "update profile failed", nil, nil, err
	}

	return http.StatusOK, "update profile success", nil, nil, nil
}

func (c User) UpdateLinkEmail(ctx *fiber.Ctx) (int, string, interface{}, interface{}, error) {
	raw := ctx.Request().Body()
	input := collections.UserLinkEmail{}
	err := json.Unmarshal([]byte(raw), &input)
	if err != nil {
		return http.StatusBadRequest, UNMARSHAL_INPUT, nil, nil, err
	}

	input.UserID = library.GetUserID(ctx)

	user, _ := c.repo.User.GetByID(input.UserID)
	if user.Email != "" {
		return http.StatusBadRequest, "you have an email", nil, nil, errors.New("you have an email")
	}

	message, err := library.ValidateInput(input.Email)
	if err != nil {
		return http.StatusBadRequest, message, nil, nil, err
	}

	exist, _ := c.repo.User.GetByEmail(input.Email)
	if exist.Id != "" {
		return http.StatusConflict, "email registered", nil, nil, errors.New("email registered")
	}

	err = c.repo.User.UpdateEmail(input)
	if err != nil {
		return http.StatusInternalServerError, "link an email failed", nil, nil, err
	}

	return http.StatusOK, "link an email success", nil, nil, nil
}

func (c User) UpdateLinkPhone(ctx *fiber.Ctx) (int, string, interface{}, interface{}, error) {
	raw := ctx.Request().Body()
	input := collections.UserLinkPhone{}
	err := json.Unmarshal([]byte(raw), &input)
	if err != nil {
		return http.StatusBadRequest, UNMARSHAL_INPUT, nil, nil, err
	}

	input.UserID = library.GetUserID(ctx)

	user, _ := c.repo.User.GetByID(input.UserID)
	if user.Phone != "" {
		return http.StatusBadRequest, "you have a phone", nil, nil, errors.New("you have a phone")
	}

	message, err := library.ValidateInput(input.Phone)
	if err != nil {
		return http.StatusBadRequest, message, nil, nil, err
	}

	exist, _ := c.repo.User.GetByPhone(input.Phone)
	if exist.Id != "" {
		return http.StatusConflict, "phone registered", nil, nil, errors.New("phone registered")
	}

	err = c.repo.User.UpdatePhone(input)
	if err != nil {
		return http.StatusInternalServerError, "link an phone failed", nil, nil, err
	}

	return http.StatusOK, "link an phone success", nil, nil, nil
}

// TODO : validate email => make a simpler function
// TODO : validate phone => make a simpler function
