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
	Accrual    float64     `json:"accrual,omitempty"`
	Sum        float64     `json:"sum,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at" format:"date-time" example:"2020-12-10T15:15:45+03:00"`
	UserLogin  string      `json:"-"`
} //@name Order

type Creds struct {
	Login    string `json:"login"`
	Password string `json:"password"`
} // @name Creds

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
} //@name Balance

type User struct {
	Creds
	Balance
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

func (u *User) GetBalance() Balance {
	return u.Balance
}

func (u *User) SetBalance(b Balance) {
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
	if err := bcrypt.CompareHashAndPassword(
		[]byte(c.Password),
		[]byte(creds.Password),
	); err != nil {
		return false
	}
	return true
}
