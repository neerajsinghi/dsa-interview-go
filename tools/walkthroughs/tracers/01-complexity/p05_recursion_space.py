from walkthrough import BORDER, MUTED, STYLES, Array, Frame, Row, Text, Widget, font

TITLE = "Space Complexity of Recursion"
PATTERN = "Call stack"


class CallStack(Widget):
    """A vertical call stack: frames[0] at the bottom, the newest call on top.

    Always reserves room for `slots` frames so the picture doesn't jump.
    """

    def __init__(self, frames, label="call stack", slots=5, hl=None, notes=None):
        self.frames = list(frames)
        self.label = label
        self.slots = slots
        self.hl = hl or {}
        self.notes = notes or {}

    def height(self, width):
        return 20 + self.slots * 32 + 6

    def draw(self, d, x, y, width):
        d.text((x, y), self.label, font=font("bold", 14), fill=MUTED)
        y += 20
        bottom = y + self.slots * 32
        f = font("mono", 14)
        for k in range(self.slots):
            y1 = bottom - (k + 1) * 32
            if k < len(self.frames):
                fill, border, ink = STYLES[self.hl.get(k)]
                d.rounded_rectangle([x, y1 + 2, x + width, y1 + 28], 6, fill=fill, outline=border,
                                    width=2 if k in self.hl else 1)
                d.text((x + 10, y1 + 15), self.frames[k], font=f, fill=ink, anchor="lm")
                if k in self.notes:
                    d.text((x + width - 10, y1 + 15), self.notes[k], font=font("bold", 13), fill=ink, anchor="rm")
            else:
                d.rounded_rectangle([x, y1 + 2, x + width, y1 + 28], 6, outline="#eaeef2", width=1)
        d.line([(x - 4, bottom + 2), (x + width + 4, bottom + 2)], fill=BORDER, width=2)


def fmt(xs):
    return "[" + ",".join(map(str, xs)) + "]"


def trace_recursive(nums, events, stack):
    """SumRecursive, recording push/return events with the stack contents."""
    stack.append(f"SumRecursive({fmt(nums)})")
    events.append(("push", list(stack), None))
    if len(nums) == 0:
        res = 0
    else:
        res = nums[0] + trace_recursive(nums[1:], events, stack)
    events.append(("ret", list(stack), res))
    stack.pop()
    return res


def trace_halves(nums, path, best):
    """SumHalves, remembering the deepest chain of pending calls in best."""
    path = path + [f"SumHalves({fmt(nums)})"]
    if len(path) > len(best):
        best[:] = path
    if len(nums) == 0:
        return 0
    if len(nums) == 1:
        return nums[0]
    mid = len(nums) // 2
    return trace_halves(nums[:mid], path, best) + trace_halves(nums[mid:], path, best)


def frames():
    nums = [1, 2, 3, 4]

    yield Frame(
        Array(nums, label="nums"),
        Text("expected: 10"),
        speaker="Interviewer",
        say="Sum the slice recursively and with a loop. Both are O(n) time. Where do they differ?",
    )

    events = []
    rec = trace_recursive(nums, events, [])
    pushes = [e for e in events if e[0] == "push"]
    rets = [e for e in events if e[0] == "ret"]
    depth = max(len(e[1]) for e in events)

    yield Frame(
        Array(nums, label="nums"),
        CallStack([], slots=depth),
        Text("SumRecursive(nums) = nums[0] + SumRecursive(nums[1:])"),
        say="The recursive version reads nicely, but every call waits for the next one to return. "
            "Each waiting call keeps a stack frame alive.",
    )

    for k, (_, st, _) in enumerate(pushes):
        if k in (0, 1) or k == len(pushes) - 1:
            if k == len(pushes) - 1:
                say = f"The empty slice is the base case. Now {len(st)} frames are alive at once: "\
                      "depth n + 1, so O(n) space even though nothing was allocated on the heap."
            elif k == 0:
                say = "The first call can't add 1 until it knows the sum of the rest, so it pushes a new frame."
            else:
                say = "Same again: one more frame per element. nums[1:] shares the backing array, so slicing is O(1)."
            yield Frame(
                Array(nums, label="nums", hl={j: "active" for j in range(k, len(nums))}),
                CallStack(st, slots=depth, hl={len(st) - 1: "active"}),
                Text(f"depth = {len(st)}"),
                say=say,
            )
        elif k == 2:
            yield Frame(
                Array(nums, label="nums", hl={j: "active" for j in range(k, len(nums))}),
                CallStack(st, slots=depth, hl={len(st) - 1: "active"}),
                Text(f"depth = {len(st)}"),
                say="It keeps going, one frame per remaining element.",
            )

    # Unwind: show the stack with the value each frame returns.
    notes, unwind, below = {}, [], None
    for _, st, res in rets:
        notes[len(st) - 1] = f"returns {res}"
        x = nums[len(st) - 1] if len(st) - 1 < len(nums) else None
        unwind.append(str(res) if x is None else f"{x}+{below}={res}")
        below = res
    yield Frame(
        Array(nums, label="nums"),
        CallStack(pushes[-1][1], slots=depth, hl={k: "match" for k in range(depth)}, notes=notes),
        Text("unwind: " + " → ".join(unwind)),
        say="Then the frames pop one by one, each adding its element. The top returns 0, "
            f"and the bottom finally returns {rec}.",
    )

    best = []
    halves = trace_halves(nums, [], best)
    yield Frame(
        Array(nums, label="nums", hl={0: "special", 1: "special"}),
        CallStack(best, label="deepest path of SumHalves", slots=depth, hl={len(best) - 1: "special"}),
        Text(f"max depth = {len(best)} = log₂ 4 + 1", "still n leaves → O(n) time"),
        say="SumHalves splits in two like merge sort. The work is still O(n), but the depth is only "
            "about log₂ n, so the stack is O(log n).",
    )

    total, running = 0, ["0"]
    for x in nums:
        total += x
        running.append(str(total))
    yield Frame(
        Array(nums, label="nums", hl={j: "match" for j in range(len(nums))}),
        CallStack(["SumIterative([1,2,3,4])"], slots=depth, hl={0: "match"}),
        Text("total = " + " → ".join(running), "one frame, one accumulator"),
        say="The loop uses one frame and one accumulator however long the slice is. That's O(1) extra space.",
    )

    assert rec == halves == total == 10 and depth == len(nums) + 1
    yield Frame(
        Text(("answer: 10", "match"),
             "SumRecursive  O(n) time  O(n) stack",
             "SumHalves     O(n) time  O(log n) stack",
             "SumIterative  O(n) time  O(1) space"),
        say="All three return 10 in O(n) time. They differ in space: n frames, log n frames, or one.",
    )
