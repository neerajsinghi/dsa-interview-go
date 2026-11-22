import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from _heaps10_util import GoHeap  # noqa: E402
from walkthrough import Array, Frame, Row, Text, Tree  # noqa: E402

TITLE = "Find Median from Data Stream"
PATTERN = "Two heaps"


def frames():
    stream = [1, 2, 3]
    expect = {2: 1.5, 3: 2.0}

    low = GoHeap(lambda a, b: a > b)   # max-heap: smaller half
    high = GoHeap(lambda a, b: a < b)  # min-heap: larger half

    def widgets(hl_low=None, hl_high=None, seen=0):
        return (
            Array(stream, label="stream", hl={**{k: "visited" for k in range(seen - 1)},
                                              **({seen - 1: "active"} if seen else {})}),
            Row(Tree(low.items, label="low: max-heap (smaller half)", hl=hl_low or {}),
                Tree(high.items, label="high: min-heap (larger half)", hl=hl_high or {})),
            Row(Array(low.items, label="low array", show_index=False, hl=hl_low or {}),
                Array(high.items, label="high array", show_index=False, hl=hl_high or {})),
        )

    def median():
        if len(low) > len(high):
            return float(low.peek())
        return (low.peek() + high.peek()) / 2

    yield Frame(
        *widgets(),
        Text("AddNum(1)  AddNum(2)  FindMedian()  AddNum(3)  FindMedian()"),
        speaker="Interviewer",
        say="Numbers arrive one at a time. Support AddNum, and FindMedian at any point. "
            "With an even count, the median is the mean of the two middle values.",
    )
    yield Frame(
        *widgets(),
        Text("brute force: keep a sorted slice; insert = binary search + shift → O(n) per AddNum"),
        say="A sorted slice gives the median in O(1), but every insert shifts elements: O(n). "
            "I only ever need the middle, so I split the data into two halves kept in two heaps.",
    )

    for i, v in enumerate(stream, 1):
        low.push(v)
        yield Frame(
            *widgets({low.items.index(v): "active"}, None, i),
            Text(f"AddNum({v}): step 1 · low.Push({v})"),
            say=(f"AddNum({v}) always goes into low first, the max-heap of the smaller half."
                 if i == 1 else f"AddNum({v}): push {v} into low."),
        )
        moved = low.pop()
        high.push(moved)
        yield Frame(
            *widgets(None, {high.items.index(moved): "active"}, i),
            Text(f"step 2 · high.Push(low.Pop()) → {moved} moves to high"),
            say=(f"Then low's largest moves to high. This keeps every value in low <= every value in high."
                 if i == 1 else f"low's largest, {moved}, moves to high, so the halves stay ordered."),
        )
        if len(high) > len(low):
            back = high.pop()
            low.push(back)
            yield Frame(
                *widgets({low.items.index(back): "active"}, None, i),
                Text(f"step 3 · high has more ({len(high) + 1} > {len(low) - 1}) → move {back} back to low"),
                say=(f"Now high has more items than low, so its smallest, {back}, goes back. "
                     "low may hold one extra item, never fewer."),
            )
        else:
            yield Frame(
                *widgets(None, None, i),
                Text(f"step 3 · sizes low={len(low)}, high={len(high)} → balanced, nothing to move"),
                say="The sizes are equal, so no rebalance is needed.",
            )
        if i in expect:
            m = median()
            assert m == expect[i], (i, m)
            if len(low) > len(high):
                calc = f"odd count → low's top = {low.peek()}"
            else:
                calc = f"even count → ({low.peek()} + {high.peek()}) / 2 = {m}"
            yield Frame(
                *widgets({0: "match"}, {0: "match"} if len(low) == len(high) else None, i),
                Text((f"FindMedian() = {m:g}", "match"), calc),
                say=(f"FindMedian: the sizes are equal, so the median is the mean of both tops, ({low.peek()} + {high.peek()}) / 2 = {m:g}."
                     if len(low) == len(high) else
                     f"FindMedian: low has one extra, so its top, {low.peek()}, is the median."),
            )

    yield Frame(
        *widgets({0: "match"}, None, len(stream)),
        Text(("medians: 1.5, then 2.0", "match"), "AddNum O(log n), FindMedian O(1), space O(n)"),
        say="AddNum does at most three heap operations, so O(log n). FindMedian only peeks at the tops: O(1). "
            "Both heaps together store every number: O(n) space.",
    )
