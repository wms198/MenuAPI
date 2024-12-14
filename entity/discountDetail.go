package entity

type DiscountDetail struct {
	OrderID  uint
	Order    Order
	DishID   uint
	Dish     Dish
	Discount float32
}
