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

	GenerateTokenOutput struct {
		AccessToken  string
		RefreshToken string
	}

	ResetPasswordInput struct {
		Email    string
		Password string
	}
)
