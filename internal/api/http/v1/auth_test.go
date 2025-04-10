package v1

import (
	"bytes"
	"context"
	"github.com/bubalync/uni-auth/internal/mocks/servicemocks"
	"github.com/bubalync/uni-auth/internal/service/auth"
	"github.com/bubalync/uni-auth/internal/service/svcErrs"
	"github.com/bubalync/uni-auth/pkg/logger"
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
		name            string
		args            args
		inputBody       string
		mockBehaviour   MockBehaviour
		wantStatusCode  int
		wantRequestBody string
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
			wantStatusCode:  201,
			wantRequestBody: `{"id":"0148edcd-e2a0-48b8-a47a-c6de5bbe4ed5"}`,
		},
		{
			name:            "Invalid password: not provided",
			args:            args{},
			inputBody:       `{"email": "test@example.com"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"errors":{"Password":"Password is a required"}}`,
		},
		{
			name:            "Invalid password: too short",
			args:            args{},
			inputBody:       `{"email": "test@example.com","password":"Qw!1"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: passwordErrMsg,
		},
		{
			name:            "Invalid password: too long",
			args:            args{},
			inputBody:       `{"email": "test@example.com","password":"Qwerty!123456789012345678901234567890"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: passwordErrMsg,
		},
		{
			name:            "Invalid password: no uppercase",
			args:            args{},
			inputBody:       `{"email": "test@example.com","password":"qwerty!1"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: passwordErrMsg,
		},
		{
			name:            "Invalid password: no lowercase",
			args:            args{},
			inputBody:       `{"email": "test@example.com","password":"QWERTY!1"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: passwordErrMsg,
		},
		{
			name:            "Invalid password: no digits",
			args:            args{},
			inputBody:       `{"email": "test@example.com","password":"Qwerty!!"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: passwordErrMsg,
		},
		{
			name:            "Invalid password: no special characters",
			args:            args{},
			inputBody:       `{"email": "test@example.com","password":"Qwerty11"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: passwordErrMsg,
		},
		{
			name:            "Invalid email: not provided",
			args:            args{},
			inputBody:       `{"password":"Qwerty!1"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"errors":{"Email":"Email is a required"}}`,
		},
		{
			name: "Invalid email: too long",
			args: args{},
			inputBody: `{"email":"testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttestte` +
				`sttesttesttesttesttesttesttesttesttesttesttesttesttesttesttetesttesttesttesttesttesttesttes` +
				`sttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest@e.com","password":"Qwerty!1"}`,
			mockBehaviour:   func(m *servicemocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"errors":{"Email":"'Email' must be shorter than 150"}}`,
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
			wantStatusCode:  422,
			wantRequestBody: `{"errors":{"message":"user already exists"}}`,
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
				m.EXPECT().CreateUser(args.ctx, args.input).Return(uuid.Nil, svcErrs.ErrInternal)
			},
			wantStatusCode:  500,
			wantRequestBody: `{"errors":{"message":"internal server error"}}`,
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

			// Log
			log := logger.New("local", "info")

			// create test server
			e := gin.New()

			cv := validator.NewCustomValidator()

			g := e.Group("/auth")
			NewAuthRoutes(g, log, cv, as)
			gin.SetMode(gin.ReleaseMode)

			// create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString(tc.inputBody))

			// execute request
			e.ServeHTTP(w, req)

			// check response
			assert.Equal(t, tc.wantStatusCode, w.Code)
			assert.Equal(t, tc.wantRequestBody, w.Body.String())
		})
	}
}
