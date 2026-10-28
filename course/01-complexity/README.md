# 01 · Complexity

> Every interview answer ends with "time O(?), space O(?)". This module makes
> those words concrete: you will *measure* the difference between O(log n),
> O(n), O(n²) and O(φⁿ) with Go benchmarks.

## What you'll learn

- What Big-O actually promises (and what it hides).
- How to read complexity off loops, recursion trees and data-structure calls.
- Amortized analysis — why `append` is O(1) even though it sometimes copies everything.
- Space complexity, including the call stack that recursion uses silently.
- How to write Go benchmarks (`BenchmarkXxx`) that make growth visible.

## Recognize the pattern

You use complexity analysis in *every* problem, but these signals mean the
interviewer is explicitly testing it:

- "Can you do better?" — your current solution has an avoidable factor of n.
- The constraints say n ≤ 10⁵ or 10⁶ — O(n²) is too slow; aim for O(n log n) or O(n).
- "Without extra space" / "in O(1) memory" — they want you to count auxiliary space.
- A recursive solution that "times out" — look for overlapping sub-problems.

## The idea

**Big-O describes how cost grows as input grows**, ignoring constants and
lower-order terms. `3n² + 10n + 7` is O(n²) because for large n the n² term
dominates.

Approximate operation counts:

| n | O(1) | O(log n) | O(n) | O(n log n) | O(n²) | O(2ⁿ) |
|---|---|---|---|---|---|---|
| 10 | 1 | 3 | 10 | 33 | 100 | 1 024 |
| 1 000 | 1 | 10 | 1 000 | 10 000 | 10⁶ | ≈ 10³⁰¹ |
| 10⁶ | 1 | 20 | 10⁶ | 2 × 10⁷ | 10¹² | — |

At ~10⁸ simple operations per second, O(n²) on 10⁶ items takes hours; O(n log n) takes a fraction of a second.

### Reading complexity from code

| Code shape | Complexity |
|---|---|
| One loop over n | O(n) |
| Nested loops over n | O(n²) |
| Loop where the range halves each step (`lo`, `hi`, `mid`) | O(log n) |
| Sort, then a loop | O(n log n) |
| Recursion with b branches and depth d | O(b^d) nodes in the call tree |
| Recursion that halves and does O(n) work per level | O(n log n) (merge sort) |

### Recursion trees

Count the nodes of the call tree and the work per node. Naive Fibonacci
branches twice per call and goes n levels deep, so it does ~φⁿ ≈ 1.618ⁿ calls
([p03](p03_fibonacci.go)). Memoization turns the tree into a path: each F(k) is
computed once, so it's O(n).

```
F(5) ─┬─ F(4) ─┬─ F(3) ─┬─ F(2) ...      naive: subtrees repeat
      │        │        └─ F(1)
      │        └─ F(2) ...               <- computed again
      └─ F(3) ...                        <- computed again

memo: F(5) → F(4) → F(3) → F(2) → F(1)   each value once
```

### Amortized analysis: `append`

A slice has a length and a capacity. When `append` runs out of capacity it
allocates a larger array (≈2× while small, ≈1.25× once large) and copies:

```
len : 1  2  3  4  5  6  7  8  9 ... 16 17
cap : 1  2  4  4  8  8  8  8 16 ... 16 32
copy: 0  1  2  0  4  0  0  0  8 ...  0 16    total copies < 2n
```

Expensive steps are rare enough that the *average* cost per append is O(1).
Inserting at the front is different: every insert shifts everything — O(n) each,
O(n²) total ([p04](p04_append_growth.go)).

### Space complexity

Count *auxiliary* space — what you allocate beyond the input — plus the
**recursion stack**. A recursion n levels deep uses O(n) stack even if it never
calls `make` ([p05](p05_recursion_space.go)). Divide-and-conquer that halves
the input is only O(log n) deep.

## Go template

A benchmark that sweeps input sizes so the growth rate is obvious:

```go
func BenchmarkThing(b *testing.B) {
    for _, n := range []int{1_000, 10_000, 100_000} {
        input := makeInput(n)              // build outside the timed loop
        b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
            b.ReportAllocs()
            for range b.N {                // Go 1.22+ range-over-int
                thing(input)
            }
        })
    }
}
```

Run it:

```sh
go test -run=^$ -bench=. -benchmem ./course/01-complexity/
```

## Complexity

