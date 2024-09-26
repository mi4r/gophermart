package storage

import (
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

	OrderCreate(login, number string) error
	OrderReadOne(number string) (storagemart.Order, error)
	OrdersReadByLogin(login string) ([]storagemart.Order, error)
}

func NewStorageGophermart(driverType string, path string) StorageGophermart {
	switch driverType {
	default:
		return drivers.NewPgxDriver(path)
	}
}

type StorageAccrualSystem interface {
	Storage
	// OrderProcessing(number string) error
}

func NewStorageAccrual(driverType string, path string) StorageAccrualSystem {
	switch driverType {
	default:
		return drivers.NewPgxDriver(path)
	}
}
