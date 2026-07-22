package mock

import (
	"context"
	"reflect"

	"github.com/aclgo/grpc-jwt/internal/models"
	"github.com/golang/mock/gomock"
)

type MockUserPGRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserPGRepoRecorder
}

type MockUserPGRepoRecorder struct {
	mock *MockUserPGRepo
}

func NewMockUserPGRepo(ctrl *gomock.Controller) *MockUserPGRepo {
	mock := &MockUserPGRepo{ctrl: ctrl}
	mock.recorder = &MockUserPGRepoRecorder{mock}

	return mock
}

func (m *MockUserPGRepo) EXPECT() *MockUserPGRepoRecorder {
	return m.recorder
}

func (m *MockUserPGRepo) Add(ctx context.Context, user *models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, user)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockUserPGRepoRecorder) Add(ctx context.Context, user any) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(m.mock, "Add", reflect.TypeOf((*MockUserPGRepo)(nil).Add), ctx, user)
}

func (m *MockUserPGRepo) FindByID(ctx context.Context, userID string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, userID)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)

	return ret0, ret1
}

func (m *MockUserPGRepoRecorder) FindByID(ctx context.Context, userID any) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(
		m.mock,
		"FindByID",
		reflect.TypeOf((*MockUserPGRepo)(nil).FindByID),
		ctx, userID,
	)
}

func (m *MockUserPGRepo) FindAll(ctx context.Context, pagination *models.Pagination) (*models.ParamsFindAllResponse, error) {
	return nil, nil
}

func (m *MockUserPGRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	m.ctrl.T.Helper()

	ret := m.ctrl.Call(m, "FindByEmail", ctx, email)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)

	return ret0, ret1
}

func (m *MockUserPGRepoRecorder) FindByEmail(ctx context.Context, email any) *gomock.Call {
	m.mock.ctrl.T.Helper()

	return m.mock.ctrl.RecordCallWithMethodType(
		m.mock,
		"FindByEmail",
		reflect.TypeOf((*MockUserPGRepo)(nil).FindByEmail),
		ctx,
		email,
	)
}

func (m *MockUserPGRepo) Update(ctx context.Context, user *models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, user)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)

	return ret0, ret1
}

func (m *MockUserPGRepoRecorder) Update(ctx context.Context, user *models.User) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(
		m.mock,
		"Update",
		reflect.TypeOf((*MockUserPGRepo)(nil).Update),
		ctx,
		user,
	)
}

func (m *MockUserPGRepo) Delete(ctx context.Context, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, userID)
	ret0, _ := ret[0].(error)

	return ret0
}

func (m *MockUserPGRepoRecorder) Delete(ctx context.Context, userID string) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(
		m.mock,
		"Delete",
		reflect.TypeOf((*MockUserPGRepo)(nil).Delete),
		ctx,
		userID,
	)
}
