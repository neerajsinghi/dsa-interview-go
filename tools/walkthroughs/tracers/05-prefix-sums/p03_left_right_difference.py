from walkthrough import Array, Frame, Text

TITLE = "Left and Right Sum Differences"
PATTERN = "Prefix sums"


def frames():
    nums = [10, 4, 8, 3]
    want = [15, 1, 11, 22]
    n = len(nums)

    yield Frame(
        Array(nums, label="nums"),
        Array([None] * n, label="out"),
        speaker="Interviewer",
        say="For each index, return the absolute difference between the sum of elements before it "
            "and the sum of elements after it.",
    )
    yield Frame(
        Array(nums, label="nums"),
        Array([None] * n, label="out"),
        Text("brute force: for each i, sum both sides", "→ O(n²)"),
        say="Summing both sides for every index is O(n²). The left sum just grows by one number "
            "each step, so I can keep it running instead.",
    )

    total = sum(nums)
    yield Frame(
        Array(nums, label="nums"),
        Array([None] * n, label="out"),
        Text(f"total = {total}", "right = total − left − nums[i]"),
        say=f"First I compute the total, {total}. Then right is total minus left minus the current "
            "number. No suffix array needed.",
    )

    out = [None] * n
    left = 0
    for i, x in enumerate(nums):
        right = total - left - x
        out[i] = abs(left - right)
        hl = {j: "visited" for j in range(i)}
        hl.update({j: "compare" for j in range(i + 1, n)})
        hl[i] = "active"
        yield Frame(
            Array(nums, label="nums", hl=hl, pointers={"i": i}),
            Array(out, label="out", hl={i: "match"}),
            Text(f"left = {left}   right = {total} − {left} − {x} = {right}",
                 (f"out[{i}] = |{left} − {right}| = {out[i]}", "match")),
            say=f"Left is {left}, right is {right}, so out[{i}] is {out[i]}. "
                f"Then left grows by {x} to {left + x} for the next index.",
        )
        left += x

    assert out == want, out
    yield Frame(
        Array(out, label="out", hl={i: "match" for i in range(n)}),
        Text(f"answer: {out}", ("time O(n) · space O(1) extra", "match")),
        say="The answer is [15, 1, 11, 22]. Two passes over the array: O(n) time and O(1) extra "
            "space besides the output.",
    )
