package v1

import (
	"errors"
	"github.com/bubalync/uni-auth/internal/lib/api/response"
	"github.com/bubalync/uni-auth/internal/service"
	"github.com/bubalync/uni-auth/internal/service/auth"
	"github.com/bubalync/uni-auth/internal/service/svcErrs"
	"github.com/bubalync/uni-auth/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type authRoutes struct {
	as service.Auth
	cv *validator.CustomValidator
}

func NewAuthRoutes(g *gin.RouterGroup, cv *validator.CustomValidator, authService service.Auth) {
	r := &authRoutes{authService, cv}

	g.POST("/sign-up", r.signUp)
	g.POST("/sign-in", r.signIn)
	g.POST("/reset-password", r.signIn)
}

type signUpRequest struct {
	Email    string `json:"email"    validate:"required,email,min=5,max=150" minLength:"5" maxLength:"150" example:"email@example.com"`
	Password string `json:"password" validate:"required,password"            minLength:"8" maxLength:"32"  example:"YourV@lidPassw0rd!"`
}

type signUpResponse struct {
	Id uuid.UUID `json:"id" example:"d13a75e2-3d21-4e57-9dc0-3a7f5bee4c25"`
}

// @Summary     Sign up
// @Description Sign up
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body signUpRequest true "Registration payload"
// @Success     201 {object} signUpResponse
// @Failure     400 {object} response.ErrResponse
// @Failure     422 {object} response.ErrResponse
// @Failure     500 {object} response.ErrResponse
// @Router      /auth/sign-up [post]
func (r *authRoutes) signUp(c *gin.Context) {
	var req signUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	if errs := r.cv.ValidateStruct(req); errs != nil {
		c.JSON(http.StatusBadRequest, response.ErrorMap(errs))
		return
	}

	id, err := r.as.CreateUser(c.Request.Context(), auth.CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, svcErrs.ErrUserAlreadyExists) {
			c.JSON(http.StatusUnprocessableEntity, response.Error(err.Error()))
			return
		}

		c.JSON(http.StatusInternalServerError, response.ErrorInternal())
		return
	}

	c.JSON(http.StatusCreated, &signUpResponse{
		Id: id,
	})
}

type signInRequest struct {
	Email    string `json:"email"    validate:"required,email,min=5,max=150" minLength:"5" maxLength:"150" example:"email@example.com"`
	Password string `json:"password" validate:"required"                     minLength:"8" maxLength:"32"  example:"YourV@lidPassw0rd!"`
}

type signInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// @Summary     Sign in
// @Description Sign in
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body signInRequest true "Sign in payload"
// @Success     200 {object} signInResponse
// @Failure     400 {object} response.ErrResponse
// @Failure     500 {object} response.ErrResponse
// @Router      /auth/sign-in [post]
func (r *authRoutes) signIn(c *gin.Context) {
	var req signInRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	if errs := r.cv.ValidateStruct(req); errs != nil {
		c.JSON(http.StatusBadRequest, response.ErrorMap(errs))
		return
	}

	tokens, err := r.as.GenerateToken(c.Request.Context(), auth.GenerateTokenInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, svcErrs.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, response.Error(err.Error()))
			return
		}

		c.JSON(http.StatusInternalServerError, response.ErrorInternal())
		return
	}

	c.JSON(http.StatusOK, signInResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})

}

func (r *userRoutes) resetPassword(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
