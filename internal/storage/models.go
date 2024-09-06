package storage

// migrations/balances
type Balance struct {
	Value   float64 `json:"value"`
	Bonuses int     `json:"bonuses"`
}

// GetValue возвращает текущее значение баланса
func (b *Balance) GetValue() float64 {
	return b.Value
}

// SetValue устанавливает новое значение баланса
func (b *Balance) SetValue(value float64) {
	b.Value = value
}

// GetBonuses возвращает текущее количество бонусов
func (b *Balance) GetBonuses() int {
	return b.Bonuses
}

// SetBonuses устанавливает новое количество бонусов
func (b *Balance) SetBonuses(bonuses int) {
	b.Bonuses = bonuses
}

// migrations/users
type User struct {
	Login string `json:"login"`

	// Hash(password)
	Password string  `json:"password"`
	Balance  Balance `json:"balance"`
}
