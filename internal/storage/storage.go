package storage

// Должен имплементировать CRUD методы
// Create
// Read
// Update
// Delete
type Storage interface {
	Open() error
	Close()

	UserCreate(user User) error
	UserReadOne(login string) (User, error)
	UserReadAll() ([]User, error)
}

func NewStorage(driverType string, path string) Storage {
	switch driverType {
	default:
		return NewPgxDriver(path)
	}
}
