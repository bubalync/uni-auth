package http

import (
	_ "github.com/bubalync/uni-auth/docs" // Swagger docs.
	"github.com/bubalync/uni-auth/internal/api/http/middleware"
	v1 "github.com/bubalync/uni-auth/internal/api/http/v1"
	"github.com/bubalync/uni-auth/internal/config"
	"github.com/bubalync/uni-auth/internal/service"
	"github.com/bubalync/uni-auth/pkg/validator"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log/slog"
)

// NewRouter -.
// Swagger spec:
// @title                        Universal authorization service API
// @description                  Authorization, registration, etc...
// @version                      1.0
// @host                         localhost:8080
// @securityDefinitions.apikey   BearerAuth
// @scheme                       bearer
// @bearerFormat                 JWT
// @in                           header
// @name                         Authorization
// @BasePath                     /
func NewRouter(handler *gin.Engine, cfg *config.Config, log *slog.Logger, services *service.Services) {
	// Middleware
	handler.Use(gin.Recovery())
	handler.Use(sloggin.New(log))

	// Swagger
	if *cfg.Swagger.Enabled {
		handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	cv := validator.NewCustomValidator()

	// Routes
	authGroup := handler.Group("/auth")
	{
		v1.NewAuthRoutes(authGroup, cv, services.Auth)
	}

	authMiddleware := middleware.NewAuthMiddleware(services.Auth)
	v1Group := handler.Group("/api/v1", authMiddleware.UserIdentity())
	{
		v1.NewUserRoutes(v1Group.Group("/users"), log, cv, services.User)
	}
}
