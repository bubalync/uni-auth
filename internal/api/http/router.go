package http

import (
	_ "github.com/bubalync/uni-auth/docs" // Swagger docs.
	v1 "github.com/bubalync/uni-auth/internal/api/http/v1"
	"github.com/bubalync/uni-auth/internal/config"
	"github.com/bubalync/uni-auth/internal/services"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log/slog"
)

// NewRouter -.
// Swagger spec:
//
//	@title			Universal authorization service API
//	@description	Authorization, registration, etc...
//	@version		1.0
//	@host			localhost:8080
//	@BasePath		/
func NewRouter(r *gin.Engine, cfg *config.Config, log *slog.Logger, us services.User) {
	// Middleware
	r.Use(gin.Recovery())
	r.Use(sloggin.New(log))

	// Swagger
	if *cfg.Swagger.Enabled {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	// Routers
	apiV1Group := r.Group("/api/v1")
	{
		v1.NewUserRoutes(log, apiV1Group, us)
	}
}
