package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrResponse struct {
	Error string `json:"error"`
}

func Error(msg string) ErrResponse {
	return ErrResponse{
		Error: msg,
	}
}

func ValidationError(errs validator.ValidationErrors) ErrResponse {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("Field '%s' is a required", err.Field()))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("Field '%s' is not a valid email format", err.Field()))
		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("Field '%s' must be longer than %s", err.Field(), err.Param()))
		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("Field '%s' must be shorter than %s", err.Field(), err.Param()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("Field '%s' is not valid", err.Field()))
		}
	}

	return ErrResponse{
		Error: strings.Join(errMsgs, "; "),
	}
}
