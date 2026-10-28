package graphs

import (
	"math/rand"
	"testing"
)

func TestLadderLength(t *testing.T) {
	tests := []struct {
		name       string
		begin, end string
		list       []string
		want       int
	}{
		{"classic", "hit", "cog", []string{"hot", "dot", "dog", "lot", "log", "cog"}, 5},
		{"end missing", "hit", "cog", []string{"hot", "dot", "dog", "lot", "log"}, 0},
		{"one step", "a", "c", []string{"a", "b", "c"}, 2},
		{"begin equals end", "hot", "hot", []string{"hot"}, 1},
		{"disconnected", "abc", "xyz", []string{"abd", "xyz"}, 0},
		{"begin in list", "hot", "dog", []string{"hot", "dog", "dot"}, 3},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func(string, string, []string) int{"bfs": LadderLength, "bidirectional": LadderLengthBidirectional} {
			if got := fn(tt.begin, tt.end, tt.list); got != tt.want {
				t.Errorf("%s/%s = %d, want %d", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestLadderLengthRandom(t *testing.T) {
	r := rand.New(rand.NewSource(30))
	word := func() string {
		return string([]byte{byte('a' + r.Intn(3)), byte('a' + r.Intn(3)), byte('a' + r.Intn(3))})
	}
	for iter := 0; iter < 300; iter++ {
		list := make([]string, r.Intn(12))
		for i := range list {
			list[i] = word()
		}
		begin, end := word(), word()
		if r.Intn(2) == 0 && len(list) > 0 {
			end = list[r.Intn(len(list))]
		}
		if a, b := LadderLength(begin, end, list), LadderLengthBidirectional(begin, end, list); a != b {
			t.Fatalf("%s->%s %v: bfs=%d bidirectional=%d", begin, end, list, a, b)
		}
	}
}
