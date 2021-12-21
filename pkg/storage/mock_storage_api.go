// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

// Package storage is a generated GoMock package.
package storage

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	types "k8s.io/apimachinery/pkg/types"
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

// CheckConfigMapEntry mocks base method.
func (m *MockStorage) CheckConfigMapEntry(arg0 string, arg1 types.NamespacedName) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckConfigMapEntry", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckConfigMapEntry indicates an expected call of CheckConfigMapEntry.
func (mr *MockStorageMockRecorder) CheckConfigMapEntry(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckConfigMapEntry", reflect.TypeOf((*MockStorage)(nil).CheckConfigMapEntry), arg0, arg1)
}

// DeleteConfigMapEntry mocks base method.
func (m *MockStorage) DeleteConfigMapEntry(arg0 string, arg1 types.NamespacedName) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteConfigMapEntry", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteConfigMapEntry indicates an expected call of DeleteConfigMapEntry.
func (mr *MockStorageMockRecorder) DeleteConfigMapEntry(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteConfigMapEntry", reflect.TypeOf((*MockStorage)(nil).DeleteConfigMapEntry), arg0, arg1)
}

// UpdateConfigMapEntry mocks base method.
func (m *MockStorage) UpdateConfigMapEntry(arg0, arg1 string, arg2 types.NamespacedName) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateConfigMapEntry", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateConfigMapEntry indicates an expected call of UpdateConfigMapEntry.
func (mr *MockStorageMockRecorder) UpdateConfigMapEntry(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateConfigMapEntry", reflect.TypeOf((*MockStorage)(nil).UpdateConfigMapEntry), arg0, arg1, arg2)
}
