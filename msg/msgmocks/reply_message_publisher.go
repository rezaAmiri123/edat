// Code generated by mockery v2.33.0. DO NOT EDIT.

package msgmocks

import (
	context "context"

	core "github.com/rezaAmiri123/edat/core"
	mock "github.com/stretchr/testify/mock"

	msg "github.com/rezaAmiri123/edat/msg"
)

// ReplyMessagePublisher is an autogenerated mock type for the ReplyMessagePublisher type
type ReplyMessagePublisher struct {
	mock.Mock
}

// PublishReply provides a mock function with given fields: ctx, reply, options
func (_m *ReplyMessagePublisher) PublishReply(ctx context.Context, reply core.Reply, options ...msg.MessageOption) error {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, reply)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, core.Reply, ...msg.MessageOption) error); ok {
		r0 = rf(ctx, reply, options...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewReplyMessagePublisher creates a new instance of ReplyMessagePublisher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReplyMessagePublisher(t interface {
	mock.TestingT
	Cleanup(func())
}) *ReplyMessagePublisher {
	mock := &ReplyMessagePublisher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
