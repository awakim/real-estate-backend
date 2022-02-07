// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/awakim/immoblock-backend/cache/redis (interfaces: Cache)

// Package mockcache is a generated GoMock package.
package mockcache

import (
	context "context"
	reflect "reflect"
	time "time"

	token "github.com/awakim/immoblock-backend/token"
	gomock "github.com/golang/mock/gomock"
)

// MockCache is a mock of Cache interface.
type MockCache struct {
	ctrl     *gomock.Controller
	recorder *MockCacheMockRecorder
}

// MockCacheMockRecorder is the mock recorder for MockCache.
type MockCacheMockRecorder struct {
	mock *MockCache
}

// NewMockCache creates a new mock instance.
func NewMockCache(ctrl *gomock.Controller) *MockCache {
	mock := &MockCache{ctrl: ctrl}
	mock.recorder = &MockCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCache) EXPECT() *MockCacheMockRecorder {
	return m.recorder
}

// DeleteRefreshToken mocks base method.
func (m *MockCache) DeleteRefreshToken(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRefreshToken", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRefreshToken indicates an expected call of DeleteRefreshToken.
func (mr *MockCacheMockRecorder) DeleteRefreshToken(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRefreshToken", reflect.TypeOf((*MockCache)(nil).DeleteRefreshToken), arg0, arg1, arg2)
}

// IsRateLimited mocks base method.
func (m *MockCache) IsRateLimited(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRateLimited", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsRateLimited indicates an expected call of IsRateLimited.
func (mr *MockCacheMockRecorder) IsRateLimited(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRateLimited", reflect.TypeOf((*MockCache)(nil).IsRateLimited), arg0, arg1)
}

// IsRevoked mocks base method.
func (m *MockCache) IsRevoked(arg0 context.Context, arg1 token.Payload) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRevoked", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsRevoked indicates an expected call of IsRevoked.
func (mr *MockCacheMockRecorder) IsRevoked(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRevoked", reflect.TypeOf((*MockCache)(nil).IsRevoked), arg0, arg1)
}

// LogoutUser mocks base method.
func (m *MockCache) LogoutUser(arg0 context.Context, arg1, arg2 token.Payload) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogoutUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// LogoutUser indicates an expected call of LogoutUser.
func (mr *MockCacheMockRecorder) LogoutUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogoutUser", reflect.TypeOf((*MockCache)(nil).LogoutUser), arg0, arg1, arg2)
}

// SetTokenData mocks base method.
func (m *MockCache) SetTokenData(arg0 context.Context, arg1 token.Payload, arg2 time.Duration, arg3 token.Payload, arg4 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetTokenData", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetTokenData indicates an expected call of SetTokenData.
func (mr *MockCacheMockRecorder) SetTokenData(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTokenData", reflect.TypeOf((*MockCache)(nil).SetTokenData), arg0, arg1, arg2, arg3, arg4)
}
