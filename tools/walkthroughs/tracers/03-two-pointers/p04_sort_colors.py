from walkthrough import Array, Frame, Text

TITLE = "Sort Colors"
PATTERN = "Dutch flag"


def regions(nums, low, mid, high, focus=None):
    """0s region green, 1s region yellow, 2s region purple, unknown plain."""
    hl = {}
    for i in range(len(nums)):
        if i < low:
            hl[i] = "match"
        elif i < mid:
            hl[i] = "window"
        elif i > high:
            hl[i] = "special"
    if focus is not None:
        hl.update(focus)
    return hl


def legend(low, mid, high):
    return Text(f"low = {low}, mid = {mid}, high = {high}",
                ("[0, low) = 0s", "match"), ("[low, mid) = 1s", "window"), ("(high, n) = 2s", "special"))


def frames():
    nums = [2, 0, 2, 1, 1, 0]

    yield Frame(
        Array(nums, label="nums"),
        Text("only 0s, 1s and 2s", "sort in place, one pass, O(1) space, no library sort"),
        speaker="Interviewer",
        say="The array holds only 0s, 1s and 2s. Sort it in place in one pass with O(1) extra space, "
            "and don't call a library sort.",
    )
    yield Frame(
        Array(nums, label="nums"),
        Text("counting: count 0s, 1s, 2s, then rewrite → two passes",
             "general sort: O(n log n)"),
        say="A general sort is O(n log n). Counting each colour and rewriting is O(n), but it takes two passes. "
            "To do it in one pass, I'll partition with three pointers.",
    )

    low, mid, high = 0, 0, len(nums) - 1
    yield Frame(
        Array(nums, label="nums", hl={mid: "active"}, pointers={"low": low, "mid": mid, "high": high}),
        legend(low, mid, high),
        say="low and mid start at 0, high at the end. Everything from mid to high is unknown. "
            "I look at nums[mid] and move it into the right region.",
    )
    while mid <= high:
        v = nums[mid]
        pre = {"low": low, "mid": mid, "high": high}
        before = list(nums)
        if v == 0:
            nums[low], nums[mid] = nums[mid], nums[low]
            focus = {mid: "active", low: "compare"}
            line = ("nums[mid] = 0 → swap with nums[low], low++, mid++", "match")
            say = ("A 0 belongs at the front. Swap it with nums[low]; what comes back is a 1 I already saw, "
                   "so both low and mid move right.") if low != mid else \
                  "A 0 belongs at the front. low equals mid, so the swap does nothing. Both low and mid move right."
            low += 1
            mid += 1
        elif v == 1:
            focus = {mid: "active"}
            line = ("nums[mid] = 1 → already in the middle region, mid++", "window")
            say = "A 1 is already in the right place, because the 1s region ends at mid. Just move mid right."
            mid += 1
        else:
            nums[mid], nums[high] = nums[high], nums[mid]
            focus = {mid: "active", high: "compare"}
            line = ("nums[mid] = 2 → swap with nums[high], high--", "special")
            say = (f"A 2 belongs at the back. Swap it with nums[high] and shrink high. "
                   f"mid stays put: the {nums[mid]} that came back hasn't been examined yet.")
            high -= 1
        yield Frame(
            Array(before, label="before", hl=regions(before, pre["low"], pre["mid"], pre["high"], focus), pointers=pre),
            Array(nums, label="after", hl=regions(nums, low, mid, high), pointers={"low": low, "mid": mid, "high": high}),
            Text(line),
            say=say,
        )

    assert nums == [0, 0, 1, 1, 2, 2]  # doc comment: [2, 0, 2, 1, 1, 0] → [0, 0, 1, 1, 2, 2]
    yield Frame(
        Array(nums, label="nums", hl=regions(nums, low, mid, high), pointers={"low": low, "mid": mid, "high": high}),
        Text((f"mid > high → done: {nums}", "match"), "time O(n): one pass, each step shrinks the unknown region",
             "space O(1): three indices"),
        say="mid passed high, so nothing is unknown any more. Every step shrank the unknown region by one, "
            "so it's one O(n) pass with O(1) space.",
    )
