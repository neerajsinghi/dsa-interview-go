from walkthrough import Array, Frame, Grid, Row, Text

TITLE = "Single Number II"
PATTERN = "Bits · state machine"

W = 8
COLS = list(range(W - 1, -1, -1))


def bits(x):
    x &= (1 << W) - 1
    return [(x >> k) & 1 for k in COLS]


def states(hl_row=None):
    rows = [[0, 0], [1, 0], [0, 1], [0, 0]]
    hl = {(hl_row, 0): "active", (hl_row, 1): "active"} if hl_row is not None else {}
    return Grid(rows, label="per-bit counter", row_labels=["seen 0×", "seen 1×", "seen 2×", "seen 3×"],
                col_labels=["ones", "twos"], hl=hl, cell=44)


def frames():
    nums = [2, 2, 3, 2]
    expected = 3

    yield Frame(
        Array(nums, label="nums"),
        speaker="Interviewer",
        say="Every value appears three times except one, which appears once. Find it in O(n) time and O(1) space.",
    )
    yield Frame(
        Array(nums, label="nums"),
        Text("map of counts → O(n) space", "count each of 64 bit positions mod 3 → O(64·n)"),
        say="A count map needs O(n) space. XOR alone fails, because three copies don't cancel. "
            "Counting each bit position mod 3 works, but takes 64 passes.",
    )
    yield Frame(
        Array(nums, label="nums"),
        Row(states(), Text("ones = (ones ^ x) &^ twos", "twos = (twos ^ x) &^ ones", size=14),
            weights=[1, 1.3]),
        say="Instead each bit position keeps a counter mod 3 in two words, ones and twos. "
            "These two lines step every bit's counter at once: 0, 1, 2, then back to 0.",
    )

    ones, twos = 0, 0
    seen = {}
    for i, x in enumerate(nums):
        ones = (ones ^ x) & ~twos
        twos = (twos ^ x) & ~ones
        seen[x] = seen.get(x, 0) + 1
        c = seen[x]
        hl = {}
        for j, b in enumerate(bits(ones)):
            if b:
                hl[(1, j)] = "match"
        for j, b in enumerate(bits(twos)):
            if b:
                hl[(2, j)] = "compare"
        if x == 2 and c == 1:
            say = "First 2: bit 1 moves to seen-once, so ones has bit 1 set."
        elif x == 2 and c == 2:
            say = "Second 2: bit 1 moves from ones to twos. It has been seen twice."
        elif x == 3:
            say = ("3 sets bits 0 and 1. Bit 0 is new, so it goes into ones. Bit 1 was in twos, "
                   "so this is its third time and it resets to 0.")
        else:
            say = ("Third 2: bit 1 has been seen once since its last reset (by the 3), so it moves to ones. "
                   "Bit 0 stays in ones.")
        yield Frame(
            Array(nums, label="nums", hl={k: "visited" for k in range(i)} | {i: "active"}, pointers={"i": i}),
            Grid([bits(x), bits(ones), bits(twos)], row_labels=[f"x = {x}", f"ones = {ones}", f"twos = {twos}"],
                 col_labels=COLS, hl=hl),
            say=say,
        )

    assert ones == expected
    yield Frame(
        Array(nums, label="nums", hl={2: "match"}),
        Grid([bits(ones)], row_labels=[f"ones = {ones}"], col_labels=COLS, hl={(0, j): "match" for j in range(W)}),
        Text((f"answer: {ones}", "match"), "O(n) time, O(1) space"),
        say=f"Bits from the triples cycle back to 0. Only the single value's bits stay in ones, so the answer is "
            f"{ones}. One pass, two words: O(n) time, O(1) space.",
    )
