// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sndzhng/gin-template/internal/usecase (interfaces: Admin)

// Package usecasemock is a generated GoMock package.
package usecasemock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/sndzhng/gin-template/internal/entity"
)

// MockAdmin is a mock of Admin interface.
type MockAdmin struct {
	ctrl     *gomock.Controller
	recorder *MockAdminMockRecorder
}

// MockAdminMockRecorder is the mock recorder for MockAdmin.
type MockAdminMockRecorder struct {
	mock *MockAdmin
}

// NewMockAdmin creates a new mock instance.
func NewMockAdmin(ctrl *gomock.Controller) *MockAdmin {
	mock := &MockAdmin{ctrl: ctrl}
	mock.recorder = &MockAdminMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdmin) EXPECT() *MockAdminMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAdmin) Create(arg0 entity.Admin) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockAdminMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAdmin)(nil).Create), arg0)
}

// Delete mocks base method.
func (m *MockAdmin) Delete(arg0 entity.Admin) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAdminMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAdmin)(nil).Delete), arg0)
}

// Get mocks base method.
func (m *MockAdmin) Get(arg0 entity.Admin) (entity.Admin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(entity.Admin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockAdminMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAdmin)(nil).Get), arg0)
}

// GetAll mocks base method.
func (m *MockAdmin) GetAll(arg0 *entity.AdminFilter, arg1 *entity.SortOrder, arg2 *entity.Pagination) ([]entity.Admin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0, arg1, arg2)
	ret0, _ := ret[0].([]entity.Admin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockAdminMockRecorder) GetAll(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockAdmin)(nil).GetAll), arg0, arg1, arg2)
}

// Initial mocks base method.
func (m *MockAdmin) Initial() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initial")
	ret0, _ := ret[0].(error)
	return ret0
}

// Initial indicates an expected call of Initial.
func (mr *MockAdminMockRecorder) Initial() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initial", reflect.TypeOf((*MockAdmin)(nil).Initial))
}

// Update mocks base method.
func (m *MockAdmin) Update(arg0 entity.Admin) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockAdminMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAdmin)(nil).Update), arg0)
}
