# 11 · Tries (prefix trees)

## What you'll learn

- How a trie stores a set of strings so that every **prefix** is a node.
- Insert / Search / StartsWith / Delete in `O(L)`, independent of how many words are stored.
- Three trie superpowers interviewers love: **wildcard search** (trie + DFS), **autocomplete** (sorted DFS from a prefix node), and **many-words-at-once grid search** (trie + backtracking).
- The **bitwise trie**: the same idea over bits, used to maximise XOR greedily.

## Recognize the pattern

Reach for a trie when you see:

- "starts with", "prefix", "autocomplete", "suggestions", "shortest root".
- Many words searched against the same text or grid (Word Search II) — share the prefix work.
- Wildcard patterns over a fixed dictionary (`.` matches any letter).
- "Maximum XOR of a pair" — a trie over the bits of the numbers.

If you only need exact lookups, a `map[string]bool` is simpler and just as fast. The trie earns its keep when the question is about **prefixes**.

## The idea

Each edge is a character; each node is the prefix spelled by the path from the root. A flag marks nodes where a stored word ends.

```
Insert: car, card, care, cat, do

            (root)
           /      \
          c        d
          |        |
          a        o*
        /   \
       r*    t*
      / \
     d*  e*          * = IsEnd (a word ends here)
```

- `Search("car")` walks c→a→r and finds `IsEnd` → true.
- `Search("ca")` walks c→a, lands on a node, but no word ends there → false.
- `StartsWith("ca")` only needs the walk to succeed → true.
- `Delete("card")` decrements a per-node `Pass` counter along the path; when a child's counter hits 0 the whole branch is unlinked.

Because a DFS that visits children in `a..z` order emits words in **sorted order**, "three smallest words with this prefix" is just a DFS that stops after three.

## Go template

```go
type Node struct {
	Children [26]*Node // fixed alphabet: fast, no map allocations
	IsEnd    bool
}

func (t *Trie) Insert(w string) {
	cur := t.root
	for i := 0; i < len(w); i++ {
		c := w[i] - 'a'
		if cur.Children[c] == nil {
			cur.Children[c] = &Node{}
		}
		cur = cur.Children[c]
	}
	cur.IsEnd = true
}

func (t *Trie) walk(s string) *Node { // shared by Search and StartsWith
	cur := t.root
	for i := 0; i < len(s) && cur != nil; i++ {
		cur = cur.Children[s[i]-'a']
	}
	return cur
}
```

The full version with `Delete`, `CountPrefix` and `WithPrefix` is in [trie.go](trie.go).

## Complexity

| Operation | Time | Notes |
|---|---|---|
| Insert / Search / StartsWith / Delete | `O(L)` | `L` = key length |
| Build from `n` words | `O(S)` | `S` = total characters |
| Wildcard search with `d` dots | up to `O(26^d · L)` | bounded by number of nodes |
| Memory | `O(S · alphabet)` with arrays, `O(S)` with maps | |

## Common pitfalls

- **Prefix vs word.** Landing on a node is not enough for `Search`; check `IsEnd`. Forgetting this is the #1 trie bug.
- **Byte vs rune.** `w[i]` indexes bytes. `w[i]-'a'` is only safe for lowercase ASCII; for Unicode use `for _, r := range w` and `map[rune]*Node` (see `MapTrie`).
- **`[26]*Node` vs `map`.** The array costs 208 bytes per node on 64-bit but needs no hashing; a map is smaller for sparse alphabets but allocates. Pick consciously and say why.
- **Duplicate inserts.** If you keep counters, inserting the same word twice must not double-count — check `Search` first (as `Trie.Insert` does).
- **Delete without pruning** leaves dead branches that make `StartsWith` lie. Prune nodes that no stored word passes through.
- **Grid DFS + trie:** restore the cell after exploring (`board[r][c] = ch`) and clear `node.word` once found, or you'll report the same word many times.
- **Shared buffer in DFS.** When building the current word in a `[]byte`, `string(buf)` copies — safe. Storing `buf` itself (or a sub-slice) would alias memory that you keep mutating.
- **Bitwise trie:** iterate bits from the most significant down, and use a fixed width (31 bits here) so every number has a path of the same length.

## Problems

| # | Problem | Difficulty | Key idea |
|---|---|---|---|
| 01 | [Implement Trie](p01_implement_trie.go) | Medium | Array-backed `Trie` + Unicode `MapTrie`; `IsEnd` separates word from prefix |
| 02 | [Longest Common Prefix](p02_longest_common_prefix.go) | Easy | Vertical scan, or walk the trie's single-child trunk |
| 03 | [Add and Search Words](p03_add_search_words.go) | Medium | `.` branches into every child (trie + DFS) |
| 04 | [Replace Words](p04_replace_words.go) | Medium | First `IsEnd` on the walk = shortest root |
| 05 | [Search Suggestions System](p05_search_suggestions.go) | Medium | Sorted DFS from the prefix node, stop at 3; or sort + binary search |
| 06 | [Maximum XOR of Two Numbers](p06_maximum_xor.go) | Medium | Bitwise trie, greedily take the opposite bit |
| 07 | [Word Search II](p07_word_search_ii.go) | Hard | One grid DFS for all words, walking the trie; prune found leaves |

## Practice tip

Close the solution file, keep only the signature, and re-solve from scratch. Then run just that problem's tests:

```bash
go test ./course/11-tries/ -run TestFindWords -v
```

Each test file also has a randomized test that compares the fast solution against a brute force — if yours disagrees, the failure message prints the exact input.

## Next

[12 · Graphs](../12-graphs/README.md) — BFS, DFS, topological sort and union-find.
