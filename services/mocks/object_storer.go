// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ObjectStorer is an autogenerated mock type for the ObjectStorer type
type ObjectStorer struct {
	mock.Mock
}

type ObjectStorer_Expecter struct {
	mock *mock.Mock
}

func (_m *ObjectStorer) EXPECT() *ObjectStorer_Expecter {
	return &ObjectStorer_Expecter{mock: &_m.Mock}
}

// Upload provides a mock function with given fields: ctx, src, bucket, remoteObjectName, srcContentType
func (_m *ObjectStorer) Upload(ctx context.Context, src []byte, bucket string, remoteObjectName string, srcContentType string) error {
	ret := _m.Called(ctx, src, bucket, remoteObjectName, srcContentType)

	if len(ret) == 0 {
		panic("no return value specified for Upload")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte, string, string, string) error); ok {
		r0 = rf(ctx, src, bucket, remoteObjectName, srcContentType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ObjectStorer_Upload_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upload'
type ObjectStorer_Upload_Call struct {
	*mock.Call
}

// Upload is a helper method to define mock.On call
//   - ctx context.Context
//   - src []byte
//   - bucket string
//   - remoteObjectName string
//   - srcContentType string
func (_e *ObjectStorer_Expecter) Upload(ctx interface{}, src interface{}, bucket interface{}, remoteObjectName interface{}, srcContentType interface{}) *ObjectStorer_Upload_Call {
	return &ObjectStorer_Upload_Call{Call: _e.mock.On("Upload", ctx, src, bucket, remoteObjectName, srcContentType)}
}

func (_c *ObjectStorer_Upload_Call) Run(run func(ctx context.Context, src []byte, bucket string, remoteObjectName string, srcContentType string)) *ObjectStorer_Upload_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]byte), args[2].(string), args[3].(string), args[4].(string))
	})
	return _c
}

func (_c *ObjectStorer_Upload_Call) Return(_a0 error) *ObjectStorer_Upload_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ObjectStorer_Upload_Call) RunAndReturn(run func(context.Context, []byte, string, string, string) error) *ObjectStorer_Upload_Call {
	_c.Call.Return(run)
	return _c
}

// NewObjectStorer creates a new instance of ObjectStorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewObjectStorer(t interface {
	mock.TestingT
	Cleanup(func())
}) *ObjectStorer {
	mock := &ObjectStorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}