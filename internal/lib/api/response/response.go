package response

type ErrResponse struct {
	Errors map[string]string `json:"errors"`
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
