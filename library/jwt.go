package library

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type JWT struct {
	SecretKey string
}

type customClaims struct {
	ID string
	jwt.StandardClaims
}

func NewJWT(secretKey string) JWT {
	return JWT{
		SecretKey: secretKey,
	}
}

func (c JWT) CreateToken(id string) (string, error) {

	expiredTime := time.Now().Add(8 * time.Hour)
	claims := &customClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime.Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(c.SecretKey))
	if err != nil {
		return "", fmt.Errorf("new with claims : %w", err)
	}

	return token, nil
}

func (c JWT) Authentication() fiber.Handler {
	return func(context *fiber.Ctx) error {
		errHandle := func(err error, message string) error {
			fmt.Println(time.Now().Format("2006-01-02 15:01:02 "), err)
			return context.Status(http.StatusUnauthorized).JSON(struct {
				Message string `json:"message"`
			}{
				Message: message,
			})
		}

		input := context.Request().Header.Peek("Authorization")

		if !strings.Contains(string(input), "Bearer") {
			return errHandle(errors.New("invalid header"), "invalid header")
		}

		authToken := ""
		arrayToken := strings.Split(string(input), " ")
		if len(arrayToken) == 2 {
			authToken = arrayToken[1]
		}

		claims := &customClaims{}
		token, err := jwt.ParseWithClaims(authToken, claims, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return method, errors.New("unexpected signing method")
			}
			return []byte(c.SecretKey), nil
		})

		if err != nil {
			return errHandle(err, "token expired")
		}

		if !token.Valid {
			return errHandle(errors.New("invalid token"), "invalid token")
		}

		SetSession(context, claims.ID)

		return context.Next()
	}
}
