from walkthrough import Array, Frame, Grid, Text

TITLE = "Min Stack / Max Stack"
PATTERN = "Value + running agg"


class AggStack:
    """Port of MinStack / MaxStack: each frame stores (val, agg)."""

    def __init__(self, pick):
        self.frames = []
        self.pick = pick

    def push(self, v):
        agg = v if not self.frames else self.pick(v, self.frames[-1][1])
        self.frames.append((v, agg))

    def pop(self):
        return self.frames.pop()[0]

    def top(self):
        return self.frames[-1][0]

    def agg(self):
        return self.frames[-1][1]

    def pop_max(self):
        target = self.agg()
        buffer = []
        while True:
            top = self.pop()
            if top == target:
                break
            buffer.append(top)
        for v in reversed(buffer):
            self.push(v)
        return target, buffer


def grid(st, aggname, hl_top=None):
    n = len(st.frames)
    if n == 0:
        return Grid([["", ""]], label="frames (bottom to top)", row_labels=["val"], cell=46)
    rows = [[f[0] for f in st.frames], [f[1] for f in st.frames]]
    hl = {}
    if hl_top:
        hl = {(0, n - 1): hl_top, (1, n - 1): hl_top}
    return Grid(rows, label="frames (bottom to top)", row_labels=["val", aggname], hl=hl,
                col_labels=["" for _ in range(n - 1)] + ["top"], cell=46)


def frames():
    s = AggStack(min)

    yield Frame(
        grid(s, "min"),
        Text("MinStack: Push, Pop, Top, GetMin — all O(1)", "MaxStack: also PopMax"),
        speaker="Interviewer",
        say="Design a stack with Push, Pop, Top and GetMin, all in O(1). Then a MaxStack that can also "
            "remove its largest element.",
    )
    yield Frame(
        grid(s, "min"),
        Text("naive: scan the stack on GetMin → O(n)", "or keep one min variable → lost after Pop"),
        say="Scanning for the minimum is O(n). A single min variable breaks on Pop: once the minimum leaves, "
            "I don't know the previous one. So I store the minimum with every frame.",
    )

    def op(call, ret, say, hl="active"):
        return Frame(grid(s, "min", hl), Text((f"{call} → {ret}" if ret is not None else call, hl)), say=say)

    s.push(-2)
    yield op("Push(-2)", None, "Push -2. The stack was empty, so this frame's min is -2 itself.")
    s.push(0)
    yield op("Push(0)", None, "Push 0. Its min is min(0, -2) = -2, copied from the frame below.")
    s.push(-3)
    yield op("Push(-3)", None, "Push -3. min(-3, -2) = -3, so this frame records a new minimum.")
    r = s.agg()
    assert r == -3
    yield op("GetMin()", r, "GetMin just reads the top frame's min: -3. No scan.", "match")
    r = s.pop()
    yield op("Pop()", r, "Pop removes -3 with its frame. The frame below still says min = -2, "
                         "so the old minimum comes back for free.")
    r = s.top()
    assert r == 0
    yield op("Top()", r, "Top reads the top value: 0.", "match")
    r = s.agg()
    assert r == -2
    yield op("GetMin()", r, "GetMin reads -2 from the top frame. Every MinStack operation was O(1).", "match")

    m = AggStack(max)
    for v in (5, 1, 5):
        m.push(v)
    yield Frame(
        grid(m, "max", "active"),
        Text("MaxStack: Push(5) Push(1) Push(5)", "PeekMax() → 5"),
        say="MaxStack uses the same frames with a running max. PeekMax is O(1). PopMax is the twist: "
            "the max may be buried under smaller values.",
    )
    target, buffer = m.pop_max()
    assert target == 5 and buffer == []
    yield Frame(
        grid(m, "max", "active"),
        Text(("PopMax() → 5", "match"), f"popped down to the top-most 5; buffered {buffer or 'nothing'}"),
        say="PopMax pops into a buffer until it hits the max, then pushes the buffer back to recompute each "
            "frame's max. Here the top-most 5 is on top, so the buffer is empty.",
    )
    r = m.top()
    assert r == 1
    yield Frame(
        grid(m, "max", "match"),
        Text(("Top() → 1", "match"), "MinStack: all O(1) · MaxStack.PopMax O(n) · space O(n)"),
        say="Top is now 1, as expected. Everything is O(1) except PopMax, which is O(n) in the worst case. "
            "A heap with lazy deletion would get O(log n).",
    )
