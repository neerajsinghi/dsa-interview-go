# 00 · Interview Playbook

> Before any pattern: a repeatable way to run a 45-minute coding interview, the
> complexity numbers you should know cold, and the Go idioms that keep you fast.

## Explain it like I'm a kid

A coding interview is like showing a friend how you solve a puzzle. They care about how you think, not only whether you finish. So follow the same seven steps every time:

1. **Clarify:** ask questions until you know the rules of the game.
2. **Examples:** try the puzzle by hand with a small example.
3. **Brute force:** say the simple, slow way out loud, like looking under every cushion for your lost toy.
4. **Optimise:** find the slow part and think of a better way.
5. **Code:** write it neatly, so a friend could read it.
6. **Test:** pretend to be the computer and follow your code step by step.
7. **Analyse:** say how much work your solution does as the puzzle gets bigger.

The other parts of this lesson are helpers:

- **Thinking out loud:** like telling a friend your plan during a board game, so they can follow along and help.
- **Big-O cheat sheet:** a chart that shows how fast each tool is.
- **Pattern recognition:** clue words that tell you which trick to use, like a detective spotting muddy footprints.

**Remember:** when you feel stuck, go to the next step. You always know what to do next.

## What you'll learn

- A seven-step framework that works for almost any coding question.
- How to think out loud so the interviewer can follow (and help).
- A Big-O cheat sheet for Go data structures and classic algorithms.
- How to map wording in a problem statement to a technique — and to the module in this course that teaches it.
- How to recognize the pattern during the interview itself: four questions, decision flowcharts, worked examples and common look-alikes.
- The Go standard-library idioms you will reach for in nearly every problem (see [`gotchas.go`](gotchas.go)).

## The framework

```mermaid
flowchart LR
    A[1. Clarify] --> B[2. Examples]
    B --> C[3. Brute force]
    C --> D[4. Optimise]
    D --> E[5. Code]
    E --> F[6. Test]
    F --> G[7. Analyse]
```

| Step | What you do | Time (45 min slot) |
|---|---|---|
| **1. Clarify** | Restate the problem. Ask about input size, value ranges (negatives? zero? overflow?), duplicates, empty input, sortedness, return format, whether you may mutate the input. | 2–4 min |
| **2. Examples** | Work one normal example and one edge case by hand. Write them down — they become your tests. | 2–3 min |
| **3. Brute force** | State the obvious solution and its complexity *out loud*, even if you won't code it. It proves you understand the problem and gives a baseline. | 2 min |
| **4. Optimise** | Find the bottleneck in the brute force ("the inner loop searches for a complement…") and remove it with a pattern (hash map, two pointers, window, prefix sums, heap…). Agree the approach with the interviewer *before* coding. | 5–10 min |
| **5. Code** | Write clean Go: good names, small helpers, no clever one-liners. Narrate intent, not syntax. | 10–15 min |
| **6. Test** | Trace your examples through the code line by line. Then the edge cases: empty, one element, all equal, negatives, max values. | 3–5 min |
| **7. Analyse** | State time and space complexity and *why*. Mention trade-offs and what you'd do with more time (or different constraints). | 1–2 min |

Constraints are a hint. Rough budget for ~10⁸ simple operations per second:

| n up to | Acceptable complexity |
|---|---|
| 10–12 | O(n!) , O(2ⁿ · n) — permutations, brute-force subsets |
| 20–25 | O(2ⁿ) — bitmask / backtracking |
| 500 | O(n³) |
| 5 000 | O(n²) |
| 10⁶ | O(n log n) — sort, heap |
| 10⁸ | O(n) or O(log n) |

## How to talk while coding

