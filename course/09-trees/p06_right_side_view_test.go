package trees

import (
	"math/rand"
	"slices"
	"testing"
)

func TestRightSideView(t *testing.T) {
	tests := []struct{ vals, want []int }{
		{[]int{1, 2, 3, Null, 5, Null, 4}, []int{1, 3, 4}},
		{[]int{1, 2, 3, 4}, []int{1, 3, 4}},
		{[]int{1, Null, 3}, []int{1, 3}},
		{[]int{1, 2}, []int{1, 2}},
		{nil, nil},
	}
	for _, tt := range tests {
		root := FromLevelOrder(tt.vals...)
		for name, fn := range map[string]func(*TreeNode) []int{"bfs": RightSideView, "dfs": RightSideViewDFS} {
			if got := fn(root); !slices.Equal(got, tt.want) {
				t.Errorf("%s(%v) = %v, want %v", name, tt.vals, got, tt.want)
			}
		}
	}
}

func TestRightSideViewRandom(t *testing.T) {
	r := rand.New(rand.NewSource(6))
	for i := 0; i < 300; i++ {
		root := randomSmallTree(r, r.Intn(30), 0, 10)
		var want []int // brute force: last element of every level
		for _, lvl := range LevelOrder(root) {
			want = append(want, lvl[len(lvl)-1])
		}
		if got := RightSideView(root); !slices.Equal(got, want) {
			t.Fatalf("RightSideView(%v) = %v, want %v", ToLevelOrder(root), got, want)
		}
		if got := RightSideViewDFS(root); !slices.Equal(got, want) {
			t.Fatalf("RightSideViewDFS(%v) = %v, want %v", ToLevelOrder(root), got, want)
		}
	}
}
