from walkthrough import Frame, Intervals, Text

TITLE = "Merge Intervals"
PATTERN = "Sort + sweep"

AXIS = (0, 18)


def frames():
    given = [(1, 3), (2, 6), (8, 10), (15, 18)]
    expected = [(1, 6), (8, 10), (15, 18)]

    yield Frame(
        Intervals(given, label="input", axis=AXIS),
        Intervals([], label="merged", axis=AXIS),
        speaker="Interviewer",
        say="Merge every group of overlapping intervals. Touching end to start counts as overlapping. "
            "Return the result sorted by start.",
    )
    yield Frame(
        Intervals(given, label="input", axis=AXIS),
        Intervals([], label="merged", axis=AXIS),
        Text("brute force: fuse any overlapping pair, repeat until none → O(n³)"),
        say="The naive way fuses any overlapping pair and repeats until nothing changes. That's O(n³). "
            "Sorting gives the intervals an order I can sweep.",
    )

    srt = sorted(given)  # sortedByStart: by start, then end
    out = [list(srt[0])]
    yield Frame(
        Intervals(srt, label="sorted by start", axis=AXIS, hl={0: "match"}),
        Intervals([tuple(x) for x in out], label="merged", axis=AXIS, hl={0: "active"}, marks={"last.end": out[-1][1]}),
        Text(f"out = [{srt[0][0]},{srt[0][1]}]"),
        say="After sorting by start, an interval can only overlap the last merged one. "
            "Everything earlier ended before it. I seed the output with the first interval.",
    )
    for k in range(1, len(srt)):
        s, e = srt[k]
        last = out[-1]
        hl_in = {j: "visited" for j in range(k)}
        if s <= last[1]:
            old = last[1]
            last[1] = max(last[1], e)
            hl_in[k] = "compare"
            yield Frame(
                Intervals(srt, label="sorted by start", axis=AXIS, hl=hl_in),
                Intervals([tuple(x) for x in out], label="merged", axis=AXIS, hl={len(out) - 1: "match"},
                          marks={"last.end": last[1]}),
                Text(f"start {s} ≤ last.end {old} → overlap", (f"last.end = max({old}, {e}) = {last[1]}", "match")),
                say=f"[{s},{e}] starts at {s}, before the last merged interval ends at {old}. "
                    f"They overlap, so I extend the end to max({old}, {e}) = {last[1]}.",
            )
        else:
            out.append([s, e])
            hl_in[k] = "active"
            yield Frame(
                Intervals(srt, label="sorted by start", axis=AXIS, hl=hl_in),
                Intervals([tuple(x) for x in out], label="merged", axis=AXIS, hl={len(out) - 1: "active"},
                          marks={"last.end": e}),
                Text(f"start {s} > last.end {last[1]} → gap", ("append as a new interval", "active")),
                say=f"[{s},{e}] starts at {s}, after the last end {last[1]}. There's a gap, "
                    "and nothing later can close it, so I start a new merged interval.",
            )

    result = [tuple(x) for x in out]
    assert result == expected
    yield Frame(
        Intervals(srt, label="sorted by start", axis=AXIS, hl={j: "visited" for j in range(len(srt))}),
        Intervals(result, label="merged", axis=AXIS, hl={j: "match" for j in range(len(result))}),
        Text(("answer: " + " ".join(f"[{a},{b}]" for a, b in result), "match"), "O(n log n) time, O(n) space"),
        say="The result is [1,6] [8,10] [15,18]. The sort costs O(n log n) and the sweep is O(n). "
            "The sorted copy and output take O(n) space.",
    )
