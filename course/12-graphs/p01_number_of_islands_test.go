package graphs

import (
	"math/rand"
	"testing"
)

func byteGrid(rows ...string) [][]byte {
	g := make([][]byte, len(rows))
	for i, r := range rows {
		g[i] = []byte(r)
	}
	return g
}

func TestNumIslands(t *testing.T) {
	tests := []struct {
		name string
		grid [][]byte
		want int
	}{
		{"classic", byteGrid("11000", "11000", "00100", "00011"), 3},
		{"one big island", byteGrid("11110", "11010", "11000", "00000"), 1},
		{"empty grid", nil, 0},
		{"all water", byteGrid("000", "000"), 0},
		{"diagonals do not connect", byteGrid("101", "010", "101"), 5},
		{"single land", byteGrid("1"), 1},
		{"ring with lake", byteGrid("111", "101", "111"), 1},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([][]byte) int{"dfs": NumIslands, "unionfind": NumIslandsUnionFind} {
			if got := fn(tt.grid); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestNumIslandsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(23))
	for iter := 0; iter < 300; iter++ {
		m, n := 1+r.Intn(7), 1+r.Intn(7)
		grid := make([][]byte, m)
		for i := range grid {
			grid[i] = make([]byte, n)
			for j := range grid[i] {
				grid[i][j] = "01"[r.Intn(2)]
			}
		}
		if a, b := NumIslands(grid), NumIslandsUnionFind(grid); a != b {
			t.Fatalf("%q: dfs=%d unionfind=%d", grid, a, b)
		}
	}
}
