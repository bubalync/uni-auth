// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/service.go
//
// Generated by this command:
//
//	mockgen -source=internal/service/service.go -destination=internal/mocks/servicemocks/service.go -package=servicemocks
//

// Package servicemocks is a generated GoMock package.
package servicemocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/bubalync/uni-auth/internal/entity"
	auth "github.com/bubalync/uni-auth/internal/service/auth"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockAuth is a mock of Auth interface.
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
	isgomock struct{}
}

// MockAuthMockRecorder is the mock recorder for MockAuth.
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance.
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockAuth) CreateUser(ctx context.Context, input auth.CreateUserInput) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, input)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthMockRecorder) CreateUser(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuth)(nil).CreateUser), ctx, input)
}

// GenerateToken mocks base method.
func (m *MockAuth) GenerateToken(ctx context.Context, input auth.GenerateTokenInput) (auth.GenerateTokenOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", ctx, input)
	ret0, _ := ret[0].(auth.GenerateTokenOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockAuthMockRecorder) GenerateToken(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockAuth)(nil).GenerateToken), ctx, input)
}

// ParseToken mocks base method.
func (m *MockAuth) ParseToken(token string) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", token)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockAuthMockRecorder) ParseToken(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockAuth)(nil).ParseToken), token)
}

// Refresh mocks base method.
func (m *MockAuth) Refresh(ctx context.Context, token string) (auth.GenerateTokenOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", ctx, token)
	ret0, _ := ret[0].(auth.GenerateTokenOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Refresh indicates an expected call of Refresh.
func (mr *MockAuthMockRecorder) Refresh(ctx, token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockAuth)(nil).Refresh), ctx, token)
}

// ResetPassword mocks base method.
func (m *MockAuth) ResetPassword(ctx context.Context, input auth.ResetPasswordInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetPassword", ctx, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetPassword indicates an expected call of ResetPassword.
func (mr *MockAuthMockRecorder) ResetPassword(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetPassword", reflect.TypeOf((*MockAuth)(nil).ResetPassword), ctx, input)
}

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
	isgomock struct{}
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockUser) Delete(ctx context.Context, u entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUserMockRecorder) Delete(ctx, u any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUser)(nil).Delete), ctx, u)
}

// Logout mocks base method.
func (m *MockUser) Logout(ctx context.Context, u entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout.
func (mr *MockUserMockRecorder) Logout(ctx, u any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockUser)(nil).Logout), ctx, u)
}

// Update mocks base method.
func (m *MockUser) Update(ctx context.Context, u entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUserMockRecorder) Update(ctx, u any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUser)(nil).Update), ctx, u)
}

// UserByEmail mocks base method.
func (m *MockUser) UserByEmail(ctx context.Context, email string) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserByEmail", ctx, email)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserByEmail indicates an expected call of UserByEmail.
func (mr *MockUserMockRecorder) UserByEmail(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserByEmail", reflect.TypeOf((*MockUser)(nil).UserByEmail), ctx, email)
}

// UserByID mocks base method.
func (m *MockUser) UserByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserByID", ctx, id)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserByID indicates an expected call of UserByID.
func (mr *MockUserMockRecorder) UserByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserByID", reflect.TypeOf((*MockUser)(nil).UserByID), ctx, id)
}
