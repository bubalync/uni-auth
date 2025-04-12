package response

import "errors"

var (
	ErrInternal          = errors.New("internal server error")
	ErrInvalidAuthHeader = errors.New("invalid auth header")
	ErrInvalidToken      = errors.New("invalid token")
)

type ErrResponse struct {
	Errors map[string]string `json:"errors"`
}

func ErrorInternal() ErrResponse {
	return ErrResponse{
		Errors: map[string]string{"message": ErrInternal.Error()},
	}
}

func Error(msg string) ErrResponse {
	return ErrResponse{
		Errors: map[string]string{"message": msg},
	}
}

func ErrorMap(errs map[string]string) ErrResponse {
	return ErrResponse{
		Errors: errs,
	}
}
