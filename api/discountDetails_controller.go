package api

import (
	"encoding/json"
	"fmt"
	"gorestserviceagain/entity"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DiscountDetailController struct {
	Repo entity.Repo
}

func (d DiscountDetailController) RegisterRoutes(r chi.Router) {
	r.Post("/discountPrice", d.CreateDiscount)
	r.Get("/{orderId}/dishes/{dishId}", d.GetPriceAfterDiscount)
}

func (d DiscountDetailController) CreateDiscount(w http.ResponseWriter, r *http.Request) {
	var price entity.DiscountDetail
	var order entity.Order
	var dish entity.Dish
	err := json.NewDecoder(r.Body).Decode(&price)
	if err != nil {
		fmt.Println(entity.ErrJson)
		SendErr(w, http.StatusUnprocessableEntity, entity.ErrJson.Error())
		return
	}

	order, err = d.Repo.GetOrder(price.OrderID)
	if err != nil {
		fmt.Println("Order does not exsits")
		SendErr(w, http.StatusUnprocessableEntity, err.Error())
		return
	} else {
		fmt.Printf("Found the order and id is %v\n", price.OrderID)
	}

	dish, err = d.Repo.GetDish(price.DishID)
	if err != nil {
		fmt.Println("Dish does not exsits")
		SendErr(w, http.StatusUnprocessableEntity, err.Error())
		return
	} else {
		fmt.Printf("Founf the dish and id is %v\n", price.DishID)
	}

	err = comparePrice(price, dish)
	if err != nil {
		SendErr(w, http.StatusUnprocessableEntity, err.Error())
	}
	price.Order = order
	price.Dish = dish

	err = d.Repo.CreateDiscount(&price)
	if err != nil {
		fmt.Println("Can not add discount price, discount with same orderId anf dishId has exist")
		SendErr(w, http.StatusUnprocessableEntity, err.Error())
		return
	} else {
		SendJson(w, http.StatusCreated, price)
		fmt.Println("Added discount price")
	}

}

func (d DiscountDetailController) GetPriceAfterDiscount(w http.ResponseWriter, r *http.Request) {
	var newPrice priceAfterDiscount
	orderId, _ := strconv.ParseUint(chi.URLParam(r, "orderId"), 10, 64)
	dishId, _ := strconv.ParseUint(chi.URLParam(r, "dishId"), 10, 64)

	price, err := d.Repo.GetPriceAfterDiscount(uint(orderId), uint(dishId))

	if err != nil {
		SendErr(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	originalPrice := price.Dish.Price
	newPrice.OrderID = price.OrderID
	newPrice.DischID = price.DishID
	newPrice.OriginalPrice = originalPrice
	newPrice.DiscountPrice = originalPrice * ((100 - price.Discount) / 100)

	SendJson(w, http.StatusOK, newPrice)
	fmt.Printf("OrderId: %d\n", newPrice.OrderID)
	fmt.Printf("DishId: %d\n", newPrice.DischID)
	fmt.Printf("Original price: %v\n", newPrice.OriginalPrice)
	fmt.Printf("Discount price: %v\n", newPrice.DiscountPrice)

}
