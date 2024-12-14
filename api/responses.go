package api

type Err struct {
	Error string
}

type priceAfterDiscount struct {
	OrderID       uint
	DischID       uint
	OriginalPrice float32
	DiscountPrice float32
}
