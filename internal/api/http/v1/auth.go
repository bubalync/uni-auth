package v1

import (
	"errors"
	"github.com/bubalync/uni-auth/internal/lib/api/response"
	"github.com/bubalync/uni-auth/internal/service"
	"github.com/bubalync/uni-auth/internal/service/auth"
	"github.com/bubalync/uni-auth/internal/service/svcErrs"
	"github.com/bubalync/uni-auth/pkg/logger/sl"
	"github.com/bubalync/uni-auth/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type authRoutes struct {
	as service.Auth
	l  *slog.Logger
	cv *validator.CustomValidator
}

func NewAuthRoutes(g *gin.RouterGroup, log *slog.Logger, cv *validator.CustomValidator, authService service.Auth) {
	r := &authRoutes{authService, log, cv}

	g.POST("/sign-up", r.signUp)
	g.POST("/sign-in", r.signIn)
	g.POST("/reset-password", r.signIn)
}

type signUpRequest struct {
	Email    string `json:"email"    validate:"required,email,min=4,max=50"     example:"email@example.com"`
	Password string `json:"password" validate:"required,password" minLength:"8" maxLength:"32" example:"YourV@lidPassw0rd!"`
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
	const op = "api.http.v1.user.register"
	log := r.l.With(slog.String("op", op))

	var req signUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	if errs := r.cv.ValidateStruct(req); errs != nil {
		log.Error("Request validation error", sl.ErrMap(errs))

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

		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, &signUpResponse{
		Id: id,
	})
}

func (r *authRoutes) signIn(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r *userRoutes) resetPassword(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
