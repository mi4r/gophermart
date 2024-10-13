package storage

import (
	"context"

	storageaccrual "github.com/mi4r/gophermart/internal/storage/accrual"
	storagedefault "github.com/mi4r/gophermart/internal/storage/default"
	"github.com/mi4r/gophermart/internal/storage/drivers"
	storagemart "github.com/mi4r/gophermart/internal/storage/gophermart"
)

type Storage interface {
	Open(ctx context.Context) error
	Close()
	Ping() error
	Migrate(path string)
}

type StorageGophermart interface {
	Storage
	UserCreate(ctx context.Context, user storagemart.User) error
	UserReadOne(ctx context.Context, login string) (storagemart.User, error)

	UserOrderCreate(ctx context.Context, login, number string) error
	UserOrderReadOne(ctx context.Context, number string) (storagemart.Order, error)
	UserOrdersReadByLogin(ctx context.Context, login string) ([]storagemart.Order, error)

	WithdrawBalance(ctx context.Context, login, order string, sum, curBalance float64) error
	GetUserWithdrawals(ctx context.Context, login string) ([]storagedefault.WithdrownOrder, error)
	UserOrderReadAllNumbers(ctx context.Context) ([]string, error)
	UserOrderUpdateStatus(ctx context.Context, number string, status storagedefault.OrderStatus) error
	UserOrderUpdateAll(ctx context.Context, orders []storagedefault.Order) error
}

func NewStorageGophermart(driverType, path string) StorageGophermart {
	switch driverType {
	default:
		return drivers.NewPgxDriver(path)
	}
}

type StorageAccrualSystem interface {
	Storage
	RewardCreate(ctx context.Context, reward storageaccrual.Reward) error
	RewardReadAll(ctx context.Context) ([]storageaccrual.Reward, error)
	OrderRegCreate(ctx context.Context, order storageaccrual.Order) error
	OrderRegReadOne(ctx context.Context, number string) (storagedefault.Order, error)
	OrderRegUpdateOne(ctx context.Context, order storagedefault.Order) error
	// Для безопасности и неизменности Accrual
	OrderRegUpdateStatus(ctx context.Context, status storagedefault.OrderStatus, number string) error
}

func NewStorageAccrual(driverType, path string) StorageAccrualSystem {
	switch driverType {
	default:
		return drivers.NewPgxDriver(path)
	}
}
