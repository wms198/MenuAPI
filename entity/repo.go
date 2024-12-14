package entity

import (
	"errors"
	"fmt"
)

type RecordNotFoundError struct {
	ID    string
	Kind  string
	Inner error
}

func (e RecordNotFoundError) Error() string {
	return fmt.Sprintf("%v with id %v not found", e.Kind, e.ID)
}

func (e RecordNotFoundError) Unwarp() error {
	return e.Inner
}

func (e RecordNotFoundError) Is(err error) bool {
	return errors.Is(err, ErrRecordNotFound)
}

func WrapRecordNotFoundError(kind string, id any, inner error) error {
	return RecordNotFoundError{
		ID:    fmt.Sprint(id),
		Kind:  kind,
		Inner: inner,
	}
}

var ErrDBNotConnected = errors.New("could not connect to DB")
var ErrInvalidData = errors.New("unsupported data")
var ErrRecordNotFound = errors.New("record not found")
var ErrJson = errors.New("can not convert object to JSON")
var ErrEntityNotFound = errors.New("entity not found")

type Repo interface {
	OrdersRepo
	DiscountDetailsRepo
	DishRepo
}

type OrdersRepo interface {
	CreateOrder(order *Order) error
	GetOrders() ([]Order, error)
	GetOrder(id uint) (Order, error)
	UpdateOrder(order *Order) error
	UpdateDiscount(discount *DiscountDetail) error
	DeleteOrder(id uint) error
}
type DiscountDetailsRepo interface {
	CreateDiscount(discount *DiscountDetail) error
	GetPriceAfterDiscount(orderId uint, dishId uint) (DiscountDetail, error)
}
type DishRepo interface {
	CreateDish(dish *Dish) error
	GetDishes() ([]Dish, error)
	GetDish(id uint) (Dish, error)
	UpdateDish(dish *Dish) error
	DeleteDish(id uint) error
}
