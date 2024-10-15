package workeraccrual

import (
	"testing"

	storageaccrual "github.com/mi4r/gophermart/internal/storage/accrual"
)

func Test_calculateReward(t *testing.T) {
	type args struct {
		price      float64
		reward     float64
		rewardType storageaccrual.RewardType
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "percent",
			args: args{
				price:      500,
				reward:     10,
				rewardType: storageaccrual.RewardTypePercent,
			},
			want: 50,
		},
		{
			name: "point",
			args: args{
				price:      500,
				reward:     10,
				rewardType: storageaccrual.RewardTypePt,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateReward(tt.args.price, tt.args.reward, tt.args.rewardType); got != tt.want {
				t.Errorf("calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}
