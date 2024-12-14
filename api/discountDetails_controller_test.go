package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"gorestserviceagain/entity"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestDiscountCreate(t *testing.T) {
	type expectations struct {
		statusCode  int
		respPayload any
	}
	order := entity.Order{Model: gorm.Model{ID: 1}, TableNumber: 2, FinalPrice: 14.}
	dish := entity.Dish{Model: gorm.Model{ID: 2}, Name: "Fish filet", Price: 10.}
	discountDetail := entity.DiscountDetail{OrderID: 1, DishID: 2, Discount: 2., Order: order, Dish: dish}
	errOrderNotFound := errors.New("order not found")
	errDishNotFound := errors.New("dish not found")

	tests := []struct {
		name             string
		payload          entity.DiscountDetail
		expected         expectations
		respPayload      any
		existingOrder    entity.Order
		existingDish     entity.Dish
		existingDiscount entity.DiscountDetail
		orderErr         error
		dishErr          error
		discountErr      error
	}{
		{
			name:          "successful creatation",
			payload:       discountDetail,
			existingOrder: order,
			existingDish:  dish,
			expected: expectations{
				statusCode: http.StatusCreated,
				respPayload: map[string]interface{}{
					"OrderID":  float64(discountDetail.OrderID),
					"DishID":   float64(discountDetail.DishID),
					"Discount": float64(discountDetail.Discount),
					"Dish": map[string]interface{}{
						"CreatedAt": "0001-01-01T00:00:00Z",
						"DeletedAt": interface{}(nil),
						"ID":        float64(dish.ID),
						"Name":      dish.Name,
						"Price":     float64(dish.Price),
						"UpdatedAt": "0001-01-01T00:00:00Z",
					},
					"Order": map[string]interface{}{
						"CreatedAt":      "0001-01-01T00:00:00Z",
						"DeletedAt":      interface{}(nil),
						"DiscountDetail": interface{}(nil),
						"FinalPrice":     float64(order.FinalPrice),
						"ID":             float64(order.ID),
						"TableNumber":    float64(order.TableNumber),
						"UpdatedAt":      "0001-01-01T00:00:00Z",
					},
				},
			},
		},
		{
			name:          "order doesn't exist",
			orderErr:      errOrderNotFound,
			payload:       discountDetail,
			existingOrder: order,
			expected: expectations{
				statusCode:  http.StatusUnprocessableEntity,
				respPayload: map[string]interface{}{"Error": errOrderNotFound.Error()},
			},
		},
		{
			name:          "dish doesn't exist",
			dishErr:       errDishNotFound,
			payload:       discountDetail,
			existingOrder: order,
			existingDish:  dish,
			expected: expectations{
				statusCode:  http.StatusUnprocessableEntity,
				respPayload: map[string]interface{}{"Error": errDishNotFound.Error()},
			},
		},
		{
			name:             "discount is too high",
			existingOrder:    order,
			existingDish:     dish,
			existingDiscount: discountDetail,
			payload: entity.DiscountDetail{
				OrderID:  order.ID,
				DishID:   dish.ID,
				Discount: 50,
				Order:    order,
				Dish:     dish,
			},
			expected: expectations{
				statusCode:  http.StatusUnprocessableEntity,
				respPayload: map[string]interface{}{"Error": "Added discount is too high"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			b := bytes.NewBuffer(nil)
			require.NoError(t, json.NewEncoder(b).Encode(tt.payload))
			r := httptest.NewRequest(http.MethodPost, "/discountPrice/", b)

			repo := new(entity.MockRepo)
			repo.On("GetOrder", tt.existingOrder.ID).Return(tt.existingOrder, tt.orderErr)
			repo.On("GetDish", tt.existingDish.ID).Return(tt.existingDish, tt.dishErr)
			repo.On("CreateDiscount", tt.payload).Return(tt.discountErr)

			DiscountDetailController{Repo: repo}.CreateDiscount(w, r)

			res := w.Result()
			assert.Equal(t, tt.expected.statusCode, res.StatusCode)
			if tt.expected.respPayload != nil {
				//parsing Body as json, puting into &tt.respPayload 从前往后
				require.NoError(t, json.NewDecoder(res.Body).Decode(&tt.respPayload))

				assert.EqualValues(t, tt.expected.respPayload, tt.respPayload)
			}
		})
	}
}

func TestPriceAfterDiscountGet(t *testing.T) {
	type expectations struct {
		statusCode  int
		respPayload any
	}
	dish := entity.Dish{Model: gorm.Model{ID: 2}, Name: "Fish filet", Price: 10.}
	discountDetail := entity.DiscountDetail{OrderID: 1, DishID: 2, Discount: 2., Dish: dish}

	tests := []struct {
		name        string
		expected    expectations
		respPayload any
		existing    entity.DiscountDetail
		err         error
	}{
		{
			name:     "successful get price after discount",
			existing: discountDetail,
			expected: expectations{
				statusCode: http.StatusOK,
				respPayload: map[string]interface{}{
					"OrderID":       float64(discountDetail.OrderID),
					"DischID":       float64(discountDetail.DishID),
					"DiscountPrice": 9.8,
					"OriginalPrice": float64(discountDetail.Dish.Price),
				},
			},
		},
		{
			name:     "failed to find the discount",
			existing: discountDetail,
			err:      entity.ErrEntityNotFound,
			expected: expectations{
				statusCode:  http.StatusUnprocessableEntity,
				respPayload: map[string]interface{}{"Error": entity.ErrEntityNotFound.Error()},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			b := bytes.NewBuffer(nil)
			r := httptest.NewRequest(http.MethodGet, "/{orderId}/dishes/{dishId}", b)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("orderId", strconv.FormatUint(uint64(tt.existing.OrderID), 10))
			rctx.URLParams.Add("dishId", strconv.FormatUint(uint64(tt.existing.DishID), 10))
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			repo := new(entity.MockRepo)
			repo.On("GetPriceAfterDiscount", tt.existing.OrderID, tt.existing.DishID).Return(tt.existing, tt.err)
			DiscountDetailController{Repo: repo}.GetPriceAfterDiscount(w, r)

			res := w.Result()
			assert.Equal(t, tt.expected.statusCode, res.StatusCode)
			if tt.expected.respPayload != nil {
				//parsing Body as json, puting into &tt.respPayload 从前往后
				require.NoError(t, json.NewDecoder(res.Body).Decode(&tt.respPayload))
				assert.EqualValues(t, tt.expected.respPayload, tt.respPayload)
			}
		})
	}
}
