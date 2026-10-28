package playbook

import (
	"math"
	"slices"
	"testing"
)

func TestSortByCountThenWord(t *testing.T) {
	items := []WordCount{{"b", 2}, {"a", 2}, {"c", 5}, {"d", 1}}
	SortByCountThenWord(items)
	want := []WordCount{{"c", 5}, {"a", 2}, {"b", 2}, {"d", 1}}
	if !slices.Equal(items, want) {
		t.Fatalf("got %v, want %v", items, want)
	}
}

func TestSortByLength(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{"empty", []string{}, []string{}},
		{"single", []string{"go"}, []string{"go"}},
		{"mixed", []string{"ccc", "a", "bb", "aa", "b"}, []string{"a", "b", "aa", "bb", "ccc"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := slices.Clone(tc.in)
			SortByLength(got)
			if !slices.Equal(got, tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSortedKeys(t *testing.T) {
	m := map[string]int{"pear": 1, "apple": 2, "fig": 3}
	// Run several times: without sorting, map order would vary between runs.
	for range 5 {
		if got, want := SortedKeys(m), []string{"apple", "fig", "pear"}; !slices.Equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestKLargest(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want []int
	}{
		{"empty", nil, 3, []int{}},
		{"k zero", []int{1, 2}, 0, []int{}},
		{"k larger than n", []int{3, 1}, 5, []int{3, 1}},
		{"duplicates", []int{5, 1, 5, 3, 5}, 2, []int{5, 5}},
		{"negatives", []int{-1, -7, -3, -2}, 3, []int{-1, -2, -3}},
		{"typical", []int{3, 2, 1, 5, 6, 4}, 2, []int{6, 5}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := KLargest(tc.nums, tc.k); !slices.Equal(got, tc.want) {
				t.Fatalf("KLargest(%v, %d) = %v, want %v", tc.nums, tc.k, got, tc.want)
			}
		})
	}
}

func TestJoinInts(t *testing.T) {
	tests := []struct {
		nums []int
		sep  string
		want string
	}{
		{nil, ",", ""},
		{[]int{7}, ",", "7"},
		{[]int{1, -2, 3}, ", ", "1, -2, 3"},
	}
	for _, tc := range tests {
		if got := JoinInts(tc.nums, tc.sep); got != tc.want {
			t.Errorf("JoinInts(%v, %q) = %q, want %q", tc.nums, tc.sep, got, tc.want)
		}
	}
}

func TestReverseString(t *testing.T) {
	tests := map[string]string{
		"":      "",
		"a":     "a",
		"abc":   "cba",
		"héllo": "olléh",
		"Go世界":  "界世oG",
	}
	for in, want := range tests {
		if got := ReverseString(in); got != want {
			t.Errorf("ReverseString(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestRuneAndByteLen(t *testing.T) {
	tests := []struct {
		s            string
		runes, bytes int
	}{
		{"", 0, 0},
		{"hello", 5, 5},
		{"héllo", 5, 6},
		{"世界", 2, 6},
	}
	for _, tc := range tests {
		r, b := RuneAndByteLen(tc.s)
		if r != tc.runes || b != tc.bytes {
			t.Errorf("RuneAndByteLen(%q) = (%d, %d), want (%d, %d)", tc.s, r, b, tc.runes, tc.bytes)
		}
	}
}

func TestLetterCounts(t *testing.T) {
	if LetterCounts("listen") != LetterCounts("silent") {
		t.Error("listen/silent should have equal letter counts")
	}
	if LetterCounts("rat") == LetterCounts("car") {
		t.Error("rat/car should differ")
	}
	if c := LetterCounts("aab"); c[0] != 2 || c[1] != 1 {
		t.Errorf("unexpected counts %v", c)
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		in, want []int
	}{
		{nil, []int{}},
		{[]int{1}, []int{1}},
		{[]int{3, 1, 3, -2, 1}, []int{-2, 1, 3}},
	}
	for _, tc := range tests {
		if got := Unique(tc.in); !slices.Equal(got, tc.want) {
			t.Errorf("Unique(%v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestAliasingDemo(t *testing.T) {
	b, c := AliasingDemo()
	// The surprising part: b was "changed" by the second append.
	if !slices.Equal(b, []int{1, 2, 20}) || !slices.Equal(c, []int{1, 2, 20}) {
		t.Fatalf("b=%v c=%v; expected both to be [1 2 20] due to aliasing", b, c)
	}
}

func TestSafeAppend(t *testing.T) {
	a := []int{1, 2, 3, 4}
	b := SafeAppend(a[:2], 10)
	c := SafeAppend(a[:2], 20)
	if !slices.Equal(b, []int{1, 2, 10}) || !slices.Equal(c, []int{1, 2, 20}) {
		t.Fatalf("b=%v c=%v", b, c)
	}
	if !slices.Equal(a, []int{1, 2, 3, 4}) {
		t.Fatalf("base mutated: %v", a)
	}
}

func TestAllPrefixes(t *testing.T) {
	got := AllPrefixes([]int{1, 2, 3})
	want := [][]int{{1}, {1, 2}, {1, 2, 3}}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if !slices.Equal(got[i], want[i]) {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
	if len(AllPrefixes(nil)) != 0 {
		t.Fatal("expected no prefixes for empty input")
	}
}

func TestMinMax(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		lo, hi int
		ok     bool
	}{
		{"empty", nil, 0, 0, false},
		{"single", []int{4}, 4, 4, true},
		{"negatives", []int{-3, -9, -1}, -9, -1, true},
		{"extremes", []int{math.MaxInt, math.MinInt, 0}, math.MinInt, math.MaxInt, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lo, hi, ok := MinMax(tc.nums)
			if lo != tc.lo || hi != tc.hi || ok != tc.ok {
				t.Fatalf("MinMax(%v) = (%d, %d, %v), want (%d, %d, %v)", tc.nums, lo, hi, ok, tc.lo, tc.hi, tc.ok)
			}
		})
	}
}

func TestAbsAndMidpoint(t *testing.T) {
	if Abs(-5) != 5 || Abs(5) != 5 || Abs(0) != 0 {
		t.Error("Abs wrong")
	}
	if got := Midpoint(math.MaxInt-2, math.MaxInt); got != math.MaxInt-1 {
		t.Errorf("Midpoint near MaxInt = %d, want %d", got, math.MaxInt-1)
	}
	if got := Midpoint(0, 9); got != 4 {
		t.Errorf("Midpoint(0, 9) = %d, want 4", got)
	}
}
