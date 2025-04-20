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
	g.POST("/refresh", r.refresh)
	g.POST("/reset-password", r.resetPassword)
	g.POST("/recovery-password", r.recoveryPassword)
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
// @Failure     401 {object} response.ErrResponse
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

type refreshRequest struct {
	// Refresh token
	Token string `json:"token" validate:"required"`
}

type refreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// @Summary     Refresh tokens
// @Description Refresh tokens by refresh-token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body refreshRequest true "Refresh payload"
// @Success     200 {object} refreshResponse
// @Failure     400 {object} response.ErrResponse
// @Failure     401 {object} response.ErrResponse
// @Failure     500 {object} response.ErrResponse
// @Router      /auth/refresh [post]
func (r *authRoutes) refresh(c *gin.Context) {
	var req refreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	if errs := r.cv.ValidateStruct(req); errs != nil {
		c.JSON(http.StatusBadRequest, response.ErrorMap(errs))
		return
	}

	tokens, err := r.as.Refresh(c.Request.Context(), req.Token)
	if err != nil {
		if errors.Is(err, svcErrs.ErrCannotParseToken) || errors.Is(err, svcErrs.ErrTokenIsExpired) {
			c.JSON(http.StatusUnauthorized, response.Error(err.Error()))
			return
		}

		c.JSON(http.StatusInternalServerError, response.ErrorInternal())
		return
	}

	c.JSON(http.StatusOK, refreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

type resetPasswordRequest struct {
	Email string `json:"email"  validate:"required,email"`
}

// @Summary     Reset password
// @Description Password reset request
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body resetPasswordRequest true "Reset password payload"
// @Success     200 {string} string
// @Failure     400 {object} response.ErrResponse
// @Failure     500 {object} response.ErrResponse
// @Router      /auth/reset-password [post]
func (r *authRoutes) resetPassword(c *gin.Context) {
	var req resetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	if errs := r.cv.ValidateStruct(req); errs != nil {
		c.JSON(http.StatusBadRequest, response.ErrorMap(errs))
		return
	}

	err := r.as.ResetPassword(c.Request.Context(), auth.ResetPasswordInput{Email: req.Email})
	if err != nil {
		if errors.Is(err, svcErrs.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
			return
		}

		c.JSON(http.StatusInternalServerError, response.ErrorInternal())
		return
	}

	c.String(http.StatusOK, "reset password email sent successfully")
}

type recoveryPasswordRequest struct {
	Token    string `json:"token"    validate:"required"`
	Password string `json:"password" validate:"required,password"  minLength:"8" maxLength:"32"  example:"YourV@lidPassw0rd!"`
}

// @Summary     Recovery password
// @Description Password recovery request
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body recoveryPasswordRequest true "Recovery password payload"
// @Success     200 {string} string
// @Failure     400 {object} response.ErrResponse
// @Failure     500 {object} response.ErrResponse
// @Router      /auth/recovery-password [post]
func (r *authRoutes) recoveryPassword(c *gin.Context) {
	var req recoveryPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	if errs := r.cv.ValidateStruct(req); errs != nil {
		c.JSON(http.StatusBadRequest, response.ErrorMap(errs))
		return
	}

	err := r.as.RecoveryPassword(c.Request.Context(), auth.RecoveryPasswordInput{
		Token:    req.Token,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, svcErrs.ErrTokenIsExpired) {
			c.JSON(http.StatusForbidden, response.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorInternal())
		return
	}

	c.String(http.StatusOK, "password updated successfully")
}
