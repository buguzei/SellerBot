// Code generated by MockGen. DO NOT EDIT.
// Source: ./repo.go

// Package mock_repo is a generated GoMock package.
package mock_repo

import (
	entities "bot/internal/entities"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserRepo is a mock of UserRepo interface.
type MockUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoMockRecorder
}

// MockUserRepoMockRecorder is the mock recorder for MockUserRepo.
type MockUserRepoMockRecorder struct {
	mock *MockUserRepo
}

// NewMockUserRepo creates a new mock instance.
func NewMockUserRepo(ctrl *gomock.Controller) *MockUserRepo {
	mock := &MockUserRepo{ctrl: ctrl}
	mock.recorder = &MockUserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepo) EXPECT() *MockUserRepoMockRecorder {
	return m.recorder
}

// GetUser mocks base method.
func (m *MockUserRepo) GetUser(arg0 int64) (*entities.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0)
	ret0, _ := ret[0].(*entities.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockUserRepoMockRecorder) GetUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUserRepo)(nil).GetUser), arg0)
}

// InsertUser mocks base method.
func (m *MockUserRepo) InsertUser(arg0 entities.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertUser", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertUser indicates an expected call of InsertUser.
func (mr *MockUserRepoMockRecorder) InsertUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertUser", reflect.TypeOf((*MockUserRepo)(nil).InsertUser), arg0)
}

// UpdateUser mocks base method.
func (m *MockUserRepo) UpdateUser(arg0 entities.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserRepoMockRecorder) UpdateUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserRepo)(nil).UpdateUser), arg0)
}

// MockOrderRepo is a mock of OrderRepo interface.
type MockOrderRepo struct {
	ctrl     *gomock.Controller
	recorder *MockOrderRepoMockRecorder
}

// MockOrderRepoMockRecorder is the mock recorder for MockOrderRepo.
type MockOrderRepoMockRecorder struct {
	mock *MockOrderRepo
}

// NewMockOrderRepo creates a new mock instance.
func NewMockOrderRepo(ctrl *gomock.Controller) *MockOrderRepo {
	mock := &MockOrderRepo{ctrl: ctrl}
	mock.recorder = &MockOrderRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderRepo) EXPECT() *MockOrderRepoMockRecorder {
	return m.recorder
}

// GetAllCurrentOrders mocks base method.
func (m *MockOrderRepo) GetAllCurrentOrders() ([]entities.CurrentOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCurrentOrders")
	ret0, _ := ret[0].([]entities.CurrentOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllCurrentOrders indicates an expected call of GetAllCurrentOrders.
func (mr *MockOrderRepoMockRecorder) GetAllCurrentOrders() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCurrentOrders", reflect.TypeOf((*MockOrderRepo)(nil).GetAllCurrentOrders))
}

// GetAllDoneOrders mocks base method.
func (m *MockOrderRepo) GetAllDoneOrders() ([]entities.DoneOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllDoneOrders")
	ret0, _ := ret[0].([]entities.DoneOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllDoneOrders indicates an expected call of GetAllDoneOrders.
func (mr *MockOrderRepoMockRecorder) GetAllDoneOrders() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDoneOrders", reflect.TypeOf((*MockOrderRepo)(nil).GetAllDoneOrders))
}

// NewCurrentOrder mocks base method.
func (m *MockOrderRepo) NewCurrentOrder(arg0 entities.CurrentOrder) (*int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewCurrentOrder", arg0)
	ret0, _ := ret[0].(*int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewCurrentOrder indicates an expected call of NewCurrentOrder.
func (mr *MockOrderRepoMockRecorder) NewCurrentOrder(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewCurrentOrder", reflect.TypeOf((*MockOrderRepo)(nil).NewCurrentOrder), arg0)
}

// NewCurrentProducts mocks base method.
func (m *MockOrderRepo) NewCurrentProducts(arg0 entities.CurrentOrder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewCurrentProducts", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewCurrentProducts indicates an expected call of NewCurrentProducts.
func (mr *MockOrderRepoMockRecorder) NewCurrentProducts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewCurrentProducts", reflect.TypeOf((*MockOrderRepo)(nil).NewCurrentProducts), arg0)
}

// NewDoneOrder mocks base method.
func (m *MockOrderRepo) NewDoneOrder(arg0 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewDoneOrder", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewDoneOrder indicates an expected call of NewDoneOrder.
func (mr *MockOrderRepoMockRecorder) NewDoneOrder(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewDoneOrder", reflect.TypeOf((*MockOrderRepo)(nil).NewDoneOrder), arg0)
}

// MockCartRepo is a mock of CartRepo interface.
type MockCartRepo struct {
	ctrl     *gomock.Controller
	recorder *MockCartRepoMockRecorder
}

// MockCartRepoMockRecorder is the mock recorder for MockCartRepo.
type MockCartRepoMockRecorder struct {
	mock *MockCartRepo
}

// NewMockCartRepo creates a new mock instance.
func NewMockCartRepo(ctrl *gomock.Controller) *MockCartRepo {
	mock := &MockCartRepo{ctrl: ctrl}
	mock.recorder = &MockCartRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCartRepo) EXPECT() *MockCartRepoMockRecorder {
	return m.recorder
}

// CartLen mocks base method.
func (m *MockCartRepo) CartLen(arg0 int64) (*int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CartLen", arg0)
	ret0, _ := ret[0].(*int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CartLen indicates an expected call of CartLen.
func (mr *MockCartRepoMockRecorder) CartLen(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CartLen", reflect.TypeOf((*MockCartRepo)(nil).CartLen), arg0)
}

// ClearCart mocks base method.
func (m *MockCartRepo) ClearCart(arg0 int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ClearCart", arg0)
}

// ClearCart indicates an expected call of ClearCart.
func (mr *MockCartRepoMockRecorder) ClearCart(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearCart", reflect.TypeOf((*MockCartRepo)(nil).ClearCart), arg0)
}

// DeleteProductFromCart mocks base method.
func (m *MockCartRepo) DeleteProductFromCart(arg0 int64, arg1 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteProductFromCart", arg0, arg1)
}

// DeleteProductFromCart indicates an expected call of DeleteProductFromCart.
func (mr *MockCartRepoMockRecorder) DeleteProductFromCart(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProductFromCart", reflect.TypeOf((*MockCartRepo)(nil).DeleteProductFromCart), arg0, arg1)
}

// GetCart mocks base method.
func (m *MockCartRepo) GetCart(arg0 int64) (map[int]entities.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCart", arg0)
	ret0, _ := ret[0].(map[int]entities.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCart indicates an expected call of GetCart.
func (mr *MockCartRepoMockRecorder) GetCart(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCart", reflect.TypeOf((*MockCartRepo)(nil).GetCart), arg0)
}

// GetCartProduct mocks base method.
func (m *MockCartRepo) GetCartProduct(arg0 int64, arg1 int) (*entities.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCartProduct", arg0, arg1)
	ret0, _ := ret[0].(*entities.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCartProduct indicates an expected call of GetCartProduct.
func (mr *MockCartRepoMockRecorder) GetCartProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCartProduct", reflect.TypeOf((*MockCartRepo)(nil).GetCartProduct), arg0, arg1)
}

// NewCartProduct mocks base method.
func (m *MockCartRepo) NewCartProduct(arg0 int64, arg1 int, arg2 entities.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewCartProduct", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewCartProduct indicates an expected call of NewCartProduct.
func (mr *MockCartRepoMockRecorder) NewCartProduct(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewCartProduct", reflect.TypeOf((*MockCartRepo)(nil).NewCartProduct), arg0, arg1, arg2)
}