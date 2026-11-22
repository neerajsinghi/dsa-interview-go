from walkthrough import Array, Frame, Text

TITLE = "Peeking Iterator & Flatten Nested List"
PATTERN = "Lookahead + frame stack"


class SliceIterator:
    def __init__(self, items):
        self.items, self.pos = items, 0

    def has_next(self):
        return self.pos < len(self.items)

    def next(self):
        v = self.items[self.pos]
        self.pos += 1
        return v


class PeekingIterator:
    def __init__(self, it):
        self.it, self.peeked, self.cached = it, False, None

    def peek(self):
        if not self.peeked:
            self.cached = self.it.next()
            self.peeked = True
        return self.cached

    def next(self):
        if self.peeked:
            self.peeked = False
            v, self.cached = self.cached, None
            return v
        return self.it.next()

    def has_next(self):
        return self.peeked or self.it.has_next()


class NestedIterator:
    """Nested lists are Python lists; integers are ints."""

    def __init__(self, items):
        self.stack = [[items, 0]]

    def has_next(self):
        while self.stack:
            top = self.stack[-1]
            if top[1] == len(top[0]):
                self.stack.pop()
                continue
            cur = top[0][top[1]]
            if isinstance(cur, int):
                return True
            top[1] += 1
            self.stack.append([cur, 0])
        return False

    def next(self):
        assert self.has_next()
        top = self.stack[-1]
        v = top[0][top[1]]
        top[1] += 1
        return v


def fmt(x):
    if isinstance(x, int):
        return str(x)
    return "[" + ",".join(fmt(y) for y in x) + "]"


def frames():
    base = SliceIterator([1, 2, 3])
    p = PeekingIterator(base)

    def peek_view(call, ret, say, hl="active"):
        hls = {i: "visited" for i in range(base.pos)}
        if p.peeked:
            hls[base.pos - 1] = "special"
        return Frame(
            Array(base.items, label="underlying SliceIterator", hl=hls, pointers={"pos": base.pos}),
            Text(f"peeked = {str(p.peeked).lower()}, cached = {p.cached if p.peeked else '—'}",
                 (f"{call} → {ret}", hl)),
            say=say,
        )

    yield Frame(
        Array([1, 2, 3], label="underlying SliceIterator", pointers={"pos": 0}),
        Text("PeekingIterator: add Peek() to any Iterator", "NestedIterator: flatten [[1,1],2,[1,1]] lazily"),
        speaker="Interviewer",
        say="Two iterators. First, wrap any iterator and add Peek, which returns the next element without "
            "consuming it. Second, flatten a nested list one integer at a time.",
    )
    yield Frame(
        Array([1, 2, 3], label="underlying SliceIterator", pointers={"pos": 0}),
        Text("naive: copy everything into a slice up front", ("O(n) memory; can't stream or stop early", "reject")),
        say="The lazy answer is to copy everything into a slice first. That costs O(n) memory and breaks for "
            "streams. I'll keep just enough state to pause between elements.",
    )

    out = []
    r = p.next(); out.append(r)
    yield peek_view("Next()", r, "Next with nothing cached passes straight through to the underlying iterator: 1.")
    r = p.peek()
    yield peek_view("Peek()", r, "Peek has to read ahead, so it calls the underlying Next and caches 2. "
                                 "The underlying pos has moved, but my caller hasn't consumed anything.", "special")
    r = p.next(); out.append(r)
    yield peek_view("Next()", r, "Next sees the cache is full, returns 2 from it and clears it. "
                                 "It doesn't touch the underlying iterator.")
    r = p.next(); out.append(r)
    hn = p.has_next()
    assert out == [1, 2, 3] and hn is False
    yield peek_view("Next() → 3, HasNext()", "false",
                    "Next returns 3. HasNext is peeked || underlying HasNext, both false now.", "match")

    nested = [[1, 1], 2, [1, 1]]
    it = NestedIterator(nested)

    def cur_hl():
        if not it.stack:
            return {}
        k = it.stack[0][1] - (1 if len(it.stack) > 1 else 0)
        return {k: "active"}

    def stack_view(call, say, hl="active"):
        lines = [f"frame {k}: list {fmt(f[0])}, pos {f[1]}" for k, f in enumerate(it.stack)]
        if not lines:
            lines = ["(stack empty)"]
        return Frame(
            Array([fmt(x) for x in nested], label="top-level list", hl=cur_hl()),
            Text(*lines, (call, hl)),
            say=say,
        )

    got = []
    says = {
        0: "HasNext does the work. The first item is [1,1], a list, so it steps past it in the parent "
           "and pushes a frame for it. Next returns 1; the frame now points at the second 1.",
        1: "Next returns the second 1. The inner list is done, so the following HasNext pops its frame, "
           "and the parent now points at 2.",
        2: "Next returns 2 straight from the parent. HasNext then meets the last [1,1] and pushes a frame for it.",
        3: "Next returns 1 from the new frame. Frames hold a list and a position, so the iterator can pause here.",
        4: "Next returns the last 1. HasNext pops the finished inner frame, then the finished parent.",
    }
    while it.has_next():
        v = it.next()
        got.append(v)
        it.has_next()  # position on the next integer so the stack shown matches HasNext's work
        yield stack_view(f"Next() → {v}   output so far: {got}", says[len(got) - 1])

    assert got == [1, 1, 2, 1, 1]
    yield stack_view("HasNext() → false   output: [1 1 2 1 1]",
                     "Output is 1 1 2 1 1. Each list is pushed and popped once, so it's O(1) amortized per "
                     "element and O(depth) space. Go 1.23's iter.Seq does the same walk recursively.", "match")
