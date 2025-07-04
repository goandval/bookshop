// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	domain "github.com/yourorg/bookshop/internal/domain"
)

// CartRepository is an autogenerated mock type for the CartRepository type
type CartRepository struct {
	mock.Mock
}

// AddItem provides a mock function with given fields: ctx, userID, bookID
func (_m *CartRepository) AddItem(ctx context.Context, userID string, bookID int) error {
	ret := _m.Called(ctx, userID, bookID)

	if len(ret) == 0 {
		panic("no return value specified for AddItem")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int) error); ok {
		r0 = rf(ctx, userID, bookID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Clear provides a mock function with given fields: ctx, userID
func (_m *CartRepository) Clear(ctx context.Context, userID string) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for Clear")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByUserID provides a mock function with given fields: ctx, userID
func (_m *CartRepository) GetByUserID(ctx context.Context, userID string) (*domain.Cart, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetByUserID")
	}

	var r0 *domain.Cart
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.Cart, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.Cart); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Cart)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListItems provides a mock function with given fields: ctx, userID
func (_m *CartRepository) ListItems(ctx context.Context, userID string) ([]*domain.CartItem, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for ListItems")
	}

	var r0 []*domain.CartItem
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*domain.CartItem, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*domain.CartItem); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.CartItem)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveItem provides a mock function with given fields: ctx, userID, bookID
func (_m *CartRepository) RemoveItem(ctx context.Context, userID string, bookID int) error {
	ret := _m.Called(ctx, userID, bookID)

	if len(ret) == 0 {
		panic("no return value specified for RemoveItem")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int) error); ok {
		r0 = rf(ctx, userID, bookID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewCartRepository creates a new instance of CartRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCartRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *CartRepository {
	mock := &CartRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
