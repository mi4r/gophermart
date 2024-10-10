package storage

import (
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

	UserOrderCreate(login, number string) error
	UserOrderReadOne(number string) (storagemart.Order, error)
	UserOrdersReadByLogin(login string) ([]storagemart.Order, error)

	WithdrawBalance(login, order string, sum, curBalance float64) error
	GetUserWithdrawals(login string) ([]storagemart.Order, error)
}

func NewStorageGophermart(driverType, path string) StorageGophermart {
	switch driverType {
	default:
		return drivers.NewPgxDriver(path)
	}
}

type StorageAccrualSystem interface {
	Storage
	RewardCreate(reward storageaccrual.Reward) error
	RewardReadAll() ([]storageaccrual.Reward, error)
	OrderRegCreate(order storageaccrual.Order) error
	OrderRegReadOne(number string) (storagedefault.Order, error)
	OrderRegUpdateOne(order storagedefault.Order) error
	// Для безопасности и неизменности Accrual
	OrderRegUpdateStatus(status storagedefault.OrderStatus, number string) error
}

func NewStorageAccrual(driverType, path string) StorageAccrualSystem {
	switch driverType {
	default:
		return drivers.NewPgxDriver(path)
	}
}