| Function | Time | Space |
|---|---|---|
| `LinearSearch` | O(n) | O(1) |
| `BinarySearch` | O(log n) | O(1) |
| `ContainsDuplicateBrute` | O(n²) | O(1) |
| `ContainsDuplicateSort` | O(n log n) | O(n) (copy) |
| `ContainsDuplicate` | O(n) | O(n) |
| `FibNaive` | O(φⁿ) | O(n) stack |
| `FibMemo` | O(n) | O(n) |
| `FibIter` | O(n) | O(1) |
| `AppendN` | O(n) total, O(1) amortized | O(n) |
| `PrependN` | O(n²) | O(n) |
| `SumRecursive` / `SumHalves` / `SumIterative` | O(n) | O(n) / O(log n) / O(1) |

### Sample benchmark results

Measured with `go test -run=^$ -bench=. -benchtime=100x` on an Apple-silicon
laptop. **Your numbers will differ — machine, Go version and load all matter.**
Look at the *ratios between rows*, not the absolute values.

| Benchmark | n | ns/op | Growth |
|---|---|---|---|
| Linear search | 1 000 → 1 000 000 | 346 → 337 627 | ~1000× for 1000× n: **O(n)** |
| Binary search | 1 000 → 1 000 000 | 15 → 31 | ~2× for 1000× n: **O(log n)** |
| Duplicate, brute | 512 → 2 048 | 47 435 → 693 971 | ~15× for 4× n: **O(n²)** |
| Duplicate, hash | 512 → 2 048 | 12 471 → 47 749 | ~4× for 4× n: **O(n)** |
| Fib naive | 20 → 25 | 46 811 → 502 800 | ~11× for +5: **O(φⁿ)** (φ⁵ ≈ 11.1) |
| Fib iterative | 20 → 25 | 17 → 12 | flat (noise): **O(n)**, tiny n |
| Append | 4 000 → 16 000 | 25 531 → 90 086 | ~3.5× for 4× n: **O(n)** |
| Prepend | 4 000 → 16 000 | 1 787 742 → 25 742 557 | ~14× for 4× n: **O(n²)** |

Two lessons hidden in those numbers:

1. **Constants matter at small n.** At n = 128 the O(n²) brute force (3.4 µs)
   beats the O(n) hash set (4.2 µs) because hashing and allocation are expensive.
   Big-O tells you who wins *eventually*.
2. **Preallocation is free speed.** `make([]int, 0, n)` did 1 allocation instead
   of ~20 and ran ~4× faster than letting `append` grow.

## Common pitfalls

- **Dropping the wrong term.** O(n + m) is not O(n) when m can be much bigger than n (e.g. edges vs vertices).
- **Forgetting the stack.** Recursive DFS on a 10⁵-node path is O(n) space. Go won't stack-overflow as early as other languages, but it still counts.
- **Hidden loops in library calls.** `slices.Contains`, `strings.Index`, `copy`, `append(s[:i], s[i+1:]...)` and string concatenation are all O(n).
- **String building with `+=`.** Strings are immutable; each `+=` copies. Use `strings.Builder` — O(n) instead of O(n²).
- **Map operations are O(1) *average*.** Worst case is O(n); mention it if asked about adversarial input.
- **Integer overflow.** `FibIter(93)` silently wraps around. Know your bounds (int64 max ≈ 9.2 × 10¹⁸).
- **Benchmarking mistakes.** Build inputs outside the `b.N` loop; make sure the compiler can't skip the work; don't compare runs from different machines.
- **Sub-slices keep big arrays alive.** `small := big[:3]` pins all of `big` in memory — a space leak that Big-O on paper won't show.

## Problems

| # | Problem | Difficulty | Key idea |
|---|---|---|---|
| 01 | [Linear vs Binary Search](p01_search.go) | Easy | Sorted input lets you discard half each step: O(log n) |
| 02 | [Contains Duplicate, three ways](p02_duplicates.go) | Easy | O(n²) pairs → O(n log n) sort → O(n) hash set |
| 03 | [Fibonacci: naive vs memo vs iterative](p03_fibonacci.go) | Easy | Recursion tree; memoization collapses it to a path |
| 04 | [Amortized Append](p04_append_growth.go) | Easy | Geometric growth ⇒ O(1) amortized; prepend is O(n²) |
| 05 | [Space Complexity of Recursion](p05_recursion_space.go) | Easy | Stack depth n vs log n vs 1 |

## Practice tip

Before running the benchmarks, write down your prediction for each ratio, then run:

```sh
go test -run=^$ -bench=Fib -benchtime=100x ./course/01-complexity/
go test ./course/01-complexity/ -run TestCountGrowthsIsLogarithmic -v
```

Then pick any problem in later modules, cover the "Time/Space" line in its doc
comment, and derive it yourself before peeking.

## Next

[02 · Arrays & Hashing →](../02-arrays-hashing/README.md)
