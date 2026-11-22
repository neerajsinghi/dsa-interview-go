import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from _trees09_util import from_level, level  # noqa: E402
from walkthrough import Frame, Text, Tree  # noqa: E402

TITLE = "Diameter of Binary Tree"
PATTERN = "Tree DFS"


def frames():
    vals = [1, 2, 3, 4, 5]
    root = from_level(vals)
    lo, idx = level(root)
    notes, done = {}, set()

    def tree(active=None, style="active", extra=None):
        hl = {i: "visited" for i in done}
        hl.update(extra or {})
        if active is not None:
            hl[active] = style
        return Tree(lo, label="tree (notes: returned height)", hl=hl, notes=dict(notes))

    yield Frame(
        tree(),
        Text("diameter = edges on the longest path between any two nodes"),
        speaker="Interviewer",
        say="Return the diameter of this tree: the number of edges on the longest path between any two nodes. It need not pass through the root.",
    )
    yield Frame(
        tree(),
        Text("brute force: at every node, height(left) + height(right)",
             "heights recomputed from scratch each time → O(n²)"),
        say="Every path turns at some top node, so I could try each node and add its left and right heights. "
            "Recomputing heights per node is O(n²). I'll compute heights once, bottom-up.",
    )

    events = []
    best = 0

    def h(n):
        nonlocal best
        if n is None:
            return 0
        l, r = h(n.left), h(n.right)
        old = best
        best = max(best, l + r)
        events.append((n, l, r, old, best, 1 + max(l, r)))
        return 1 + max(l, r)

    h(root)

    first_leaf = True
    for n, l, r, old, new, ret in events:
        i = idx[id(n)]
        notes[i] = f"h={ret}"
        leaf = n.left is None and n.right is None
        if leaf:
            say = (f"Post-order, so leaves go first. Leaf {n.val}: both sides are 0, the path through it has 0 edges, "
                   "and it returns height 1." if first_leaf else
                   f"Leaf {n.val} also returns height 1. A path through it alone has 0 edges.")
            first_leaf = False
        elif new > old:
            say = (f"At {n.val}: left height {l}, right height {r}. A path turning here uses {l} + {r} = {l + r} edges, "
                   f"which beats best. It returns {ret} to its parent.")
        else:
            say = f"At {n.val}: {l} + {r} = {l + r} edges, no better than best = {new}. It still returns height {ret}."
        yield Frame(
            tree(i, "match"),
            Text(f"node {n.val}: l={l}, r={r}  →  best = max({old}, {l}+{r}) = {new}",
                 f"returns 1 + max({l}, {r}) = {ret}"),
            say=say,
        )
        done.add(i)

    assert best == 3
    path = {idx[id(root.left.left)]: "match", idx[id(root.left)]: "match",
            idx[id(root)]: "match", idx[id(root.right)]: "match"}
    yield Frame(
        tree(extra=path),
        Text(("answer: 3  (4-2-1-3)", "match"), "time O(n), space O(h)"),
        say="The best is 3 edges, for example 4-2-1-3. The helper returns a height but records a path length: "
            "one pass, O(n) time and O(h) stack space.",
    )
