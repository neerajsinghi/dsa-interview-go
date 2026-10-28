package trees

import (
	"math"
	"math/rand"
	"testing"
)

func TestSerialize(t *testing.T) {
	tests := []struct {
		vals []int
		want string
	}{
		{nil, "#"},
		{[]int{1}, "1,#,#"},
		{[]int{1, 2, 3, Null, Null, 4, 5}, "1,2,#,#,3,4,#,#,5,#,#"},
		{[]int{-7, Null, -7}, "-7,#,-7,#,#"},
	}
	for _, tt := range tests {
		root := FromLevelOrder(tt.vals...)
		got := Serialize(root)
		if got != tt.want {
			t.Errorf("Serialize(%v) = %q, want %q", tt.vals, got, tt.want)
		}
		back, err := Deserialize(got)
		if err != nil || !Equal(back, root) {
			t.Errorf("Deserialize(%q) = %v, %v; want %v", got, ToLevelOrder(back), err, tt.vals)
		}
	}
}

func TestDeserializeErrors(t *testing.T) {
	for _, bad := range []string{"", "1", "1,#", "x,#,#", "1,#,#,#", "#,#"} {
		if _, err := Deserialize(bad); err == nil {
			t.Errorf("Deserialize(%q) should fail", bad)
		}
	}
}

func TestSerializeRandom(t *testing.T) {
	r := rand.New(rand.NewSource(8))
	for i := 0; i < 300; i++ {
		root := randomSmallTree(r, r.Intn(40), -3, 3) // duplicates and negatives
		if i == 0 {
			root = &TreeNode{Val: math.MaxInt, Left: &TreeNode{Val: math.MinInt}}
		}
		s := Serialize(root)
		back, err := Deserialize(s)
		if err != nil || !Equal(back, root) {
			t.Fatalf("round trip of %q failed: %v", s, err)
		}
	}
}
