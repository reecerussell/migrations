// Code generated by MockGen. DO NOT EDIT.
// Source: ../provider.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	migrations "github.com/reecerussell/migrations"
	reflect "reflect"
)

// MockProvider is a mock of Provider interface
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return m.recorder
}

// GetAppliedMigrations mocks base method
func (m *MockProvider) GetAppliedMigrations(ctx context.Context) ([]*migrations.Migration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAppliedMigrations", ctx)
	ret0, _ := ret[0].([]*migrations.Migration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAppliedMigrations indicates an expected call of GetAppliedMigrations
func (mr *MockProviderMockRecorder) GetAppliedMigrations(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAppliedMigrations", reflect.TypeOf((*MockProvider)(nil).GetAppliedMigrations), ctx)
}

// Apply mocks base method
func (m_2 *MockProvider) Apply(ctx context.Context, m *migrations.Migration) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Apply", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Apply indicates an expected call of Apply
func (mr *MockProviderMockRecorder) Apply(ctx, m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Apply", reflect.TypeOf((*MockProvider)(nil).Apply), ctx, m)
}

// Rollback mocks base method
func (m_2 *MockProvider) Rollback(ctx context.Context, m *migrations.Migration) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Rollback", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Rollback indicates an expected call of Rollback
func (mr *MockProviderMockRecorder) Rollback(ctx, m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockProvider)(nil).Rollback), ctx, m)
}