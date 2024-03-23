package library

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	validation "github.com/go-playground/validator/v10"
)

const TIME_FORMAT = "2006-01-02 15:01:02 "

var validationError = validation.New()

// Using validator
func ValidateInput(data interface{}) (string, error) {

	// create new validationa and check the struct

	err := validationError.Struct(data)

	if err != nil {
		var errors []string
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validation.InvalidValidationError); ok {
			fmt.Println(time.Now().Format(TIME_FORMAT), err)
			return "", nil
		}

		for n, e := range err.(validation.ValidationErrors) {
			message := fmt.Sprintf("%s must %s,", e.Field(), e.Tag())
			if n == (len(err.(validation.ValidationErrors)) - 1) {
				if e.Tag() == "email" {
					message = "Please input correct email format"
				} else {
					message = fmt.Sprintf("%s must %s", e.Field(), e.Tag())
				}
			}
			errors = append(errors, message)
		}
		return fmt.Sprint(errors), err
	}

	if err != nil {
		arrayOfErrors := []string{err.Error()}
		return fmt.Sprint(arrayOfErrors), err
	}

	return "", err
}

func IsEmail(value string) bool {

	err := validationError.Var(value, "required,email")
	if err != nil {
		fmt.Println(time.Now().Format(TIME_FORMAT), err)
		return false
	}

	return true
}

func IsPhone(value string) bool {

	err := validationError.Var(value, "required,startswith=+,min=7,max=13")
	if err != nil {
		fmt.Println(time.Now().Format(TIME_FORMAT), err)
		return false
	}

	return true
}

func IsHaveExt(value string) bool {
	u, err := url.Parse(value)
	if err != nil {
		return false
	}

	pos := strings.LastIndex(u.Path, ".")
	return pos != -1
}
