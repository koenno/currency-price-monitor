// Code generated by mockery v2.32.0. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"

	request "github.com/koenno/currency-price-monitor/request"
)

// Requester is an autogenerated mock type for the Requester type
type Requester[T interface{}] struct {
	mock.Mock
}

type Requester_Expecter[T interface{}] struct {
	mock *mock.Mock
}

func (_m *Requester[T]) EXPECT() *Requester_Expecter[T] {
	return &Requester_Expecter[T]{mock: &_m.Mock}
}

// Process provides a mock function with given fields: req
func (_m *Requester[T]) Process(req *http.Request) (request.Descriptor[T], error) {
	ret := _m.Called(req)

	var r0 request.Descriptor[T]
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request) (request.Descriptor[T], error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*http.Request) request.Descriptor[T]); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(request.Descriptor[T])
	}

	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Requester_Process_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Process'
type Requester_Process_Call[T interface{}] struct {
	*mock.Call
}

// Process is a helper method to define mock.On call
//   - req *http.Request
func (_e *Requester_Expecter[T]) Process(req interface{}) *Requester_Process_Call[T] {
	return &Requester_Process_Call[T]{Call: _e.mock.On("Process", req)}
}

func (_c *Requester_Process_Call[T]) Run(run func(req *http.Request)) *Requester_Process_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*http.Request))
	})
	return _c
}

func (_c *Requester_Process_Call[T]) Return(_a0 request.Descriptor[T], _a1 error) *Requester_Process_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Requester_Process_Call[T]) RunAndReturn(run func(*http.Request) (request.Descriptor[T], error)) *Requester_Process_Call[T] {
	_c.Call.Return(run)
	return _c
}

// NewRequester creates a new instance of Requester. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRequester[T interface{}](t interface {
	mock.TestingT
	Cleanup(func())
}) *Requester[T] {
	mock := &Requester[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}