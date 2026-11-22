from walkthrough import Frame, Intervals, Text

TITLE = "Non-overlapping Intervals"
PATTERN = "Greedy · earliest end"

AXIS = (0, 5)


def frames():
    given = [(1, 2), (2, 3), (3, 4), (1, 3)]
    expected = 1

    yield Frame(
        Intervals(given, label="intervals", axis=AXIS),
        speaker="Interviewer",
        say="Return the fewest intervals to delete so the rest don't overlap. Intervals that only touch, "
            "like [1,2] and [2,3], are fine.",
    )
    yield Frame(
        Intervals(given, label="intervals", axis=AXIS),
        Text("brute force: try every subset to keep → O(2ⁿ · n log n)"),
        say="Deleting the fewest is the same as keeping the most. Trying every subset is exponential. "
            "This is interval scheduling, which a greedy solves.",
    )

    srt = sorted(given, key=lambda iv: (iv[1], iv[0]))  # sortedByEnd
    yield Frame(
        Intervals(srt, label="sorted by end", axis=AXIS),
        Text("order: " + " ".join(f"[{a},{b}]" for a, b in srt), "kept = 0   lastEnd = -∞"),
        say="I sort by end. The interval that ends first leaves the most room for the rest, "
            "so I always keep the earliest finisher that fits.",
    )

    kept, last_end = 0, float("-inf")
    hl = {}
    for k, (s, e) in enumerate(srt):
        cur = dict(hl)
        marks = {"lastEnd": last_end} if last_end != float("-inf") else {}
        if s >= last_end:
            kept += 1
            old = last_end
            last_end = e
            cur[k] = "match"
            hl[k] = "match"
            yield Frame(
                Intervals(srt, label="sorted by end", axis=AXIS, hl=cur, marks=marks),
                Text(f"start {s} ≥ lastEnd {old:g} → keep" if old != float("-inf") else "start 1 ≥ -∞ → keep",
                     (f"kept = {kept}   lastEnd = {last_end}", "match")),
                say=(f"[{s},{e}] starts at or after the last kept end, so it fits. I keep it and lastEnd becomes {e}."
                     if k else f"[{s},{e}] ends first, so I keep it. lastEnd becomes {e}."),
            )
        else:
            cur[k] = "reject"
            hl[k] = "reject"
            yield Frame(
                Intervals(srt, label="sorted by end", axis=AXIS, hl=cur, marks=marks),
                Text(f"start {s} < lastEnd {last_end} → overlap", (f"drop [{s},{e}]", "reject")),
                say=f"[{s},{e}] starts at {s}, before lastEnd {last_end}. It overlaps what I kept. "
                    f"It also ends no earlier, so dropping it is never worse.",
            )

    removed = len(given) - kept
    assert removed == expected
    yield Frame(
        Intervals(srt, label="sorted by end", axis=AXIS, hl=hl),
        Text(f"kept {kept} of {len(given)}", (f"answer: {removed} removal  (drop [1,3])", "match"),
             "O(n log n) time, O(n) space"),
        say=f"I kept {kept} of {len(given)}, so the answer is {removed}. Sorting is O(n log n), "
            "the sweep is O(n), and the sorted copy is O(n) space.",
    )
