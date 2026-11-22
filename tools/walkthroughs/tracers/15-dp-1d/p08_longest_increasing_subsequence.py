import bisect

from walkthrough import Array, Frame, Text

TITLE = "Longest Increasing Subsequence"
PATTERN = "1-D DP · patience sorting"


def frames():
    nums = [10, 9, 2, 5, 3, 7, 101, 18]
    tails = []

    def widgets(i=None, hl_t=None, hl_n=None):
        hl = dict(hl_n or {})
        if i is not None:
            hl[i] = "active"
        return (
            Array(nums, label="nums", hl=hl, pointers={"x": i} if i is not None else None),
            Array(tails if tails else [""], label="tails[k] = smallest tail of an increasing run of length k+1",
                  hl=hl_t, cell=52),
        )

    yield Frame(
        *widgets(),
        speaker="Interviewer",
        say="Return the length of the longest strictly increasing subsequence. It doesn't have to be contiguous.",
    )
    yield Frame(
        *widgets(),
        Text("lis[i] = 1 + max(lis[j]) for j < i, nums[j] < nums[i]", ("O(n²) table", "compare")),
        say="The classic DP stores the longest run ending at each index and scans all earlier ones: O(n²). "
            "I can do O(n log n) by keeping only the best tails.",
    )

    says = {
        0: "10 starts the first run. tails is [10].",
        1: "9 is smaller than 10. A run of length 1 ending in 9 is better than one ending in 10, so I replace it.",
        2: "2 replaces 9 for the same reason: a smaller tail can only help later numbers.",
        3: "5 is bigger than every tail, so it extends the longest run. Length 2.",
        4: "3 lands on 5. A length-2 run ending in 3 (2, 3) beats one ending in 5, so replace.",
        5: "7 is bigger than every tail: append. Length 3.",
        6: "101 extends again: length 4.",
        7: "18 lands on 101. Same length 4, but a smaller tail.",
    }
    for i, x in enumerate(nums):
        k = bisect.bisect_left(tails, x)
        if k == len(tails):
            tails.append(x)
            line = (f"x = {x} > every tail → append (length {len(tails)})", "match")
            st = "match"
        else:
            old = tails[k]
            tails[k] = x
            line = (f"first tail ≥ {x} is tails[{k}] = {old} → replace with {x}", "compare")
            st = "compare"
        yield Frame(*widgets(i, {k: st}), Text(f"binary search for {x} in tails", line), say=says[i])

    assert len(tails) == 4
    yield Frame(
        *widgets(None, {k: "match" for k in range(len(tails))}),
        Text(f"answer: len(tails) = {len(tails)}", "note: tails is not itself a subsequence",
             ("O(n log n) time · O(n) space", "match")),
        say="The answer is the length of tails: 4. tails isn't a real subsequence, only the best tail per length. "
            "n binary searches gives O(n log n).",
    )
