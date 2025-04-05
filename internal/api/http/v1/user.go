package v1

import (
	"github.com/bubalync/uni-auth/internal/service"
	"github.com/bubalync/uni-auth/pkg/validator"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type userRoutes struct {
	us service.User
	l  *slog.Logger
	cv *validator.CustomValidator
}

func NewUserRoutes(g *gin.RouterGroup, log *slog.Logger, cv *validator.CustomValidator, us service.User) {
	r := &userRoutes{us, log, cv}

	g.GET("/profile", r.user)
	g.PUT("/profile", r.update)
	g.DELETE("/profile", r.delete)
	g.POST("/logout", r.logout)

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
