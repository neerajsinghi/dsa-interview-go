from walkthrough import Frame, Grid, Text

TITLE = "Power of Two"
PATTERN = "Bits · n & (n-1)"

W = 8
COLS = list(range(W - 1, -1, -1))


def bits(x):
    x &= (1 << W) - 1  # two's complement view for negatives
    return [(x >> k) & 1 for k in COLS]


def is_power_of_two(n):
    return n > 0 and n & (n - 1) == 0


def frames():
    # Every example from the doc comment, in order.
    cases = [(1, True), (16, True), (6, False), (0, False), (-8, False)]

    yield Frame(
        Grid([bits(16), bits(6)], row_labels=["16", "6"], col_labels=COLS),
        speaker="Interviewer",
        say="Decide whether n equals 2 to the power k for some k ≥ 0. n is a Go int, so it can be zero or negative.",
    )
    yield Frame(
        Grid([bits(16), bits(6)], row_labels=["16", "6"], col_labels=COLS,
             hl={(0, 3): "match", (1, 5): "reject", (1, 6): "reject"}),
        Text("loop: while n is even, n /= 2; check n == 1", "O(log n) steps"),
        say="I could divide by 2 while n is even and check I end at 1. That's O(log n). "
            "But a power of two has exactly one 1 bit, and I can test that directly.",
    )

    for n, want in cases:
        got = is_power_of_two(n)
        assert got == want
        m = n - 1
        new = n & m
        if n <= 0:
            rows = [bits(n)]
            labels = [f"n = {n}"]
            hl = {(0, j): "reject" for j, b in enumerate(bits(n)) if b}
            if n == 0:
                say = "0 has no 1 bits at all. 0 & -1 is 0, which would wrongly pass, so the n > 0 guard rejects it."
                lines = ["n > 0?  no", ("false", "reject")]
            else:
                say = ("-8 is negative, and its two's complement has several 1s. No negative number is a power of two, "
                       "so n > 0 rejects it first.")
                lines = ["n > 0?  no  (two's complement shown in 8 bits)", ("false", "reject")]
        else:
            rows = [bits(n), bits(m), bits(new)]
            labels = [f"n = {n}", f"n-1 = {m}", f"n & (n-1) = {new}"]
            hl = {}
            low = (n & -n).bit_length() - 1
            hl[(0, COLS.index(low))] = "active"
            for j, b in enumerate(bits(new)):
                if new == 0:
                    hl[(2, j)] = "match"
                elif b:
                    hl[(2, j)] = "reject"
            if got:
                say = (f"{n} has one 1 bit. n - 1 turns it off and sets the bits below, "
                       f"so n & (n-1) is 0. {n} is a power of two.")
                lines = [f"n > 0 and n & (n-1) == 0", ("true", "match")]
            else:
                say = (f"{n} has two 1 bits. n & (n-1) clears only the lowest one, so a 1 is left. "
                       f"{n} is not a power of two.")
                lines = [f"n & (n-1) = {new} ≠ 0", ("false", "reject")]
        yield Frame(Grid(rows, row_labels=labels, col_labels=COLS, hl=hl), Text(*lines), say=say)

    yield Frame(
        Grid([bits(16), bits(15), bits(0)], row_labels=["n = 16", "n-1 = 15", "n & (n-1)"], col_labels=COLS,
             hl={(2, j): "match" for j in range(W)}),
        Text(("return n > 0 && n&(n-1) == 0", "match"), "O(1) time, O(1) space"),
        say="So the whole test is n > 0 and n & (n-1) == 0. It's one AND and two compares: "
            "O(1) time and O(1) space.",
    )
