package backtracking

import (
	"reflect"
	"strings"
	"testing"
)

// validBoard checks that a rendered board has n queens, none attacking.
func validBoard(rows []string) bool {
	n := len(rows)
	var qs [][2]int
	for r, row := range rows {
		if len(row) != n || strings.Count(row, "Q") != 1 {
			return false
		}
		qs = append(qs, [2]int{r, strings.IndexByte(row, 'Q')})
	}
	for i := range qs {
		for j := i + 1; j < len(qs); j++ {
			dr, dc := qs[i][0]-qs[j][0], qs[i][1]-qs[j][1]
			if dc == 0 || dr == dc || dr == -dc {
				return false
			}
		}
	}
	return true
}

func TestSolveNQueens(t *testing.T) {
	want4 := [][]string{{".Q..", "...Q", "Q...", "..Q."}, {"..Q.", "Q...", "...Q", ".Q.."}}
	if got := SolveNQueens(4); !reflect.DeepEqual(got, want4) {
		t.Errorf("SolveNQueens(4) = %q, want %q", got, want4)
	}
	if got := SolveNQueens(1); !reflect.DeepEqual(got, [][]string{{"Q"}}) {
		t.Errorf("SolveNQueens(1) = %q", got)
	}
	known := []int{1, 1, 0, 0, 2, 10, 4, 40, 92} // index = n (n=0: one empty board)
	for n := 1; n < len(known); n++ {
		boards := SolveNQueens(n)
		if len(boards) != known[n] {
			t.Errorf("SolveNQueens(%d) found %d boards, want %d", n, len(boards), known[n])
		}
		seen := map[string]bool{}
		for _, b := range boards {
			if !validBoard(b) {
				t.Fatalf("n=%d: invalid board %q", n, b)
			}
			key := strings.Join(b, "/")
			if seen[key] {
				t.Fatalf("n=%d: duplicate board %q", n, b)
			}
			seen[key] = true
		}
	}
}

func TestTotalNQueens(t *testing.T) {
	known := []int{1, 1, 0, 0, 2, 10, 4, 40, 92, 352}
	for n, want := range known {
		if got := TotalNQueens(n); got != want {
			t.Errorf("TotalNQueens(%d) = %d, want %d", n, got, want)
		}
	}
}
