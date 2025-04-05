package v1

import (
	"errors"
	"github.com/bubalync/uni-auth/internal/lib/api/response"
	"github.com/bubalync/uni-auth/internal/services"
	"github.com/bubalync/uni-auth/pkg/logger/sl"
	"github.com/bubalync/uni-auth/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type userRoutes struct {
	us services.User
	l  *slog.Logger
	cv *validator.CustomValidator
}

func NewUserRoutes(l *slog.Logger, v1Group *gin.RouterGroup, cv *validator.CustomValidator, us services.User) {
	r := &userRoutes{us, l, cv}

	v1Group.POST("/users/register", r.register)
	v1Group.POST("/users/login", r.login)
	v1Group.POST("/users/reset-password", r.resetPassword)

	protected := v1Group.Group("/")
	//protected.Use(authMiddleware.Authenticate)
	{
		protected.GET("/users/profile", r.user)
		protected.PUT("/users/profile", r.update)
		protected.DELETE("/users/profile", r.delete)
		protected.POST("/users/logout", r.logout)
	}
}

type registerRequest struct {
	Email    string `json:"email"     validate:"required,email"  example:"email@example.com"`
	Password string `json:"password"  validate:"required,password" minLength:"8" maxLength:"30" example:"YourV@lidPassw0rd!"`
}

type registerResponse struct {
	Id uuid.UUID `json:"id" example:"d13a75e2-3d21-4e57-9dc0-3a7f5bee4c25"`
}

// @Summary     Register
// @Description Register a new user
// @Tags  	    Users
// @Accept      json
// @Produce     json
// @Param       request body registerRequest true "Registration payload"
// @Success     200 {object} registerResponse
// @Failure     400 {object} response.ErrResponse
// @Failure     422 {object} response.ErrResponse
// @Failure     500 {object} response.ErrResponse
// @Router      /api/v1/users/register [post]
func (r *userRoutes) register(c *gin.Context) {
	const op = "api.http.v1.user.register"
	log := r.l.With(slog.String("op", op))

	var req registerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	if err := r.cv.ValidateStruct(req); err != nil {
		log.Error("Request validation error", sl.ErrMap(err))

		c.JSON(http.StatusBadRequest, response.ErrorMap(err))
		return
	}

	id, err := r.us.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			c.JSON(http.StatusUnprocessableEntity, response.Error(err.Error()))
			return
		}

		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))

		return
	}

	c.JSON(http.StatusCreated, &registerResponse{Id: id})
}

func (r *userRoutes) login(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r *userRoutes) resetPassword(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r *userRoutes) update(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r *userRoutes) user(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r *userRoutes) delete(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r *userRoutes) logout(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
