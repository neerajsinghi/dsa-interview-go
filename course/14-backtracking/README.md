# 14 · Backtracking

## What you'll learn

- The single **choose → explore → unchoose** template behind subsets, permutations, combinations, partitions, grid search and puzzles.
- The three enumeration shapes: **subsets/combinations** (start index), **permutations** (used[] array), **positions** (one decision per slot).
- How to **deduplicate** when the input has repeats — sort, then skip equal siblings.
- How **pruning** (sorted break, constraint sets, bitmasks, most-constrained-first) turns hopeless searches into fast ones.

## Recognize the pattern

- "Return **all** …" — all subsets, all permutations, all combinations, all partitions, all boards.
- "Is there **any** arrangement that …" when no greedy/DP structure is apparent (Word Search, Sudoku).
- Input sizes are tiny (n ≤ 10–20, a 9×9 grid) — the interviewer expects exponential time.
- If the question asks for a **count** or an **optimum** and subproblems repeat, it is probably DP instead (module 15); backtracking enumerates, DP aggregates.

## The idea

Backtracking walks a **decision tree** depth-first. Each level makes one decision, the path from the root is the partial answer, and leaves (or every node, for subsets) are complete answers. When a partial answer can't possibly succeed, you return early — the whole subtree is pruned.

```
Subsets of [1 2 3] (record at every node):

                      []
            /          |         \
          [1]         [2]        [3]
        /     \        |
     [1 2]   [1 3]   [2 3]
       |
    [1 2 3]

Permutations of [1 2 3] (record at leaves; any unused element can come next):

                         []
           /             |             \
         [1]            [2]            [3]
       /     \        /     \        /     \
   [1 2]   [1 3]  [2 1]   [2 3]  [3 1]   [3 2]
     |       |      |       |      |       |
  [1 2 3] [1 3 2] [2 1 3] [2 3 1] [3 1 2] [3 2 1]
```

| Shape | Loop | Recurse with | Record |
|---|---|---|---|
| Subsets | `i := start..n-1` | `i + 1` | every call |
| Combinations (k of n) | `i := start..n-1` | `i + 1` | when `len == k` |
| Combination sum (reuse) | `i := start..n-1` | `i` | when `remaining == 0` |
| Permutations | `i := 0..n-1`, skip `used[i]` | — | when `len == n` |
| Positions (phone letters, queens) | each option for slot `i` | `i + 1` | when `i == n` |

**Duplicates in the input:** sort first, then skip `nums[i]` when it equals `nums[i-1]` *at the same depth* (`i > start` for subsets, `!used[i-1]` for permutations).

## Go template

```go
var res [][]int
var path []int
var backtrack func(start int)
backtrack = func(start int) {
	if done(path) {
		res = append(res, append([]int(nil), path...)) // COPY the path
		return
	}
	for i := start; i < len(nums); i++ {
		if !valid(i) {
			continue // prune
		}
		path = append(path, nums[i]) // choose
		backtrack(i + 1)             // explore
		path = path[:len(path)-1]    // unchoose
	}
}
backtrack(0)
```

A worked example (`Combine(n, k)`) is in [backtracking.go](backtracking.go).

## Complexity

| Problem | Number of answers | Time |
|---|---|---|
| Subsets | `2^n` | `O(n · 2^n)` |
| Permutations | `n!` | `O(n · n!)` |
| k-combinations | `C(n, k)` | `O(k · C(n, k))` |
| Phone letters | up to `4^n` | `O(n · 4^n)` |
| Palindrome partitions | up to `2^(n-1)` | `O(n · 2^n)` |
| N-Queens | — | `O(n!)` |

The `O(n)` factor is the cost of copying each answer. Space is `O(n)` recursion depth plus the output.

## Common pitfalls

- **Appending a shared slice.** `res = append(res, path)` stores a slice header that points at the *same* backing array you keep mutating — all your answers end up identical or garbage. Always copy: `append([]int(nil), path...)` (or a `clone` helper).
- **`append([]int(nil), empty...)` is `nil`, not `[]int{}`.** Harmless for printing, but `reflect.DeepEqual(nil, []int{})` is false — matters in tests. The course's `clone` uses `make` + `copy` to always return non-nil.
- **Forgetting to unchoose.** Every mutation before the recursive call needs its inverse after it: pop the path, clear `used[i]`, restore the grid cell, clear the bitmask bit.
- **Capturing loop variables.** Since Go 1.22 each iteration has its own variable, but closures passed around in older code bases may still see the last value. Don't rely on it — pass indices explicitly.
- **Mutating the caller's input.** Sorting `nums` in place for dedup changes the caller's slice. Copy first, as these solutions do.
- **Skipping duplicates with the wrong condition.** For subsets it's `i > start` (same depth), not `i > 0` (which would also forbid `[2 2]`).
- **Not pruning a sorted search.** In Combination Sum, `break` (not `continue`) once a candidate is too big — everything after it is bigger too.
- **Deep recursion is fine here** — depth is at most n or the number of empty cells — but exponential *breadth* is the real cost. Prune early.

## Problems

| # | Problem | Difficulty | Key idea |
|---|---|---|---|
| 01 | [Letter Combinations of a Phone Number](p01_letter_combinations.go) | Medium | One decision per digit; iterative product as an alternative |
| 02 | [Subsets I & II](p02_subsets.go) | Medium | Record at every node; sort + skip equal siblings for duplicates |
| 03 | [Permutations I & II](p03_permutations.go) | Medium | used[] array or in-place swaps; `!used[i-1]` rule for duplicates |
| 04 | [Combination Sum](p04_combination_sum.go) | Medium | Recurse with `i` to allow reuse; sorted `break` prune |
| 05 | [Palindrome Partitioning](p05_palindrome_partitioning.go) | Medium | Cut positions + precomputed `isPal` table |
| 06 | [Word Search](p06_word_search.go) | Medium | Grid DFS, mark/restore cell; letter-count and reverse-word pruning |
| 07 | [N-Queens I & II](p07_n_queens.go) | Hard | One queen per row; column / `r-c` / `r+c` sets; bitmask counting |
| 08 | [Sudoku Solver](p08_sudoku_solver.go) | Hard | Bitmask candidates + fill the most constrained cell first |

## Practice tip

Before coding, draw the first two levels of the decision tree and decide: *what is one decision? when is the path an answer? what makes a branch invalid?* Then hide the solution and re-solve:

```bash
go test ./course/14-backtracking/ -run TestPermuteUnique -v
```

The randomized tests compare against brute-force enumerations (bitmasks, dedupe-by-map), so they catch both missing and duplicated answers.

## Next

[15 · Dynamic programming (1-D)](../15-dp-1d/README.md) — when the recursion tree repeats itself, cache it.
