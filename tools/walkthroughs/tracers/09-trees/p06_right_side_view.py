import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from _trees09_util import from_level, level  # noqa: E402
from walkthrough import Array, Frame, Text, Tree  # noqa: E402

TITLE = "Right Side View"
PATTERN = "Tree BFS"


def frames():
    vals = [1, 2, 3, None, 5, None, 4]
    root = from_level(vals)
    lo, idx = level(root)
    seen = set()
    view = []

    def show(queue, level_hl=None, hl_q=None):
        hl = {i: "visited" for i in seen}
        hl.update(level_hl or {})
        return (Tree(lo, label="tree", hl=hl),
                Array([n.val for n in queue], label="queue", hl=hl_q or {}, show_index=False, cell=50),
                Array(view, label="view", show_index=False, cell=50))

    yield Frame(
        *show([]),
        speaker="Interviewer",
        say="Stand to the right of this tree. Return the values you can see from top to bottom.",
    )
    yield Frame(
        *show([]),
        Text("brute force: for each depth, search the whole tree for its rightmost node → O(n·h)"),
        say="Following right children isn't enough: 4 hangs off a left branch. I could search the tree once per depth, "
            "but that's O(n·h). What I want is the last node on each level, which BFS gives me.",
    )

    queue = [root]
    depth = 0
    while queue:
        size = len(queue)
        last = queue[size - 1]
        view.append(last.val)
        level_hl = {idx[id(n)]: "active" for n in queue[:size]}
        level_hl[idx[id(last)]] = "match"
        hl_q = {k: "active" for k in range(size)}
        hl_q[size - 1] = "match"
        vals_here = ", ".join(str(n.val) for n in queue[:size])
        if depth == 0:
            say = "The queue starts with the root. I snapshot size = len(queue); the last of those size nodes is the one seen from the right."
        elif depth == 1:
            say = f"Level 1 holds {vals_here}. The last one, {last.val}, is visible. Then I enqueue each node's children, left before right."
        else:
            say = f"Level 2 is just {vals_here}, so {last.val} is visible even though it is a left child: nothing is to its right."
        yield Frame(
            *show(queue, level_hl, hl_q),
            Text(f"level {depth}: size = {size}, last = queue[{size - 1}] = {last.val} → view += {last.val}"),
            say=say,
        )
        for n in queue[:size]:
            if n.left:
                queue.append(n.left)
            if n.right:
                queue.append(n.right)
            seen.add(idx[id(n)])
        queue = queue[size:]
        depth += 1

    assert view == [1, 3, 4]
    yield Frame(
        *show([]),
        Text(("answer: [1, 3, 4]", "match"), "time O(n), space O(width) for the queue"),
        say="The view is [1, 3, 4]. Each node enters and leaves the queue once: O(n) time, O(width) space. "
            "The DFS version visits right first and records the first node at each new depth.",
    )
