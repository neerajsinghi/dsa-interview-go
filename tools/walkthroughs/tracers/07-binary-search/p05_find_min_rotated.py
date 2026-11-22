from walkthrough import Array, Frame, Text

TITLE = "Find Minimum in Rotated Sorted Array"
PATTERN = "Binary search"


def view(nums, lo, hi, mid=None, mid_style="compare", label="nums"):
    hl = {i: "visited" for i in range(len(nums)) if i < lo or i > hi}
    ptr = {"lo": lo, "hi": hi}
    if mid is not None:
        hl[mid] = mid_style
        ptr["mid"] = mid
    hl.setdefault(hi, "active")
    return Array(nums, label=label, hl=hl, pointers=ptr, window=(lo, hi))


def find_min(nums):
    """Port of FindMin. Yields (lo, hi, mid, nums[mid] > nums[hi]) then ("done", lo)."""
    lo, hi = 0, len(nums) - 1
    while lo < hi:
        mid = lo + (hi - lo) // 2
        bigger = nums[mid] > nums[hi]
        yield lo, hi, mid, bigger
        if bigger:
            lo = mid + 1
        else:
            hi = mid
    yield "done", lo


def frames():
    nums = [3, 4, 5, 1, 2]

    yield Frame(
        Array(nums, label="nums (sorted, then rotated)"),
        Text("[1, 2, 3, 4, 5] rotated → [3, 4, 5, 1, 2]"),
        speaker="Interviewer",
        say="A sorted array of distinct values was rotated at an unknown point. Return the minimum in O(log n).",
    )
    yield Frame(
        Array(nums, label="nums (sorted, then rotated)"),
        Text("linear scan for the min → O(n)"),
        say="A scan finds the min in O(n). The array isn't fully sorted, so I need a yes/no question that is still monotone.",
    )
    last = nums[-1]
    pred = ["T" if v <= last else "F" for v in nums]
    yield Frame(
        Array(nums, label="nums", hl={i: ("match" if p == "T" else "reject") for i, p in enumerate(pred)}),
        Text(f"pred(i) = nums[i] <= last ({last}):  " + " ".join(pred), "the minimum is the first T"),
        say=f"Ask 'is nums[i] <= the last value?' The big left segment says no, the small right segment says yes. "
            "The minimum is the first yes, so this is a first-true search.",
    )

    answer = None
    for step in find_min(nums):
        if step[0] == "done":
            answer = nums[step[1]]
            yield Frame(
                Array(nums, label="nums", hl={step[1]: "match"}, pointers={"lo": step[1], "hi": step[1]}),
                Text(f"lo == hi == {step[1]}", (f"minimum = nums[{step[1]}] = {answer}", "match")),
                say=f"lo meets hi at index {step[1]}, so the minimum is {answer}.",
            )
            break
        lo, hi, mid, bigger = step
        if bigger:
            yield Frame(view(nums, lo, hi, mid),
                        Text(f"mid = {mid}: nums[mid] = {nums[mid]} > nums[hi] = {nums[hi]}", (f"mid is in the left segment → lo = {mid + 1}", "reject")),
                        say=f"{nums[mid]} is bigger than nums[hi] = {nums[hi]}, so mid is in the big left segment. "
                            f"The minimum is strictly to its right: lo = {mid + 1}.")
        else:
            yield Frame(view(nums, lo, hi, mid),
                        Text(f"mid = {mid}: nums[mid] = {nums[mid]} <= nums[hi] = {nums[hi]}", (f"mid is in the right segment → hi = {mid}", "active")),
                        say=f"{nums[mid]} is not bigger than nums[hi] = {nums[hi]}, so mid is in the small right segment. "
                            f"It could be the minimum itself, so I keep it: hi = {mid}.")
    assert answer == 1

    # Follow-up: the unrotated second example.
    nums2 = [11, 13, 15, 17]
    lines = []
    for step in find_min(nums2):
        if step[0] == "done":
            ans2 = nums2[step[1]]
            break
        lo, hi, mid, bigger = step
        lines.append(f"[{lo},{hi}] mid={mid}: {nums2[mid]} {'>' if bigger else '<='} {nums2[hi]} → "
                     + (f"lo={mid + 1}" if bigger else f"hi={mid}"))
    assert ans2 == 11
    yield Frame(
        Array(nums2, label="nums (not rotated)"),
        Text("compare mid with nums[hi], or with nums[lo]?"),
        speaker="Interviewer",
        say="Why compare with nums[hi] and not nums[lo]? Try [11, 13, 15, 17].",
    )
    yield Frame(
        Array(nums2, label="nums (not rotated)", hl={0: "match"}, pointers={"lo": 0, "hi": 0}),
        Text(*lines, (f"minimum = {ans2}", "match"), "vs nums[lo]: 13 > 11 looks like 'left segment' → moves right, wrong", ("O(log n) time, O(1) space", "match")),
        say="Against nums[hi], every mid is <= it, so hi slides left to index 0. Against nums[lo], an unrotated array "
            "looks like the left segment and I'd move the wrong way. O(log n), O(1).",
    )
