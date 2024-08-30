// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	entity "kururen/entity"

	mock "github.com/stretchr/testify/mock"
)

// Reader is an autogenerated mock type for the Reader type
type Reader struct {
	mock.Mock
}

// GetCarByID provides a mock function with given fields: _a0
func (_m *Reader) GetCarByID(_a0 uint) (*entity.Car, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetCarByID")
	}

	var r0 *entity.Car
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) (*entity.Car, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(uint) *entity.Car); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Car)
		}
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCars provides a mock function with given fields:
func (_m *Reader) GetCars() ([]entity.Car, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetCars")
	}

	var r0 []entity.Car
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]entity.Car, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []entity.Car); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Car)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewReader creates a new instance of Reader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *Reader {
	mock := &Reader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
