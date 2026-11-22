from walkthrough import Array, Frame, Text

TITLE = "Search Insert Position"
PATTERN = "Lower bound"


def view(nums, lo, hi, mid=None, mid_style="compare", label="nums"):
    # Half-open [lo, hi): index hi itself is "the answer could be here" but is not a cell to test.
    hl = {i: "visited" for i in range(len(nums)) if i < lo or i >= hi}
    ptr = {"lo": lo, "hi": hi}
    if mid is not None:
        hl[mid] = mid_style
        ptr["mid"] = mid
    return Array(nums, label=label, hl=hl, pointers=ptr, window=(lo, hi - 1) if lo < hi else None)


def search_insert(nums, target):
    """Port of SearchInsert: half-open [lo, hi). Yields (event, lo, hi, mid)."""
    lo, hi = 0, len(nums)
    while lo < hi:
        mid = lo + (hi - lo) // 2
        if nums[mid] < target:
            yield "right", lo, hi, mid
            lo = mid + 1
        else:
            yield "left", lo, hi, mid
            hi = mid
    yield "done", lo, hi, None


def run(nums, target):
    out = []
    for ev, lo, hi, mid in search_insert(nums, target):
        if ev == "done":
            return lo, out
        calc = f"[{lo},{hi})  mid = {mid}  nums[{mid}] = {nums[mid]}"
        if ev == "right":
            out.append((lo, hi, mid, calc, f"{nums[mid]} < {target} → lo = {mid + 1}", "reject"))
        else:
            out.append((lo, hi, mid, calc, f"{nums[mid]} >= {target} → hi = {mid}", "active"))


def frames():
    nums = [1, 3, 5, 6]

    yield Frame(
        Array(nums, label="nums (sorted, distinct)"),
        Text("target = 5 → 2      target = 2 → 1      target = 7 → 4"),
        speaker="Interviewer",
        say="Return the index of the target if it's there. Otherwise return where it would be inserted to keep the array sorted.",
    )
    yield Frame(
        Array(nums, label="nums (sorted, distinct)"),
        Text("scan: first i with nums[i] >= target → O(n)"),
        say="Scanning for the first value that is at least the target answers both cases, but it's O(n). "
            "That 'first index where nums[i] >= target' is the lower bound.",
    )
    yield Frame(
        view(nums, 0, len(nums)),
        Text("lo = 0, hi = 4 = len(nums)", "pred(i) = nums[i] >= target:  F F ... T T"),
        say="The predicate goes false, false, then true. I search the half-open range [0, 4). "
            "hi starts at len(nums) because 'insert at the end' is a valid answer.",
    )

    target = 5
    ans, steps = run(nums, target)
    for lo, hi, mid, calc, verdict, st in steps:
        say = (f"nums[{mid}] = {nums[mid]} is at least {target}, so mid might be the answer. "
               f"I keep it: hi = {mid}." if st == "active" else
               f"nums[{mid}] = {nums[mid]} is below {target}, so mid and everything left of it is too small: lo = {mid + 1}.")
        yield Frame(view(nums, lo, hi, mid), Text(f"target = {target}", calc, (verdict, st)), say=say)
    assert ans == 2
    yield Frame(
        view(nums, ans, ans, label="nums"),
        Text(f"lo == hi == {ans}", (f"nums[{ans}] = {nums[ans]} == {target} → found at {ans}", "match")),
        say=f"lo and hi meet at {ans}. The loop stops and lo is the first index with nums[i] >= {target}. Here it holds {target} itself.",
    )

    target = 2
    ans2, steps = run(nums, target)
    lines = [f"target = {target}"] + [f"{c}: {v}" for _, _, _, c, v, _ in steps]
    yield Frame(
        Array(nums, label="nums"),
        Text(f"target = {target}"),
        speaker="Interviewer",
        say="And when the target is missing, like 2?",
    )
    assert ans2 == 1
    ans3, _ = run(nums, 7)
    assert ans3 == 4
    yield Frame(
        view(nums, ans2, ans2),
        Text(*lines, (f"return {ans2}: insert 2 before 3", "match"), "target 7 → 4 (append at the end)",
             ("O(log n) time, O(1) space", "match")),
        say=f"Same loop, no special case. It stops at {ans2}, the first value bigger than 2, which is exactly where 2 goes. "
            "For 7 it returns 4, the end. O(log n) time, O(1) space.",
    )
