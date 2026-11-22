from walkthrough import Frame, LinkedList, Text

TITLE = "Add Two Numbers"
PATTERN = "Dummy head + carry"


def add(A, B):
    """Port of AddTwoNumbers on value lists. Yields one dict per loop iteration."""
    i = j = 0
    carry = 0
    out = []
    while i < len(A) or j < len(B) or carry > 0:
        s = carry
        parts = [str(carry)]
        ai, bj = i, j
        if i < len(A):
            s += A[i]
            parts.append(str(A[i]))
            i += 1
        if j < len(B):
            s += B[j]
            parts.append(str(B[j]))
            j += 1
        c_in = carry
        carry = s // 10
        out.append(s % 10)
        yield dict(ai=ai, bj=bj, s=s, parts=parts, c_in=c_in, carry=carry, out=list(out),
                   a_ok=ai < len(A), b_ok=bj < len(B))


def views(A, B, st):
    ha = {st["ai"]: "active"} if st["a_ok"] else {}
    hb = {st["bj"]: "active"} if st["b_ok"] else {}
    out = ["D"] + st["out"]
    ho = {0: "special", len(out) - 1: "match"}
    return (LinkedList(A, label="a (least-significant digit first)", hl=ha, pointers={"a": st["ai"]}),
            LinkedList(B, label="b", hl=hb, pointers={"b": st["bj"]}),
            LinkedList(out, label="result (D = dummy)", hl=ho, pointers={"tail": len(out) - 1}))


def frames():
    A, B = [2, 4, 3], [5, 6, 4]

    yield Frame(
        LinkedList(A, label="a = 342"), LinkedList(B, label="b = 465"),
        Text("digits are stored in reverse: 342 is 2 → 4 → 3", "expected: 7 → 0 → 8  (807)"),
        speaker="Interviewer",
        say="Two numbers are stored as digit lists, least significant digit first. Return their sum in the same form.",
    )
    yield Frame(
        LinkedList(A, label="a = 342"), LinkedList(B, label="b = 465"),
        Text("convert to int, add, convert back", "→ breaks once a number has more than ~19 digits"),
        say="Turning each list into an int and back fails for long numbers, since int64 overflows past 19 digits. "
            "I'll add digit by digit, like on paper.",
    )
    yield Frame(
        LinkedList(A, label="a (least-significant digit first)", pointers={"a": 0}),
        LinkedList(B, label="b", pointers={"b": 0}),
        LinkedList(["D"], label="result (D = dummy)", hl={0: "special"}, pointers={"tail": 0}),
        Text("carry = 0", "loop while a != nil || b != nil || carry > 0"),
        say="The lowest digits come first, which is exactly the order I add in. I keep a carry and a dummy head for the result.",
    )

    last = None
    for k, st in enumerate(add(A, B)):
        last = st
        expr = " + ".join(st["parts"]) + f" = {st['s']}"
        yield Frame(
            *views(A, B, st),
            Text(f"sum = carry + a + b = {expr}", f"digit = {st['s']} % 10 = {st['s'] % 10},  carry = {st['s']} / 10 = {st['carry']}"),
            say=(f"{st['c_in']} + {A[st['ai']]} + {B[st['bj']]} = {st['s']}. I append {st['s'] % 10} and carry {st['carry']}"
                 + (" into the next column." if st["carry"] else ".")),
        )
    assert last["out"] == [7, 0, 8]

    # Second doc example: the final carry adds a digit.
    A2, B2 = [9, 9], [1]
    steps = list(add(A2, B2))
    lines = []
    for st in steps:
        lines.append(" + ".join(st["parts"]) + f" = {st['s']} → digit {st['s'] % 10}, carry {st['carry']}")
    assert steps[-1]["out"] == [0, 0, 1]
    yield Frame(
        LinkedList(A2, label="a = 99"), LinkedList(B2, label="b = 1"),
        Text("lists of different lengths, carry left at the end"),
        speaker="Interviewer",
        say="What about 99 + 1? The lists have different lengths, and there's a carry left at the end.",
    )
    yield Frame(
        LinkedList(A2, label="a = 99", hl={0: "visited", 1: "visited"}, pointers={"a": 2}),
        LinkedList(B2, label="b = 1", hl={0: "visited"}, pointers={"b": 1}),
        LinkedList(["D"] + steps[-1]["out"], label="result", hl={0: "special", 3: "match"}),
        Text(*lines, ("0 → 0 → 1  (100);  O(max(m, n)) time and output", "match")),
        say="A missing digit counts as 0. When both lists end with carry 1, the 'carry > 0' condition runs one more "
            "time and adds the final 1. O(max(m, n)) time.",
    )
