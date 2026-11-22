import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from _heaps10_util import GoHeap, heap_view  # noqa: E402
from walkthrough import Array, Frame, Text  # noqa: E402

TITLE = "K Closest Points to Origin"
PATTERN = "Size-k max-heap"


def dist2(p):
    return p[0] * p[0] + p[1] * p[1]


def pt(p):
    return f"({p[0]},{p[1]})"


def frames():
    # The first doc example has only two points; the second shows the eviction better.
    points, k = [[3, 3], [5, -1], [-2, 4]], 2
    labels = [pt(p) for p in points]
    d2 = [dist2(p) for p in points]

    def heap(h, hl=None):
        return heap_view(h.items, "max-heap by dist²", fmt=pt, tree_fmt=lambda p: str(dist2(p)), hl=hl)

    yield Frame(
        Array(labels, label="points"),
        Text(f"k = {k}"),
        speaker="Interviewer",
        say="Return the k = 2 points closest to the origin, in any order.",
    )
    yield Frame(
        Array(labels, label="points"),
        Array(d2, label="squared distance x² + y²"),
        Text("brute force: sort by distance, take the first k → O(n log n)"),
        say="Sorting all points by distance is O(n log n). I compare squared distances, so there's no sqrt "
            "and no floats. And I only need the best k, not a full sort.",
    )

    h = GoHeap(lambda a, b: dist2(a) > dist2(b))  # farthest on top
    for i, p in enumerate(points):
        h.push(p)
        evicted = h.pop() if len(h) > k else None
        if evicted is None:
            line = f"push {pt(p)} (dist² {dist2(p)}) → size {len(h)}"
            say = ("For the k smallest I use a max-heap of size k: the farthest of my current picks sits on top, "
                   "ready to be thrown out." if i == 0 else
                   f"Push {pt(p)} with distance² {dist2(p)}. It's the farthest so far, so it rises to the top.")
        else:
            line = f"push {pt(p)} (dist² {dist2(p)}) → size {k + 1} > k, pop {pt(evicted)} (dist² {dist2(evicted)})"
            say = (f"Push {pt(p)} with distance² {dist2(p)}. Now there are 3, so I pop the top, the farthest: "
                   f"{pt(evicted)} at {dist2(evicted)}. It can't be among the 2 closest.")
        yield Frame(
            Array(labels, label="points", hl={**{j: "visited" for j in range(i)}, i: "active"}),
            Array(d2, label="squared distance x² + y²", hl={i: "active"}),
            heap(h, hl={0: "compare"}),
            Text(line),
            say=say,
        )

    out = []
    while len(h):
        out.append(h.pop())
    assert sorted(out) == sorted([[3, 3], [-2, 4]]), out
    yield Frame(
        Array(labels, label="points", hl={j: ("match" if points[j] in out else "reject") for j in range(len(points))}),
        Array(d2, label="squared distance x² + y²"),
        heap(GoHeap(lambda a, b: dist2(a) > dist2(b), out), hl={0: "match", 1: "match"}),
        Text((f"answer: {' '.join(pt(p) for p in out)}  (any order)", "match"), "time O(n log k), space O(k)"),
        say="Popping what's left gives (-2,4) and (3,3). Each point costs one push and maybe one pop on a heap "
            "of size k: O(n log k) time, O(k) space.",
    )
