import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from _heaps10_util import GoHeap, heap_view  # noqa: E402
from walkthrough import Array, Frame, Grid, Text  # noqa: E402

TITLE = "Merge k Sorted Lists"
PATTERN = "K-way merge"


class ListNode:
    def __init__(self, val, nxt=None, where=None):
        self.val = val
        self.next = nxt
        self.where = where  # (list, position) for drawing only


def frames():
    lists_vals = [[1, 4, 5], [1, 3, 4], [2, 6]]
    heads = []
    for r, vals in enumerate(lists_vals):
        head = None
        for c in range(len(vals) - 1, -1, -1):
            head = ListNode(vals[c], head, (r, c))
        heads.append(head)
    total = sum(len(v) for v in lists_vals)
    width = max(len(v) for v in lists_vals)
    rows = [v + [""] * (width - len(v)) for v in lists_vals]
    labels = [f"list {r}" for r in range(len(lists_vals))]

    out = []
    used = set()

    def widgets(h, hl_extra=None, out_hl=None):
        hl = {pos: "visited" for pos in used}
        for node in (h.items if h else []):
            hl[node.where] = "active"
        hl.update(hl_extra or {})
        return (
            Grid(rows, label="input lists (blue = in the heap)", hl=hl, row_labels=labels, cell=46),
            heap_view([x.val for x in h.items] if h else [], "min-heap of heads", weights=(1, 1.2)),
            Array(out + [""] * (total - len(out)), label="merged", hl=out_hl or {}, cell=46),
        )

    yield Frame(
        *widgets(None),
        speaker="Interviewer",
        say="Merge these k sorted linked lists into one sorted list, reusing the nodes.",
    )
    yield Frame(
        *widgets(None),
        Text("brute force: merge list 0 with list 1, the result with list 2, ... → O(N·k)"),
        say="Merging one list at a time re-walks the growing result k times: O(N·k). "
            "The smallest remaining value is always one of the k current heads, so I keep the heads in a min-heap.",
    )

    h = GoHeap(lambda a, b: a.val < b.val)
    for l in heads:
        if l is not None:
            h.push(l)
    yield Frame(
        *widgets(h),
        Text(f"push the {len(h)} heads → heap top = {h.peek().val}"),
        say="I push every non-nil head. The heap never holds more than k nodes, one per list.",
    )

    step = 0
    while len(h):
        node = h.pop()
        out.append(node.val)
        used.add(node.where)
        step += 1
        pushed = None
        if node.next is not None:
            h.push(node.next)
            pushed = node.next
        r = node.where[0]
        line = f"pop {node.val} (list {r}) → append" + (f"; push its next, {pushed.val}" if pushed else "; list ends, push nothing")
        if step == 1:
            say = (f"Pop the smallest head, {node.val} from list {r}, and append it. "
                   f"Its successor {pushed.val} takes its place in the heap.")
        elif step == 2:
            say = f"Pop {node.val} from list {r}. Ties don't matter; either 1 can go first. Push {pushed.val}."
        elif step == 3:
            say = f"Pop {node.val} and push {pushed.val}. Each step is one pop and at most one push, both O(log k)."
        elif step < total:
            if step == 4:
                say = "Same again for the rest: pop the smallest head, append it, push its successor."
            else:
                say = (f"Pop {node.val}. List {r} has nothing left, so nothing is pushed and the heap shrinks."
                       if pushed is None else f"Pop {node.val}, push {pushed.val}.")
        else:
            say = f"Pop {node.val}, the last node. The heap is empty, so the merge is done."
        yield Frame(
            *widgets(h, {node.where: "match"}, {len(out) - 1: "match"}),
            Text(line),
            say=say,
        )

    assert out == [1, 1, 2, 3, 4, 4, 5, 6], out
    yield Frame(
        *widgets(h, None, {k: "match" for k in range(total)}),
        Text(("answer: 1→1→2→3→4→4→5→6", "match"), "time O(N log k), space O(k)"),
        say="Every one of the N nodes is pushed and popped once on a heap of at most k items: O(N log k) time, "
            "O(k) extra space. Pairwise merging also gets O(N log k).",
    )
