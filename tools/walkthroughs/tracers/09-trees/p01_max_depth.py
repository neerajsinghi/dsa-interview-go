import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from _trees09_util import from_level, level  # noqa: E402
from walkthrough import Frame, Text, Tree  # noqa: E402

TITLE = "Maximum Depth"
PATTERN = "Tree DFS"


def frames():
    vals = [3, 9, 20, None, None, 15, 7]
    root = from_level(vals)
    lo, idx = level(root)
    notes, done = {}, set()

    def tree(active=None, style="active"):
        hl = {i: "visited" for i in done}
        if active is not None:
            hl[active] = style
        return Tree(lo, label="tree", hl=hl, notes=dict(notes))

    yield Frame(
        tree(),
        Text("depth = nodes on the longest root-to-leaf path"),
        speaker="Interviewer",
        say="Return the maximum depth of this binary tree: the number of nodes on the longest path from the root down to a leaf.",
    )
    yield Frame(
        tree(),
        Text("brute force: list every root-to-leaf path, measure each",
             "→ O(n·h) time, and the paths take extra memory"),
        say="I could list every root-to-leaf path and take the longest, but that copies paths around. "
            "A node's depth only needs its children's depths, so one recursive pass is enough.",
    )

    events = []

    def depth(n):
        if n is None:
            return 0
        l, r = depth(n.left), depth(n.right)
        d = 1 + max(l, r)
        events.append((n, l, r, d))
        return d

    answer = depth(root)

    yield Frame(
        tree(0, "compare"),
        Text("MaxDepth(3) = 1 + max(MaxDepth(9), MaxDepth(20))"),
        say="At the root I can't answer yet. I ask the left subtree, then the right subtree, "
            "and add one for the root itself. A nil child returns 0.",
    )
    for k, (n, l, r, d) in enumerate(events):
        i = idx[id(n)]
        notes[i] = f"d={d}"
        leaf = n.left is None and n.right is None
        if leaf:
            say = (f"{n.val} is a leaf. Both children are nil and return 0, so its depth is 1 + max(0, 0) = 1."
                   if k == 0 else f"{n.val} is also a leaf, so it returns 1.")
        elif n is root:
            say = f"Back at the root: left gave {l}, right gave {r}. The deeper side wins, so the depth is 1 + {max(l, r)} = {d}."
        else:
            say = f"Node {n.val} now has both answers: left {l}, right {r}. It returns 1 + {max(l, r)} = {d} to its parent."
        yield Frame(
            tree(i, "match"),
            Text(f"MaxDepth({n.val}) = 1 + max({l}, {r}) = {d}"),
            say=say,
        )
        done.add(i)

    assert answer == 3
    yield Frame(
        tree(0, "match"),
        Text(("answer: 3", "match"), "time O(n), space O(h) for the recursion"),
        say="The answer is 3. Every node is visited once, so it's O(n) time. "
            "The recursion is as deep as the tree, so O(h) space, which is O(n) for a skewed tree.",
    )
