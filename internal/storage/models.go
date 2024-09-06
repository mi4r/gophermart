package storage

// migrations/wallet
type Wallet struct {
	Balance float64 `json:"balance"`
	Scores  int     `json:"scores"`
	IsLock  bool    `json:"is_lock"`
}

// GetBalance возвращает текущее значение баланса
func (w *Wallet) GetBalance() float64 {
	return w.Balance
}

// SetBalance устанавливает новое значение баланса
func (w *Wallet) SetBalance(balance float64) {
	w.Balance = balance
}

// GetGScores возвращает текущее количество бонусов
func (w *Wallet) GetScores() int {
	return w.Scores
}

// SetGScores устанавливает новое количество бонусов
func (w *Wallet) SetScores(scores int) {
	w.Scores = scores
}

// migrations/users
type User struct {
	Login string `json:"login"`

	// Hash(password)
	Password string `json:"password"`
	Wallet   Wallet `json:"wallet"`
}
