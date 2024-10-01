package storage

import (
	storageaccrual "github.com/mi4r/gophermart/internal/storage/accrual"
	"github.com/mi4r/gophermart/internal/storage/drivers"
	storagemart "github.com/mi4r/gophermart/internal/storage/gophermart"
)

type Storage interface {
	Open() error
	Close()
	Ping() error
	Migrate(path string)
}

type StorageGophermart interface {
	Storage
	UserCreate(user storagemart.User) error
	UserReadOne(login string) (storagemart.User, error)
	UserReadAll() ([]storagemart.User, error)

	UserOrderCreate(login, number string) error
	UserOrderReadOne(number string) (storagemart.Order, error)
	UserOrdersReadByLogin(login string) ([]storagemart.Order, error)

	WithdrawBalance(login, order string, sum, curBalance float64) error
}

func NewStorageGophermart(driverType string, path string) StorageGophermart {
	switch driverType {
	default:
		return drivers.NewPgxDriver(path)
	}
}

type StorageAccrualSystem interface {
	Storage
	RewardCreate(reward storageaccrual.Reward) error
	OrderRegCreate(order storageaccrual.Order) error
}

func NewStorageAccrual(driverType string, path string) StorageAccrualSystem {
	switch driverType {
	default:
		return drivers.NewPgxDriver(path)
	}
}
