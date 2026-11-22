from walkthrough import Array, Frame, KV, Text

TITLE = "Contiguous Array"
PATTERN = "Prefix + first index"


def frames():
    # The first doc example ([0, 1]) is only two steps; the second one shows the
    # idea properly, and the doc comment draws its prefix table.
    nums = [0, 1, 1, 1, 0, 0]
    want = 6
    n = len(nums)
    signs = [1 if x == 1 else -1 for x in nums]

    yield Frame(
        Array(nums, label="nums"),
        speaker="Interviewer",
        say="The array holds only 0s and 1s. Return the length of the longest contiguous subarray "
            "with as many 0s as 1s.",
    )
    yield Frame(
        Array(nums, label="nums"),
        Text("brute force: every start, track 1s − 0s while extending → O(n²)"),
        say="I could fix each start and extend, tracking ones minus zeros, and record the length "
            "whenever it hits zero. That's O(n²).",
    )
    yield Frame(
        Array(nums, label="nums"),
        Array([f"{s:+d}" for s in signs], label="as ±1  (0 counts as −1)"),
        Text("equal 0s and 1s  ⇔  sum of ±1 values is 0  ⇔  two equal balances"),
        say="If I count a 0 as −1, equal counts means the sum is 0. A zero-sum stretch sits between "
            "two points with the same running balance.",
    )

    first = {0: -1}
    yield Frame(
        Array(nums, label="nums"),
        KV(first, label="first  (balance: first index)", hl={0: "special"}),
        Text("best = 0   balance = 0"),
        say="For the longest stretch I keep the first index of each balance. Balance 0 is seeded at "
            "index −1, so a stretch can start at index 0.",
    )

    best, balance = 0, 0
    bal_row = [None] * n
    for i, x in enumerate(nums):
        balance += 1 if x == 1 else -1
        bal_row[i] = balance
        before = dict(first)
        if balance in first:
            j = first[balance]
            length = i - j
            improved = length > best
            best = max(best, length)
            yield Frame(
                Array(nums, label="nums", hl={t: "window" for t in range(j + 1, i + 1)}, window=(j + 1, i),
                      pointers={"i": i}),
                Array(bal_row, label="balance", hl={i: "match"}),
                KV(before, label="first", hl={balance: "match"}),
                Text(f"balance {balance} first seen at {j} → length {i} − ({j}) = {length}",
                     (f"best = {best}", "match" if improved else None)),
                say=f"Balance {balance} was first seen at index {j}, so indices {j + 1}..{i} are balanced: "
                    f"length {length}. I keep the old first index; overwriting it would shorten later spans.",
            )
        else:
            first[balance] = i
            yield Frame(
                Array(nums, label="nums", hl={i: "active"}, pointers={"i": i}),
                Array(bal_row, label="balance", hl={i: "active"}),
                KV(first, label="first", hl={balance: "active"}),
                Text(f"balance {balance} is new → first[{balance}] = {i}", f"best = {best}"),
                say=f"Balance {balance} is new, so I record index {i} as its first appearance. "
                    "Nothing balanced ends here yet.",
            )

    assert best == want
    yield Frame(
        Array(nums, label="nums", hl={t: "match" for t in range(n)}, window=(0, n - 1)),
        Text(f"answer: {best}", ("time O(n) · space O(n)", "match")),
        say="The answer is 6: the whole array. One pass with a map of at most n + 1 balances, "
            "so O(n) time and O(n) space.",
    )
