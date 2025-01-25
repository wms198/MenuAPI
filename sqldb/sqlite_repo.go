package sqldb

import (
	"errors"
	"fmt"
	"gorestserviceagain/entity"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Sqlite *gorm.DB

type SqliteDB struct {
	db *gorm.DB
}

func newConnection(path string) error {
	if Sqlite != nil {
		return nil
	}
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return entity.ErrDBNotConnected
	}
	db.Exec("PRAGMA foreign_keys = ON")
	fmt.Printf("DB connects to '%s'\n", path)
	Sqlite = db
	return nil
}

func NewSqlite(path string) (SqliteDB, error) {
	var r SqliteDB
	err := newConnection(path)
	if err != nil {
		return r, err
	}
	r.db = Sqlite
	return r, nil
}

func (r SqliteDB) Migrate() error {
	return errors.Join(
		r.db.AutoMigrate(&entity.Order{}, &entity.Dish{}, &entity.DiscountDetail{}),
		errors.New("error migrating db schema"),
	)
}

func (r SqliteDB) CreateOrder(order *entity.Order) error {
	result := r.db.Create(&order)
	if errors.Is(result.Error, gorm.ErrInvalidData) {
		return fmt.Errorf("%w:%w", entity.ErrInvalidData, result.Error)
	}
	return nil
}

func (r SqliteDB) GetOrders() (o []entity.Order, err error) {
	result := r.db.Find(&o)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w:%w", errors.New("oders not found"), result.Error)
	}
	return o, nil
}

func (r SqliteDB) GetOrder(id uint) (o entity.Order, err error) {
	o.ID = id
	result := r.db.First(&o, o.ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return o, entity.WrapRecordNotFoundError("Order", id, result.Error)
	}
	return o, nil
}
func (r SqliteDB) UpdateOrder(o *entity.Order) error {
	result := r.db.Model(o).Updates(*o)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return entity.WrapRecordNotFoundError("Order", o.ID, result.Error)
	}
	return nil
}

func (r SqliteDB) UpdateDiscount(d *entity.DiscountDetail) error {
	result := r.db.Model(d).Updates(*d)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("OrderId or dishId can not be found, can not update")
	}
	return nil
}
func (r SqliteDB) DeleteOrder(id uint) error {
	var o entity.Order
	result := r.db.Delete(&o, id)
	if result.RowsAffected == 0 {
		return entity.WrapRecordNotFoundError("Order", id, result.Error)
	}
	return nil
}

func (r SqliteDB) CreateDish(dish *entity.Dish) error {
	result := r.db.Create(&dish)
	if errors.Is(result.Error, gorm.ErrInvalidData) {
		return fmt.Errorf("%w:%w", entity.ErrInvalidData, result.Error)
	}
	return nil
}

func (r SqliteDB) GetDishes() (d []entity.Dish, err error) {
	result := r.db.Find(&d)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return d, fmt.Errorf("%w:%w", entity.ErrRecordNotFound, result.Error)
	}
	return d, nil
}

func (r SqliteDB) GetDish(id uint) (d entity.Dish, err error) {
	d.ID = id
	result := r.db.First(&d, d.ID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return d, entity.WrapRecordNotFoundError("Dish", id, result.Error)
	}
	return d, nil
}

func (r SqliteDB) UpdateDish(dish *entity.Dish) error {
	result := r.db.Model(dish).Updates(*dish)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return entity.WrapRecordNotFoundError("Dish", dish.ID, result.Error)
	}
	return nil
}

func (r SqliteDB) DeleteDish(id uint) error {
	var dish entity.Dish
	result := r.db.Delete(&dish, id)
	if result.RowsAffected == 0 {
		return entity.WrapRecordNotFoundError("Dish", id, result.Error)
	}
	return nil
}

func (r SqliteDB) CreateDiscount(price *entity.DiscountDetail) error {
	result := r.db.Create(&price)
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return fmt.Errorf("%s:%w", "the discount price with the same orderId and dishIs has existed", result.Error)
	}
	return nil
}

func (r SqliteDB) GetPriceAfterDiscount(orderId uint, dishId uint) (discountDetail entity.DiscountDetail, err error) {
	result := r.db.Joins("Dish").Joins("Order").Where(&discountDetail).First(&discountDetail)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return discountDetail, fmt.Errorf("orderId or itemId does not exist, just ne delete")
	}
	return discountDetail, nil
}
