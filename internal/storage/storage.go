package storage

// Должен имплементировать CRUD методы
// Create
// Read
// Update
// Delete
type Storage interface {
}

func NewStorage(driverType string) Storage {
	switch driverType {
	// TODO
	}
	return nil
}
