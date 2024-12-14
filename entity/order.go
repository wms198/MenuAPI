package entity

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	TableNumber    int
	FinalPrice     float32
	DiscountDetail []DiscountDetail
}
