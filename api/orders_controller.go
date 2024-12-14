package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorestserviceagain/entity"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type OrdersController struct {
	Repo entity.Repo
}

func (o OrdersController) RegisterRoutes(r chi.Router) {
	r.Post("/", o.CreateOrder)
	r.Get("/", o.ReadAllOrders)
	r.Get("/{id}", o.ReadOrderById)
	r.Put("/{id}", o.UpdateOderById)
	r.Put("/{orderId}/dishes/{dishId} ", o.UpdateDiscountById)
	r.Delete("/{id}", o.DeleteOrderById)
}

func (o OrdersController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order entity.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		fmt.Println(entity.ErrJson)
		SendErr(w, http.StatusUnprocessableEntity, entity.ErrJson.Error())
		return
	}
	err = o.Repo.CreateOrder(&order)
	if err != nil {
		SendErr(w, http.StatusUnprocessableEntity, err.Error())
		fmt.Println("Can not add order")
	} else {
		SendJson(w, http.StatusCreated, order)
		fmt.Println("Added order")
	}
}
func (o OrdersController) ReadAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := o.Repo.GetOrders()
	if err != nil {
		SendErr(w, http.StatusUnprocessableEntity, err.Error())
		fmt.Println("Can not find orders")
	} else {
		SendJson(w, http.StatusOK, orders)
		fmt.Println("Found order")
	}
}
func (o OrdersController) ReadOrderById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	fmt.Printf("id: %+v\n", id)

	order, err := o.Repo.GetOrder(uint(id))
	if err != nil {
		notFoundErr := entity.RecordNotFoundError{}
		if errors.As(err, &notFoundErr) {
			SendErr(w, http.StatusNotFound, err.Error())
			fmt.Println("Can not find order")
		} else {
			SendErr(w, http.StatusInternalServerError, "Unknown error")
			fmt.Println("Inner issue, can not find order")
		}
		return
	}
	SendJson(w, http.StatusOK, order)
	fmt.Println("Found order")
}

func (o OrdersController) UpdateOderById(w http.ResponseWriter, r *http.Request) {
	var order entity.Order
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		SendErr(w, http.StatusUnprocessableEntity, err.Error())
		fmt.Println(entity.ErrJson)
		return
	}
	order.ID = uint(id)
	err = o.Repo.UpdateOrder(&order)
	if err != nil {
		notFoundErr := entity.RecordNotFoundError{}
		if errors.As(err, &notFoundErr) {
			SendErr(w, http.StatusNotFound, err.Error())
			fmt.Println("Can not update the order")
		} else {
			SendErr(w, http.StatusInternalServerError, "Unknown error")
			fmt.Println("Inner error, can not update order")
		}
		return
	}
	SendJson(w, http.StatusNoContent, nil)
	fmt.Println("Order is updated")
}

func (o OrdersController) UpdateDiscountById(w http.ResponseWriter, r *http.Request) {
	var discount entity.DiscountDetail
	orderId, _ := strconv.ParseUint(chi.URLParam(r, "orderId"), 10, 64)
	dishId, _ := strconv.ParseUint(chi.URLParam(r, "dishId"), 10, 64)
	err := json.NewDecoder(r.Body).Decode(&discount)
	if err != nil {
		fmt.Println("Can not convert json to object")
		SendErr(w, http.StatusUnprocessableEntity, err.Error())
	}

	discount.OrderID = uint(orderId)
	discount.DishID = uint(dishId)
	_, err = o.Repo.GetOrder(discount.OrderID)
	if err != nil {
		fmt.Println("Customer does not exsist")
		SendErr(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	fmt.Printf("Found the custome and id is %v\n", discount.OrderID)

	_, err = o.Repo.GetDish(discount.DishID)
	if err != nil {
		fmt.Println("Dish does not exsist")
		SendErr(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	fmt.Printf("Found the dish and id is %v\n", discount.DishID)

	err = o.Repo.UpdateDiscount(&discount)
	if err != nil {
		notFoundErr := entity.RecordNotFoundError{}
		if errors.As(err, &notFoundErr) {
			SendErr(w, http.StatusNotFound, err.Error())
		} else {
			SendErr(w, http.StatusInternalServerError, "Uknown error")
		}
		return
	}
	SendJson(w, http.StatusNoContent, nil)
	fmt.Println("Discount is updated")
}

func (o OrdersController) DeleteOrderById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	fmt.Printf("id: %#v\n", id)
	err := o.Repo.DeleteOrder(uint(id))
	if err != nil {
		notFoundErr := entity.RecordNotFoundError{}
		if errors.As(err, &notFoundErr) {
			SendErr(w, http.StatusNotFound, err.Error())
			fmt.Println("Can not find the id, so can not delete the order")
		} else {
			SendErr(w, http.StatusInternalServerError, err.Error())
			fmt.Println("Inner error, can not delete oreder")
		}
		return
	}
	SendJson(w, http.StatusNoContent, nil)
	fmt.Println("Deleted order")
}
