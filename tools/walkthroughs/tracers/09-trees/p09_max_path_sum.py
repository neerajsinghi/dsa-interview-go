import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from _trees09_util import from_level, level  # noqa: E402
from walkthrough import Frame, Text, Tree  # noqa: E402

TITLE = "Maximum Path Sum"
PATTERN = "Tree DFS"


def frames():
    vals = [-10, 9, 20, None, None, 15, 7]
    root = from_level(vals)
    lo, idx = level(root)
    notes, done = {}, set()

    def tree(active=None, style="active", extra=None):
        hl = {i: "visited" for i in done}
        hl.update(extra or {})
        if active is not None:
            hl[active] = style
        return Tree(lo, label="tree (notes: returned gain)", hl=hl, notes=dict(notes))

    yield Frame(
        tree(),
        Text("path: connected nodes, no repeats, at least one node"),
        speaker="Interviewer",
        say="Values can be negative. Find the largest sum of any path in this tree. "
            "A path can start and end anywhere but can't reuse a node.",
    )
    yield Frame(
        tree(),
        Text("brute force: DFS from every node over every simple path → O(n²)"),
        say="I could start a DFS from every node and sum every path, which is O(n²). "
            "Better: like Diameter, return one thing and record another in a single post-order pass.",
    )

    best = None
    events = []

    def gain(n):
        nonlocal best
        if n is None:
            return 0
        gl, gr = gain(n.left), gain(n.right)
        l, r = max(gl, 0), max(gr, 0)
        old = best
        best = n.val + l + r if best is None else max(best, n.val + l + r)
        ret = n.val + max(l, r)
        events.append((n, gl, gr, l, r, old, best, ret))
        return ret

    gain(root)

    for k, (n, gl, gr, l, r, old, new, ret) in enumerate(events):
        i = idx[id(n)]
        notes[i] = f"g={ret}"
        old_s = "-∞" if old is None else str(old)
        if k == 0:
            say = (f"Leaf {n.val}: both gains are 0. The path that turns here is just {n.val}, so best = {new}. "
                   f"It offers its parent a downward gain of {ret}.")
        elif n.val == 20:
            say = (f"At 20: left gain {l}, right gain {r}. The path that turns here, 15-20-7, sums to {n.val + l + r}. "
                   f"But a parent can only extend one side, so it returns {n.val} + {max(l, r)} = {ret}.")
        elif n is root:
            say = (f"At the root: -10 + {l} + {r} = {n.val + l + r}, less than {new}. "
                   f"The negative root makes paths through it worse, so the best stays at 15-20-7.")
        else:
            say = f"Leaf {n.val} returns gain {ret}. The path of just {n.val} doesn't beat best = {new}."
        yield Frame(
            tree(i, "match"),
            Text(f"node {n.val}: l = max({gl},0) = {l}, r = max({gr},0) = {r}",
                 f"best = max({old_s}, {n.val}+{l}+{r}) = {new};  returns {n.val}+max({l},{r}) = {ret}"),
            say=say,
        )
        done.add(i)

    assert best == 42
    path = {idx[id(root.right)]: "match", idx[id(root.right.left)]: "match", idx[id(root.right.right)]: "match"}
    yield Frame(
        tree(extra=path),
        Text(("answer: 42  (15-20-7)", "match"), "time O(n), space O(h)"),
        say="The answer is 42. Each node is visited once: O(n) time, O(h) stack. "
            "Negative gains are clamped to 0, and best starts at minus infinity so an all-negative tree still works.",
    )
