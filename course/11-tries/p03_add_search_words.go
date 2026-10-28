package tries

// Problem: Design Add and Search Words Data Structure
// Difficulty: Medium
//
// Build a dictionary that can add words and answer whether a pattern matches
// any stored word. A pattern is made of lowercase letters and '.', where '.'
// stands for exactly one arbitrary letter. The match must cover the whole
// word, not just a prefix.
//
// Examples:
//
//	AddWord("bad"); AddWord("dad"); AddWord("mad")
//	Search("pad") -> false   Search("bad") -> true
//	Search(".ad") -> true    Search("b..") -> true    Search("b.") -> false
//
// Reference: https://leetcode.com/problems/design-add-and-search-words-data-structure/
//
// Pattern: trie + DFS. A literal letter follows exactly one edge; a '.'
// branches into every existing child. Branching only happens on wildcards and
// dies as soon as a path falls off the trie, so typical searches touch far
// fewer nodes than scanning the whole word list.
//
// Time: AddWord O(L). Search O(L) with no dots; worst case O(26^d · L) with
// d dots (bounded by the number of trie nodes). Space: O(total characters).

// WordDictionary stores words and answers wildcard queries.
type WordDictionary struct {
	root *Node // reuses the trie Node from trie.go
}

// NewWordDictionary returns an empty dictionary.
func NewWordDictionary() *WordDictionary {
	return &WordDictionary{root: &Node{}}
}

// AddWord stores a lowercase word.
func (d *WordDictionary) AddWord(word string) {
	cur := d.root
	for i := 0; i < len(word); i++ {
		c := word[i] - 'a'
		if cur.Children[c] == nil {
			cur.Children[c] = &Node{}
		}
		cur = cur.Children[c]
	}
	cur.IsEnd = true
}

// Search reports whether some stored word matches pattern.
func (d *WordDictionary) Search(pattern string) bool {
	var match func(n *Node, i int) bool
	match = func(n *Node, i int) bool {
		if i == len(pattern) {
			return n.IsEnd // whole pattern consumed: must sit on a word end
		}
		if pattern[i] != '.' {
			next := n.Children[pattern[i]-'a']
			return next != nil && match(next, i+1)
		}
		for _, child := range n.Children {
			if child != nil && match(child, i+1) {
				return true // first successful branch is enough
			}
		}
		return false
	}
	return match(d.root, 0)
}

// wordListDictionary is the brute force: keep a slice and test every word.
// O(N·L) per search. Used by the tests as the reference.
type wordListDictionary struct {
	words []string
}

func (d *wordListDictionary) AddWord(word string) { d.words = append(d.words, word) }

func (d *wordListDictionary) Search(pattern string) bool {
	for _, w := range d.words {
		if len(w) != len(pattern) {
			continue
		}
		ok := true
		for i := 0; i < len(w) && ok; i++ {
			ok = pattern[i] == '.' || pattern[i] == w[i]
		}
		if ok {
			return true
		}
	}
	return false
}
