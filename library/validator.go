package library

import (
	"fmt"
	"net/url"
	"strings"

	validation "github.com/go-playground/validator/v10"
)

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
			fmt.Println(err)
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
		fmt.Println(err.Error())
		return false
	}

	return true
}

func IsPhone(value string) bool {

	err := validationError.Var(value, "required,startswith=+,min=7,max=13")
	if err != nil {
		fmt.Println(err.Error())
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
