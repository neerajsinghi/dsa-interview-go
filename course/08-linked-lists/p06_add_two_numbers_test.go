package linkedlists

import (
	"math/big"
	"math/rand"
	"slices"
	"testing"
)

func TestAddTwoNumbers(t *testing.T) {
	tests := []struct{ a, b, want []int }{
		{[]int{2, 4, 3}, []int{5, 6, 4}, []int{7, 0, 8}},
		{[]int{0}, []int{0}, []int{0}},
		{[]int{9, 9}, []int{1}, []int{0, 0, 1}},
		{[]int{9, 9, 9, 9, 9, 9, 9}, []int{9, 9, 9, 9}, []int{8, 9, 9, 9, 0, 0, 0, 1}},
		{[]int{5}, []int{5}, []int{0, 1}},
	}
	for _, tt := range tests {
		if got := ToSlice(AddTwoNumbers(FromSlice(tt.a), FromSlice(tt.b))); !slices.Equal(got, tt.want) {
			t.Errorf("AddTwoNumbers(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

// digitsToBig reads least-significant-first digits as a big.Int.
func digitsToBig(d []int) *big.Int {
	n := new(big.Int)
	for i := len(d) - 1; i >= 0; i-- {
		n.Mul(n, big.NewInt(10))
		n.Add(n, big.NewInt(int64(d[i])))
	}
	return n
}

// TestAddTwoNumbersRandom checks numbers far beyond int64 against math/big.
func TestAddTwoNumbersRandom(t *testing.T) {
	r := rand.New(rand.NewSource(6))
	digits := func() []int {
		d := make([]int, 1+r.Intn(40))
		for i := range d {
			d[i] = r.Intn(10)
		}
		if len(d) > 1 && d[len(d)-1] == 0 {
			d[len(d)-1] = 1 + r.Intn(9) // no leading zero
		}
		return d
	}
	for iter := 0; iter < 300; iter++ {
		a, b := digits(), digits()
		got := digitsToBig(ToSlice(AddTwoNumbers(FromSlice(a), FromSlice(b))))
		want := new(big.Int).Add(digitsToBig(a), digitsToBig(b))
		if got.Cmp(want) != 0 {
			t.Fatalf("AddTwoNumbers(%v, %v) = %v, want %v", a, b, got, want)
		}
	}
}
