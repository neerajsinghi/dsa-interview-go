// Package playbook collects small, tested Go idioms that come up again and
// again in coding interviews. Each function is deliberately tiny: the point is
// to have the idiom at your fingertips so you spend interview time on the
// algorithm, not on remembering the standard library.
package playbook

import (
	"cmp"
	"container/heap"
	"math"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ---------------------------------------------------------------------------
// Sorting
// ---------------------------------------------------------------------------

// WordCount pairs a word with how often it appears.
type WordCount struct {
	Word  string
	Count int
}

// SortByCountThenWord sorts by count descending, breaking ties alphabetically.
// It uses sort.Slice, which takes a "less" function over indices.
//
// Gotcha: sort.Slice is NOT stable. If you need ties to keep their input order,
// use sort.SliceStable — or, as here, make the comparator total by adding a
// tie-breaker so the result is deterministic either way.
func SortByCountThenWord(items []WordCount) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count != items[j].Count {
			return items[i].Count > items[j].Count
		}
		return items[i].Word < items[j].Word
	})
}

// SortByLength sorts words by length, then lexicographically, using the
// generic slices.SortFunc (Go 1.21+). The comparator returns a negative, zero
// or positive int — cmp.Compare and cmp.Or make multi-key comparisons tidy.
func SortByLength(words []string) {
	slices.SortFunc(words, func(a, b string) int {
		return cmp.Or(
			cmp.Compare(len(a), len(b)),
			strings.Compare(a, b),
		)
	})
}

// SortedKeys returns the keys of m in ascending order.
//
// Gotcha: map iteration order in Go is deliberately randomised. Any time your
// output depends on map order (tests, "return in any order" problems you want
// to compare), collect the keys and sort them.
func SortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// ---------------------------------------------------------------------------
// container/heap
// ---------------------------------------------------------------------------

// MinHeap is an int min-heap usable with container/heap. You implement five
// methods: Len, Less, Swap (from sort.Interface) plus Push and Pop. Flip Less
// to get a max-heap.
type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push is called by heap.Push — never call h.Push directly.
func (h *MinHeap) Push(x any) { *h = append(*h, x.(int)) }

// Pop is called by heap.Pop. It removes the LAST element; heap.Pop has
// already swapped the minimum there.
func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// KLargest returns the k largest values in descending order. It keeps a
// min-heap of size k: the root is the smallest of the current top-k, so any
// bigger value evicts it. O(n log k) time, O(k) space.
func KLargest(nums []int, k int) []int {
	if k <= 0 {
		return []int{}
	}
	h := &MinHeap{}
	for _, x := range nums {
		heap.Push(h, x)
		if h.Len() > k {
			heap.Pop(h) // drop the smallest
		}
	}
	out := make([]int, h.Len())
	for i := len(out) - 1; i >= 0; i-- {
		out[i] = heap.Pop(h).(int)
	}
	return out
}

// ---------------------------------------------------------------------------
// Strings
// ---------------------------------------------------------------------------

// JoinInts renders nums separated by sep. strings.Builder avoids the O(n²)
// cost of repeated string concatenation (strings are immutable, so s += x
// copies the whole string every time).
func JoinInts(nums []int, sep string) string {
	var b strings.Builder
	for i, x := range nums {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(strconv.Itoa(x))
	}
	return b.String()
}

// ReverseString reverses s by Unicode code point.
//
// Gotcha: indexing a string (s[i]) gives you a byte, not a character.
// Reversing bytes would corrupt multi-byte UTF-8 characters like "é" or "世".
// Convert to []rune when the input may contain non-ASCII text.
func ReverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// RuneAndByteLen reports the number of code points and the number of bytes
// in s. For ASCII they match; for "héllo" they don't (5 runes, 6 bytes).
// Note that `for i, r := range s` yields byte offsets i and runes r.
func RuneAndByteLen(s string) (runes, bytes int) {
	return utf8.RuneCountInString(s), len(s)
}

// LetterCounts counts lowercase ASCII letters with a fixed [26]int array —
// faster than a map and, being an array (not a slice), comparable with ==.
// That makes "are these two strings anagrams?" a one-liner.
func LetterCounts(s string) [26]int {
	var c [26]int
	for i := 0; i < len(s); i++ {
		if ch := s[i]; ch >= 'a' && ch <= 'z' {
			c[ch-'a']++
		}
	}
	return c
}

// ---------------------------------------------------------------------------
// Maps as sets
// ---------------------------------------------------------------------------

// Unique returns the distinct values of nums in ascending order. Go has no
// set type; map[T]struct{} is the idiom (struct{} takes zero bytes).
func Unique(nums []int) []int {
	seen := make(map[int]struct{}, len(nums))
	out := []int{}
	for _, x := range nums {
		if _, ok := seen[x]; !ok {
			seen[x] = struct{}{}
			out = append(out, x)
		}
	}
	slices.Sort(out)
	return out
}

// ---------------------------------------------------------------------------
// Slices: copying and aliasing
// ---------------------------------------------------------------------------

// AliasingDemo shows the classic append trap. b and c are both appended to
// the same sub-slice a[:2], which still has spare capacity — so both appends
// write into a's backing array and c overwrites b's last element.
// It returns b and c so the test can show they ended up equal.
func AliasingDemo() (b, c []int) {
	a := []int{1, 2, 3, 4}
	b = append(a[:2], 10) // writes a[2] = 10
	c = append(a[:2], 20) // writes a[2] = 20 — b sees it too!
	return b, c
}

// SafeAppend appends x to base without ever touching base's backing array.
// slices.Clip trims capacity to length, forcing append to allocate.
func SafeAppend(base []int, x int) []int {
	return append(slices.Clip(base), x)
}

// AllPrefixes returns [nums[:1], nums[:2], ...] as independent slices.
// In backtracking you hit the same bug: appending your working slice to the
// result stores a VIEW that later mutations will change. Always copy.
func AllPrefixes(nums []int) [][]int {
	res := [][]int{}
	cur := []int{}
	for _, x := range nums {
		cur = append(cur, x)
		res = append(res, slices.Clone(cur)) // or: append([]int(nil), cur...)
	}
	return res
}

// ---------------------------------------------------------------------------
// Numbers
// ---------------------------------------------------------------------------

// MinMax returns the smallest and largest values in nums, or ok=false for an
// empty slice. Seed with math.MaxInt / math.MinInt and use the min/max
// builtins (Go 1.21+) instead of hand-written helpers.
func MinMax(nums []int) (lo, hi int, ok bool) {
	if len(nums) == 0 {
		return 0, 0, false
	}
	lo, hi = math.MaxInt, math.MinInt
	for _, x := range nums {
		lo = min(lo, x)
		hi = max(hi, x)
	}
	return lo, hi, true
}

// Abs returns |x|. Go has no integer abs in the standard library
// (math.Abs is for float64). Note: Abs(math.MinInt) overflows back to itself.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Midpoint returns the midpoint of lo and hi without the overflow that
// (lo+hi)/2 can hit when both are large.
func Midpoint(lo, hi int) int {
	return lo + (hi-lo)/2
}
