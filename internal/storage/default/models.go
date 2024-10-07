package storagedefault

type OrderStatus string

const (
	StatusNew        OrderStatus = "NEW"
	StatusRegistered OrderStatus = "REGISTERED"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	Number  string      `json:"number" example:"12345678903"`
	Status  OrderStatus `json:"status"`
	Accrual float64     `json:"accrual,omitempty"`
}
