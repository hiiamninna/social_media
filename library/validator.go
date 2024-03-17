package library

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const (
	New    string = "new"
	Second string = "second"
)

func Validate(item interface{}) error {
	val := reflect.ValueOf(item)
	for i := 0; i < val.NumField(); i++ {
		var err error
		fieldType := val.Type().Field(i)
		rule := fieldType.Tag.Get("validate")
		rules := strings.Split(rule, ";")
		//no rules defined, then exit
		if len(rules) == 0 {
			continue
		}

		field := val.Field(i)
		value := field.Interface()

		//get name of variable
		name := fieldType.Tag.Get("json")
		//get value of variable
		//check type of validation
		switch field.Type().Kind() {
		case reflect.String:
			err = validateString(rules, name, value)
		case reflect.Int:
			err = validateInt(rules, name, value)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func validateString(rules []string, name string, v interface{}) error {
	value := v.(string)
	for _, rule := range rules {
		if err := strRequired(rule, name, value); err != nil {
			return err
		}
		if err := strMinLength(rule, name, value); err != nil {
			return err
		}
		if err := strMaxLength(rule, name, value); err != nil {
			return err
		}
		if err := strIsValidUrl(rule, name, value); err != nil {
			return err
		}
		if err := strCheckEnum(rule, name, value); err != nil {
			return err
		}
	}

	return nil
}

func strRequired(rule, name, value string) error {
	if !strings.Contains(rule, "required") {
		return nil
	}
	if value == "" {
		return fmt.Errorf("field %s must be filled", name)
	}

	return nil
}

func strIsValidUrl(rule, name, value string) error {
	if !strings.Contains(rule, "url") {
		return nil
	}

	_, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("fiels %s must be a valid url", name)
	}

	return nil
}

func strCheckEnum(rule, name, value string) error {
	if !strings.Contains(rule, "enum") {
		return nil
	}

	r := strings.Split(rule, ":")
	if len(r) < 2 {
		return fmt.Errorf("enum invalid rule definition, name : " + name)
	}

	enumType := strings.TrimSpace(r[1])

	switch enumType {
	case "condition":
		checkEnumCondition(name, value)
	default:
		return fmt.Errorf("enum invalid rule definition, name : " + name)
	}

	return nil
}

func checkEnumCondition(name, value string) error {
	switch value {
	case New, Second:
		return nil
	default:
		return fmt.Errorf("invalid value for field %s", name)
	}
}

func strMinLength(rule, name, value string) error {
	if !strings.Contains(rule, "min") {
		return nil
	}

	r := strings.Split(rule, ":")
	if len(r) < 2 {
		return fmt.Errorf("min-length invalid rule definition, name : " + name)
	}

	limit, err := strconv.Atoi(strings.TrimSpace(r[1]))
	if err != nil {
		return fmt.Errorf("invalid rule:(%s) %w", name, err)
	}

	if len(value) < limit {
		return fmt.Errorf("field %s must have at least %d character(s)", name, limit)
	}

	return nil
}

func strMaxLength(rule, name, value string) error {
	if !strings.Contains(rule, "max") {
		return nil
	}

	r := strings.Split(rule, ":")
	if len(r) < 2 {
		return fmt.Errorf("max-length invalid rule definition, name : " + name)
	}

	limit, err := strconv.Atoi(strings.TrimSpace(r[1]))
	if err != nil {
		return fmt.Errorf("invalid rule:(%s) %w", name, err)
	}

	if len(value) > limit {
		return fmt.Errorf("total characters for field %s must be less or same than %d character(s)", name, limit)
	}

	return nil
}

func validateInt(rules []string, name string, v interface{}) error {
	value, _ := v.(int)
	for _, rule := range rules {
		if err := intRequired(rule, name, value); err != nil {
			return err
		}
		if err := intMinValue(rule, name, value); err != nil {
			return err
		}
		if err := intMaxValue(rule, name, value); err != nil {
			return err
		}
	}

	return nil
}

func intRequired(rule, name string, value int) error {
	if !strings.Contains(rule, "required") {
		return nil
	}
	if value == 0 {
		return fmt.Errorf("field %s must be filled", name)
	}

	return nil
}

func intMinValue(rule, name string, value int) error {
	if !strings.Contains(rule, "min") {
		return nil
	}
	r := strings.Split(rule, ":")
	if len(r) < 2 {
		return fmt.Errorf("min-value invalid rule definition, name : " + name)
	}

	limit, err := strconv.Atoi(strings.TrimSpace(r[1]))
	if err != nil {
		return fmt.Errorf("invalid rule:(%s) %w", name, err)
	}

	if value < limit {
		return fmt.Errorf("field %s must not less than %d", name, limit)
	}

	return nil
}

func intMaxValue(rule, name string, value int) error {
	if !strings.Contains(rule, "max") {
		return nil
	}
	r := strings.Split(rule, ":")
	if len(r) < 2 {
		return fmt.Errorf("max-value invalid rule definition, name : " + name)
	}

	limit, err := strconv.Atoi(strings.TrimSpace(r[1]))
	if err != nil {
		return fmt.Errorf("invalid rule:(%s) %w", name, err)
	}

	if value > limit {
		return fmt.Errorf("field %s must not greater than %d", name, limit)
	}

	return nil
}
