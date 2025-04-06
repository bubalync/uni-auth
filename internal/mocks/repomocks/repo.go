// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repo/repo.go
//
// Generated by this command:
//
//	mockgen -source=internal/repo/repo.go -destination=internal/mocks/repomocks/repo.go -package=repomocks
//

// Package repomocks is a generated GoMock package.
package repomocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/bubalync/uni-auth/internal/entity"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

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

// Create mocks base method.
func (m *MockUser) Create(ctx context.Context, u *entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockUserMockRecorder) Create(ctx, u any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUser)(nil).Create), ctx, u)
}

// Delete mocks base method.
func (m *MockUser) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUserMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUser)(nil).Delete), ctx, id)
}

// ResetPassword mocks base method.
func (m *MockUser) ResetPassword(ctx context.Context, passwordHash []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetPassword", ctx, passwordHash)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetPassword indicates an expected call of ResetPassword.
func (mr *MockUserMockRecorder) ResetPassword(ctx, passwordHash any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetPassword", reflect.TypeOf((*MockUser)(nil).ResetPassword), ctx, passwordHash)
}

// Update mocks base method.
func (m *MockUser) Update(ctx context.Context, u *entity.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, u)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockUserMockRecorder) Update(ctx, u any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUser)(nil).Update), ctx, u)
}

// UserByEmail mocks base method.
func (m *MockUser) UserByEmail(ctx context.Context, email string) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserByEmail", ctx, email)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserByEmail indicates an expected call of UserByEmail.
func (mr *MockUserMockRecorder) UserByEmail(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserByEmail", reflect.TypeOf((*MockUser)(nil).UserByEmail), ctx, email)
}

// UserByEmailIsExists mocks base method.
func (m *MockUser) UserByEmailIsExists(ctx context.Context, email string) (*bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserByEmailIsExists", ctx, email)
	ret0, _ := ret[0].(*bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserByEmailIsExists indicates an expected call of UserByEmailIsExists.
func (mr *MockUserMockRecorder) UserByEmailIsExists(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserByEmailIsExists", reflect.TypeOf((*MockUser)(nil).UserByEmailIsExists), ctx, email)
}

// UserByID mocks base method.
func (m *MockUser) UserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserByID", ctx, id)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserByID indicates an expected call of UserByID.
func (mr *MockUserMockRecorder) UserByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserByID", reflect.TypeOf((*MockUser)(nil).UserByID), ctx, id)
}
