import math

from walkthrough import MUTED, STYLES, Frame, Text, Widget, arrow, centered, fit_font, font, _label

TITLE = "Copy List with Random Pointer"
PATTERN = "Interleave clones"

RANDOM_ORIG = "#bc4c00"   # random pointers of original nodes
RANDOM_CLONE = "#1a7f37"  # random pointers of clones


class Node:
    def __init__(self, val, slot, clone=False):
        self.val, self.slot, self.clone = val, slot, clone
        self.next = None
        self.random = None


class RandList(Widget):
    """Nodes in fixed slots. Next arrows run above the boxes (straight when the
    target is the next visible box, arced otherwise); Random arrows hang below."""

    def __init__(self, nodes, nslots, label=None, hl=None, show_random=True):
        # Snapshot now: frames are rendered after the algorithm has finished
        # mutating the real nodes.
        snap = {}
        for n in nodes:
            if n is not None:
                snap[id(n)] = Node(n.val, n.slot, n.clone)
        for n in nodes:
            if n is not None:
                c = snap[id(n)]
                c.next = Node(None, n.next.slot) if n.next is not None else None
                c.random = Node(None, n.random.slot) if n.random is not None else None
        self.nodes = list(snap.values())
        self.nslots = nslots
        self.label = label
        self.hl = hl or {}
        self.show_random = show_random

    def height(self, width):
        return (20 if self.label else 0) + 30 + 44 + 62

    def draw(self, d, x, y, width):
        _label(d, x, y, self.label)
        y += 20 if self.label else 0
        step = min(80, (width - 30) / self.nslots)
        box = max(28, min(46, step - 20))
        top = y + 30
        cy = top + 22
        by = top + 44
        cx = lambda s: x + s * step + box / 2
        visible = sorted(n.slot for n in self.nodes)

        for n in self.nodes:
            if n.next is None:
                d.text((cx(n.slot) + box / 2 + 3, cy), "∅", font=font("mono", 13), fill=MUTED, anchor="lm")
                continue
            j = n.next.slot
            later = [s for s in visible if s > n.slot]
            if later and later[0] == j:
                arrow(d, cx(n.slot) + box / 2, cy, cx(j) - box / 2 - 1, cy, "#57606a")
            else:
                x1, x2 = cx(n.slot), cx(j)
                h = 12 + 4 * min(4, abs(j - n.slot))
                pts = [(x1 + (x2 - x1) * t / 20, top - 2 - h * math.sin(math.pi * t / 20)) for t in range(21)]
                d.line(pts, fill="#57606a", width=2)
                arrow(d, pts[-2][0], pts[-2][1], x2, top - 1, "#57606a", head=7)

        if self.show_random:
            for n in self.nodes:
                if n.random is None:
                    continue
                col = RANDOM_CLONE if n.clone else RANDOM_ORIG
                x1, x2 = cx(n.slot) + (4 if n.clone else -4), cx(n.random.slot) + (4 if n.clone else -4)
                dist = abs(n.random.slot - n.slot)
                h = 10 + 5 * min(9, dist)
                pts = [(x1 + (x2 - x1) * t / 24, by + 2 + h * math.sin(math.pi * t / 24)) for t in range(25)]
                d.line(pts, fill=col, width=2)
                arrow(d, pts[-2][0], pts[-2][1], x2, by + 1, col, head=6)

        for n in self.nodes:
            st = self.hl.get(n.slot, "special" if n.clone else None)
            fill, border, ink = STYLES[st]
            d.rounded_rectangle([cx(n.slot) - box / 2, top, cx(n.slot) + box / 2, by], 8, fill=fill, outline=border,
                                width=2 if n.slot in self.hl else 1)
            s = f"{n.val}'" if n.clone else str(n.val)
            centered(d, cx(n.slot), cy, s, fit_font(s, "mono", 15, box - 4), ink)


SPEC = [(7, None), (13, 0), (11, 4), (10, 2), (1, 0)]


def build():
    nodes = [Node(v, 2 * i) for i, (v, _) in enumerate(SPEC)]
    for i, (_, r) in enumerate(SPEC):
        nodes[i].next = nodes[i + 1] if i + 1 < len(nodes) else None
        nodes[i].random = nodes[r] if r is not None else None
    return nodes


def compact(nodes):
    """Re-slot a separate list 0, 1, 2, ... for drawing on its own row."""
    out, k = [], 0
    copies = {}
    cur = nodes
    while cur is not None:
        c = Node(cur.val, k, cur.clone)
        copies[id(cur)] = c
        out.append((cur, c))
        cur = cur.next
        k += 1
    for old, c in out:
        c.next = copies.get(id(old.next)) if old.next else None
        c.random = copies.get(id(old.random)) if old.random else None
    return [c for _, c in out]


def all_nodes(head):
    out, n = [], head
    while n is not None:
        out.append(n)
        n = n.next
    return out


