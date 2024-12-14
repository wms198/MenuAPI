package entity

import "github.com/stretchr/testify/mock"

type MockRepo struct {
	mock.Mock
	Orders              []Order
	Dishes              []Dish
	DiscountDetails     []DiscountDetail
	OrderError          error
	DishError           error
	DiscountDetailError error
}

func (m *MockRepo) CreateOrder(order *Order) error {
	args := m.Called(*order)
	return args.Error(0)
}

func (m *MockRepo) GetOrders() (c []Order, err error) {
	args := m.Called()
	if result := args.Get(0); result != nil {
		return result.([]Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepo) GetOrder(id uint) (Order, error) {
	args := m.Called(id)
	if result := args.Get(0); result != nil {
		return result.(Order), args.Error(1)
	}
	return Order{}, args.Error(1)
}

func (m *MockRepo) UpdateOrder(order *Order) error {
	args := m.Called(*order)
	return args.Error(0)
}

func (m *MockRepo) UpdateDiscount(discount *DiscountDetail) error {
	args := m.Called(*discount)
	return args.Error(0)
}

func (m *MockRepo) DeleteOrder(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepo) CreateDish(dish *Dish) error {
	args := m.Called(*dish)
	return args.Error(0)
}

func (m *MockRepo) GetDishes() ([]Dish, error) {
	args := m.Called()
	if result := args.Get(0); result != nil {
		return result.([]Dish), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepo) GetDish(id uint) (Dish, error) {
	args := m.Called(id)
	if result := args.Get(0); result != nil {
		return result.(Dish), args.Error(1)
	}
	return Dish{}, args.Error(1)
}

func (m *MockRepo) UpdateDish(dish *Dish) error {
	args := m.Called(*dish)
	return args.Error(0)
}

func (m *MockRepo) DeleteDish(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepo) CreateDiscount(discount *DiscountDetail) error {
	args := m.Called(*discount)
	return args.Error(0)
}

func (m *MockRepo) GetPriceAfterDiscount(orderId uint, dishId uint) (DiscountDetail, error) {
	args := m.Called(orderId, dishId)
	if result := args.Get(0); result != nil {
		return result.(DiscountDetail), args.Error(1)
	}
	return DiscountDetail{}, args.Error(1)
}
