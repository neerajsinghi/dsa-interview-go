from walkthrough import Array, Frame, Text

TITLE = "First and Last Position"
PATTERN = "Lower / upper bound"


def view(nums, lo, hi, mid=None, mid_style="compare", extra=None):
    hl = {i: "visited" for i in range(len(nums)) if i < lo or i >= hi}
    hl.update(extra or {})
    ptr = {"lo": lo, "hi": hi}
    if mid is not None:
        hl[mid] = mid_style
        ptr["mid"] = mid
    return Array(nums, label="nums", hl=hl, pointers=ptr, window=(lo, hi - 1) if lo < hi else None)


def bound(nums, go_right):
    """Half-open first-true search shared by LowerBound and UpperBound.

    go_right(v) is the "mid is too small" test: v < target for LowerBound,
    v <= target for UpperBound. Yields (lo, hi, mid, went_right) then ("done", lo).
    """
    lo, hi = 0, len(nums)
    while lo < hi:
        mid = lo + (hi - lo) // 2
        right = go_right(nums[mid])
        yield lo, hi, mid, right
        if right:
            lo = mid + 1
        else:
            hi = mid
    yield "done", lo


def frames():
    nums, target = [5, 7, 7, 8, 8, 10], 8

    yield Frame(
        Array(nums, label="nums (sorted, may repeat)"),
        Text(f"target = {target}"),
        speaker="Interviewer",
        say="The sorted array may contain duplicates. Return the first and last index of 8, or [-1, -1]. Aim for O(log n).",
    )
    yield Frame(
        Array(nums, label="nums (sorted, may repeat)"),
        Text("scan all cells → O(n)", "find any 8, then walk outward → O(n) with many copies"),
        say="A scan is O(n). Finding one 8 and walking outward is also O(n) if the whole array is 8s. "
            "I'll do two boundary searches instead.",
    )

    # Search 1: lower bound, first index with nums[i] >= target.
    first = None
    for step in bound(nums, lambda v: v < target):
        if step[0] == "done":
            first = step[1]
            break
        lo, hi, mid, right = step
        verdict = (f"{nums[mid]} < {target} → lo = {mid + 1}", "reject") if right else \
                  (f"{nums[mid]} >= {target} → hi = {mid}", "active")
        say = (f"Lower bound: {nums[mid]} is at least {target}, so the first 8 is at mid or to its left. hi = {mid}."
               if not right else
               f"Lower bound: {nums[mid]} is below {target}, so the first 8 is to the right. lo = {mid + 1}.")
        yield Frame(view(nums, lo, hi, mid), Text("first = lowerBound: first i with nums[i] >= 8",
                                                   f"[{lo},{hi})  mid = {mid}", verdict), say=say)
    assert first < len(nums) and nums[first] == target
    yield Frame(
        view(nums, first, first, extra={first: "match"}),
        Text(f"lo == hi == {first}", (f"nums[{first}] = {nums[first]} == {target} → first = {first}", "match")),
        say=f"The lower bound is {first}, and nums[{first}] really is {target}. If it weren't, I'd return [-1, -1] here.",
    )

    # Search 2: upper bound, first index with nums[i] > target.
    up = None
    for step in bound(nums, lambda v: v <= target):
        if step[0] == "done":
            up = step[1]
            break
        lo, hi, mid, right = step
        verdict = (f"{nums[mid]} <= {target} → lo = {mid + 1}", "reject") if right else \
                  (f"{nums[mid]} > {target} → hi = {mid}", "active")
        say = (f"Upper bound: {nums[mid]} is not bigger than {target}, so the first value past the 8s is to the right. lo = {mid + 1}."
               if right else
               f"Upper bound: {nums[mid]} is bigger than {target}. It could be the first value past the 8s, so hi = {mid}.")
        yield Frame(view(nums, lo, hi, mid, extra={first: "match"}),
                    Text("upperBound: first i with nums[i] > 8", f"[{lo},{hi})  mid = {mid}", verdict), say=say)

    last = up - 1
    result = [first, last]
    assert result == [3, 4]
    hl = {first: "match", last: "match"}
    if up < len(nums):
        hl[up] = "window"
    yield Frame(
        Array(nums, label="nums", hl=hl, pointers={"first": first, "last": last, "upper": up}),
        Text(f"upperBound = {up}, last = {up} - 1 = {last}", (f"answer: [{first}, {last}]", "match"),
             "two searches: O(log n) time, O(1) space"),
        say=f"The upper bound is {up}, the first value past the 8s, so the last 8 sits at {last}. "
            f"The answer is [{first}, {last}]. Two binary searches: O(log n) time, O(1) space.",
    )
