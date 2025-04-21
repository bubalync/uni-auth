package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
)

const (
	passwordMinLength = 8
	passwordMaxLength = 32
	passwordMinLower  = 1
	passwordMinUpper  = 1
	passwordMinDigit  = 1
	passwordMinSymbol = 1
)

var (
	lengthRegexp    = regexp.MustCompile(fmt.Sprintf(`^.{%d,%d}$`, passwordMinLength, passwordMaxLength))
	lowerCaseRegexp = regexp.MustCompile(fmt.Sprintf(`[a-z]{%d,}`, passwordMinLower))
	upperCaseRegexp = regexp.MustCompile(fmt.Sprintf(`[A-Z]{%d,}`, passwordMinUpper))
	digitRegexp     = regexp.MustCompile(fmt.Sprintf(`[0-9]{%d,}`, passwordMinDigit))
	symbolRegexp    = regexp.MustCompile(fmt.Sprintf(`[!@#$%%^&*]{%d,}`, passwordMinSymbol))
)

type CustomValidator struct {
	V *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	v := validator.New(validator.WithRequiredStructEnabled())
	cv := &CustomValidator{V: v}

	err := v.RegisterValidation("password", cv.passwordValidate)
	if err != nil {
		panic(err)
	}

	return cv
}

func (cv *CustomValidator) ValidateStruct(data interface{}) map[string]string {
	err := cv.V.Struct(data)
	if err != nil {
		return parseErrors(err)
	}

	return nil
}

func parseErrors(err error) map[string]string {
	errors := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		switch err.Tag() {
		case "required":
			errors[field] = "Is a required"
		case "email":
			errors[field] = fmt.Sprintf("Invalid email format")
		case "password":
			errors[field] = fmt.Sprintf(
				"%s must be between %d and %d in length"+
					", contain at least %d lowercase"+
					", %d uppercase"+
					", %d digits"+
					", and %d special characters (!@#$%%^&*)",
				field,
				passwordMinLength, passwordMaxLength,
				passwordMinLower,
				passwordMinUpper,
				passwordMinDigit,
				passwordMinSymbol,
			)
		case "min":
			errors[field] = fmt.Sprintf("Must be longer than %s", err.Param())
		case "max":
			errors[field] = fmt.Sprintf("Must be shorter than %s", err.Param())
		default:
			errors[field] = "Is not valid"
		}
	}

	return errors
}

func (cv *CustomValidator) passwordValidate(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		return false
	}

	fieldValue := fl.Field().String()

	if ok := lengthRegexp.MatchString(fieldValue); !ok {
		return false
	} else if ok = lowerCaseRegexp.MatchString(fieldValue); !ok {
		return false
	} else if ok = upperCaseRegexp.MatchString(fieldValue); !ok {
		return false
	} else if ok = digitRegexp.MatchString(fieldValue); !ok {
		return false
	} else if ok = symbolRegexp.MatchString(fieldValue); !ok {
		return false
	}

	return true
}
