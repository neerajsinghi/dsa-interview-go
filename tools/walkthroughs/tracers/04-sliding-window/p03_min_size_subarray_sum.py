from walkthrough import Array, Frame, Text

TITLE = "Minimum Size Subarray Sum"
PATTERN = "Variable window"


def view(nums, left, right, extra=None):
    hl = {i: "window" for i in range(left, right + 1)}
    if extra:
        hl.update(extra)
    return Array(nums, label="nums", hl=hl, window=(left, right) if left <= right else None,
                 pointers={"left": left, "right": right})


def frames():
    target, nums = 7, [2, 3, 1, 2, 4, 3]

    yield Frame(
        Array(nums, label="nums (all positive)"),
        Text(f"target = {target}", "length of the shortest subarray with sum >= target, or 0"),
        speaker="Interviewer",
        say="All numbers are positive. Return the length of the shortest contiguous subarray "
            "whose sum is at least 7, or 0 if there isn't one.",
    )
    yield Frame(
        Array(nums, label="nums (all positive)"),
        Text("brute force: from every start, extend until sum >= target", "O(n²)"),
        say="From every start I could extend until the sum reaches 7. That's O(n²). "
            "Because the values are positive, growing a window only raises its sum and shrinking only lowers it.",
    )

    best = 0
    total, left = 0, 0
    pending = []  # right-edge additions that didn't make the window valid
    for right, x in enumerate(nums):
        total += x
        if total < target:
            pending.append(x)
            continue
        first = True
        while total >= target:
            length = right - left + 1
            improved = best == 0 or length < best
            if improved:
                best = length
                best_range = (left, right)
            if first and pending:
                head = (f"grow: add {' + '.join(map(str, pending + [x]))} → sum = {total} >= {target}")
                say = (f"I grow the right edge until the window is valid. After adding {x} the sum is {total}, "
                       f"which reaches {target}. Length {length} is my first answer.")
            elif first:
                head = f"grow: add {x} → sum = {total} >= {target}"
                say = (f"Adding {x} makes the sum {total}, valid again. Length {length}"
                       + (" is a new best. " if improved else f" doesn't beat {best}. ")
                       + "Now I try shrinking.")
            else:
                head = f"after shrinking: sum = {total} >= {target}, still valid"
                say = (f"The window is still valid with sum {total}, so length {length} counts"
                       + (", a new best. " if improved else ". ")
                       + f"Drop {nums[left]} from the left and keep going.")
            pending = []
            first = False
            removed = nums[left]
            yield Frame(
                view(nums, left, right, {right: "active", left: "reject"}),
                Text(head, (f"length = {right} - {left} + 1 = {length}   best = {best}", "match" if improved else None),
                     (f"valid → shrink: sum -= {removed}, left++", "reject")),
                say=say,
            )
            total -= removed
            left += 1

    assert best == 2  # doc comment: target = 7, [2, 3, 1, 2, 4, 3] → 2
    bl, br = best_range
    yield Frame(
        Array(nums, label="nums", hl={i: "match" for i in range(bl, br + 1)}, window=(bl, br)),
        Text((f"answer: {best}  ({nums[bl:br + 1]})", "match"),
             "time O(n): each index enters once (right) and leaves once (left)", "space O(1)"),
        say=f"right hit the end, so the answer is {best}, from {nums[bl:br + 1]}. Each index enters the window once and "
            "leaves at most once, so even with the inner loop it's O(n) time and O(1) space.",
    )
