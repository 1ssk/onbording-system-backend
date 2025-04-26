package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/glebkasilov/grpc-manager/internal/domain/models"
	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddMeeting mocks base method.
func (m *MockStorage) AddMeeting(arg0 context.Context, arg1 *models.Meeting) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMeeting", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddMeeting indicates an expected call of AddMeeting.
func (mr *MockStorageMockRecorder) AddMeeting(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(
		mr.mock,
		"AddMeeting",
		reflect.TypeOf((*MockStorage)(nil).AddMeeting),
		arg0,
		arg1,
	)
}

// AddUser mocks base method.
func (m *MockStorage) AddUser(arg0 context.Context, arg1 *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUser indicates an expected call of AddUser.
func (mr *MockStorageMockRecorder) AddUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(
		mr.mock,
		"AddUser",
		reflect.TypeOf((*MockStorage)(nil).AddUser),
		arg0,
		arg1,
	)
}

// ... (аналогично для всех методов интерфейса Storage)
