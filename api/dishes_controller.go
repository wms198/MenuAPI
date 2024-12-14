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

type DishesController struct {
	Repo entity.Repo
}

func (d DishesController) RegisterRoutes(r chi.Router) {
	r.Post("/", d.CreateDish)
	r.Get("/", d.ReadAllDishes)
	r.Get("/{id}", d.ReadDishById)
	r.Put("/{id}", d.UpdateDishById)
	r.Delete("/{id}", d.DeleteDishById)
}

func (d DishesController) CreateDish(w http.ResponseWriter, r *http.Request) {
	var dish entity.Dish
	err := json.NewDecoder(r.Body).Decode(&dish)
	if err != nil {
		fmt.Println(entity.ErrJson)
		SendErr(w, http.StatusUnprocessableEntity, entity.ErrJson.Error())
		return
	}
	err = d.Repo.CreateDish(&dish)
	if err != nil {
		SendErr(w, http.StatusUnprocessableEntity, err.Error())
		fmt.Println("Can not create dish")
	} else {
		SendJson(w, http.StatusCreated, dish)
		fmt.Println("Added dish")
	}
}

func (d DishesController) ReadAllDishes(w http.ResponseWriter, r *http.Request) {
	dishes, err := d.Repo.GetDishes()
	if err != nil {
		SendErr(w, http.StatusUnprocessableEntity, err.Error())
		fmt.Println("Can not find dishes")
	} else {
		SendJson(w, http.StatusOK, dishes)
		fmt.Println("Found dishes")
	}
}

func (d DishesController) ReadDishById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	fmt.Printf("id: %+v\n", id)

	dish, err := d.Repo.GetDish(uint(id))
	if err != nil {
		notFoundErr := entity.RecordNotFoundError{}
		if errors.As(err, &notFoundErr) {
			SendErr(w, http.StatusNotFound, err.Error())
			fmt.Println("Can not find dish")
		} else {
			SendErr(w, http.StatusInternalServerError, "Unkown err")
			fmt.Printf("Inner error, can not find dish")
		}
		return
	}
	SendJson(w, http.StatusOK, dish)
	fmt.Println("Found dish")
}

func (d DishesController) UpdateDishById(w http.ResponseWriter, r *http.Request) {
	var dish entity.Dish
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	err := json.NewDecoder(r.Body).Decode(&dish)
	if err != nil {
		fmt.Println(entity.ErrJson)
		SendErr(w, http.StatusUnprocessableEntity, entity.ErrInvalidData.Error())
		return
	}
	dish.ID = uint(id)
	err = d.Repo.UpdateDish(&dish)
	if err != nil {
		msg := err.Error()
		notFoundErr := entity.RecordNotFoundError{}
		if errors.As(err, &notFoundErr) {
			SendErr(w, http.StatusNotFound, msg)
			fmt.Println("Can not update dish")
		} else {
			SendErr(w, http.StatusInternalServerError, "Unkown error")
			fmt.Println("Inner error", msg)
		}
		return
	}
	SendJson(w, http.StatusNoContent, nil)
	fmt.Println("Updated dish")
}

func (d DishesController) DeleteDishById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	fmt.Printf("id: %+v\n", id)
	err := d.Repo.DeleteDish(uint(id))
	if err != nil {
		msg := err.Error()
		notFoundErr := entity.RecordNotFoundError{}
		if errors.As(err, &notFoundErr) {
			SendErr(w, http.StatusNotFound, msg)
			fmt.Println("Can not found id, can not delete dish")
		} else {
			SendErr(w, http.StatusInternalServerError, "Unkown error")
			fmt.Println("Internal error", msg)
		}
		return
	}
	SendJson(w, http.StatusNoContent, nil)
	fmt.Println("Delete dish")
}
