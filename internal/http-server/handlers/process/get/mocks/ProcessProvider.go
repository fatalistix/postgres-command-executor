// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	models "github.com/fatalistix/postgres-command-executor/internal/domain/models"
	uuid "github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

// ProcessProvider is an autogenerated mock type for the ProcessProvider type
type ProcessProvider struct {
	mock.Mock
}

// Process provides a mock function with given fields: id
func (_m *ProcessProvider) Process(id uuid.UUID) (*models.Process, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Process")
	}

	var r0 *models.Process
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*models.Process, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Process); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Process)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProcessProvider creates a new instance of ProcessProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProcessProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProcessProvider {
	mock := &ProcessProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}