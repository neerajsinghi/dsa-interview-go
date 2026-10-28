# DSA for Coding Interviews, in Go

[![CI](https://github.com/neerajsinghi/dsa-interview-go/actions/workflows/ci.yml/badge.svg)](https://github.com/neerajsinghi/dsa-interview-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/neerajsinghi/dsa-interview-go.svg)](https://pkg.go.dev/github.com/neerajsinghi/dsa-interview-go)

**A 21-module course on data structures and algorithms for coding interviews: a lesson per pattern, 150+ interview problems, and a tested Go solution for every one.**

Most interview prep teaches problems. This course teaches **patterns**: how to spot that a problem is "sliding window" or "binary search on the answer", the template that solves it, and the traps that cost you the round. Each module is a short lesson followed by problems ordered from easy to hard, so you practise the pattern until it sticks.

- **Every solution is tested.** Table-driven tests cover edge cases, and randomized tests check the optimal solution against a brute force, so you can trust the code and see *why* the fast version is right.
- **Brute force first, then optimal.** Where it helps, problems show both, mirroring how you should talk through a solution in an interview.
- **Written for Go interviews.** Idiomatic Go, plus the Go-specific traps interviewers notice: slice aliasing in backtracking, `rune` vs `byte`, integer overflow, `container/heap` boilerplate.
- **Standard library only.** `go test ./...` is all you need.

Companion to [system-design-patterns-go](https://github.com/neerajsinghi/system-design-patterns-go), which covers the system design round the same way.

## Quickstart

```bash
git clone https://github.com/neerajsinghi/dsa-interview-go
cd dsa-interview-go

go test ./...                                              # run every solution's tests
go test ./course/04-sliding-window/ -run TestMinWindow -v  # one problem
go test -bench=. ./course/01-complexity/                   # see Big-O in real numbers
```

Requires Go 1.23 or later.

## Course map

| # | Module | What you'll learn | Problems |
|---|---|---|---|
| 00 | [Interview playbook](course/00-interview-playbook) | A 7-step interview framework, Big-O and pattern-recognition cheat sheets, Go idioms for interviews | — |
| 01 | [Complexity](course/01-complexity) | Big-O, amortized cost, recursion trees, with benchmarks that make growth visible | 5 |
| 02 | [Arrays & hashing](course/02-arrays-hashing) | Hash maps and sets for O(1) lookups, counting, bucket sort | 8 |
| 03 | [Two pointers](course/03-two-pointers) | Converging and same-direction pointers on sorted data | 7 |
| 04 | [Sliding window](course/04-sliding-window) | Fixed and variable windows over arrays and strings, monotonic deque | 8 |
| 05 | [Prefix sums](course/05-prefix-sums) | Range sums in O(1), prefix sum + hash map, difference arrays, 2D prefix sums | 8 |
| 06 | [Stacks & monotonic stacks](course/06-stacks-monotonic) | Matching brackets, expression evaluation, next greater element | 8 |
| 07 | [Binary search](course/07-binary-search) | Lower/upper bound templates, rotated arrays, binary search on the answer | 8 |
| 08 | [Linked lists](course/08-linked-lists) | Reversal, fast/slow pointers, dummy heads, in-place reordering | 8 |
| 09 | [Trees](course/09-trees) | Traversals, recursion on trees, BSTs, construction and serialization | 9 |
| 10 | [Heaps](course/10-heaps) | Top-k, k-way merge, two heaps for a running median, quickselect | 7 |
| 11 | [Tries](course/11-tries) | Prefix trees for word search, autocomplete and bitwise XOR | 7 |
| 12 | [Graphs](course/12-graphs) | BFS/DFS, grids, topological sort, union-find | 9 |
| 13 | [Shortest paths & MST](course/13-shortest-paths-mst) | Dijkstra, Bellman-Ford, 0-1 BFS, Floyd-Warshall, Kruskal and Prim | 7 |
| 14 | [Backtracking](course/14-backtracking) | Subsets, permutations, combinations, constraint search (N-Queens, Sudoku) | 8 |
| 15 | [Dynamic programming 1D](course/15-dp-1d) | Framing DP: state, transition, base case; memo vs table vs O(1) space | 9 |
| 16 | [Dynamic programming 2D](course/16-dp-2d) | Grids, two-string DP (LCS, edit distance), knapsack, interval DP | 9 |
| 17 | [Greedy & intervals](course/17-greedy-intervals) | When greedy is provably right, Kadane, merging and scheduling intervals | 9 |
| 18 | [Bit manipulation](course/18-bit-manipulation) | Masks, two's complement, XOR tricks, `math/bits` | 9 |
| 19 | [String algorithms](course/19-strings) | KMP, Z-function, Rabin-Karp, Manacher, palindromes | 9 |
| 20 | [Design problems](course/20-design-problems) | LRU/LFU caches, rate limiter, iterators, a simplified Twitter feed | 9 |

## Study plan

Eight weeks at about an hour a day. Already comfortable with the basics? Skim 01–03 and spend the extra time on graphs and DP.

| Week | Modules | Goal |
|---|---|---|
| 1 | 00, 01, 02, 03 | The interview framework, Big-O, and the most common easy/medium patterns |
| 2 | 04, 05, 06 | Window and prefix techniques, stack-based problems |
| 3 | 07, 08 | Binary search in all its forms, pointer manipulation |
| 4 | 09, 10, 11 | Tree recursion, heaps, tries |
| 5 | 12, 13 | Graph traversal and shortest paths |
| 6 | 14, 15 | Backtracking and the DP framework |
| 7 | 16, 17 | Harder DP, greedy and intervals |
| 8 | 18, 19, 20 + review | Bits and strings, design rounds, then re-solve problems you got wrong |

## How to practise

1. **Read the lesson**, especially *Recognize the pattern*. Spotting the pattern is most of the battle.
2. **Solve before you look.** Read the problem's doc comment, then write your own solution in a scratch file with the same function signature, or delete the body of the provided one (`git checkout` brings it back).
3. **Run the tests** for that problem: `go test ./course/<module>/ -run <TestName> -v`. The randomized tests will catch the edge cases you missed.
4. **Compare** with the reference solution and its complexity analysis.
5. **Say it out loud.** Explain the brute force, the optimisation and the complexity as you would in an interview. The [playbook](course/00-interview-playbook) has a script.
6. **Revisit** the problems you got wrong after 3 days and again after a week.

## Layout

```
course/NN-topic/
├── README.md           the lesson: recognize the pattern, the idea, a Go template, pitfalls, problem list
├── <topic>.go          the data structure or reusable templates for the pattern
├── pNN_<name>.go       one problem: statement, examples, key insight, complexity, solution(s)
└── pNN_<name>_test.go  table-driven and randomized tests
```

Each module is its own Go package, named after the topic (for example `course/04-sliding-window` is package `slidingwindow`).

## A note on problem statements

Problems are classic interview questions, restated in our own words. Where a problem is well known, its doc comment links to the original so you can also submit there.

## Contributing

Issues and pull requests are welcome: a new problem, a clearer explanation, or a bug fix. Please follow the module layout, use only the standard library, and make sure `gofmt`, `go vet` and `go test -race ./...` pass.

## License

[MIT](LICENSE) © 2026 Neeraj Singhi
