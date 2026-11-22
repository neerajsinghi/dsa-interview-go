import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from trieview import TNode, TrieView  # noqa: E402
from walkthrough import Frame, KV, Text  # noqa: E402

TITLE = "Implement Trie"
PATTERN = "Trie"


class Trie:
    """Port of the array-backed Trie in trie.go (Insert / Search / StartsWith)."""

    def __init__(self):
        self.root = TNode()

    def walk(self, s):
        cur = self.root
        for ch in s:
            cur = cur.children.get(ch)
            if cur is None:
                return None
        return cur

    def search(self, w):
        n = self.walk(w)
        return n is not None and n.end

    def starts_with(self, p):
        return self.walk(p) is not None

    def insert(self, w):
        if self.search(w):
            return
        cur = self.root
        cur.passes += 1
        for ch in w:
            if ch not in cur.children:
                cur.children[ch] = TNode()
            cur = cur.children[ch]
            cur.passes += 1
        cur.end = True


def path_hl(word, last_style, upto=None):
    upto = len(word) if upto is None else upto
    hl = {word[:i]: "active" for i in range(upto)}
    hl[word[:upto]] = last_style
    return hl


def frames():
    ops = [("Insert", "apple", None), ("Search", "apple", True), ("Search", "app", False),
           ("StartsWith", "app", True), ("Insert", "app", None), ("Search", "app", True)]
    t = Trie()
    log = []

    def ops_text():
        return KV(log or ["(no calls yet)"], label="calls so far")

    yield Frame(
        Text(*[f"{op}(\"{w}\")" + (f"  -> {str(exp).lower()}" if exp is not None else "") for op, w, exp in ops]),
        speaker="Interviewer",
        say="Build a structure with Insert, Search for an exact word, and StartsWith for a prefix. "
            "Here's a sequence of calls and what each should return.",
    )
    yield Frame(
        Text("brute force: keep a []string",
             "Search / StartsWith scan every word -> O(N · L) per call",
             ("cost grows with N, the number of stored words", "reject")),
        say="I could keep a slice of words and scan it on each call. That's O(N·L) per query, "
            "so it slows down as the dictionary grows. I want a cost that depends only on the key length.",
    )

    for op, w, exp in ops:
        if op == "Insert":
            existed = t.walk(w) is not None
            t.insert(w)
            notes = {}
            n, p = t.root, ""
            for ch in w:
                n = n.children[ch]
                p += ch
                notes[p] = f"pass {n.passes}"
            log.append(f'Insert("{w}")')
            if not existed:
                say = (f"Insert walks one letter at a time and creates a node for each missing letter. "
                       f"The last node gets IsEnd, the double ring, meaning a word ends here.")
            else:
                say = (f"Inserting \"{w}\" creates nothing: every node is already there. "
                       f"I only set IsEnd on the last node. Pass counts words through each node.")
            yield Frame(
                TrieView(t.root, label="trie", hl=path_hl(w, "match"), notes=notes),
                Text(f'Insert("{w}")', ("IsEnd set on \"%s\"" % w, "match")),
                ops_text(),
                say=say,
            )
        elif op == "Search":
            got = t.search(w)
            assert got == exp, (op, w, got)
            log.append(f'Search("{w}")={str(got).lower()}')
            if got:
                say = f"Search walks {' → '.join(w)}. It lands on a node with IsEnd set, so \"{w}\" is a stored word: true."
                last = "match"
            else:
                say = (f"Search walks {' → '.join(w)} and lands on a real node, but IsEnd is false. "
                       f"\"{w}\" is only a prefix of apple, so the answer is false.")
                last = "reject"
            yield Frame(
                TrieView(t.root, label="trie", hl=path_hl(w, last)),
                Text(f'Search("{w}")', (f"node found, IsEnd = {str(got).lower()} -> {str(got).lower()}", last)),
                ops_text(),
                say=say,
            )
        else:
            got = t.starts_with(w)
            assert got == exp, (op, w, got)
            log.append(f'StartsWith("{w}")={str(got).lower()}')
            yield Frame(
                TrieView(t.root, label="trie", hl=path_hl(w, "match")),
                Text(f'StartsWith("{w}")', ("walk did not fall off -> true", "match")),
                ops_text(),
                say="StartsWith uses the same walk but skips the IsEnd check. "
                    "If the walk never falls off the tree, some word has this prefix: true.",
            )

    yield Frame(
        TrieView(t.root, label="trie"),
        Text(("results: true, false, true, true", "match"),
             "Insert / Search / StartsWith: O(L) time",
             "space: O(total characters inserted)"),
        say="Every call touches at most L nodes, one per letter, no matter how many words are stored. "
            "So each operation is O(L), and memory is one node per distinct prefix.",
    )
