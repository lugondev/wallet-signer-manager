// Code generated by MockGen. DO NOT EDIT.
// Source: keys.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	entities "github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockStore is a mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Info mocks base method
func (m *MockStore) Info(arg0 context.Context) (*entities.StoreInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Info", arg0)
	ret0, _ := ret[0].(*entities.StoreInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Info indicates an expected call of Info
func (mr *MockStoreMockRecorder) Info(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockStore)(nil).Info), arg0)
}

// Create mocks base method
func (m *MockStore) Create(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, id, alg, attr)
	ret0, _ := ret[0].(*entities.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockStoreMockRecorder) Create(ctx, id, alg, attr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockStore)(nil).Create), ctx, id, alg, attr)
}

// Import mocks base method
func (m *MockStore) Import(ctx context.Context, id, privKey string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Import", ctx, id, privKey, alg, attr)
	ret0, _ := ret[0].(*entities.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Import indicates an expected call of Import
func (mr *MockStoreMockRecorder) Import(ctx, id, privKey, alg, attr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Import", reflect.TypeOf((*MockStore)(nil).Import), ctx, id, privKey, alg, attr)
}

// Get mocks base method
func (m *MockStore) Get(ctx context.Context, id, version string) (*entities.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id, version)
	ret0, _ := ret[0].(*entities.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockStoreMockRecorder) Get(ctx, id, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStore)(nil).Get), ctx, id, version)
}

// List mocks base method
func (m *MockStore) List(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockStoreMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockStore)(nil).List), ctx)
}

// Update mocks base method
func (m *MockStore) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, attr)
	ret0, _ := ret[0].(*entities.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockStoreMockRecorder) Update(ctx, id, attr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockStore)(nil).Update), ctx, id, attr)
}

// Refresh mocks base method
func (m *MockStore) Refresh(ctx context.Context, id string, expirationDate time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", ctx, id, expirationDate)
	ret0, _ := ret[0].(error)
	return ret0
}

// Refresh indicates an expected call of Refresh
func (mr *MockStoreMockRecorder) Refresh(ctx, id, expirationDate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockStore)(nil).Refresh), ctx, id, expirationDate)
}

// Delete mocks base method
func (m *MockStore) Delete(ctx context.Context, id string, versions ...string) (*entities.Key, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, id}
	for _, a := range versions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(*entities.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete
func (mr *MockStoreMockRecorder) Delete(ctx, id interface{}, versions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, id}, versions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStore)(nil).Delete), varargs...)
}

// GetDeleted mocks base method
func (m *MockStore) GetDeleted(ctx context.Context, id string) (*entities.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeleted", ctx, id)
	ret0, _ := ret[0].(*entities.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeleted indicates an expected call of GetDeleted
func (mr *MockStoreMockRecorder) GetDeleted(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeleted", reflect.TypeOf((*MockStore)(nil).GetDeleted), ctx, id)
}

// ListDeleted mocks base method
func (m *MockStore) ListDeleted(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListDeleted", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListDeleted indicates an expected call of ListDeleted
func (mr *MockStoreMockRecorder) ListDeleted(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListDeleted", reflect.TypeOf((*MockStore)(nil).ListDeleted), ctx)
}

// Undelete mocks base method
func (m *MockStore) Undelete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Undelete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Undelete indicates an expected call of Undelete
func (mr *MockStoreMockRecorder) Undelete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Undelete", reflect.TypeOf((*MockStore)(nil).Undelete), ctx, id)
}

// Destroy mocks base method
func (m *MockStore) Destroy(ctx context.Context, id string, versions ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, id}
	for _, a := range versions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Destroy", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy
func (mr *MockStoreMockRecorder) Destroy(ctx, id interface{}, versions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, id}, versions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockStore)(nil).Destroy), varargs...)
}

// Sign mocks base method
func (m *MockStore) Sign(ctx context.Context, id, data, version string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", ctx, id, data, version)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sign indicates an expected call of Sign
func (mr *MockStoreMockRecorder) Sign(ctx, id, data, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockStore)(nil).Sign), ctx, id, data, version)
}

// Encrypt mocks base method
func (m *MockStore) Encrypt(ctx context.Context, id, data string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Encrypt", ctx, id, data)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Encrypt indicates an expected call of Encrypt
func (mr *MockStoreMockRecorder) Encrypt(ctx, id, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Encrypt", reflect.TypeOf((*MockStore)(nil).Encrypt), ctx, id, data)
}

// Decrypt mocks base method
func (m *MockStore) Decrypt(ctx context.Context, id, data string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decrypt", ctx, id, data)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Decrypt indicates an expected call of Decrypt
func (mr *MockStoreMockRecorder) Decrypt(ctx, id, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decrypt", reflect.TypeOf((*MockStore)(nil).Decrypt), ctx, id, data)
}