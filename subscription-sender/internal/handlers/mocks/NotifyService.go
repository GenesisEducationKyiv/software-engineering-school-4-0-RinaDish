// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	services "github.com/RinaDish/subscription-sender/internal/services"
)

// NotifyService is an autogenerated mock type for the NotifyService type
type NotifyService struct {
	mock.Mock
}

// NotifySubscribers provides a mock function with given fields: ctx, notification
func (_m *NotifyService) NotifySubscribers(ctx context.Context, notification services.Notification) {
	_m.Called(ctx, notification)
}

// NewNotifyService creates a new instance of NotifyService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNotifyService(t interface {
	mock.TestingT
	Cleanup(func())
}) *NotifyService {
	mock := &NotifyService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}