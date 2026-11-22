from walkthrough import Array, Frame, Grid, Text

TITLE = "Missing Number"
PATTERN = "Bits · XOR"

W = 8
COLS = list(range(W - 1, -1, -1))


def bits(x):
    x &= (1 << W) - 1
    return [(x >> k) & 1 for k in COLS]


def frames():
    nums = [3, 0, 1]
    expected = 2
    n = len(nums)

    yield Frame(
        Array(nums, label="nums"),
        speaker="Interviewer",
        say="The slice holds n distinct numbers from 0 to n, so exactly one value is missing. Return it.",
    )
    yield Frame(
        Array(nums, label="nums"),
        Text("seen table: mark each value, find the gap → O(n) extra space",
             "sum: n(n+1)/2 - Σ nums → works, but can overflow"),
        say="A seen table needs O(n) space. The Gauss sum is O(1) space, but the sum can overflow in fixed-width ints. "
            "XOR can't overflow.",
    )

    acc = n
    yield Frame(
        Array(nums, label="nums"),
        Grid([bits(acc)], row_labels=[f"acc = n = {acc}"], col_labels=COLS),
        Text(f"indices 0..{n - 1} plus n cover 0..{n}"),
        say=f"I XOR every index with every value. The loop only covers indices 0 to {n - 1}, "
            f"so I start acc at n = {n}. Each present number then appears twice and cancels.",
    )
    for i, x in enumerate(nums):
        old = acc
        acc ^= i ^ x
        hl = {(3, j): "active" for j, b in enumerate(bits(acc)) if b}
        yield Frame(
            Array(nums, label="nums", hl={k: "visited" for k in range(i)} | {i: "active"}, pointers={"i": i}),
            Grid([bits(old), bits(i), bits(x), bits(acc)],
                 row_labels=[f"acc = {old}", f"i = {i}", f"x = {x}", f"acc^i^x = {acc}"], col_labels=COLS, hl=hl),
            say=(f"Index {i}, value {x}. acc becomes {old} ^ {i} ^ {x} = {acc}." +
                 (" The 3 from n cancels the value 3." if i == 0 else
                  " The values seen so far have cancelled their indices." if i == n - 1 else "")),
        )

    assert acc == expected
    yield Frame(
        Array(nums, label="nums"),
        Grid([bits(acc)], row_labels=[f"acc = {acc}"], col_labels=COLS,
             hl={(0, j): "match" for j in range(W)}),
        Text(f"indices + n: {' '.join(map(str, range(n + 1)))}    values: {' '.join(map(str, nums))}", (f"answer: {acc}", "match"), "O(n) time, O(1) space"),
        say=f"Every number from 0 to {n} appears twice except {acc}, which only came in as an index. "
            f"The answer is {acc}. O(n) time, O(1) space.",
    )
