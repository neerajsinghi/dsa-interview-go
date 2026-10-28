// Package tries implements the trie (prefix tree) and the interview problems
// built on top of it: prefix queries, wildcard search, autocomplete, grid word
// search and the bitwise trie for XOR maximisation.
package tries

// alphabet is the number of distinct letters this trie supports ('a'..'z').
// A fixed-size array of children is faster and allocation-friendlier than a
// map; switch to map[rune]*Node when the alphabet is large or unknown (see
// MapTrie in p01_implement_trie.go).
const alphabet = 26

// Node is one trie node. The path of edges from the root to a node spells a
// prefix; the node itself stores no character.
type Node struct {
	Children [alphabet]*Node
	IsEnd    bool // a stored word ends exactly here
	Pass     int  // how many stored words pass through (or end at) this node
}

// Trie is a prefix tree over lowercase ASCII words. Every word passed to it
// must consist only of 'a'..'z'; other bytes panic with an index error.
// The zero value is not usable: create one with New.
type Trie struct {
	root *Node
}

// New returns an empty trie.
func New() *Trie {
	return &Trie{root: &Node{}}
}

// Insert adds word to the trie. Inserting a word that is already present is
// a no-op, which keeps the Pass counters equal to the number of distinct
// words below each node.
//
// Time: O(L) for a word of length L.
func (t *Trie) Insert(word string) {
	if t.Search(word) {
		return
	}
	cur := t.root
	cur.Pass++
	for i := 0; i < len(word); i++ {
		c := word[i] - 'a'
		if cur.Children[c] == nil {
			cur.Children[c] = &Node{}
		}
		cur = cur.Children[c]
		cur.Pass++
	}
	cur.IsEnd = true
}

// Search reports whether word was inserted (an exact match, not a prefix).
func (t *Trie) Search(word string) bool {
	n := t.walk(word)
	return n != nil && n.IsEnd
}

// StartsWith reports whether any stored word begins with prefix.
// Every trie "starts with" the empty prefix, even an empty trie, matching
// the convention of strings.HasPrefix(s, "").
func (t *Trie) StartsWith(prefix string) bool {
	return t.walk(prefix) != nil
}

// CountPrefix returns how many stored words begin with prefix.
func (t *Trie) CountPrefix(prefix string) int {
	n := t.walk(prefix)
	if n == nil {
		return 0
	}
	return n.Pass
}

// Len returns the number of distinct words stored.
func (t *Trie) Len() int {
	return t.root.Pass
}

// Delete removes word and reports whether it was present. Nodes that no
// longer lie on any stored word are unlinked so the memory can be reclaimed.
//
// Time: O(L).
func (t *Trie) Delete(word string) bool {
	if !t.Search(word) {
		return false
	}
	cur := t.root
	cur.Pass--
	for i := 0; i < len(word); i++ {
		c := word[i] - 'a'
		next := cur.Children[c]
		next.Pass--
		if next.Pass == 0 {
			// Only this word used the rest of the path: drop the whole branch.
			cur.Children[c] = nil
			return true
		}
		cur = next
	}
	cur.IsEnd = false
	return true
}

// WithPrefix returns up to limit stored words that start with prefix, in
// lexicographic order (limit <= 0 means no limit). A pre-order DFS that
// visits children 'a'..'z' yields words in sorted order for free.
func (t *Trie) WithPrefix(prefix string, limit int) []string {
	start := t.walk(prefix)
	if start == nil {
		return nil
	}
	var out []string
	buf := []byte(prefix)
	var dfs func(n *Node)
	dfs = func(n *Node) {
		if limit > 0 && len(out) == limit {
			return
		}
		if n.IsEnd {
			out = append(out, string(buf))
		}
		for c, child := range n.Children {
			if child == nil {
				continue
			}
			buf = append(buf, byte('a'+c)) // choose
			dfs(child)                     // explore
			buf = buf[:len(buf)-1]         // unchoose
		}
	}
	dfs(start)
	return out
}

// walk follows s from the root and returns the node it ends at, or nil if
// the path falls off the trie.
func (t *Trie) walk(s string) *Node {
	cur := t.root
	for i := 0; i < len(s); i++ {
		cur = cur.Children[s[i]-'a']
		if cur == nil {
			return nil
		}
	}
	return cur
}
