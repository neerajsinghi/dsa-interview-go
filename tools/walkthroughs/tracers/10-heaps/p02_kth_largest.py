import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from _heaps10_util import GoHeap, heap_view  # noqa: E402
from walkthrough import Array, Frame, Text  # noqa: E402

TITLE = "Kth Largest Element"
PATTERN = "Size-k min-heap"


def frames():
    nums, k = [3, 2, 1, 5, 6, 4], 2

    yield Frame(
        Array(nums, label="nums"),
        Text(f"k = {k}"),
        speaker="Interviewer",
        say="Return the k-th largest value in this unsorted array, with k = 2. Duplicates count separately.",
    )
    yield Frame(
        Array(sorted(nums), label="sorted copy", hl={len(nums) - k: "match"}),
        Text(f"brute force: sort, answer = s[n-k] = {sorted(nums)[len(nums) - k]}", "O(n log n) time, O(n) copy"),
        say="Sorting and reading index n - k works in O(n log n). But I only need the top k values, "
            "not the order of everything.",
    )

    h = GoHeap(lambda a, b: a < b)
    for i, v in enumerate(nums):
        h.push(v)
        evicted = None
        if len(h) > k:
            evicted = h.pop()
        if evicted is None:
            line = f"push {v} → size {len(h)} <= k, keep"
            say = (f"I keep a min-heap of the {k} largest values seen so far. Push {v}; the heap isn't full yet."
                   if i == 0 else f"Push {v}. Now the heap holds {k} values and its top, {h.peek()}, is the smallest of them.")
        elif evicted == v:
            line = f"push {v} → size {k + 1} > k, pop {evicted}"
            say = f"{v} goes in and comes straight back out: it's smaller than everything in the top {k}."
        else:
            line = f"push {v} → size {k + 1} > k, pop {evicted}"
            say = (f"Push {v}, so the heap has {k + 1} items. I pop the smallest, {evicted}: "
                   f"{k} values beat it, so it can't be the {k}nd largest.")
        yield Frame(
            Array(nums, label="nums", hl={**{j: "visited" for j in range(i)}, i: "active"}),
            heap_view(h.items, "min-heap", hl={0: "compare"}),
            Text(line, f"top = {h.peek()}"),
            say=say,
        )

    answer = h.peek()
    assert answer == 5
    yield Frame(
        Array(nums, label="nums", hl={j: "visited" for j in range(len(nums))}),
        heap_view(h.items, "min-heap", hl={0: "match"}),
        Text(("answer: top = 5", "match"), "heap: O(n log k) time, O(k) space · quickselect: O(n) expected"),
        say="After the last value the heap holds the 2 largest, 5 and 6, and the top is the answer: 5. "
            "That's O(n log k) time and O(k) space. Quickselect gets O(n) expected.",
    )
