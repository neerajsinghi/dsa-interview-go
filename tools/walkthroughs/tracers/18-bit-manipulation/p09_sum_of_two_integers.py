from walkthrough import Frame, Grid, Text

TITLE = "Sum of Two Integers"
PATTERN = "Bits · half adder"

W = 8
COLS = list(range(W - 1, -1, -1))


def bits(x):
    x &= (1 << W) - 1
    return [(x >> k) & 1 for k in COLS]


def i32(x):
    """Wrap to Go's int32: two's complement, overflow wraps."""
    x &= 0xFFFFFFFF
    return x - (1 << 32) if x >= 1 << 31 else x


def get_sum(a, b, trace=None):
    while b != 0:
        carry = i32((a & b) << 1)
        if trace is not None:
            trace.append((a, b, i32(a ^ b), carry))
        a = i32(a ^ b)
        b = carry
    return a


def frames():
    # The doc comment's 5 + 3 diagram. 1 + 2 (the first example) finishes in one round.
    a0, b0 = 5, 3
    expected = 8
    assert get_sum(1, 2) == 3 and get_sum(-2, 3) == 1 and get_sum(2**31 - 1, 1) == -2**31

    yield Frame(
        Grid([bits(a0), bits(b0)], row_labels=[f"a = {a0}", f"b = {b0}"], col_labels=COLS),
        speaker="Interviewer",
        say="Add two signed 32-bit integers without using + or -. Let's trace 5 + 3.",
    )
    yield Frame(
        Grid([bits(a0), bits(b0)], row_labels=[f"a = {a0}", f"b = {b0}"], col_labels=COLS),
        Text("brute: increment a one at a time, b times → O(|b|)", "and increment itself needs +"),
        say="Counting up one at a time is O(|b|), and still needs an increment. "
            "I'll do what the hardware does: add column by column with logic gates.",
    )
    yield Frame(
        Grid([[0, 1, 1, 0], [0, 0, 0, 1]], row_labels=["a ^ b (sum)", "a & b (carry)"],
             col_labels=["0+0", "0+1", "1+0", "1+1"], cell=60, hl={(0, 1): "match", (0, 2): "match", (1, 3): "compare"}),
        Text("a ^ b = sum without carries", "(a & b) << 1 = carries, moved one column left"),
        say="For one column, XOR gives the sum bit and AND says whether it carries. "
            "So a ^ b adds without carrying, and (a & b) << 1 is the carry. Repeat until no carry is left.",
    )

    trace = []
    result = get_sum(a0, b0, trace)
    for r, (a, b, s, carry) in enumerate(trace, 1):
        hl = {}
        for j, bit in enumerate(bits(carry)):
            if bit:
                hl[(3, j)] = "compare"
        for j, bit in enumerate(bits(s)):
            if bit:
                hl[(2, j)] = "match"
        both = [k for k in range(W) if (a & b) >> k & 1]
        where = ", ".join(map(str, both))
        if carry == 0:
            say = f"Round {r}: no bit is set in both, so the carry is 0. a ^ b = {s} is the answer."
        elif r == 1:
            say = (f"Round 1: a ^ b = {s}. Bit {where} is set in both, so it carries one column left: "
                   f"carry = {carry}. Now I need {s} + {carry}.")
        else:
            say = f"Round {r}: bit {where} is set in both, so it carries again. Now a = {s} and b = {carry}."
        yield Frame(
            Grid([bits(a), bits(b), bits(s), bits(carry)],
                 row_labels=[f"a = {a}", f"b = {b}", f"a ^ b = {s}", f"(a&b)<<1 = {carry}"], col_labels=COLS, hl=hl),
            Text(f"round {r}: a = {s}, b = {carry}" + ("  → stop" if carry == 0 else "")),
            say=say,
        )

    assert result == expected
    yield Frame(
        Grid([bits(result)], row_labels=[f"a = {result}"], col_labels=COLS, hl={(0, COLS.index(3)): "match"}),
        Text((f"answer: {a0} + {b0} = {result}", "match"), "≤ 32 rounds for int32, O(1) space"),
        say=f"The answer is {result}. Each round pushes the carry one column left, so int32 needs at most 32 rounds. "
            "O(1) time and space.",
    )
