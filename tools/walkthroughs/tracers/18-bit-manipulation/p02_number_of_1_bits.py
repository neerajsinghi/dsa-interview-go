from walkthrough import Frame, Grid, Text

TITLE = "Number of 1 Bits"
PATTERN = "Bits · Kernighan"

W = 8
COLS = list(range(W - 1, -1, -1))


def bits(x):
    x &= (1 << W) - 1
    return [(x >> k) & 1 for k in COLS]


def frames():
    n0 = 11
    expected = 3

    yield Frame(
        Grid([bits(n0)], row_labels=[f"n = {n0}"], col_labels=COLS),
        speaker="Interviewer",
        say="Count how many bits are 1 in an unsigned 32-bit integer. For example, 11.",
    )
    yield Frame(
        Grid([bits(n0)], row_labels=[f"n = {n0}"], col_labels=COLS, hl={(0, W - 1): "active"}),
        Text("shift loop: add n & 1, then n >>= 1, until n == 0", "one step per bit position (up to 32)"),
        say="The simple way checks the low bit and shifts right, once per bit position. "
            "That's up to 32 steps, even when only a few bits are set.",
    )

    n, count = n0, 0
    while n != 0:
        m = n - 1
        new = n & m
        low = (n & -n).bit_length() - 1
        j = COLS.index(low)
        hl = {(0, j): "reject", (2, j): "reject"}
        for k in range(j + 1, W):
            hl[(1, k)] = "compare"
        hl[(1, j)] = "compare"
        count += 1
        say = (f"n - 1 flips the lowest 1 (bit {low}) to 0 and the zeros below it to 1. "
               f"ANDing with n clears just that bit. count = {count}.")
        if low == 0 and count == 1:
            say = (f"n - 1 turns off the lowest 1, bit 0 here (and would turn any zeros below it into 1s). "
                   f"ANDing with n clears just that bit. count = {count}.")
        if count == 2:
            say = f"Same again: the lowest 1 is now bit {low}. n & (n-1) clears it. count = {count}."
        if new == 0:
            say = f"The last 1 is bit {low}. Clearing it leaves 0, so the loop stops. count = {count}."
        yield Frame(
            Grid([bits(n), bits(m), bits(new)], row_labels=[f"n = {n}", f"n-1 = {m}", f"n & (n-1) = {new}"],
                 col_labels=COLS, hl=hl),
            Text(f"count = {count}"),
            say=say,
        )
        n = new

    assert count == expected
    yield Frame(
        Grid([bits(n0)], row_labels=[f"n = {n0}"], col_labels=COLS,
             hl={(0, j): "match" for j, b in enumerate(bits(n0)) if b}),
        Text((f"answer: {count}", "match"), "O(k) time for k set bits (≤ 32), O(1) space"),
        say=f"The answer is {count}. The loop ran once per set bit, not once per position. "
            "That's O(k) for k set bits, at most 32, and O(1) space.",
    )