- **Announce the plan, then the step.** "I'll keep a map from value to index. For each number I check if its complement is already there."
- **Name invariants.** "Everything left of `slow` is deduplicated." Invariants are what interviewers listen for.
- **Flag shortcuts honestly.** "I'm assuming the input fits in memory; if not, I'd stream it."
- **Don't go silent when stuck.** Say what you're trying: "Sorting would make this two pointers but costs O(n log n) — is that OK?"
- **Accept hints gracefully.** A hint is not a failure; ignoring one is.
- **Leave time to test.** Finding your own bug is a positive signal. Running out of time before testing is not.

## Big-O cheat sheet

### Data structures (Go)

| Structure | Access | Search | Insert | Delete | Notes |
|---|---|---|---|---|---|
| Slice `[]T` | O(1) | O(n) | O(1)* append, O(n) middle | O(n) middle, O(1) end | *amortized; see [01-complexity](../01-complexity/README.md) |
| Sorted slice | O(1) | O(log n) `slices.BinarySearch` | O(n) | O(n) | |
| Map `map[K]V` | — | O(1) avg | O(1) avg | O(1) avg | random iteration order |
| Set `map[K]struct{}` | — | O(1) avg | O(1) avg | O(1) avg | |
| Stack (slice) | O(1) top | O(n) | O(1) push | O(1) pop | |
| Queue / deque (slice or ring) | O(1) ends | O(n) | O(1) | O(1) | `q = q[1:]` leaks until GC; fine for interviews |
| Linked list (`container/list`) | O(n) | O(n) | O(1) at node | O(1) at node | |
| Binary heap (`container/heap`) | O(1) min | O(n) | O(log n) | O(log n) pop | `heap.Init` is O(n) |
| Balanced BST | O(log n) | O(log n) | O(log n) | O(log n) | not in stdlib; use sorted slice or heap |
| Trie | — | O(L) | O(L) | O(L) | L = key length |
| Union-Find | — | ~O(α(n)) | ~O(α(n)) | — | with path compression + rank |

### Algorithms

| Algorithm | Time | Space |
|---|---|---|
| Sort (`slices.Sort`, pdqsort) | O(n log n) | O(log n) |
| Binary search | O(log n) | O(1) |
| BFS / DFS on graph | O(V + E) | O(V) |
| Dijkstra (binary heap) | O((V + E) log V) | O(V) |
| Topological sort | O(V + E) | O(V) |
| Kruskal / Prim MST | O(E log E) | O(V) |
| Subsets / permutations | O(2ⁿ · n) / O(n! · n) | O(n) stack |
| 1-D DP over n states | O(n · transitions) | O(n), often O(1) |
| 2-D DP (two strings) | O(n · m) | O(n · m), often O(min(n, m)) |

## Pattern recognition

