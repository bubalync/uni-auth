package middleware

import (
	"errors"
	"github.com/bubalync/uni-auth/internal/lib/jwtgen"
	"github.com/bubalync/uni-auth/internal/mocks/servicemocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type MockBehaviour func(m *servicemocks.MockAuth)

	testCases := []struct {
		name             string
		accessToken      string
		mockBehaviour    MockBehaviour
		wantStatusCode   int
		wantResponseBody string
	}{
		{
			name:        "OK",
			accessToken: `Bearer valid_access_token`,
			mockBehaviour: func(a *servicemocks.MockAuth) {
				claims := &jwtgen.Claims{UserId: uuid.MustParse("0148edcd-e2a0-48b8-a47a-c6de5bbe4ed5")}
				a.EXPECT().ParseToken("valid_access_token").Return(claims, nil)
			},
			wantStatusCode:   200,
			wantResponseBody: `{"user_id":"0148edcd-e2a0-48b8-a47a-c6de5bbe4ed5"}`,
		},
		{
			name:             "authorization header empty",
			accessToken:      "",
			mockBehaviour:    func(a *servicemocks.MockAuth) {},
			wantStatusCode:   401,
			wantResponseBody: `{"errors":{"message":"invalid auth header"}}`,
		},
		{
			name:             "authorization header invalid",
			accessToken:      `_Bearer_ invalid_access_token`,
			mockBehaviour:    func(a *servicemocks.MockAuth) {},
			wantStatusCode:   401,
			wantResponseBody: `{"errors":{"message":"invalid auth header"}}`,
		},
		{
			name:        "auth service some error",
			accessToken: `Bearer valid_access_token`,
			mockBehaviour: func(a *servicemocks.MockAuth) {
				a.EXPECT().ParseToken("valid_access_token").Return(nil, errors.New("some error"))
			},
			wantStatusCode:   401,
			wantResponseBody: `{"errors":{"message":"invalid token"}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// init service mock
			as := servicemocks.NewMockAuth(ctrl)
			tc.mockBehaviour(as)

			// create middleware
			middleware := NewAuthMiddleware(as)

			// create test server
			gin.SetMode(gin.TestMode)
			r := gin.New()

			r.Use(middleware.UserIdentity())
			r.GET("/protected", func(c *gin.Context) {
				userID := c.MustGet(UserIdKey).(uuid.UUID)

				c.JSON(200, gin.H{UserIdKey: userID})
			})

			// create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", tc.accessToken)

			// execute request
			r.ServeHTTP(w, req)

			// check response
			assert.Equal(t, tc.wantStatusCode, w.Code)
			assert.Equal(t, tc.wantResponseBody, w.Body.String())
		})

	}
}
