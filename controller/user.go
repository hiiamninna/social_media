package controller

import (
	"encoding/json"
	"errors"
	"fmt"
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
		return http.StatusBadRequest, "unmarshal input", nil, nil, err
	}

	// TODO : validation

	generated, err := bcrypt.GenerateFromPassword([]byte(input.Password), c.bcryptSalt)
	if err != nil {
		return http.StatusInternalServerError, "failed generate", nil, nil, err
	}

	input.Password = string(generated)

	if input.CredType == "phone" {
		input.Phone = input.CredValue
	} else if input.CredType == "email" {
		input.Email = input.CredValue
	} else {
		// TODO : return error enum or get into validation
	}

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
		return http.StatusBadRequest, "unmarshal input", nil, nil, err
	}

	// TODO : validation

	user := collections.User{}
	if input.CredType == "phone" {
		user, err = c.repo.User.GetByPhone(input.CredValue)
		if err != nil {
			return http.StatusNotFound, "User not found", nil, nil, err
		}
	} else if input.CredType == "email" {
		user, err = c.repo.User.GetByEmail(input.CredValue)
		if err != nil {
			return http.StatusNotFound, "User not found", nil, nil, err
		}
	} else {
		// TODO : return error enum or get into validation
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return http.StatusBadRequest, "Check your password again", nil, nil, err
	}

	token, err := c.jwt.CreateToken(user.Id)
	if err != nil {
		fmt.Println(err)
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
		return http.StatusBadRequest, "unmarshal input", nil, nil, err
	}

	input.UserID, _ = library.GetUserID(ctx)
	if input.UserID == "" {
		return http.StatusForbidden, "please check your credential", nil, nil, errors.New("not login")
	}

	// TODO : validation

	err = c.repo.User.UpdateProfile(input)
	if err != nil {
		return http.StatusInternalServerError, "update profile failed", nil, nil, err
	}

	return http.StatusOK, "update profile success", nil, nil, nil
}

func (c User) UpdateLinkEmail(ctx *fiber.Ctx) (int, string, interface{}, interface{}, error) {
	raw := ctx.Request().Body()
	input := collections.UserLinkInput{}
	err := json.Unmarshal([]byte(raw), &input)
	if err != nil {
		return http.StatusBadRequest, "unmarshal input", nil, nil, err
	}

	input.UserID, _ = library.GetUserID(ctx)
	if input.UserID == "" {
		return http.StatusForbidden, "please check your credential", nil, nil, errors.New("not login")
	}

	// TODO : validation
	user, _ := c.repo.User.GetByID(input.UserID)
	if user.Email != "" {
		return http.StatusBadRequest, "you have an email", nil, nil, errors.New("you have an email")
	}

	err = c.repo.User.UpdateEmail(input)
	if err != nil {
		return http.StatusInternalServerError, "link an email failed", nil, nil, err
	}

	return http.StatusOK, "link an email success", nil, nil, nil
}

func (c User) UpdateLinkPhone(ctx *fiber.Ctx) (int, string, interface{}, interface{}, error) {
	raw := ctx.Request().Body()
	input := collections.UserLinkInput{}
	err := json.Unmarshal([]byte(raw), &input)
	if err != nil {
		return http.StatusBadRequest, "unmarshal input", nil, nil, err
	}

	input.UserID, _ = library.GetUserID(ctx)
	if input.UserID == "" {
		return http.StatusForbidden, "please check your credential", nil, nil, errors.New("not login")
	}

	// TODO : validation
	user, _ := c.repo.User.GetByID(input.UserID)
	if user.Phone != "" {
		return http.StatusBadRequest, "you have a phone", nil, nil, errors.New("you have a phone")
	}

	err = c.repo.User.UpdatePhone(input)
	if err != nil {
		return http.StatusInternalServerError, "link an phone failed", nil, nil, err
	}

	return http.StatusOK, "link an phone success", nil, nil, nil
}