| Signal in the problem | Technique | Module |
|---|---|---|
| "Can this be faster?", constraints table, comparing approaches | Complexity analysis | [01-complexity](../01-complexity/README.md) |
| "Find a pair / count occurrences / seen before / group by" | Hash map / set, counting | [02-arrays-hashing](../02-arrays-hashing/README.md) |
| Sorted array, pair/triplet with target, in-place, palindrome | Two pointers | [03-two-pointers](../03-two-pointers/README.md) |
| "Longest/shortest **contiguous** subarray/substring such that…" | Sliding window | [04-sliding-window](../04-sliding-window/README.md) |
| "Sum of a range", "number of subarrays summing to k", many range queries | Prefix sums (+ hash map) | [05-prefix-sums](../05-prefix-sums/README.md) |
| "Next greater/smaller", matching brackets, evaluate expression | Stack / monotonic stack | [06-stacks-monotonic](../06-stacks-monotonic/README.md) |
| Sorted input, "minimum X such that…", O(log n) required | Binary search (on index or answer) | [07-binary-search](../07-binary-search/README.md) |
| Reverse/merge/detect cycle in a list, O(1) extra space | Linked-list pointers, fast/slow | [08-linked-lists](../08-linked-lists/README.md) |
| Hierarchy, BST, "path from root", level order | Tree DFS / BFS | [09-trees](../09-trees/README.md) |
| "Top k", "k-th largest", merge k sorted, running median | Heap | [10-heaps](../10-heaps/README.md) |
| Prefix matching, autocomplete, word search on many words | Trie | [11-tries](../11-tries/README.md) |
| Grid, islands, dependencies, "connected", "can reach" | Graph BFS/DFS, topo sort, union-find | [12-graphs](../12-graphs/README.md) |
| Weighted edges, cheapest path, connect all nodes cheaply | Dijkstra, Bellman-Ford, MST | [13-shortest-paths-mst](../13-shortest-paths-mst/README.md) |
| "All combinations / permutations / subsets", n ≤ 20 | Backtracking | [14-backtracking](../14-backtracking/README.md) |
| "Number of ways", "min cost", choices along a line | 1-D dynamic programming | [15-dp-1d](../15-dp-1d/README.md) |
| Two strings, grid paths, knapsack | 2-D dynamic programming | [16-dp-2d](../16-dp-2d/README.md) |
| Intervals, scheduling, "locally best choice works" | Greedy, sort by start/end | [17-greedy-intervals](../17-greedy-intervals/README.md) |
| "Single number", powers of two, subsets as masks, no extra space | Bit manipulation | [18-bit-manipulation](../18-bit-manipulation/README.md) |
| Pattern matching, palindromic substrings, string hashing | String algorithms | [19-strings](../19-strings/README.md) |
| "Design a cache / rate limiter / data structure with O(1) ops" | Composite data structures | [20-design-problems](../20-design-problems/README.md) |

## How to recognize the pattern in an interview

The table above lists the clues. This section is about using them under pressure: the first two or three minutes of an interview, when you have a problem statement and no idea yet which trick it needs.

Think like a detective. Don't guess the answer from the title. Collect clues from the statement, then let them point to a technique.

### Ask four questions, in order

```mermaid
flowchart LR
    A["1. What is the input?<br/>array, string, list, tree,<br/>grid, intervals, stream"] --> B["2. What is asked?<br/>one answer, count,<br/>min/max, all options"]
    B --> C["3. What do the<br/>constraints allow?<br/>n tells you the Big-O"]
    C --> D["4. Where is the brute<br/>force slow?<br/>that step names the fix"]
```

1. **What is the input?** The shape of the data narrows the choice to a few modules.
2. **What is asked?** "Longest stretch", "number of ways", "all subsets" and "k-th largest" each point somewhere different.
3. **What do the constraints allow?** The size of `n` tells you how fast the solution must be, and so which techniques are in or out.
4. **Where is the brute force slow?** Say the simple solution out loud, then find the line that repeats work. The tool that removes that repeated work is usually the pattern.

### Visual 1: start from the input

```mermaid
flowchart TD
    S["Read the problem"] --> Q{"What is the input?"}
    Q -->|"array or string"| A["Use Visual 2"]
    Q -->|"linked list"| L["08 Linked lists<br/>dummy head, fast and slow pointers"]
    Q -->|"tree"| T["09 Trees<br/>DFS or level-order BFS"]
    Q -->|"grid, network, or tasks<br/>that depend on each other"| G{"Do edges have<br/>different costs?"}
    G -->|"no"| G1["12 Graphs<br/>BFS, DFS, topological sort, union-find"]
    G -->|"yes"| G2["13 Shortest paths and MST<br/>Dijkstra, Bellman-Ford, Kruskal, Prim"]
    Q -->|"intervals with start and end"| I["17 Greedy and intervals<br/>sort by start or by end"]
    Q -->|"a stream of values"| H["10 Heaps<br/>top k, running median"]
    Q -->|"a dictionary of words"| W["11 Tries<br/>prefixes, autocomplete"]
    Q -->|"design a class with fast operations"| D["20 Design problems"]
```

Grids, prerequisites and "who knows whom" are graphs even when the word "graph" never appears.

### Visual 2: arrays and strings

Most interview problems take an array or a string, so this is the branch you'll use most.

