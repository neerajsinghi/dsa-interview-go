import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from _trees09_util import from_level, level  # noqa: E402
from walkthrough import Frame, Row, Text, Tree  # noqa: E402

TITLE = "Invert Binary Tree"
PATTERN = "Tree DFS"


def frames():
    vals = [4, 2, 7, 1, 3, 6, 9]
    root = from_level(vals)
    done = set()  # ids of nodes whose subtree is already mirrored

    def view(active=None, style="active", note=None):
        lo, idx = level(root)
        hl = {idx[k]: "visited" for k in done}
        notes = {}
        if active is not None:
            hl[idx[id(active)]] = style
            if note:
                notes[idx[id(active)]] = note
        return Row(Tree(vals, label="input"), Tree(lo, label="tree now", hl=hl, notes=notes))

    yield Frame(
        view(),
        Text("mirror: swap left and right at every node"),
        speaker="Interviewer",
        say="Invert this binary tree: mirror it so that every node's left and right subtrees trade places. Return the root.",
    )
    yield Frame(
        view(),
        Text("naive: read each level, reverse it, rebuild a new tree",
             "→ O(n) time but O(n) extra copies and fiddly null slots"),
        say="I could copy each level out, reverse it and rebuild, but that allocates a second tree. "
            "Mirroring is just a local swap at each node, done in place.",
    )

    def invert(n):
        """InvertTree as a generator: yields each node right after its swap."""
        if n is None:
            return None
        # Go evaluates the right-hand side left to right: right subtree first.
        new_left = yield from invert(n.right)
        new_right = yield from invert(n.left)
        n.left, n.right = new_left, new_right
        yield n
        return n

    yield Frame(
        view(root, "compare"),
        Text("InvertTree(4): left, right = InvertTree(7), InvertTree(2)"),
        say="At 4 the code says left, right = InvertTree(right), InvertTree(left). "
            "Go evaluates both calls first, right subtree first, and only then assigns.",
    )

    leaves_said = False
    for n in invert(root):
        v = n.val
        leaf = n.left is None and n.right is None
        done.add(id(n))
        if leaf:
            if leaves_said:
                continue
            leaves_said = True
            yield Frame(
                view(n, "match", "leaf"),
                Text(f"InvertTree({v}): swap nil, nil → returns {v}"),
                say=f"{v} is a leaf. Swapping two nils changes nothing, so it returns itself. "
                    "The other leaves behave the same way.",
            )
            continue
        kids = f"{n.left.val}, {n.right.val}"
        yield Frame(
            view(n, "match"),
            Text(f"InvertTree({v}): children are now ({kids})"),
            say=(f"Both of {v}'s subtrees are mirrored, so now the assignment runs and its children swap. "
                 f"{v}'s children are now {kids}."
                 if n is not root else
                 f"Finally the root: 7's and 2's subtrees are mirrored, so 4 swaps them. The whole tree is now a mirror image."),
        )

    lo, _ = level(root)
    assert lo == [4, 7, 2, 9, 6, 3, 1], lo
    yield Frame(
        view(),
        Text(("answer: [4,7,2,9,6,3,1]", "match"), "time O(n), space O(h)"),
        say="The result is [4,7,2,9,6,3,1]. Each node is swapped once: O(n) time. "
            "The recursion stack is O(h); the BFS version uses a queue instead.",
    )
