import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from trieview import TNode, TrieView  # noqa: E402
from walkthrough import KV, Frame, Text  # noqa: E402

TITLE = "Add and Search Words"
PATTERN = "Trie + DFS"


class WordDictionary:
    def __init__(self):
        self.root = TNode()

    def add_word(self, w):
        cur = self.root
        for ch in w:
            cur = cur.children.setdefault(ch, TNode())
        cur.end = True

    def search(self, pattern, trace):
        """Port of Search; trace collects (path, result) for every node the DFS enters."""

        def match(n, path, i):
            if i == len(pattern):
                trace.append((path, n.end))
                return n.end
            if pattern[i] != ".":
                nxt = n.children.get(pattern[i])
                if nxt is None:
                    trace.append((path + pattern[i], None))  # fell off
                    return False
                ok = match(nxt, path + pattern[i], i + 1)
                trace.append((path, ok))
                return ok
            for ch, child in n.kids():
                if match(child, path + ch, i + 1):
                    trace.append((path, True))
                    return True
            trace.append((path, False))
            return False

        return match(self.root, "", 0)


def frames():
    words = ["bad", "dad", "mad"]
    queries = [("pad", False), ("bad", True), (".ad", True), ("b..", True), ("b.", False)]
    d = WordDictionary()

    yield Frame(
        Text('AddWord("bad"), AddWord("dad"), AddWord("mad")',
             *[f'Search("{q}") -> {str(e).lower()}' for q, e in queries]),
        speaker="Interviewer",
        say="Store words, then answer patterns where a dot matches exactly one letter. "
            "The match must cover the whole word.",
    )
    yield Frame(
        Text("brute force: keep a []string",
             "Search compares the pattern with every word of the same length",
             ("O(N · L) per search", "reject")),
        say="The simple version keeps a list and compares the pattern against every word. "
            "That's O(N·L) per search. A trie lets a literal letter follow just one edge.",
    )
    for w in words:
        d.add_word(w)
    yield Frame(
        TrieView(d.root, label="trie"),
        KV(words, label="added"),
        say="AddWord is a normal trie insert. Three words give three branches from the root, "
            "and each last letter is marked as a word end.",
    )

    results = []
    for q, exp in queries:
        trace = []
        got = d.search(q, trace)
        assert got == exp, (q, got)
        results.append(got)
        hl, fell = {}, None
        for path, ok in trace:
            if ok is None:
                fell = path
            elif ok:
                hl[path] = "match"
            else:
                hl.setdefault(path, "reject")
        hl[""] = "match" if got else "reject"
        tried = [p for p, ok in trace if ok is not None and len(p) == len(q)]
        status = (f"Search(\"{q}\") -> {str(got).lower()}", "match" if got else "reject")
        if fell is not None:
            say = (f"\"{q}\" starts with '{q[0]}', but the root has no '{q[0]}' child. "
                   "The walk falls off right away, so the answer is false.")
            lines = [f"root has no child '{fell[-1]}'", status]
        elif "." in q and got:
            say = (f"A dot tries every child in turn. For \"{q}\" the first branch I try, "
                   f"\"{tried[0]}\", already ends on a word, so I stop and return true.")
            lines = [f"dot -> try children in order; first full path: {tried[0]}", status]
        elif got:
            say = "A literal letter follows exactly one edge: b, a, d. The walk ends on a word end, so the answer is true."
            lines = ["b -> a -> d, IsEnd = true", status]
        else:
            say = (f"\"{q}\" walks to node \"{trace[0][0]}\" and the pattern runs out there. "
                   "That node isn't a word end, so the whole-word match fails: false.")
            lines = [f"pattern ended on \"{trace[0][0]}\", IsEnd = false", status]
        yield Frame(
            TrieView(d.root, label="trie", hl=hl),
            Text(*lines),
            say=say,
        )

    yield Frame(
        TrieView(d.root, label="trie"),
        Text(("results: " + ", ".join(str(r).lower() for r in results), "match"),
             "AddWord: O(L)   Search: O(L) without dots",
             "worst case with d dots: O(26^d · L), capped by the trie size"),
        say="AddWord is O(L). Search is O(L) with no dots. Each dot can branch into up to 26 children, "
            "but a branch dies as soon as it falls off the trie.",
    )
