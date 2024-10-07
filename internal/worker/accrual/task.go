package workeraccrual

import (
	storageaccrual "github.com/mi4r/gophermart/internal/storage/accrual"
	storagedefault "github.com/mi4r/gophermart/internal/storage/default"
)

type Task struct {
	Order storageaccrual.Order
}

type TaskResult struct {
	ResultOrder storagedefault.Order
}

func NewTask(order storageaccrual.Order) Task {
	return Task{
		Order: order,
	}
}