```mermaid
flowchart TD
    A{"What does it ask for?"} -->|"a pair or group<br/>that matches a target"| P{"Sorted, or allowed<br/>to sort?"}
    P -->|"yes"| P1["03 Two pointers"]
    P -->|"no"| P2["02 Hash map<br/>have I seen the partner?"]
    A -->|"longest or shortest<br/>contiguous stretch"| W{"All values non-negative,<br/>or a string?"}
    W -->|"yes"| W1["04 Sliding window"]
    W -->|"no, negatives allowed"| W2["05 Prefix sums + hash map"]
    A -->|"sums over many ranges"| R["05 Prefix sums"]
    A -->|"next greater or smaller,<br/>matching brackets"| S["06 Stack or monotonic stack"]
    A -->|"find in sorted data, or<br/>the smallest X that works"| B["07 Binary search"]
    A -->|"k-th largest or top k"| H["10 Heap or quickselect"]
    A -->|"a pattern inside text,<br/>palindromes"| T["19 String algorithms"]
    A -->|"every value twice but one,<br/>no extra space"| X["18 Bit manipulation"]
```

### Visual 3: what kind of answer?

When the input doesn't settle it, look at the answer they want.

```mermaid
flowchart TD
    A{"What kind of answer?"} -->|"list every combination,<br/>subset or order"| B["14 Backtracking"]
    A -->|"count the ways, or<br/>the min or max cost"| C{"Is the best move now<br/>always safe?"}
    C -->|"yes, and you can prove it"| G["17 Greedy"]
    C -->|"no, choices affect later ones"| D{"How many numbers<br/>describe a smaller problem?"}
    D -->|"one: a position or amount"| D1["15 DP in 1-D"]
    D -->|"two: two strings, a grid cell,<br/>or a range i to j"| D2["16 DP in 2-D"]
    A -->|"fewest steps or moves"| E{"Does every step<br/>cost the same?"}
    E -->|"yes"| E1["12 BFS"]
    E -->|"no"| E2["13 Dijkstra"]
```

### Let the constraints narrow it down

