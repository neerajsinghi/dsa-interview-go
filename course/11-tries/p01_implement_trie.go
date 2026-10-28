package tries

// Problem: Implement Trie (Prefix Tree)
// Difficulty: Medium
//
// Build a container of words that supports three operations: add a word,
// ask whether an exact word was added, and ask whether any added word begins
// with a given prefix. (This course's Trie in trie.go also supports Delete,
// CountPrefix and WithPrefix.)
//
// Examples:
//
//	Insert("apple"); Search("apple") -> true;  Search("app") -> false
//	StartsWith("app") -> true;  Insert("app"); Search("app") -> true
//
// Reference: https://leetcode.com/problems/implement-trie-prefix-tree/
//
// Pattern: one node per prefix. Walking a string from the root either falls
// off the tree (no word has this prefix) or lands on a node; an IsEnd flag on
// that node separates "is a word" from "is only a prefix".
//
// Time: O(L) per operation for a key of length L, independent of how many
// words are stored. Space: O(total characters inserted) nodes.
//
// The array-backed Trie lives in trie.go. MapTrie below is the variant you
// reach for when keys are arbitrary Unicode: children live in a map, so each
// node only pays for the edges it actually has.

// mapNode is a trie node whose children are keyed by rune.
type mapNode struct {
	children map[rune]*mapNode
	isEnd    bool
}

// MapTrie is a trie over arbitrary Unicode strings.
type MapTrie struct {
	root *mapNode
}

// NewMapTrie returns an empty map-backed trie.
func NewMapTrie() *MapTrie {
	return &MapTrie{root: &mapNode{children: map[rune]*mapNode{}}}
}

// Insert adds word.
func (t *MapTrie) Insert(word string) {
	cur := t.root
	for _, r := range word { // range over a string yields runes, not bytes
		next, ok := cur.children[r]
		if !ok {
			next = &mapNode{children: map[rune]*mapNode{}}
			cur.children[r] = next
		}
		cur = next
	}
	cur.isEnd = true
}

// Search reports whether word was inserted.
func (t *MapTrie) Search(word string) bool {
	n := t.walk(word)
	return n != nil && n.isEnd
}

// StartsWith reports whether some inserted word starts with prefix.
func (t *MapTrie) StartsWith(prefix string) bool {
	return t.walk(prefix) != nil
}

func (t *MapTrie) walk(s string) *mapNode {
	cur := t.root
	for _, r := range s {
		cur = cur.children[r] // a missing key yields nil: no ok-check needed
		if cur == nil {
			return nil
		}
	}
	return cur
}
