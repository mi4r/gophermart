package storage

type OrderStatus string

const (
	StatusNew        OrderStatus = "NEW"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessed  OrderStatus = "PROCESSED"
)
