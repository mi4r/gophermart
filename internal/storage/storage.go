package storage

import (
	"context"

	storageaccrual "github.com/mi4r/gophermart/internal/storage/accrual"
	storagedefault "github.com/mi4r/gophermart/internal/storage/default"
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
	RewardReadAll(ctx context.Context) ([]storageaccrual.Reward, error)
	OrderRegCreate(order storageaccrual.Order) error
	OrderRegReadOne(ctx context.Context, number string) (storagedefault.Order, error)
	OrderRegUpdateOne(ctx context.Context, order storagedefault.Order) error
	// Для безопасности и неизменности Accrual
	OrderRegUpdateStatus(ctx context.Context, status storagedefault.OrderStatus, number string) error
}

func NewStorageAccrual(driverType string, path string) StorageAccrualSystem {
	switch driverType {
	default:
		return drivers.NewPgxDriver(path)
	}
}
