package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorestserviceagain/entity"
	"net/http"
)

func SendJson(w http.ResponseWriter, status int, body any) {
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	if body != nil {
		err := json.NewEncoder(w).Encode(body)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SendErr(w http.ResponseWriter, status int, err string) {
	SendJson(w, status, Err{Error: err})
}

func comparePrice(DiscountDetail entity.DiscountDetail, dish entity.Dish) error {
	var msg = "Added discount is too high"
	discount := DiscountDetail.Discount
	price := dish.Price
	priceAfterDiscount := price * ((100 - discount) / 100)

	if priceAfterDiscount < price*0.8 {
		fmt.Println(msg)
		return errors.New(msg)
	}
	return nil
}
