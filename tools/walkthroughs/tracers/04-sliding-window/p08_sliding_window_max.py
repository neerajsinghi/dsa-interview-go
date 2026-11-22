from walkthrough import Array, Frame, Row, Text

TITLE = "Sliding Window Maximum"
PATTERN = "Monotonic deque"


def deque_view(nums, dq, hl=None):
    # Show values with their indices, front on the left.
    vals = [f"{nums[i]}@{i}" for i in dq]
    return Array(vals, label="deque (value@index, front to back)", hl=hl or {}, show_index=False, cell=64)


def frames():
    nums, k = [1, 3, -1, -3, 5, 3, 6, 7], 3

    yield Frame(
        Array(nums, label="nums"),
        Text(f"k = {k}", "report the maximum of every window of k consecutive elements"),
        speaker="Interviewer",
        say="For every window of 3 consecutive elements, left to right, report the window's maximum.",
    )
    yield Frame(
        Array(nums, label="nums"),
        Text("brute force: scan each window", f"{len(nums) - k + 1} windows × {k} → O(n·k)"),
        say="Scanning each window is O(n·k). A heap gives O(n log n). A deque of candidates gives O(n): "
            "it holds indices whose values decrease from front to back.",
    )

    dq, out = [], []
    for i, x in enumerate(nums):
        expired = None
        if dq and dq[0] <= i - k:
            expired = dq.pop(0)
        popped = []
        while dq and nums[dq[-1]] <= x:
            popped.append(dq.pop())
        dq.append(i)
        lo = max(0, i - k + 1)
        hl = {j: "window" for j in range(lo, i + 1)}
        hl[i] = "active"
        for j in popped:
            hl[j] = "reject"
        if expired is not None:
            hl[expired] = "visited"
        hl[dq[0]] = "match" if i >= k - 1 else hl[dq[0]]
        lines = []
        parts = []
        if expired is not None:
            lines.append((f"front index {expired} ≤ {i} - {k} → expired, pop front", "reject"))
            parts.append(f"index {expired} slid out, so it leaves the front")
        if popped:
            vals = ", ".join(str(nums[j]) for j in popped)
            lines.append((f"pop back {vals} (≤ {x}): older and not bigger", "reject"))
            one = len(popped) == 1
            parts.append(f"{vals} {'is' if one else 'are'} older and not bigger than {x}, "
                         f"so {'it' if one else 'they'} can never be a max again")
        lines.append(f"push {x}@{i}")
        if i >= k - 1:
            out.append(nums[dq[0]])
            lines.append((f"window [{i - k + 1}..{i}] max = front = {nums[dq[0]]}   out = {out}", "match"))
        if parts:
            say = f"{x} arrives. " + ". ".join(p[0].upper() + p[1:] for p in parts) + "."
        elif len(dq) == 1:
            say = f"{x} arrives. The deque is empty, so {x} just joins it."
        else:
            say = f"{x} arrives. Nothing behind it is smaller, so it joins the back as a future candidate."
        if i == k - 1:
            say += f" The first full window is ready: the front, {nums[dq[0]]}, is its max."
        elif i >= k:
            say += f" The front, {nums[dq[0]]}, is the max."
        yield Frame(
            Array(nums, label="nums", hl=hl, window=(lo, i), pointers={"i": i}),
            deque_view(nums, dq, {len(dq) - 1: "active", 0: "match"} if i >= k - 1 else {len(dq) - 1: "active"}),
            Text(*lines),
            say=say,
        )

    assert out == [3, 3, 5, 5, 6, 7]  # doc comment: [1, 3, -1, -3, 5, 3, 6, 7], k = 3 → [3, 3, 5, 5, 6, 7]
    yield Frame(
        Array(nums, label="nums"),
        Array(out, label="output: max of each window"),
        Text((f"answer: {out}", "match"), "time O(n): each index is pushed once and popped at most once",
             "space O(k): the deque never holds more than k indices"),
        say=f"The answer is {out}. Each index enters the deque once and leaves at most once, "
            "so it's O(n) amortised time and O(k) space.",
    )
