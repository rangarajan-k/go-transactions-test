// Code generated by MockGen. DO NOT EDIT.
// Source: datastore/postgress_client.go

// Package datastore is a generated GoMock package.
package datastore

import (
	config "go-transactions-test/config"
	reflect "reflect"

	pg "github.com/go-pg/pg/v10"
	gomock "github.com/golang/mock/gomock"
)

// MockIPgClientInterface is a mock of IPgClientInterface interface.
type MockIPgClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockIPgClientInterfaceMockRecorder
}

// MockIPgClientInterfaceMockRecorder is the mock recorder for MockIPgClientInterface.
type MockIPgClientInterfaceMockRecorder struct {
	mock *MockIPgClientInterface
}

// NewMockIPgClientInterface creates a new mock instance.
func NewMockIPgClientInterface(ctrl *gomock.Controller) *MockIPgClientInterface {
	mock := &MockIPgClientInterface{ctrl: ctrl}
	mock.recorder = &MockIPgClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPgClientInterface) EXPECT() *MockIPgClientInterfaceMockRecorder {
	return m.recorder
}

// CreateSchema mocks base method.
func (m *MockIPgClientInterface) CreateSchema(db *pg.DB) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSchema", db)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSchema indicates an expected call of CreateSchema.
func (mr *MockIPgClientInterfaceMockRecorder) CreateSchema(db interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSchema", reflect.TypeOf((*MockIPgClientInterface)(nil).CreateSchema), db)
}

// GetDbClient mocks base method.
func (m *MockIPgClientInterface) GetDbClient() *pg.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDbClient")
	ret0, _ := ret[0].(*pg.DB)
	return ret0
}

// GetDbClient indicates an expected call of GetDbClient.
func (mr *MockIPgClientInterfaceMockRecorder) GetDbClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDbClient", reflect.TypeOf((*MockIPgClientInterface)(nil).GetDbClient))
}

// NewPgClient mocks base method.
func (m *MockIPgClientInterface) NewPgClient(config config.DBConfig) *pg.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewPgClient", config)
	ret0, _ := ret[0].(*pg.DB)
	return ret0
}

// NewPgClient indicates an expected call of NewPgClient.
func (mr *MockIPgClientInterfaceMockRecorder) NewPgClient(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewPgClient", reflect.TypeOf((*MockIPgClientInterface)(nil).NewPgClient), config)
}
