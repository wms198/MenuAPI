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

func TestOrderCreate(t *testing.T) {
	type expectations struct {
		statusCode  int
		respPayload any
	}
	order := entity.Order{Model: gorm.Model{ID: 1}, TableNumber: 2, FinalPrice: 14.}

	tests := []struct {
		name        string
		payload     entity.Order
		expected    expectations
		respPayload any
		err         error
	}{
		{
			name:    "successful creatation",
			payload: order,
			expected: expectations{
				statusCode: http.StatusCreated,
				respPayload: map[string]interface{}{
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
		{
			name: "failed creatation",
			err:  entity.ErrInvalidData,
			expected: expectations{
				statusCode:  http.StatusUnprocessableEntity,
				respPayload: map[string]interface{}{"Error": entity.ErrInvalidData.Error()},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			b := bytes.NewBuffer(nil)
			require.NoError(t, json.NewEncoder(b).Encode(tt.payload))
			r := httptest.NewRequest(http.MethodPost, "/orders/", b)

			repo := new(entity.MockRepo)
			// function in Mock_Repo file
			repo.On("CreateOrder", tt.payload).Return(tt.err)

			OrdersController{Repo: repo}.CreateOrder(w, r)

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

func TestOrdersReadAll(t *testing.T) {
	type expectations struct {
		statusCode  int
		respPayload any
	}
	order := entity.Order{Model: gorm.Model{ID: 1}, TableNumber: 2, FinalPrice: 14.}
	tests := []struct {
		name        string
		payload     entity.Order
		existing    []entity.Order
		expected    expectations
		respPayload any
		err         error
	}{
		{
			name: "successful get orders",
			existing: []entity.Order{
				order,
			},
			expected: expectations{
				statusCode: http.StatusOK,
				respPayload: []interface{}{map[string]interface{}{
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
			name: "can not find orders",
			err:  entity.ErrEntityNotFound,
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
			require.NoError(t, json.NewEncoder(b).Encode(tt.payload))
			r := httptest.NewRequest(http.MethodGet, "/orders/", b)

			repo := new(entity.MockRepo)
			repo.On("GetOrders").Return(tt.existing, tt.err)
			OrdersController{Repo: repo}.ReadAllOrders(w, r)

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

func TestOrderReadById(t *testing.T) {
	type expectations struct {
		statusCode  int
		respPayload any
	}
	order := entity.Order{Model: gorm.Model{ID: 1}, TableNumber: 2, FinalPrice: 14.}
	notFoundErr := entity.RecordNotFoundError{
		Kind:  "Order",
		ID:    strconv.FormatInt(int64(order.ID), 10),
		Inner: errors.New("mock repo says no"),
	}

	tests := []struct {
		name        string
		payload     entity.Order
		respPayload any
		existing    entity.Order
		expected    expectations
		err         error
	}{
		{
			name:     "successful get order",
			existing: order,
			expected: expectations{
				statusCode: http.StatusOK,
				respPayload: map[string]interface{}{
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
		{
			name: "can not find the order",
			err:  notFoundErr,
			expected: expectations{
				statusCode:  http.StatusNotFound,
				respPayload: map[string]interface{}{"Error": notFoundErr.Error()},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			b := bytes.NewBuffer(nil)
			r := httptest.NewRequest(http.MethodGet, "/orders/{id}", b)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strconv.FormatUint(uint64(tt.existing.ID), 10))
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			repo := new(entity.MockRepo)
			repo.On("GetOrder", tt.existing.ID).Return(tt.existing, tt.err)

			OrdersController{Repo: repo}.ReadOrderById(w, r)

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

func TestOderUpdateById(t *testing.T) {
	type expectations struct {
		statusCode  int
		respPayload any
	}
	order := entity.Order{Model: gorm.Model{ID: 1}, TableNumber: 2, FinalPrice: 14.}
	notFoundErr := entity.RecordNotFoundError{
		Kind:  "Order",
		ID:    strconv.FormatInt(int64(order.ID), 10),
		Inner: errors.New("mock repo says no"),
	}
	tests := []struct {
		name        string
		payload     entity.Order
		respPayload any
		existing    entity.Order
		expected    expectations
		err         error
	}{
		{
			name:     "successful update",
			existing: order,
			payload: entity.Order{
				FinalPrice: 16.,
				Model:      gorm.Model{ID: 1},
			},
			expected: expectations{
				statusCode:  http.StatusNoContent,
				respPayload: nil,
			},
		},
		{
			name: "failed to update order",
			err:  notFoundErr,
			expected: expectations{
				statusCode:  http.StatusNotFound,
				respPayload: map[string]interface{}{"Error": notFoundErr.Error()},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			b := bytes.NewBuffer(nil)
			require.NoError(t, json.NewEncoder(b).Encode(tt.payload))
			r := httptest.NewRequest(http.MethodPut, "/orders/{id}", b)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strconv.FormatUint(uint64(tt.existing.ID), 10))
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			repo := new(entity.MockRepo)
			repo.On("UpdateOrder", tt.payload).Return(tt.err)
			OrdersController{Repo: repo}.UpdateOderById(w, r)

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
func TestDiscountUpdateById(t *testing.T) {
	type expectations struct {
		statusCode  int
		respPayload any
	}
	order := entity.Order{Model: gorm.Model{ID: 1}, TableNumber: 2, FinalPrice: 14.}
	dish := entity.Dish{Model: gorm.Model{ID: 2}, Name: "Fish filet", Price: 10.}
	discountDetail := entity.DiscountDetail{OrderID: 1, DishID: 2, Discount: 2., Dish: dish, Order: order}
	errOrderNotFound := errors.New("order not found")
	errDishNotFound := errors.New("dish not found")
	errDiscount := entity.RecordNotFoundError{}

	tests := []struct {
		name             string
		payload          entity.DiscountDetail
		respPayload      any
		expected         expectations
		existingOrder    entity.Order
		existingDish     entity.Dish
		existingDiscount entity.DiscountDetail
		orderErr         error
		dishErr          error
		discountErr      error
	}{
		{
			name: "successful update the order",
			payload: entity.DiscountDetail{
				Discount: 1.8,
				DishID:   dish.ID,
				OrderID:  order.ID,
			},
			existingOrder:    order,
			existingDish:     dish,
			existingDiscount: discountDetail,
			expected: expectations{
				statusCode:  http.StatusNoContent,
				respPayload: nil,
			},
		},
		{
			name:          "order doesn't exist",
			orderErr:      errOrderNotFound,
			payload:       discountDetail,
			existingOrder: order,
			//existingDiscount: discountDetail,
			expected: expectations{
				statusCode:  http.StatusUnprocessableEntity,
				respPayload: map[string]interface{}{"Error": errOrderNotFound.Error()},
			},
		},
		{
			name:          "dish doesn't exist",
			dishErr:       errDishNotFound,
			existingDish:  dish,
			existingOrder: order,
			//existingDiscount: discountDetail,
			payload: discountDetail,
			expected: expectations{
				statusCode:  http.StatusUnprocessableEntity,
				respPayload: map[string]interface{}{"Error": errDishNotFound.Error()},
			},
		},
		{
			name:        "failed to update the discount",
			discountErr: errDiscount,
			expected: expectations{
				statusCode:  http.StatusNotFound,
				respPayload: map[string]interface{}{"Error": errDiscount.Error()},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			b := bytes.NewBuffer(nil)
			require.NoError(t, json.NewEncoder(b).Encode(tt.payload))
			r := httptest.NewRequest(http.MethodPost, "/orders/{orderId}/dishes/{dishId}", b)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("orderId", strconv.FormatUint(uint64(tt.existingOrder.ID), 10))
			rctx.URLParams.Add("dishId", strconv.FormatUint(uint64(tt.existingDish.ID), 10))
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			repo := new(entity.MockRepo)
			repo.On("GetOrder", tt.existingOrder.ID).Return(tt.existingOrder, tt.orderErr)
			repo.On("GetDish", tt.existingDish.ID).Return(tt.existingDish, tt.dishErr)
			repo.On("UpdateDiscount", tt.payload).Return(tt.discountErr)

			OrdersController{Repo: repo}.UpdateDiscountById(w, r)

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

func TestOrderDeleteById(t *testing.T) {
	type expectations struct {
		statusCode  int
		respPayload any
	}
	order := entity.Order{Model: gorm.Model{ID: 1}, TableNumber: 2, FinalPrice: 14.}
	notFoundErr := entity.RecordNotFoundError{
		Kind:  "Order",
		ID:    strconv.FormatInt(int64(order.ID), 10),
		Inner: errors.New("mock says no"),
	}

	tests := []struct {
		name        string
		payload     entity.Order
		respPayload any
		existing    entity.Order
		err         error
		expected    expectations
	}{
		{
			name:     "successful deleted order",
			existing: order,
			expected: expectations{
				statusCode:  http.StatusNoContent, 
				respPayload: nil,
			},
		},
		{
			name: "failed to delete order",
			err:  notFoundErr,
			expected: expectations{
				statusCode:  http.StatusNotFound,
				respPayload: map[string]interface{}{"Error": notFoundErr.Error()},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			b := bytes.NewBuffer(nil)
			r := httptest.NewRequest(http.MethodDelete, "/orders/{id}", b)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strconv.FormatUint(uint64(tt.existing.ID), 10))
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			repo := new(entity.MockRepo)
			repo.On("DeleteOrder", tt.existing.ID).Return(tt.err)
			OrdersController{Repo: repo}.DeleteOrderById(w, r)

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
