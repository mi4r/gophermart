package storage

type Storage interface {
	Open() error
	Close()
	Ping() error

	UserCreate(user User) error
	UserReadOne(login string) (User, error)
	UserReadAll() ([]User, error)

	OrderCreate(login, number string) error
	OrderReadOne(number string) (Order, error)
	OrdersReadByLogin(login string) ([]Order, error)
}

func NewStorage(driverType string, path string) Storage {
	switch driverType {
	default:
		return NewPgxDriver(path)
	}
}
