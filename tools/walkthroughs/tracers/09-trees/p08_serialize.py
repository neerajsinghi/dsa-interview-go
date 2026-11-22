import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from _trees09_util import Node, from_level, level  # noqa: E402
from walkthrough import Array, Frame, Text, Tree  # noqa: E402

TITLE = "Serialize and Deserialize"
PATTERN = "Preorder + null markers"


def frames():
    vals = [1, 2, 3, None, None, 4, 5]
    root = from_level(vals)
    lo, idx = level(root)

    yield Frame(
        Tree(lo, label="tree"),
        Array([""] * 11, label="tokens", show_index=True, cell=56),
        speaker="Interviewer",
        say="Design Serialize, which turns any binary tree into a string, and Deserialize, which turns that string "
            "back into the same tree. Values can repeat or be negative.",
    )
    yield Frame(
        Tree(lo, label="tree"),
        Array([""] * 11, label="tokens", show_index=True, cell=56),
        Text("plain preorder [1,2,3,4,5] fits many shapes; preorder + inorder needs two lists",
             "and breaks when values repeat"),
        say="One plain traversal is ambiguous, and preorder plus inorder fails with duplicate values. "
            "If I also write every nil child as '#', preorder alone pins down the shape.",
    )

    # --- Serialize: preorder walk, '#' for nil
    tokens = []
    seen = set()

    def walk(n):
        if n is None:
            tokens.append("#")
            return
        tokens.append(str(n.val))
        seen.add(idx[id(n)])
        yield n
        yield from walk(n.left)
        yield from walk(n.right)

    def tok_array(hl=None, ptr=None):
        return Array(tokens + [""] * (11 - len(tokens)), label="tokens", show_index=True, cell=56,
                     hl=hl or {}, pointers=ptr or {})

    prev = 0
    shown = 0
    for n in walk(root):
        i = idx[id(n)]
        hl = {k: "visited" for k in range(len(tokens) - 1)}
        for k in range(prev, len(tokens) - 1):
            if tokens[k] == "#":
                hl[k] = "special"
        hl[len(tokens) - 1] = "active"
        if shown == 0:
            say = "Serialize walks in preorder: write the node's value, then its left subtree, then its right."
        elif n.val == 3:
            say = "Back up at 1, its right child is 3. Before reaching it, 2's two nil children each wrote '#'."
        elif n.val == 2:
            say = "Next comes 2, the left child of 1. Its children are both nil, so the walk will write '#' twice."
        else:
            say = (f"Then 4, the left child of 3." if n.val == 4 else
                   "4's two nil children wrote '#', '#', and then 5 is written. The '#' markers record exactly where each subtree stops.")
        prev = len(tokens)
        shown += 1
        yield Frame(
            Tree(lo, label="tree", hl={**{k: "visited" for k in seen}, i: "active"}),
            tok_array(hl),
            Text(f"visit {n.val} → write \"{n.val}\""),
            say=say,
        )

    data = ",".join(tokens)
    assert data == "1,2,#,#,3,4,#,#,5,#,#", data
    yield Frame(
        Tree(lo, label="tree", hl={k: "visited" for k in seen}),
        tok_array({k: ("special" if t == "#" else "visited") for k, t in enumerate(tokens)}),
        Text(("Serialize → \"" + data + "\"", "match")),
        say="The full string has 5 values and 6 '#' markers. A tree with n nodes always has n + 1 nil children, "
            "so the output is O(n).",
    )

    # --- Deserialize: read tokens in the same order with a shared cursor
    toks = data.split(",")
    cur = 0
    holder = {"root": None}

    def build(parent=None, side=None):
        nonlocal cur
        tok = toks[cur]
        at = cur
        cur += 1
        if tok == "#":
            return None
        n = Node(int(tok))
        if parent is None:
            holder["root"] = n
        else:
            setattr(parent, side, n)
        yield n, at
        n.left = yield from build(n, "left")
        n.right = yield from build(n, "right")
        return n

    for n, at in build():
        if n.val not in (1, 2, 3):
            continue
        lv, bidx = level(holder["root"])
        hl = {k: "visited" for k in range(at)}
        hl[at] = "active"
        if n.val == 1:
            say = "Deserialize reads the tokens in the same order with a cursor i. \"1\" becomes the root; then I build its left subtree, then its right."
        elif n.val == 2:
            say = "\"2\" is 1's left child. The next two tokens are '#', so both of 2's children are nil and that subtree is done."
        else:
            say = "The cursor is now at \"3\", which becomes 1's right child. The same rule builds 4 and 5 from the rest."
        yield Frame(
            Tree(lv, label="rebuilt so far", hl={bidx[id(n)]: "match"}),
            Array(toks, label="tokens", show_index=True, cell=56, hl=hl, pointers={"i": at}),
            Text(f"token[{at}] = \"{toks[at]}\" → new node {n.val}"),
            say=say,
        )

    assert cur == len(toks)
    lv, _ = level(holder["root"])
    assert lv == lo
    yield Frame(
        Tree(lv, label="rebuilt tree", hl={k: "match" for k in range(len(lv)) if lv[k] is not None}),
        Array(toks, label="tokens", show_index=True, cell=56, hl={k: "visited" for k in range(len(toks))}),
        Text(("Deserialize(data) = [1,2,3,null,null,4,5] ✓", "match"), "both directions: O(n) time, O(n) space"),
        say="The rebuilt tree matches the original, and the cursor ended exactly at the last token. "
            "Both directions touch each token once: O(n) time and space.",
    )
