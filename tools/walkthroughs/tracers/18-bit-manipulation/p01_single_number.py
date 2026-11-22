from walkthrough import Array, Frame, Grid, Text

TITLE = "Single Number"
PATTERN = "Bits · XOR"

W = 8
COLS = list(range(W - 1, -1, -1))


def bits(x):
    x &= (1 << W) - 1
    return [(x >> k) & 1 for k in COLS]


def frames():
    nums = [2, 2, 1]
    expected = 1

    yield Frame(
        Array(nums, label="nums"),
        speaker="Interviewer",
        say="Every value appears twice except one. Find the single one in linear time and constant extra memory.",
    )
    yield Frame(
        Array(nums, label="nums"),
        Text("map: count each value, return the one with count 1", "O(n) time, but O(n) extra space"),
        say="A map of counts works in O(n) time, but it needs O(n) extra space. "
            "I need something that cancels pairs without remembering them.",
    )
    yield Frame(
        Array(nums, label="nums"),
        Grid([bits(5), bits(5), bits(0)], row_labels=["x = 5", "x = 5", "x ^ x"], col_labels=COLS,
             hl={(2, j): "match" for j in range(W)}),
        Text("x ^ x = 0    x ^ 0 = x    order doesn't matter"),
        say="XOR does that. A value XORed with itself is 0, and XOR with 0 changes nothing. "
            "XOR is also order-free, so pairs cancel wherever they sit.",
    )

    acc = 0
    for i, x in enumerate(nums):
        old = acc
        acc ^= x
        hl = {(2, j): "match" for j, b in enumerate(bits(acc)) if b}
        for j, (a, b) in enumerate(zip(bits(old), bits(x))):
            if a and b:
                hl[(2, j)] = "reject"
        if i == 0:
            say = f"acc starts at 0. XOR with {x} just copies it: acc = {acc}."
        elif acc == 0:
            say = f"The second {x} cancels the first. Its 1 bit meets a 1, and 1 ^ 1 = 0, so acc is back to 0."
        else:
            say = f"{x} has no partner left to cancel it, so acc becomes {acc}."
        yield Frame(
            Array(nums, label="nums", hl={k: "visited" for k in range(i)} | {i: "active"}, pointers={"i": i}),
            Grid([bits(old), bits(x), bits(acc)], row_labels=[f"acc = {old}", f"x = {x}", f"acc ^ x = {acc}"],
                 col_labels=COLS, hl=hl),
            say=say,
        )

    assert acc == expected
    yield Frame(
        Array(nums, label="nums", hl={2: "match"}),
        Grid([bits(acc)], row_labels=[f"acc = {acc}"], col_labels=COLS, hl={(0, j): "match" for j in range(W)}),
        Text((f"answer: {acc}", "match"), "O(n) time, O(1) space"),
        say=f"Only the unpaired value survives, so the answer is {acc}. One pass and one integer: "
            "O(n) time, O(1) space.",
    )
