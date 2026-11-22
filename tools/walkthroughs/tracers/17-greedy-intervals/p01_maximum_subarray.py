from walkthrough import Array, Frame, Text

TITLE = "Maximum Subarray"
PATTERN = "Greedy · Kadane"


def frames():
    nums = [-2, 1, -3, 4, -1, 2, 1, -5, 4]
    expected = 6

    yield Frame(
        Array(nums, label="nums"),
        speaker="Interviewer",
        say="Return the largest sum of any contiguous, non-empty run of elements. Negatives are allowed.",
    )
    yield Frame(
        Array(nums, label="nums"),
        Text("brute force: every start i, extend to every end j → O(n²)"),
        say="I could try every start and extend it to every end, keeping a running sum. "
            "That's O(n²). I want one pass.",
    )

    # Kadane, as in MaxSubarray. `start` is tracked only to draw the current run.
    cur, best, start = nums[0], nums[0], 0
    best_lo, best_hi = 0, 0
    yield Frame(
        Array(nums, label="nums", hl={0: "window"}, pointers={"i": 0}, window=(0, 0)),
        Text(f"cur = {cur}", f"best = {best}"),
        say="I start with cur and best equal to the first element. cur is the best sum of a run that ends exactly here.",
    )
    for i in range(1, len(nums)):
        x = nums[i]
        extend = cur + x
        restart = extend < x
        old = cur
        cur = max(x, extend)
        if restart:
            start = i
        improved = cur > best
        best = max(best, cur)
        if improved:
            best_lo, best_hi = start, i
        hl = {k: "window" for k in range(start, i + 1)}
        hl[i] = "reject" if restart else ("match" if improved else "active")
        if restart:
            why = f"cur + x = {old} + {x} = {extend} < {x}"
            line = ("restart at x", "reject")
            say = (f"The run so far sums to {old}, which is negative. Carrying it only lowers {x}, "
                   f"so I drop it and restart at {x}.")
        else:
            why = f"cur + x = {old} + {x} = {extend} ≥ {x}"
            line = ("extend the run", "active")
            if improved:
                say = f"Extending gives {extend}, a new best. A non-negative prefix never hurts, so I keep it."
            elif i < 5:
                say = f"Extending gives {extend}. That's below best, but cur might grow again, so I keep extending."
            else:
                say = f"Extending gives {extend}. Best stays {best}."
        yield Frame(
            Array(nums, label="nums", hl=hl, pointers={"i": i}, window=(start, i)),
            Text(why, line, f"cur = {cur}    best = {best}"),
            say=say,
        )

    assert best == expected
    yield Frame(
        Array(nums, label="nums", hl={k: "match" for k in range(best_lo, best_hi + 1)},
              window=(best_lo, best_hi)),
        Text((f"answer: {best}  (run {nums[best_lo:best_hi + 1]})", "match"), "O(n) time, O(1) space"),
        say=f"The answer is {best}, from the run 4, -1, 2, 1. One pass, two variables: O(n) time and O(1) space.",
    )