The [constraints table](#the-framework) shows how much work fits in time. Read it the other way round, too: the limit on `n` rules techniques in and out.

| If n is up to | You can afford | Which suggests |
|---|---|---|
| 10–12 | O(n!) | trying every order: [14 backtracking](../14-backtracking/README.md) |
| 20–25 | O(2ⁿ) | trying every subset: [14 backtracking](../14-backtracking/README.md), bitmasks from [18](../18-bit-manipulation/README.md) |
| about 500 | O(n³) | interval DP in [16](../16-dp-2d/README.md), Floyd-Warshall in [13](../13-shortest-paths-mst/README.md) |
| about 5,000 | O(n²) | DP over pairs, or a careful brute force |
| up to 10⁶ | O(n log n) | sorting, heaps, binary search |
| beyond 10⁶ | O(n) or O(log n) | one pass: hash maps, two pointers, windows, prefix sums, stacks |

If `n` is small, that's a hint. Interviewers don't write "n ≤ 20" by accident.

### Find the slow step in the brute force

Say the simple solution first. Then point at the line that repeats work and swap in the tool that removes it.

| The brute force keeps… | Replace it with | Module |
|---|---|---|
| searching the whole array for a partner value | a hash map lookup | [02](../02-arrays-hashing/README.md) |
| checking every pair in sorted data | two pointers moving inward | [03](../03-two-pointers/README.md) |
| recounting a stretch every time it moves by one | a sliding window update | [04](../04-sliding-window/README.md) |
| adding up the same range again and again | prefix sums | [05](../05-prefix-sums/README.md) |
| scanning right for the next bigger value | a monotonic stack | [06](../06-stacks-monotonic/README.md) |
| trying every possible answer from 1 upward | binary search on the answer | [07](../07-binary-search/README.md) |
| sorting everything just to read the top k | a heap of size k | [10](../10-heaps/README.md) |
| comparing each word letter by letter with many others | a trie | [11](../11-tries/README.md) |
| solving the same smaller problem many times | dynamic programming | [15](../15-dp-1d/README.md), [16](../16-dp-2d/README.md) |
| restarting a text match from scratch after a mismatch | KMP or a rolling hash | [19](../19-strings/README.md) |

### Four worked examples

Each example puts the clue words in bold, then follows them to a technique.

**[Longest Substring Without Repeating Characters](../04-sliding-window/p04_longest_substring_no_repeat.go)**

> Given a **string** s, find the length of the **longest substring** **without repeating** characters.

```
string  →  "longest"  →  substring means contiguous  →  a rule to keep  →  SLIDING WINDOW (04)
```

The brute force checks every substring. The window keeps the stretch that obeys the rule and only moves its edges.

**[Koko Eating Bananas](../07-binary-search/p07_koko_bananas.go)**

> Return the **minimum** eating speed k **such that** Koko can eat all the bananas **within h hours**.

```
"minimum k such that"  →  if speed k works, every faster speed works too  →  BINARY SEARCH ON THE ANSWER (07)
```

Whenever the answers look like "no, no, no, yes, yes, yes", you can binary search for the first "yes".

**[Course Schedule](../12-graphs/p04_course_schedule.go)**

> There are n courses. **prerequisites**[i] = [a, b] means you must take b **before** a. **Can you finish** them all?

```
"before" means an arrow b → a  →  a directed graph  →  stuck only if there's a cycle  →  TOPOLOGICAL SORT (12)
```

The word "graph" never appears. "Prerequisites", "dependencies" and "build order" are the disguise.

**[Coin Change](../15-dp-1d/p05_coin_change.go)**

> Given coin values and an amount, return the **fewest coins** that make up the amount.

```
"fewest"  →  min cost  →  is taking the biggest coin always safe?  →  no, coins 1, 3, 4 for 6  →  DP (15)
```

A min or max question tempts you towards greedy. Test it on a small example first. When one counterexample breaks it, switch to DP.

### Look-alikes that trip people up

| It looks like… | But check… | If so, use… |
|---|---|---|
| a sliding window | Can numbers be negative? Then shrinking the window may not lower the sum. | prefix sums + hash map ([05](../05-prefix-sums/README.md)) |
| greedy | Can you find one small example where the greedy choice loses? | DP ([15](../15-dp-1d/README.md)) |
| BFS for the shortest path | Do the edges cost different amounts? | Dijkstra ([13](../13-shortest-paths-mst/README.md)) |
| backtracking | Do they want a *count* rather than the full *list*? | often DP, because a count doesn't need every answer listed |
| sorting | Is the data arriving as a stream, or do they only need the top k? | a heap ([10](../10-heaps/README.md)) |
| binary search on an index | Is the data unsorted, but the answer has a yes/no boundary? | binary search on the answer ([07](../07-binary-search/README.md)) |
| a hash map | Do they ask for O(1) extra space? | two pointers ([03](../03-two-pointers/README.md)) or bit tricks ([18](../18-bit-manipulation/README.md)) |

### Say your reasoning out loud

Interviewers want to hear how you found the pattern, not only the name of it. A useful script:

> "The input is an unsorted array and I need the longest stretch whose sum is k. Brute force checks every subarray, which is O(n²). The values can be negative, so a sliding window won't work. I'll use prefix sums with a hash map, which is O(n) time and O(n) space."

That's one sentence for each of the four questions, followed by the complexity.

### When nothing fits

1. Solve two tiny examples by hand and watch what you do. You'll often find yourself using the pattern without noticing.
2. Write the brute force and look for the repeated work.
3. Ask whether sorting the input would help. Many patterns only appear after a sort.
4. Ask whether a smaller version of the problem would help you solve the bigger one. That's recursion, which leads to DP or backtracking.
5. Tell the interviewer which ideas you're weighing. They'll often nudge you.

### Practise recognizing on its own

Spotting the pattern is a separate skill from coding it, so practise it separately. Pick problems at random from the module lists, read only the statement, and give yourself 60 seconds to name the technique and say why. Then check the problem's *Key idea* column. Ten a day for a week makes a big difference.

## Go for coding interviews

Every idiom below is implemented and tested in [`gotchas.go`](gotchas.go).

```go
// Sorting
sort.Slice(xs, func(i, j int) bool { return xs[i].Count > xs[j].Count }) // not stable!
slices.SortFunc(ws, func(a, b string) int {
    return cmp.Or(cmp.Compare(len(a), len(b)), strings.Compare(a, b))
})

// Set
seen := map[int]struct{}{}
seen[x] = struct{}{}
if _, ok := seen[x]; ok { /* ... */ }

// Counting lowercase letters: arrays are comparable, maps are not
var c [26]int
c[s[i]-'a']++

// Building strings
var b strings.Builder
b.WriteString("x"); b.WriteByte(','); b.WriteRune('é')
s := b.String()

// Safe copies
cp := slices.Clone(cur)          // before appending `cur` to results
grid2 := make([][]int, len(grid)) // 2-D copy needs a loop over rows

// Numbers
best := math.MaxInt
best = min(best, x)              // builtins since Go 1.21
mid := lo + (hi-lo)/2

// Iterate code points, not bytes
for i, r := range s { _ = i; _ = r } // i = byte offset, r = rune

// Heap: implement Len/Less/Swap/Push/Pop, then use heap.Push / heap.Pop
h := &MinHeap{}
heap.Push(h, 3)
x := heap.Pop(h).(int)
```

## Common pitfalls

- **Map iteration order is random.** Never rely on it for output; sort keys first ([`SortedKeys`](gotchas.go)).
- **Sub-slices share memory.** `append(a[:2], x)` can overwrite `a[2]`. Use `slices.Clone` or `slices.Clip` ([`AliasingDemo`](gotchas.go)).
- **Backtracking results mutate.** `res = append(res, path)` stores a view; append `slices.Clone(path)`.
- **`s[i]` is a byte.** Use `[]rune(s)` or `range s` for Unicode; `len(s)` counts bytes.
- **Integer overflow is silent.** `int` is 64-bit on interview machines, but `(lo+hi)/2`, products and `-math.MinInt` still overflow.
- **`sort.Slice` is unstable.** Add a tie-breaker or use `sort.SliceStable` / `slices.SortStableFunc`.
- **Loop variables.** Since Go 1.22 each iteration gets a fresh variable, so closures capturing `i` are safe — but say so if your interviewer uses older Go.
- **`heap.Pop` vs `h.Pop`.** Always call the package functions; your methods are the plumbing.

## Problems

This module has no numbered problems — it is the reference you come back to. The
exercises are the idioms in [`gotchas.go`](gotchas.go):

| # | Idiom (link to file) | Key idea |
|---|---|---|
| 1 | [Sorting with comparators](gotchas.go) | `sort.Slice`, `slices.SortFunc`, `cmp.Or` |
| 2 | [Min-heap with `container/heap`](gotchas.go) | five methods; keep size-k heap for top-k |
| 3 | [Strings and runes](gotchas.go) | `strings.Builder`, `[]rune`, `[26]int` counts |
| 4 | [Map as set](gotchas.go) | `map[T]struct{}` |
| 5 | [Slice aliasing](gotchas.go) | `slices.Clone`, `slices.Clip` |
| 6 | [Numbers](gotchas.go) | `math.MaxInt`, `min`/`max`, overflow-safe midpoint |

## Practice tip

Before each mock interview, re-type `MinHeap` and `KLargest` from memory, then run
`go test ./course/00-interview-playbook/ -run TestKLargest`. If it takes more than
three minutes, do it again tomorrow. Read the framework table aloud once — it
makes step 1 (clarify) automatic.

## Next

[01 · Complexity →](../01-complexity/README.md)
