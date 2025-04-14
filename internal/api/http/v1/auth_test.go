package v1

import (
	"bytes"
	"context"
	"github.com/bubalync/uni-auth/internal/mocks/servicemocks"
	"github.com/bubalync/uni-auth/internal/service/auth"
	"github.com/bubalync/uni-auth/internal/service/svcErrs"
	"github.com/bubalync/uni-auth/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthRoutes_SignUp(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.CreateUserInput
	}

	const passwordErrMsg = `{"errors":{"Password":"Password must be between 8 and 32 in length, contain at least 1 lowercase, 1 uppercase, 1 digits, and 1 special characters (!@#$%^\u0026*)"}}`

	type MockBehaviour func(m *servicemocks.MockAuth, args args)

	testCases := []struct {
		name             string
		args             args
		inputBody        string
		mockBehaviour    MockBehaviour
		wantStatusCode   int
		wantResponseBody string
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				input: auth.CreateUserInput{
					Email:    "test@example.com",
					Password: "Qwerty!1",
				},
			},
			inputBody: `{"email":"test@example.com","password":"Qwerty!1"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().CreateUser(args.ctx, args.input).Return(uuid.MustParse("0148edcd-e2a0-48b8-a47a-c6de5bbe4ed5"), nil)
			},
			wantStatusCode:   201,
			wantResponseBody: `{"id":"0148edcd-e2a0-48b8-a47a-c6de5bbe4ed5"}`,
		},
		{
			name:             "Invalid password: not provided",
			args:             args{},
			inputBody:        `{"email": "test@example.com"}`,
			mockBehaviour:    func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:   400,
			wantResponseBody: `{"errors":{"Password":"Password is a required"}}`,
		},
		{
			name:             "Invalid password: too short",
			args:             args{},
			inputBody:        `{"email": "test@example.com","password":"Qw!1"}`,
			mockBehaviour:    func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:   400,
			wantResponseBody: passwordErrMsg,
		},
		{
			name:             "Invalid password: too long",
			args:             args{},
			inputBody:        `{"email": "test@example.com","password":"Qwerty!123456789012345678901234567890"}`,
			mockBehaviour:    func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:   400,
			wantResponseBody: passwordErrMsg,
		},
		{
			name:             "Invalid password: no uppercase",
			args:             args{},
			inputBody:        `{"email": "test@example.com","password":"qwerty!1"}`,
			mockBehaviour:    func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:   400,
			wantResponseBody: passwordErrMsg,
		},
		{
			name:             "Invalid password: no lowercase",
			args:             args{},
			inputBody:        `{"email": "test@example.com","password":"QWERTY!1"}`,
			mockBehaviour:    func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:   400,
			wantResponseBody: passwordErrMsg,
		},
		{
			name:             "Invalid password: no digits",
			args:             args{},
			inputBody:        `{"email": "test@example.com","password":"Qwerty!!"}`,
			mockBehaviour:    func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:   400,
			wantResponseBody: passwordErrMsg,
		},
		{
			name:             "Invalid password: no special characters",
			args:             args{},
			inputBody:        `{"email": "test@example.com","password":"Qwerty11"}`,
			mockBehaviour:    func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:   400,
			wantResponseBody: passwordErrMsg,
		},
		{
			name:             "Invalid email: not provided",
			args:             args{},
			inputBody:        `{"password":"Qwerty!1"}`,
			mockBehaviour:    func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:   400,
			wantResponseBody: `{"errors":{"Email":"Email is a required"}}`,
		},
		{
			name: "Invalid email: too long",
			args: args{},
			inputBody: `{"email":"testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttestte` +
				`sttesttesttesttesttesttesttesttesttesttesttesttesttesttesttetesttesttesttesttesttesttesttes` +
				`sttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest@e.com","password":"Qwerty!1"}`,
			mockBehaviour:    func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:   400,
			wantResponseBody: `{"errors":{"Email":"'Email' must be shorter than 150"}}`,
		},
		{
			name: "Auth service error: already exists",
			args: args{
				ctx: context.Background(),
				input: auth.CreateUserInput{
					Email:    "test@example.com",
					Password: "Qwerty!1",
				},
			},
			inputBody: `{"email": "test@example.com","password":"Qwerty!1"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().CreateUser(args.ctx, args.input).Return(uuid.Nil, svcErrs.ErrUserAlreadyExists)
			},
			wantStatusCode:   422,
			wantResponseBody: `{"errors":{"message":"user already exists"}}`,
		},
		{
			name: "Internal server error",
			args: args{
				ctx: context.Background(),
				input: auth.CreateUserInput{
					Email:    "test@example.com",
					Password: "Qwerty!1",
				},
			},
			inputBody: `{"email": "test@example.com", "password":"Qwerty!1"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().CreateUser(args.ctx, args.input).Return(uuid.Nil, svcErrs.ErrCannotCreateUser)
			},
			wantStatusCode:   500,
			wantResponseBody: `{"errors":{"message":"internal server error"}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// init service mock
			as := servicemocks.NewMockAuth(ctrl)
			tc.mockBehaviour(as, tc.args)

			// create test server
			e := gin.New()

			cv := validator.NewCustomValidator()

			g := e.Group("/auth")
			NewAuthRoutes(g, cv, as)
			gin.SetMode(gin.ReleaseMode)

			// create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString(tc.inputBody))

			// execute request
			e.ServeHTTP(w, req)

			// check response
			assert.Equal(t, tc.wantStatusCode, w.Code)
			assert.Equal(t, tc.wantResponseBody, w.Body.String())
		})
	}
}

func TestAuthRoutes_SignIn(t *testing.T) {
	type args struct {
		ctx   context.Context
		input auth.GenerateTokenInput
	}

	type MockBehaviour func(m *servicemocks.MockAuth, args args)

	testCases := []struct {
		name             string
		args             args
		inputBody        string
		mockBehaviour    MockBehaviour
		wantStatusCode   int
		wantResponseBody string
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				input: auth.GenerateTokenInput{
					Email:    "test@example.com",
					Password: "Qwerty!1",
				},
			},
			inputBody: `{"email":"test@example.com","password":"Qwerty!1"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(args.ctx, args.input).
					Return(auth.GenerateTokenOutput{AccessToken: "1", RefreshToken: "2"}, nil)
			},
			wantStatusCode:   200,
			wantResponseBody: `{"access_token":"1","refresh_token":"2"}`,
		},
		{
			name:             "Invalid password: not provided",
			args:             args{},
			inputBody:        `{"email": "test@example.com"}`,
			mockBehaviour:    func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:   400,
			wantResponseBody: `{"errors":{"Password":"Password is a required"}}`,
		},
		{
			name:             "Invalid email: not provided",
			args:             args{},
			inputBody:        `{"password":"Qwerty!1"}`,
			mockBehaviour:    func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:   400,
			wantResponseBody: `{"errors":{"Email":"Email is a required"}}`,
		},
		{
			name: "Invalid email: too long",
			args: args{},
			inputBody: `{"email":"testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttestte` +
				`sttesttesttesttesttesttesttesttesttesttesttesttesttesttesttetesttesttesttesttesttesttesttes` +
				`sttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest@e.com","password":"Qwerty!1"}`,
			mockBehaviour:    func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:   400,
			wantResponseBody: `{"errors":{"Email":"'Email' must be shorter than 150"}}`,
		},
		{
			name: "Auth service error: invalid credentials",
			args: args{
				ctx: context.Background(),
				input: auth.GenerateTokenInput{
					Email:    "test@example.com",
					Password: "123",
				},
			},
			inputBody: `{"email": "test@example.com","password":"123"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(args.ctx, args.input).Return(auth.GenerateTokenOutput{}, svcErrs.ErrInvalidCredentials)
			},
			wantStatusCode:   401,
			wantResponseBody: `{"errors":{"message":"invalid credentials"}}`,
		},
		{
			name: "Internal server error",
			args: args{
				ctx: context.Background(),
				input: auth.GenerateTokenInput{
					Email:    "test@example.com",
					Password: "Qwerty!1",
				},
			},
			inputBody: `{"email": "test@example.com", "password":"Qwerty!1"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(args.ctx, args.input).Return(auth.GenerateTokenOutput{}, svcErrs.ErrCannotCreateUser)
			},
			wantStatusCode:   500,
			wantResponseBody: `{"errors":{"message":"internal server error"}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// init service mock
			as := servicemocks.NewMockAuth(ctrl)
			tc.mockBehaviour(as, tc.args)

			// create test server
			gin.SetMode(gin.TestMode)
			e := gin.New()

			cv := validator.NewCustomValidator()

			g := e.Group("/auth")
			NewAuthRoutes(g, cv, as)

			// create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBufferString(tc.inputBody))

			// execute request
			e.ServeHTTP(w, req)

			// check response
			assert.Equal(t, tc.wantStatusCode, w.Code)
			assert.Equal(t, tc.wantResponseBody, w.Body.String())
		})
	}
}

func TestAuthRoutes_Refresh(t *testing.T) {
	type args struct {
		ctx   context.Context
		token string
	}

	type MockBehaviour func(m *servicemocks.MockAuth, args args)

	testCases := []struct {
		name             string
		args             args
		inputBody        string
		mockBehaviour    MockBehaviour
		wantStatusCode   int
		wantResponseBody string
	}{
		{
			name: "OK",
			args: args{
				ctx:   context.Background(),
				token: "valid_token",
			},
			inputBody: `{"token":"valid_token"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().Refresh(args.ctx, args.token).
					Return(auth.GenerateTokenOutput{AccessToken: "111", RefreshToken: "222"}, nil)
			},
			wantStatusCode:   200,
			wantResponseBody: `{"access_token":"111","refresh_token":"222"}`,
		},
		{
			name:             "Invalid token: not provided",
			args:             args{},
			inputBody:        `{"token": ""}`,
			mockBehaviour:    func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:   400,
			wantResponseBody: `{"errors":{"Token":"Token is a required"}}`,
		},
		{
			name: "Auth service error: invalid token",
			args: args{
				ctx:   context.Background(),
				token: "invalid_token",
			},
			inputBody: `{"token": "invalid_token"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().Refresh(args.ctx, args.token).Return(auth.GenerateTokenOutput{}, svcErrs.ErrCannotParseToken)
			},
			wantStatusCode:   401,
			wantResponseBody: `{"errors":{"message":"cannot parse token"}}`,
		},
		{
			name: "Auth service error: token is expired",
			args: args{
				ctx:   context.Background(),
				token: "valid_token",
			},
			inputBody: `{"token": "valid_token"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().Refresh(args.ctx, args.token).Return(auth.GenerateTokenOutput{}, svcErrs.ErrTokenIsExpired)
			},
			wantStatusCode:   401,
			wantResponseBody: `{"errors":{"message":"token is expired"}}`,
		},
		{
			name: "Auth service error: Internal server error",
			args: args{
				ctx:   context.Background(),
				token: "valid_token",
			},
			inputBody: `{"token": "valid_token"}`,
			mockBehaviour: func(m *servicemocks.MockAuth, args args) {
				m.EXPECT().Refresh(args.ctx, args.token).Return(auth.GenerateTokenOutput{}, svcErrs.ErrCannotSignToken)
			},
			wantStatusCode:   500,
			wantResponseBody: `{"errors":{"message":"internal server error"}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// init service mock
			as := servicemocks.NewMockAuth(ctrl)
			tc.mockBehaviour(as, tc.args)

			// create test server
			gin.SetMode(gin.TestMode)
			e := gin.New()

			cv := validator.NewCustomValidator()

			g := e.Group("/auth")
			NewAuthRoutes(g, cv, as)

			// create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBufferString(tc.inputBody))

			// execute request
			e.ServeHTTP(w, req)

			// check response
			assert.Equal(t, tc.wantStatusCode, w.Code)
			assert.Equal(t, tc.wantResponseBody, w.Body.String())
		})
	}
}
