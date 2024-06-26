// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	uuid "github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

// CommandExecutionStarter is an autogenerated mock type for the CommandExecutionStarter type
type CommandExecutionStarter struct {
	mock.Mock
}

// StartCommandExecution provides a mock function with given fields: id
func (_m *CommandExecutionStarter) StartCommandExecution(id int64) (uuid.UUID, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for StartCommandExecution")
	}

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (uuid.UUID, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int64) uuid.UUID); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCommandExecutionStarter creates a new instance of CommandExecutionStarter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCommandExecutionStarter(t interface {
	mock.TestingT
	Cleanup(func())
}) *CommandExecutionStarter {
	mock := &CommandExecutionStarter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
