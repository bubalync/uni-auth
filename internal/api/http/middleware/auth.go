package middleware

import (
	"github.com/bubalync/uni-auth/internal/lib/api/response"
	"github.com/bubalync/uni-auth/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	UserIdKey = "user_id"
)

type AuthMiddleware struct {
	authService service.Auth
}

func NewAuthMiddleware(authService service.Auth) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) UserIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(response.ErrInvalidAuthHeader.Error()))
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		uid, err := m.authService.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(response.ErrInvalidToken.Error()))
			return
		}
		c.Set(UserIdKey, uid)
		c.Next()
	}
}
