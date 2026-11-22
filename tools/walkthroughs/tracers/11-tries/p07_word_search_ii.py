import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from trieview import TNode, TrieView  # noqa: E402
from walkthrough import Frame, Grid, KV, Row, Text  # noqa: E402

TITLE = "Word Search II"
PATTERN = "Trie + grid DFS"

DIRS = [(1, 0), (-1, 0), (0, 1), (0, -1)]


def frames():
    rows = ["oaan", "etae", "ihkr", "iflv"]
    words = ["oath", "pea", "eat", "rain"]
    expected = ["eat", "oath"]
    board = [list(r) for r in rows]
    shown = [list(r) for r in rows]  # letters for drawing (board gets '#' marks)

    def grid(path, last=None):
        hl = {cell: "active" for cell in path}
        if path and last:
            hl[path[-1]] = last
        return Grid(shown, label="board", hl=hl, cell=40)

    def tview(root, p, last="active"):
        hl = {p[:i]: "active" for i in range(len(p))}
        hl[p] = last
        return TrieView(root, label="trie of words", hl=hl, level_h=36, r=12)

    yield Frame(
        Row(Grid(shown, label="board", cell=40), KV(words, label="words"), weights=[1, 1]),
        speaker="Interviewer",
        say="Return every word from the list that can be traced through adjacent cells, "
            "using each cell at most once per word. Here the answer is eat and oath.",
    )
    yield Frame(
        Row(Grid(shown, label="board", cell=40), KV(words, label="words"), weights=[1, 1]),
        Text("brute force: run a full Word Search for every word",
             ("O(words · m · n · 4 · 3^(L-1)): shared prefixes are explored again and again", "reject")),
        say="The brute force runs a separate grid search per word. Words that share a prefix "
            "repeat the same exploration, and the cost multiplies by the number of words.",
    )

    root = TNode()
    for w in words:
        cur = root
        for ch in w:
            cur = cur.children.setdefault(ch, TNode())
        cur.word = w
        cur.end = True
    yield Frame(
        Row(Grid(shown, label="board", cell=40), tview(root, "", None), weights=[1, 1.3]),
        say="I put all words in a trie and store each word at its end node. Then one DFS from each cell "
            "walks the grid and the trie together.",
    )

    found = []
    out = []   # frames collected during the DFS
    pruned = []

    def dfs(r, c, parent, ppath, path):
        ch = board[r][c]
        node = parent.children.get(ch)
        if node is None:
            return
        tp = ppath + ch
        path = path + [(r, c)]
        board[r][c] = "#"
        nxt = []
        for dr, dc in DIRS:
            nr, nc = r + dr, c + dc
            if 0 <= nr < len(board) and 0 <= nc < len(board[nr]) and board[nr][nc] in node.children:
                nxt.append((nr, nc))
        board[r][c] = ch
        if node.word:
            w = node.word
            found.append(w)
            node.word = ""
            out.append(Frame(
                Row(grid(path, "match"), tview(root, tp, "match"), weights=[1, 1.3]),
                Text((f'node "{tp}" holds the word "{w}" -> record it, then clear it', "match"),
                     "found: " + ", ".join(found)),
                say=f"The node for \"{tp}\" stores \"{w}\", so I record it. I clear the stored word "
                    "so another path can't report it twice.",
            ))
        else:
            if len(path) == 1:
                lead = f"Cell {r, c} is '{ch}', and the root has a child '{ch}', so a DFS starts here."
            else:
                lead = f"\"{tp}\" is still a prefix in the trie, so I keep going."
            if nxt:
                tail = f" Next letter '{board[nxt[0][0]][nxt[0][1]]}' is at {nxt[0]}."
                st = "active"
            else:
                tail = f" But no neighbour holds a letter that continues \"{tp}\", so I back up."
                st = "reject"
            out.append(Frame(
                Row(grid(path, st), tview(root, tp, st), weights=[1, 1.3]),
                Text(f'path "{tp}" at {(r, c)}', ("next: " + ", ".join(str(x) for x in nxt) if nxt else "dead end", st)),
                say=lead + tail,
            ))
        board[r][c] = "#"
        for dr, dc in DIRS:
            nr, nc = r + dr, c + dc
            if 0 <= nr < len(board) and 0 <= nc < len(board[nr]) and board[nr][nc] != "#":
                dfs(nr, nc, node, tp, path)
        board[r][c] = ch
        if node.word == "" and not node.children:
            del parent.children[ch]
            pruned.append(tp)

    for r in range(len(board)):
        for c in range(len(board[r])):
            before = len(pruned)
            dfs(r, c, root, "", [])
            yield from out
            out.clear()
            if len(pruned) > before:
                cut = pruned[before:]
                yield Frame(
                    Row(Grid(shown, label="board", cell=40), TrieView(root, label="trie after pruning", level_h=36, r=12),
                        weights=[1, 1.3]),
                    Text("pruned: " + ", ".join(f'"{p}"' for p in cut), ("found: " + ", ".join(found), "match")),
                    say="On the way back I cut every node that holds no word and has no children. "
                        f"The \"{cut[-1]}\" branch is gone, so later starts skip it immediately.",
                )

    result = sorted(found)
    assert result == expected, result
    yield Frame(
        Row(Grid(shown, label="board", cell=40), TrieView(root, label="trie at the end", level_h=36, r=12),
            weights=[1, 1.3]),
        Text((f"answer: {result}".replace("'", '"'), "match"),
             "time: O(m · n · 4 · 3^(L-1)) worst case, far less with pruning",
             "space: O(total word length) trie + O(L) recursion"),
        say="Sorted, the answer is eat and oath. One DFS per cell serves all words at once, "
            "and pruning shrinks the trie as words are found.",
    )
