package storage

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type OrderStatus string
type Orders []Order

const (
	StatusNew        OrderStatus = "NEW"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	Number      string      `json:"number" example:"12345678903"`
	Status      OrderStatus `json:"status"`
	Accrual     float64     `json:"accrual,omitempty"`
	Sum         float64     `json:"sum,omitempty"`
	UploadedAt  time.Time   `json:"uploaded_at" format:"date-time" example:"2020-12-10T15:15:45+03:00"`
	UserLogin   string      `json:"-"`
	IsWithdrawn bool        `json:"is_withdrawn"`
} //@name Order

type Creds struct {
	Login    string `json:"login"`
	Password string `json:"password"`
} // @name Creds

type Wallet struct {
	Balance   float64 `json:"balance"`
	Withdrawn float64 `json:"withdrawn"`
} //@name Wallet

type User struct {
	Creds
	Wallet
} //@name User

func NewUserFromCreds(creds Creds) (User, error) {
	hashedPassword, err := creds.Password2Hash()
	if err != nil {
		return User{}, err
	}
	return User{
		Creds: Creds{
			Login:    creds.Login,
			Password: hashedPassword,
		},
	}, nil
}

func (u *User) GetBalance() Wallet {
	return u.Wallet
}

func (c *Creds) IsEmpty() bool {
	if c.Login == "" || c.Password == "" {
		return true
	}
	return false
}

func (c *Creds) Password2Hash() (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (c *Creds) PasswordCompare(creds Creds) bool {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(c.Password),
		[]byte(creds.Password),
	); err != nil {
		return false
	}
	return true
}
