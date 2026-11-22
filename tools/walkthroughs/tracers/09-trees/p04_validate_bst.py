import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from _trees09_util import from_level, level  # noqa: E402
from walkthrough import Frame, Text, Tree  # noqa: E402

TITLE = "Validate BST"
PATTERN = "Tree DFS + bounds"


def _b(lo, hi):
    return f"({'-∞' if lo is None else lo}, {'∞' if hi is None else hi})"


def trace(vals):
    """IsValidBST as a generator: yields (node, lo, hi, ok) at each check."""
    root = from_level(vals)

    def valid(n, lo, hi):
        if n is None:
            return True
        ok = not ((lo is not None and n.val <= lo) or (hi is not None and n.val >= hi))
        yield n, lo, hi, ok
        if not ok:
            return False
        return (yield from valid(n.left, lo, n.val)) and (yield from valid(n.right, n.val, hi))

    return root, valid(root, None, None)


def frames():
    ex1 = [2, 1, 3]
    yield Frame(
        Tree(ex1, label="example 1"),
        Text("left subtree < node < right subtree, for EVERY node"),
        speaker="Interviewer",
        say="Is this a valid binary search tree? Everything in a node's left subtree must be strictly smaller, "
            "and everything in its right subtree strictly larger.",
    )
    yield Frame(
        Tree(ex1, label="example 1"),
        Text("brute force: compare each node with every node below it → O(n²)",
             "parent/child checks only → fast but WRONG"),
        say="Brute force compares each node with its whole subtree, which is O(n²). "
            "Checking only parent and child is fast but misses grandchildren. "
            "Instead I'll pass down the allowed range.",
    )

    for name, vals, want in (("example 1", ex1, True), ("example 3", [5, 4, 6, None, None, 3, 7], False)):
        root, gen = trace(vals)
        lo_list, idx = level(root)
        notes, hl = {}, {}
        result = True
        first = True
        if name == "example 3":
            yield Frame(
                Tree(lo_list, label=name),
                Text("every parent/child pair looks fine: 4<5, 6>5, 3<6, 7>6"),
                speaker="Interviewer",
                say="That one was easy. What about this tree? Every parent and child pair looks correct.",
            )
        for n, lo, hi, ok in gen:
            i = idx[id(n)]
            notes[i] = _b(lo, hi)
            hl[i] = "match" if ok else "reject"
            if n is root:
                say = (f"The root {n.val} has no limits, so it passes. Going left, everything must be below {n.val}; "
                       f"going right, everything must be above {n.val}." if first else
                       f"Root {n.val}: no limits, passes. Its right subtree must stay above {n.val}.")
            elif ok:
                say = f"{n.val} must lie in {_b(lo, hi)}. It does, so it passes and hands tighter bounds to its children."
            else:
                say = (f"3 is 6's left child, so it must be below 6. But it is also in 5's right subtree, so it must be above 5. "
                       f"3 is not in {_b(lo, hi)}, so the tree is invalid.")
            yield Frame(
                Tree(lo_list, label=name + "  (notes: allowed range)", hl=dict(hl), notes=dict(notes)),
                Text(f"check {n.val} in {_b(lo, hi)} → {'ok' if ok else 'FAIL'}",
                     ("valid so far" if ok else ("return false", "reject"))),
                say=say,
            )
            first = False
            if not ok:
                result = False
        assert result == want
        if want:
            yield Frame(
                Tree(lo_list, label=name, hl=dict(hl), notes=dict(notes)),
                Text(("example 1 → true", "match")),
                say="All three nodes sit inside their ranges, so [2,1,3] is a valid BST.",
            )

    yield Frame(
        Tree([5, 4, 6, None, None, 3, 7], label="example 3",
             hl={0: "compare", 5: "reject"}, notes={5: "(5, 6)"}),
        Text(("example 1 → true, example 3 → false", "match"), "time O(n), space O(h)"),
        say="Each node is checked once against two numbers: O(n) time and O(h) stack. "
            "The bounds are pointers, so nil means no limit and no sentinel value can clash with real data.",
    )
