package v1

import (
	"errors"
	"github.com/bubalync/uni-auth/internal/api/http/middleware"
	"github.com/bubalync/uni-auth/internal/lib/api/response"
	"github.com/bubalync/uni-auth/internal/service"
	"github.com/bubalync/uni-auth/internal/service/svcErrs"
	"github.com/bubalync/uni-auth/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type userRoutes struct {
	us service.User
	l  *slog.Logger
	cv *validator.CustomValidator
}

func NewUserRoutes(g *gin.RouterGroup, log *slog.Logger, cv *validator.CustomValidator, us service.User) {
	r := &userRoutes{us, log, cv}

	g.GET("/", r.user)
	g.GET("/:user_id", r.userById)
	g.PUT("/", r.update)
	g.DELETE("/", r.delete)
	g.POST("/logout", r.logout)
}

func userIdFromContext(c *gin.Context) uuid.UUID {
	return c.MustGet(middleware.UserIdKey).(uuid.UUID)
}

func (r *userRoutes) update(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

// @Summary     Current user info
// @Description Get information about the current user
// @Tags        users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} entity.User
// @Failure     400 {object} response.ErrResponse
// @Failure     500 {object} response.ErrResponse
// @Router      /api/v1/users [get]
func (r *userRoutes) user(c *gin.Context) {
	uid := userIdFromContext(c)

	user, err := r.us.UserById(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorInternal())
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary     User info by id
// @Description Get information about user by id
// @Tags        users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       user_id path string true "User id (UUID)"
// @Success     200 {object} entity.User
// @Failure     400 {object} response.ErrResponse
// @Failure     500 {object} response.ErrResponse
// @Router      /api/v1/users/{user_id} [get]
func (r *userRoutes) userById(c *gin.Context) {
	id := c.Param("user_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Error("user_id is required"))
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("user_id is invalid uuid"))
		return
	}

	user, err := r.us.UserById(c.Request.Context(), uid)
	if err != nil {
		if errors.Is(err, svcErrs.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorInternal())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (r *userRoutes) delete(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r *userRoutes) logout(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
