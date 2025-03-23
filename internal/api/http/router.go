package http

import (
	_ "github.com/bubalync/uni-auth/docs"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log/slog"
	"net/http"
)

// FillRouter -.
// Swagger spec:
//
//	@title			Universal authorization service API
//	@description	Authorization, registration, etc...
//	@version		1.0
//	@host			localhost:8080
//	@BasePath		/
func FillRouter(r *gin.Engine, log *slog.Logger) *gin.Engine {
	// Middleware
	r.Use(gin.Recovery())
	r.Use(sloggin.New(log))

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Routers
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/ping", pongHandler())
	}

	return r
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Hello world
// @Router /api/ping [get]
func pongHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	}
}
