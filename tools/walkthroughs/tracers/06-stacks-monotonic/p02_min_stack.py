from walkthrough import Array, Frame, Row, Text

TITLE = "Min Stack"
PATTERN = "Stack + running min"


class MinStack:
    """Port of the Go MinStack: each entry stores (val, min of itself and below)."""

    def __init__(self):
        self.items = []

    def push(self, v):
        m = v
        if self.items:
            m = min(m, self.items[-1][1])
        self.items.append((v, m))

    def pop(self):
        self.items.pop()

    def get_min(self):
        return self.items[-1][1]


def view(s, hl_style=None):
    vals = [v for v, _ in s.items]
    mins = [m for _, m in s.items]
    top = len(vals) - 1
    hl = {top: hl_style} if hl_style and vals else None
    ptr = {"top": top} if vals else {}
    return Row(
        Array(vals, label="val", hl=hl, pointers=ptr, show_index=False, cell=52),
        Array(mins, label="min (of this and below)", hl=hl, pointers=ptr, show_index=False, cell=52),
    )


def frames():
    ops = [("Push", 5), ("Push", 2), ("Push", 7), ("GetMin", 2), ("Pop", None), ("GetMin", 2),
           ("Pop", None), ("GetMin", 5)]
    log = "  ".join(f"{op}({a})" if op == "Push" else f"{op}()" for op, a in ops)

    s = MinStack()
    yield Frame(
        view(s),
        Text(log, size=13),
        speaker="Interviewer",
        say="Design a stack with Push, Pop, Top and GetMin, where GetMin returns the smallest element. "
            "Every operation must be O(1).",
    )
    yield Frame(
        view(s),
        Text("plain stack + scan for the min on GetMin → O(n)",
             ("one 'min' variable breaks on Pop: what was the previous min?", "reject")),
        say="Scanning for the min is O(n). Keeping one min variable fails too: after popping the min, "
            "I'd have no idea what the previous min was.",
    )

    seen_getmin = 0
    for op, arg in ops:
        if op == "Push":
            below = s.items[-1][1] if s.items else None
            s.push(arg)
            m = s.items[-1][1]
            calc = f"Push({arg}): min = {arg}" if below is None else f"Push({arg}): min = min({arg}, {below}) = {m}"
            yield Frame(
                view(s, "active"),
                Text(calc),
                say=(f"Each entry stores its value and the min of itself and everything below. "
                     f"The stack was empty, so 5's min is 5." if below is None else
                     f"Push {arg}: its min is the smaller of {arg} and the min below it, {below}. So {m}."),
            )
        elif op == "Pop":
            gone = s.items[-1]
            s.pop()
            yield Frame(
                view(s, "special"),
                Text(f"Pop(): remove ({gone[0]}, min {gone[1]})", f"new top remembers min = {s.items[-1][1]}"),
                say=f"Pop removes {gone[0]}. The entry below still carries its own min, "
                    f"{s.items[-1][1]}, so the old minimum comes back for free.",
            )
        else:
            got = s.get_min()
            assert got == arg, (got, arg)
            seen_getmin += 1
            say = {1: "GetMin just reads the min stored in the top entry: 2. No scan.",
                   2: "GetMin is still 2: the 2 is still in the stack, and the top entry says so.",
                   3: "GetMin is 5 now. The 2 is gone and the remaining entry already knew its own min."}[seen_getmin]
            yield Frame(
                view(s, "match"),
                Text(f"GetMin() → {got}", ("read top.min", "match")),
                say=say,
            )

    yield Frame(
        view(s),
        Text("GetMin results: 2, 2, 5", ("every operation O(1) · space O(n)", "match")),
        say="GetMin returned 2, 2, 5 as expected. Every operation touches only the top, so all are O(1). "
            "The extra min field costs O(n) space.",
    )
