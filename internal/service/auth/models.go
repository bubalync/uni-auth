package auth

type (
	CreateUserInput struct {
		Email    string
		Password string
	}

	GenerateTokenInput struct {
		Email    string
		Password string
	}

	ResetPasswordInput struct {
		Email    string
		Password string
	}
)
