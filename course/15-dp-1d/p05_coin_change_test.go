package dp1d

import (
	"math/rand"
	"testing"
)

var coinChangeImpls = map[string]func([]int, int) int{
	"recursive": CoinChangeRecursive, "memo": CoinChangeMemo, "tab": CoinChange,
}

var changeImpls = map[string]func(int, []int) int{
	"recursive": ChangeRecursive, "memo": ChangeMemo, "tab": ChangeTab, "optimized": Change,
}

func TestCoinChange(t *testing.T) {
	tests := []struct {
		name   string
		coins  []int
		amount int
		want   int
	}{
		{"classic", []int{1, 2, 5}, 11, 3},
		{"impossible", []int{2}, 3, -1},
		{"zero amount", []int{1}, 0, 0},
		{"greedy fails", []int{1, 3, 4}, 6, 2}, // 3+3, not 4+1+1
		{"no coins", nil, 5, -1},
		{"exact single coin", []int{7}, 14, 2},
	}
	for _, tt := range tests {
		for name, fn := range coinChangeImpls {
			if got := fn(tt.coins, tt.amount); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestChange(t *testing.T) {
	tests := []struct {
		name   string
		amount int
		coins  []int
		want   int
	}{
		{"classic", 5, []int{1, 2, 5}, 4},
		{"impossible", 3, []int{2}, 0},
		{"one coin exact", 10, []int{10}, 1},
		{"zero amount", 0, []int{1, 2}, 1}, // the empty combination
		{"no coins", 4, nil, 0},
	}
	for _, tt := range tests {
		for name, fn := range changeImpls {
			if got := fn(tt.amount, tt.coins); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
	// Loop order matters: amount 3 with coins {1,2} has 2 combinations
	// (1+1+1, 1+2) but 3 orderings (1+1+1, 1+2, 2+1).
	if got := Change(3, []int{1, 2}); got != 2 {
		t.Errorf("Change(3, [1 2]) = %d, want 2", got)
	}
	if got := CountOrderings(3, []int{1, 2}); got != 3 {
		t.Errorf("CountOrderings(3, [1 2]) = %d, want 3", got)
	}
}

func TestCoinChangeRandom(t *testing.T) {
	r := rand.New(rand.NewSource(74))
	for iter := 0; iter < 200; iter++ {
		var coins []int
		for _, v := range r.Perm(9)[:1+r.Intn(4)] {
			coins = append(coins, v+1) // distinct denominations 1..9
		}
		amount := r.Intn(25)
		want := CoinChangeRecursive(coins, amount)
		for name, fn := range coinChangeImpls {
			if got := fn(coins, amount); got != want {
				t.Fatalf("coins %v amount %d: %s=%d brute=%d", coins, amount, name, got, want)
			}
		}
		wantWays := ChangeRecursive(amount, coins)
		for name, fn := range changeImpls {
			if got := fn(amount, coins); got != wantWays {
				t.Fatalf("coins %v amount %d: %s=%d brute=%d", coins, amount, name, got, wantWays)
			}
		}
	}
}
