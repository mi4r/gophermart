package storageaccrual

const (
	RewardTypePt      RewardType = "pt"
	RewardTypePercent RewardType = "%"
)

type RewardType string

type Order struct {
	Order string `json:"order"`
	Goods []Good
} // @name Order

type Good struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
} // @name Good

type Reward struct {
	Match      string     `json:"match"`
	Reward     float64    `json:"reward"`
	RewardType RewardType `json:"reward_type"`
} // @name Reward

func (r *Reward) IsEmptyMatch() bool {
	return r.Match == ""
}

func (r *Reward) IsNegative() bool {
	return r.Reward < 0
}

func (r *Reward) IsValidType() bool {
	return r.RewardType == RewardTypePercent ||
		r.RewardType == RewardTypePt
}
