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
	Number     string      `json:"number" example:"12345678903"`
	Status     OrderStatus `json:"status"`
	Accrual    int64       `json:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at" format:"date-time" example:"2020-12-10T15:15:45+03:00"`
} //@name Order

type Creds struct {
	Login    string `json:"login"`
	Password string `json:"password"`
} // @name Creds

type User struct {
	Creds
	Balance float64 `json:"balance"`
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

func (u *User) GetBalance() float64 {
	return u.Balance
}

func (u *User) SetBalance(b float64) {
	u.Balance = b
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
	oldHash, err := c.Password2Hash()
	if err != nil {
		return false
	}
	newHash, err := creds.Password2Hash()
	if err != nil {
		return false
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(oldHash),
		[]byte(newHash),
	); err != nil {
		return false
	}
	return true
}
