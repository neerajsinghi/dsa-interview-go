from walkthrough import Array, Frame, Text

TITLE = "Maximum Product Subarray"
PATTERN = "1-D DP · Kadane"


def frames():
    nums = [2, 3, -2, 4]
    n = len(nums)
    his = [None] * n
    los = [None] * n

    def widgets(i=None, hl_hi=None, hl_lo=None, best=None):
        ptr = {"i": i} if i is not None else None
        hl = {i: "active"} if i is not None else None
        return (
            Array(nums, label="nums", hl=hl, pointers=ptr),
            Array(his, label="hi = max product of a subarray ending here", hl=hl_hi),
            Array(los, label="lo = min product of a subarray ending here", hl=hl_lo),
        )

    yield Frame(
        *widgets(),
        speaker="Interviewer",
        say="Find the contiguous, non-empty subarray with the largest product. Numbers can be negative or zero.",
    )
    yield Frame(
        *widgets(),
        Text("try every start i, extend the product to every end j", ("O(n²)", "reject")),
        say="The brute force tries every start and extends the product to every end. That's O(n²). "
            "I want one pass, like Kadane for sums.",
    )

    hi = lo = best = nums[0]
    his[0], los[0] = hi, lo
    yield Frame(
        *widgets(0, {0: "active"}, {0: "active"}),
        Text(f"hi = lo = best = nums[0] = {nums[0]}"),
        say="But a negative number flips order: the smallest product can become the largest. "
            "So I track both hi and lo for subarrays ending here. Both start at 2.",
    )
    for i in range(1, n):
        x = nums[i]
        a, b = hi * x, lo * x
        hi, lo = max(x, a, b), min(x, a, b)
        best = max(best, hi)
        his[i], los[i] = hi, lo
        if i == 1:
            say = "3 extends the run: hi = 2 × 3 = 6, while lo keeps the smallest choice, 3 alone. best becomes 6."
        elif i == 2:
            say = ("-2 flips everything. The old hi times -2 is now the smallest, -12. "
                   "Starting fresh at -2 is the largest. best stays 6.")
        else:
            say = ("4: the largest is starting fresh at 4. lo grows to -48, which a later negative "
                   "could turn into a big win.")
        yield Frame(
            *widgets(i, {i: "active", i - 1: "compare"}, {i: "active", i - 1: "compare"}),
            Text(f"candidates: x = {x}, hi·x = {a}, lo·x = {b}",
                 (f"hi = max = {hi}, lo = min = {lo}, best = {best}", "active")),
            say=say,
        )

    assert best == 6
    yield Frame(
        *widgets(None, {1: "match"}, None),
        Text(f"answer: best = {best}  (subarray [2, 3])", ("O(n) time · O(1) space", "match")),
        say="The answer is 6, from [2, 3]. I compute hi and lo together from the old values, "
            "in one tuple assignment. One pass, O(n) time, O(1) space.",
    )
