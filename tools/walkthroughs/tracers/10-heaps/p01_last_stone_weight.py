import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from _heaps10_util import GoHeap, heap_view  # noqa: E402
from walkthrough import Array, Frame, Text  # noqa: E402

TITLE = "Last Stone Weight"
PATTERN = "Max-heap"


def frames():
    stones = [2, 7, 4, 1, 8, 1]

    yield Frame(
        Array(stones, label="stones"),
        Text("each turn: smash the two heaviest, y >= x → leaves y - x (or nothing)"),
        speaker="Interviewer",
        say="Each turn, take the two heaviest stones and smash them. Equal stones both vanish; otherwise the "
            "difference stays. Return the last stone's weight, or 0.",
    )
    yield Frame(
        Array(sorted(stones), label="sorted copy"),
        Text("brute force: re-sort every turn → O(n · n log n)"),
        say="I could sort, take the last two, put the difference back and sort again. "
            "That's a sort per turn, O(n² log n). I only ever need the largest, so a max-heap fits.",
    )

    h = GoHeap(lambda a, b: a > b, stones)  # FromSlice: heap.Init, O(n)
    yield Frame(
        heap_view(h.items, "max-heap", hl={0: "active"}),
        Text(f"FromSlice({stones}) → {h.items}", "heapify is O(n); the largest is at index 0"),
        say="FromSlice heapifies the copy in O(n). Every parent is now at least as big as its children, "
            "so the heaviest stone, 8, is at the root.",
    )

    turn = 0
    while len(h) > 1:
        turn += 1
        y = h.pop()
        x = h.pop()
        if y != x:
            h.push(y - x)
            line = f"turn {turn}: pop {y}, pop {x} → push {y} - {x} = {y - x}"
            if turn == 1:
                say = (f"Pop gives {y}, then {x}. They differ, so I push {y - x}. "
                       "Each pop and push is O(log n): the heap fixes itself along one root-to-leaf path.")
            else:
                say = f"The two heaviest are now {y} and {x}. I push back {y - x}."
            hl = {}
        else:
            line = f"turn {turn}: pop {y}, pop {x} → equal, both vanish"
            say = f"{y} and {x} are equal, so both disappear and nothing is pushed."
            hl = {}
        yield Frame(heap_view(h.items, "max-heap", hl=hl), Text(line, f"stones left: {len(h)}"), say=say)

    answer = h.peek() if len(h) else 0
    assert answer == 1
    yield Frame(
        heap_view(h.items, "max-heap", hl={0: "match"}),
        Text(("answer: 1", "match"), "time O(n log n), space O(n)"),
        say="One stone of weight 1 is left, so the answer is 1. There are at most n - 1 turns of O(log n) each: "
            "O(n log n) time and O(n) space for the heap.",
    )
