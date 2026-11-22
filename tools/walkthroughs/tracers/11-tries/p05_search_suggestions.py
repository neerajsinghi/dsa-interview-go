import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from trieview import TNode, TrieView  # noqa: E402
from walkthrough import KV, Frame, Text  # noqa: E402

TITLE = "Search Suggestions System"
PATTERN = "Trie autocomplete"


def insert(root, w):
    cur = root
    for ch in w:
        cur = cur.children.setdefault(ch, TNode())
    cur.end = True


def with_prefix(root, prefix, limit, visited):
    """Port of Trie.WithPrefix: pre-order DFS, children a..z, stop at limit."""
    cur = root
    for ch in prefix:
        cur = cur.children.get(ch)
        if cur is None:
            return []
    out = []
    buf = list(prefix)

    def dfs(n):
        if limit > 0 and len(out) == limit:
            return
        visited.append("".join(buf))
        if n.end:
            out.append("".join(buf))
        for ch, child in n.kids():
            buf.append(ch)
            dfs(child)
            buf.pop()

    dfs(cur)
    return out


def frames():
    products = ["mobile", "mouse", "moneypot", "monitor", "mousepad"]
    search = "mouse"
    expected = [["mobile", "moneypot", "monitor"], ["mobile", "moneypot", "monitor"],
                ["mouse", "mousepad"], ["mouse", "mousepad"], ["mouse", "mousepad"]]

    yield Frame(
        KV(products, label="products"),
        Text(f'searchWord = "{search}"'),
        speaker="Interviewer",
        say="As the user types \"mouse\" one letter at a time, suggest up to three products "
            "with the typed prefix, smallest first. Return one list per letter.",
    )
    yield Frame(
        KV(products, label="products"),
        Text("brute force: for each typed prefix, filter all products and sort",
             ("O(L · n · log n) with string compares", "reject")),
        say="The direct way filters every product for each prefix and sorts the matches. "
            "It redoes the same work for every letter typed.",
    )

    root = TNode()
    for p in products:
        insert(root, p)
    yield Frame(
        TrieView(root, label="trie", level_h=31, r=11),
        KV(products, label="products"),
        say="I insert all products into a trie. A DFS that visits children from a to z "
            "produces words in sorted order, so the first three it finds are the answer.",
    )

    res = []
    cur = root
    for i in range(len(search)):
        prefix = search[: i + 1]
        if cur is not None:
            cur = cur.children.get(search[i])
        if cur is None:
            res.append([])
            continue
        visited = []
        got = with_prefix(root, prefix, 3, visited)
        res.append(got)
        hl = {prefix[:j]: "active" for j in range(len(prefix) + 1)}
        for v in visited:
            if len(v) > len(prefix):
                hl[v] = "visited"
        for w in got:
            hl[w] = "match"
        if i == 0:
            say = ("After \"m\" I step to node m and run the DFS from there. It meets mobile, "
                   "then moneypot, then monitor, and stops at three.")
        elif got == res[-2]:
            say = f"\"{prefix}\": one more edge down, then the same DFS. The answer matches the previous letter."
        else:
            say = (f"\"{prefix}\": now I'm in the u branch. The DFS finds mouse first, because it ends "
                   "higher up, then mousepad. Only two words exist here.")
        yield Frame(
            TrieView(root, label="trie", hl=hl, level_h=31, r=11),
            Text(f'typed "{prefix}" -> ' + str(got).replace("'", "")),
            say=say,
        )

    assert res == expected, res
    yield Frame(
        Text(*[(f'{search[:i + 1]:<6}-> ' + str(r).replace("'", ""), "match") for i, r in enumerate(res)]),
        Text("trie: O(P) to build, one edge + a short DFS per letter",
             "sort + binary search: O(P log n + L · log n)"),
        say="Building the trie is O(P) for P total characters. Each letter costs one step plus a DFS "
            "that stops after three words. Sorting plus binary search is the other common answer.",
    )
