package workermart

type Task struct {
	OrderNumber string
}

func NewTask(orderNumber string) Task {
	return Task{
		OrderNumber: orderNumber,
	}
}
