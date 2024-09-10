package storage

import (
	"time"
)

type OrderStatus string

const (
	StatusNew        OrderStatus = "NEW"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	Number     string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    int64       `json:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at"`
} //@name Order

type Creds struct {
	Login    string `json:"login"`
	Password string `json:"password"`
} // @name Creds

type User struct {
	Creds
	Balance float64 `json:"balance"`
} //@name User

func (u *User) GetBalance() float64 {
	return u.Balance
}

func (u *User) SetBalance(b float64) {
	u.Balance = b
}
