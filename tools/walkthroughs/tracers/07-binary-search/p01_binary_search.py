from walkthrough import Array, Frame, Text

TITLE = "Binary Search"
PATTERN = "Binary search"


def view(nums, lo, hi, mid=None, mid_style="compare"):
    hl = {i: "visited" for i in range(len(nums)) if i < lo or i > hi}
    ptr = {"lo": lo, "hi": hi}
    if mid is not None:
        hl[mid] = mid_style
        ptr["mid"] = mid
    return Array(nums, label="nums", hl=hl, pointers=ptr, window=(lo, hi) if lo <= hi else None)


def search(nums, target):
    """Port of Search: closed interval [lo, hi]. Yields (event, lo, hi, mid)."""
    lo, hi = 0, len(nums) - 1
    while lo <= hi:
        mid = lo + (hi - lo) // 2
        if nums[mid] == target:
            yield "found", lo, hi, mid
            return
        if nums[mid] < target:
            yield "right", lo, hi, mid
            lo = mid + 1
        else:
            yield "left", lo, hi, mid
            hi = mid - 1
    yield "missing", lo, hi, None


def frames():
    nums, target = [-1, 0, 3, 5, 9, 12], 9

    yield Frame(
        Array(nums, label="nums (sorted, distinct)"),
        Text(f"target = {target}"),
        speaker="Interviewer",
        say="The array is sorted with distinct values. Return the index of 9, or -1 if it isn't there. Aim for O(log n).",
    )
    yield Frame(
        Array(nums, label="nums (sorted, distinct)"),
        Text(f"target = {target}", "linear scan: check every cell → O(n)"),
        say="A linear scan works but is O(n). Because the array is sorted, one comparison with the middle "
            "tells me which half can't hold the target.",
    )
    yield Frame(
        view(nums, 0, len(nums) - 1),
        Text("lo = 0, hi = 5   closed interval [lo, hi]", "loop while lo <= hi"),
        say="I keep a closed interval [lo, hi] of indices that could still hold the target. "
            "At the start that is the whole array.",
    )

    answer = None
    for ev, lo, hi, mid in search(nums, target):
        calc = f"mid = {lo} + ({hi} - {lo}) / 2 = {mid}"
        if ev == "found":
            answer = mid
            yield Frame(
                view(nums, lo, hi, mid, "match"),
                Text(calc, (f"nums[{mid}] = {nums[mid]} == {target} → return {mid}", "match")),
                say=f"nums[{mid}] is {nums[mid]}, which is the target, so I return {mid}.",
            )
        elif ev == "right":
            yield Frame(
                view(nums, lo, hi, mid),
                Text(calc, (f"nums[{mid}] = {nums[mid]} < {target} → lo = {mid + 1}", "reject")),
                say=f"nums[{mid}] is {nums[mid]}, smaller than {target}. Everything left of mid is even smaller, "
                    f"so I drop indices {lo}..{mid} and set lo = {mid + 1}.",
            )
            yield Frame(
                view(nums, mid + 1, hi),
                Text(f"lo = {mid + 1}, hi = {hi}   {hi - mid} candidates left"),
                say="One comparison threw away half the array. The target, if it exists, is still inside [lo, hi].",
            )
    assert answer == 4

    # The second doc example shows how the loop ends when the target is absent.
    target2 = 2
    trace = list(search(nums, target2))
    steps = []
    for ev, lo, hi, mid in trace[:-1]:
        op = "<" if ev == "right" else ">"
        steps.append(f"[{lo},{hi}] mid={mid}: {nums[mid]} {op} {target2}")
    _, lo, hi, _ = trace[-1]
    yield Frame(
        Array(nums, label="nums"),
        Text(f"target = {target2}"),
        speaker="Interviewer",
        say="What happens if the target is 2, which isn't in the array?",
    )
    assert trace[-1][0] == "missing"
    yield Frame(
        view(nums, lo, hi),
        Text(*steps, (f"lo = {lo} > hi = {hi} → return -1", "reject"), ("O(log n) time, O(1) space", "match")),
        say=f"The interval keeps shrinking until lo passes hi, here lo = {lo} and hi = {hi}. It's empty, so I return -1. "
            "Each step halves the range: O(log n) time, O(1) space.",
    )
