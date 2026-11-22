from walkthrough import Array, Frame, Text

TITLE = "Maximum Sum Subarray of Size K"
PATTERN = "Fixed window"


def win(nums, lo, hi, extra=None):
    hl = {i: "window" for i in range(lo, hi + 1)}
    if extra:
        hl.update(extra)
    return Array(nums, label="nums", hl=hl, window=(lo, hi), pointers={"left": lo, "right": hi})


def frames():
    nums, k = [2, 1, 5, 1, 3, 2], 3

    yield Frame(
        Array(nums, label="nums"),
        Text(f"k = {k}", "largest sum of any k consecutive elements"),
        speaker="Interviewer",
        say="Return the largest sum of any 3 consecutive elements. If k doesn't fit the array, report that it's not ok.",
    )
    yield Frame(
        Array(nums, label="nums"),
        Text("brute force: sum every window from scratch", f"{len(nums) - k + 1} windows × {k} additions → O(n·k)"),
        say="I could add up every window from scratch. That's O(n·k). But neighbouring windows share k-1 elements, "
            "so recomputing them is wasted work.",
    )

    total = sum(nums[:k])
    best = total
    best_lo = 0
    yield Frame(
        win(nums, 0, k - 1),
        Text(f"first window: {' + '.join(map(str, nums[:k]))} = {total}", (f"best = {best}", "match")),
        say=f"I sum the first window once: {total}. I seed best with it, not 0, because every value could be negative.",
    )
    for right in range(k, len(nums)):
        out = nums[right - k]
        old = total
        total += nums[right] - out
        improved = total > best
        best = max(best, total)
        lo = right - k + 1
        if improved:
            best_lo = lo
        yield Frame(
            win(nums, lo, right, {right: "match", right - k: "reject"}),
            Text(f"sum = {old} + {nums[right]} - {out} = {total}",
                 (f"best = {best}" + ("  (new best)" if improved else ""), "match" if improved else None)),
            say=(f"Slide right by one: {nums[right]} enters and {out} leaves, so the sum is {old} + {nums[right]} - {out} = {total}. "
                 + ("That beats the old best." if improved else f"That doesn't beat {best}.")),
        )

    assert best == 9  # doc comment: [2, 1, 5, 1, 3, 2], k = 3 → 9
    yield Frame(
        Array(nums, label="nums", hl={i: "match" for i in range(best_lo, best_lo + k)}, window=(best_lo, best_lo + k - 1)),
        Text((f"answer: {best}  ({' + '.join(map(str, nums[best_lo:best_lo + k]))})", "match"), "time O(n): one add and one subtract per step",
             "space O(1)"),
        say=f"The answer is {best}, from {' + '.join(map(str, nums[best_lo:best_lo + k]))}. Each slide costs one add and one subtract, "
            "so it's O(n) time and O(1) space.",
    )
