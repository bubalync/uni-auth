package v1

import (
	"github.com/bubalync/uni-auth/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log/slog"
)

type userRoutes struct {
	us services.User
	l  *slog.Logger
	v  *validator.Validate
}

func NewUserRoutes(l *slog.Logger, v1Group *gin.RouterGroup, us services.User) {
	r := &userRoutes{us, l, validator.New(validator.WithRequiredStructEnabled())}

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

func (r *userRoutes) register(c *gin.Context) {
	//TODO implement me
	panic("implement me")
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
