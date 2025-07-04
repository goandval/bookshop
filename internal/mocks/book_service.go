// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	domain "github.com/yourorg/bookshop/internal/domain"
)

// BookService is an autogenerated mock type for the BookService type
type BookService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, book
func (_m *BookService) Create(ctx context.Context, book *domain.Book) error {
	ret := _m.Called(ctx, book)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Book) error); ok {
		r0 = rf(ctx, book)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id
func (_m *BookService) Delete(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *BookService) GetByID(ctx context.Context, id int) (*domain.Book, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *domain.Book
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (*domain.Book, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) *domain.Book); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Book)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, categoryIDs, limit, offset
func (_m *BookService) List(ctx context.Context, categoryIDs []int, limit int, offset int) ([]*domain.Book, error) {
	ret := _m.Called(ctx, categoryIDs, limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []*domain.Book
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []int, int, int) ([]*domain.Book, error)); ok {
		return rf(ctx, categoryIDs, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []int, int, int) []*domain.Book); ok {
		r0 = rf(ctx, categoryIDs, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Book)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []int, int, int) error); ok {
		r1 = rf(ctx, categoryIDs, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, book
func (_m *BookService) Update(ctx context.Context, book *domain.Book) error {
	ret := _m.Called(ctx, book)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Book) error); ok {
		r0 = rf(ctx, book)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewBookService creates a new instance of BookService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBookService(t interface {
	mock.TestingT
	Cleanup(func())
}) *BookService {
	mock := &BookService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
