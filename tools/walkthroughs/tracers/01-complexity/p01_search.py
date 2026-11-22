from walkthrough import Array, Frame, Text

TITLE = "Linear vs Binary Search"
PATTERN = "Binary search"


def frames():
    nums, target = [1, 3, 5, 7, 9], 7
    label = "nums (sorted)"

    yield Frame(
        Array(nums, label=label),
        Text(f"target = {target}"),
        speaker="Interviewer",
        say="The slice is sorted. Return the index of 7, or -1 if it isn't there.",
    )

    # LinearSearch, run for real to count its comparisons.
    linear_steps, linear_ans = 0, -1
    for i, x in enumerate(nums):
        linear_steps += 1
        if x == target:
            linear_ans = i
            break
    yield Frame(
        Array(nums, label=label, hl={i: "visited" for i in range(linear_ans)} | {linear_ans: "match"}),
        Text(f"linear scan: {linear_steps} comparisons to reach index {linear_ans}",
             "worst case: n comparisons → O(n)"),
        say="A linear scan checks every value from the left. It works on any slice, "
            "but it's O(n) and it ignores that the data is sorted.",
    )

    lo, hi = 0, len(nums) - 1
    yield Frame(
        Array(nums, label=label, pointers={"lo": lo, "hi": hi}, window=(lo, hi)),
        Text(f"target = {target}", "invariant: if 7 exists, it is in nums[lo..hi]"),
        say="Sorted order lets me throw away half the range per comparison. "
            "I keep lo and hi so the target, if present, is always between them.",
    )

    steps, ans = 0, -1
    while lo <= hi:
        mid = lo + (hi - lo) // 2
        steps += 1
        if nums[mid] == target:
            ans = mid
            yield Frame(
                Array(nums, label=label, hl={mid: "match"}, pointers={"lo": lo, "mid": mid, "hi": hi},
                      window=(lo, hi)),
                Text(f"mid = {lo} + ({hi}-{lo})/2 = {mid}", (f"nums[{mid}] = {nums[mid]} == {target} ✓ found", "match")),
                say=f"mid is {mid} and nums[{mid}] is exactly {target}, so I return {mid}. "
                    f"That took {steps} comparisons instead of {linear_steps}.",
            )
            break
        if nums[mid] < target:
            yield Frame(
                Array(nums, label=label, hl={i: "reject" for i in range(lo, mid + 1)},
                      pointers={"lo": lo, "mid": mid, "hi": hi}, window=(lo, hi)),
                Text(f"mid = {lo} + ({hi}-{lo})/2 = {mid}",
                     (f"nums[{mid}] = {nums[mid]} < {target} → lo = mid + 1 = {mid + 1}", "reject")),
                say=f"nums[{mid}] is {nums[mid]}, smaller than {target}. Everything left of mid is even smaller, "
                    f"so the whole left half goes. lo moves to {mid + 1}.",
            )
            lo = mid + 1
        else:
            yield Frame(
                Array(nums, label=label, hl={i: "reject" for i in range(mid, hi + 1)},
                      pointers={"lo": lo, "mid": mid, "hi": hi}, window=(lo, hi)),
                Text(f"mid = {lo} + ({hi}-{lo})/2 = {mid}",
                     (f"nums[{mid}] = {nums[mid]} > {target} → hi = mid - 1 = {mid - 1}", "reject")),
                say=f"nums[{mid}] is {nums[mid]}, bigger than {target}. Everything right of mid is bigger too, "
                    f"so the right half goes. hi moves to {mid - 1}.",
            )
            hi = mid - 1

    assert ans == 3 and linear_ans == 3
    yield Frame(
        Array(nums, label=label, hl={ans: "match"}),
        Text(("answer: 3", "match"), "LinearSearch: O(n) time, O(1) space",
             "BinarySearch: O(log n) time, O(1) space  (10⁶ items → ~20 steps)"),
        say="The answer is 3. Each step halves the range, so binary search needs about log₂ n steps: "
            "O(log n) time and O(1) space.",
    )