def frames():
    orig = build()
    head = orig[0]
    legend = Text("orange = Random of an original, green = Random of a clone (X')")

    yield Frame(
        RandList(compact(head), 5, label="original: arrows above = Next, below = Random"),
        Text("[(7,nil) (13,0) (11,4) (10,2) (1,0)]   (value, index of Random)"),
        speaker="Interviewer",
        say="Each node has Next and a Random pointer to any node or nil. Build a deep copy: new nodes only, same structure. "
            "Leave the original unchanged.",
    )
    yield Frame(
        RandList(compact(head), 5, label="original: arrows above = Next, below = Random"),
        Text("pass 1: clone[old] = new node", "pass 2: new.Next = clone[old.Next], new.Random = clone[old.Random]",
             "O(n) time, O(n) map"),
        say="The hard part is Random: it can point ahead to a node I haven't copied yet. An old-to-new map in two passes "
            "solves it in O(n) space. I can store that map inside the list instead.",
    )

    # 1. Weave: A -> A' -> B -> B' ...
    n = head
    first = True
    while n is not None:
        c = Node(n.val, n.slot + 1, clone=True)
        c.next = n.next
        n.next = c
        if first:
            yield Frame(
                RandList(all_nodes(head), 10, label="step 1: insert each clone right after its original", hl={c.slot: "special"}),
                Text(f"{n.val}.Next = new {n.val}'  and  {n.val}'.Next = old {n.val}.Next"),
                say=f"Step 1: put a clone right after each original. {n.val}' goes between {n.val} and its old Next. "
                    "Now X.Next is X's clone, which is my map.",
            )
            first = False
        n = n.next.next
    yield Frame(
        RandList(all_nodes(head), 10, label="step 1 done: 7, 7', 13, 13', ... in one list"),
        Text("every original X is followed by its clone X'", "clones have no Random yet"),
        say="Same for every node. The list is now twice as long, with each clone right after its original.",
    )

    # 2. Random: clone of X.Random is X.Random.Next.
    n = head
    while n is not None:
        if n.random is not None:
            n.next.random = n.random.next
        n = n.next.next
    # Show one worked case, then the result.
    eleven = next(x for x in all_nodes(head) if x.val == 11 and not x.clone)
    view_nodes = all_nodes(head)
    partial = []
    for x in view_nodes:
        y = Node(x.val, x.slot, x.clone)
        partial.append(y)
    by_slot = {y.slot: y for y in partial}
    for x, y in zip(view_nodes, partial):
        y.next = by_slot[x.next.slot] if x.next else None
        if x.random is not None and (not x.clone or x.slot == eleven.slot + 1):
            y.random = by_slot[x.random.slot]
    yield Frame(
        RandList(partial, 10, label="step 2: X'.Random = X.Random.Next",
                 hl={eleven.slot: "compare", eleven.slot + 1: "match", eleven.random.slot: "compare",
                     eleven.random.slot + 1: "match"}),
        legend,
        Text(f"11.Random = 1   →   11'.Random = 1.Next = 1'"),
        say="Step 2: 11's Random is 1. The clone of 1 is 1.Next, so 11'.Random = 11.Random.Next = 1'. No map needed.",
    )
    yield Frame(
        RandList(view_nodes, 10, label="step 2 done: every clone's Random points at a clone"),
        legend,
        Text("7 has Random nil, so 7' keeps nil"),
        say="I do that for every original that has a Random. Each green arrow lands on a clone, never on an original.",
    )

    # 3. Unweave, restoring the originals.
    new_head = head.next
    n = head
    while n is not None:
        c = n.next
        n.next = c.next
        if c.next is not None:
            c.next = c.next.next
        n = n.next

    orig_nodes = all_nodes(head)
    copy_nodes = all_nodes(new_head)
    idx = {id(x): i for i, x in enumerate(copy_nodes)}
    got = [(x.val, idx[id(x.random)] if x.random else None) for x in copy_nodes]
    assert got == SPEC
    assert not any(id(x) in {id(o) for o in orig_nodes} for x in copy_nodes)
    assert [(o.val) for o in orig_nodes] == [v for v, _ in SPEC]
    assert all(o.next is None or not o.next.clone for o in orig_nodes)

    yield Frame(
        RandList(compact(head), 5, label="original (restored)"),
        RandList(compact(new_head), 5, label="copy (all new nodes)"),
        Text("n.Next = c.Next;  c.Next = c.Next.Next"),
        say="Step 3: split them. Each original skips its clone, and each clone skips the next original. "
            "The original list is back to how it was.",
    )
    yield Frame(
        RandList(compact(head), 5, label="original (restored)"),
        RandList(compact(new_head), 5, label="copy (all new nodes)"),
        Text(("copy matches [(7,nil) (13,0) (11,4) (10,2) (1,0)]", "match"), "O(n) time, O(1) extra space (besides the copy)"),
        say="I return the head of the copy. Three linear passes, and the only memory is the copy itself: O(n) time, O(1) extra space.",
    )
