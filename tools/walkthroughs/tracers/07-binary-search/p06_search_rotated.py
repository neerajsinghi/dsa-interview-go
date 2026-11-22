from walkthrough import Array, Frame, Text

TITLE = "Search in Rotated Sorted Array"
PATTERN = "Binary search"


def view(nums, lo, hi, mid=None, sorted_half=None, mid_style="compare"):
    hl = {i: "visited" for i in range(len(nums)) if i < lo or i > hi}
    if sorted_half:
        for i in range(sorted_half[0], sorted_half[1] + 1):
            hl[i] = "window"
    ptr = {"lo": lo, "hi": hi}
    if mid is not None:
        hl[mid] = mid_style
        ptr["mid"] = mid
    return Array(nums, label="nums", hl=hl, pointers=ptr, window=(lo, hi) if lo <= hi else None)


def frames():
    nums, target = [4, 5, 6, 7, 0, 1, 2], 0

    yield Frame(
        Array(nums, label="nums (sorted, then rotated)"),
        Text(f"target = {target}"),
        speaker="Interviewer",
        say="A sorted array of distinct values was rotated at an unknown pivot. Return the index of 0, or -1. Aim for O(log n).",
    )
    yield Frame(
        Array(nums, label="nums (sorted, then rotated)"),
        Text("linear scan → O(n)", "find the pivot first, then search one side → two passes"),
        say="A scan is O(n). I could find the rotation point and then search one side, but one loop is enough "
            "if I use one fact about rotated arrays.",
    )
    yield Frame(
        view(nums, 0, len(nums) - 1),
        Text("cut at mid → at least one half is fully sorted", "sorted half: I can test 'is target inside?' with its two ends"),
        say="If I cut a rotated array at mid, at least one half is fully sorted. For a sorted half, "
            "checking its first and last value tells me if the target can be in it.",
    )

    answer = -1
    lo, hi = 0, len(nums) - 1
    while lo <= hi:
        mid = lo + (hi - lo) // 2
        if nums[mid] == target:
            answer = mid
            yield Frame(
                view(nums, lo, hi, mid, mid_style="match"),
                Text(f"lo = {lo}, hi = {hi}, mid = {mid}", (f"nums[{mid}] = {nums[mid]} == {target} → return {mid}", "match"),
                     "O(log n) time, O(1) space"),
                say=f"nums[{mid}] is {target}, so I return {mid}. Each step halved the range: O(log n) time, O(1) space.",
            )
            break
        if nums[lo] <= nums[mid]:
            half = (lo, mid)
            name = f"left half [{lo}..{mid}] sorted: {nums[lo]} <= {nums[mid]}"
            inside = nums[lo] <= target < nums[mid]
            test = f"{nums[lo]} <= {target} < {nums[mid]}?"
            if inside:
                nxt, move = (lo, mid - 1), f"yes → hi = {mid - 1}"
            else:
                nxt, move = (mid + 1, hi), f"no → lo = {mid + 1}"
        else:
            half = (mid, hi)
            name = f"right half [{mid}..{hi}] sorted: nums[lo] = {nums[lo]} > {nums[mid]}"
            inside = nums[mid] < target <= nums[hi]
            test = f"{nums[mid]} < {target} <= {nums[hi]}?"
            if inside:
                nxt, move = (mid + 1, hi), f"yes → lo = {mid + 1}"
            else:
                nxt, move = (lo, mid - 1), f"no → hi = {mid - 1}"
        side = "left" if half[0] == lo else "right"
        yield Frame(
            view(nums, lo, hi, mid, sorted_half=half),
            Text(f"mid = {mid}: nums[mid] = {nums[mid]} ≠ {target}", name),
            say=f"nums[{mid}] is {nums[mid]}, not the target. nums[lo] = {nums[lo]} "
                + (f"is <= {nums[mid]}, so the left half {nums[lo:mid + 1]} is sorted."
                   if side == "left" else f"is > {nums[mid]}, so the right half {nums[mid:hi + 1]} is sorted."),
        )
        yield Frame(
            view(nums, nxt[0], nxt[1]),
            Text(f"is {target} in the sorted {side} half?  {test}", (move, "active" if inside else "reject")),
            say=(f"{target} lies inside the sorted half's range, so I search there." if inside else
                 f"{target} is outside the sorted half's range {nums[half[0]]}..{nums[half[1]]}, so it can only be in the other half.")
                + f" Now lo = {nxt[0]}, hi = {nxt[1]}.",
        )
        lo, hi = nxt
    assert answer == 4
