from walkthrough import MUTED, font, Frame, Intervals, Text


class MarkedIntervals(Intervals):
    """Intervals with a little extra room above the rows, so a marker's name
    never covers the widget label."""

    def height(self, width):
        return super().height(width) + 16

    def draw(self, d, x, y, width):
        label = self.label
        if label:
            d.text((x, y), label, font=font("bold", 14), fill=MUTED)
        self.label = None
        super().draw(d, x, y + (36 if label else 16), width)
        self.label = label


TITLE = "Insert Interval"
PATTERN = "Linear sweep"

AXIS = (0, 10)


def frames():
    intervals = [(1, 3), (6, 9)]
    add = [2, 5]
    expected = [(1, 5), (6, 9)]

    def view(hl=None, out=(), out_hl=None, marks=None, new=None):
        nw = tuple(new or add)
        return (MarkedIntervals(list(intervals) + [nw], label="intervals  +  new (last row)", axis=AXIS,
                          hl=hl or {}, marks=marks),
                MarkedIntervals(list(out), label="out", axis=AXIS, hl=out_hl or {}))

    yield Frame(
        *view(hl={2: "special"}),
        speaker="Interviewer",
        say="The intervals are sorted by start and don't overlap. Insert [2,5], merging where needed, "
            "and keep the list sorted and non-overlapping.",
    )
    yield Frame(
        *view(hl={2: "special"}),
        Text("brute force: append, sort, run Merge Intervals → O(n log n)"),
        say="I could append it, sort, and merge. That's O(n log n). But the input is already sorted, "
            "so one left-to-right pass in three phases is enough.",
    )

    out = []
    i, n = 0, len(intervals)
    a = list(add)
    # Phase 1: strictly before.
    while i < n and intervals[i][1] < a[0]:
        out.append(intervals[i])
        yield Frame(*view(hl={i: "visited", 2: "special"}, out=out, out_hl={len(out) - 1: "active"}),
                    Text(f"phase 1: end {intervals[i][1]} < new.start {a[0]} → copy"),
                    say="This interval ends before the new one starts, so I copy it as is.")
        i += 1
    yield Frame(
        *view(hl={i: "compare", 2: "special"}, out=out, marks={"new.start": a[0]}),
        Text(f"phase 1: end {intervals[i][1]} < new.start {a[0]}?  No → phase 1 ends"),
        say=f"Phase 1 copies intervals that end before {a[0]}. [{intervals[i][0]},{intervals[i][1]}] "
            f"ends at {intervals[i][1]}, which is not before {a[0]}, so phase 1 copies nothing.",
    )
    # Phase 2: overlapping.
    while i < n and intervals[i][0] <= a[1]:
        old = list(a)
        a[0] = min(a[0], intervals[i][0])
        a[1] = max(a[1], intervals[i][1])
        yield Frame(
            *view(hl={i: "compare", 2: "match"}, out=out, marks={"new.end": old[1]}, new=a),
            Text(f"phase 2: start {intervals[i][0]} ≤ new.end {old[1]} → overlap",
                 (f"new = [min({old[0]},{intervals[i][0]}), max({old[1]},{intervals[i][1]})] = [{a[0]},{a[1]}]",
                  "match")),
            say=f"[{intervals[i][0]},{intervals[i][1]}] starts at or before {old[1]}, where the new interval ends. "
                f"They overlap, so the new interval grows to cover it: [{a[0]},{a[1]}].",
        )
        i += 1
    if i < n:
        yield Frame(
            *view(hl={0: "visited", i: "compare", 2: "match"}, out=out, marks={"new.end": a[1]}, new=a),
            Text(f"phase 2: start {intervals[i][0]} ≤ new.end {a[1]}?  No → phase 2 ends"),
            say=f"[{intervals[i][0]},{intervals[i][1]}] starts at {intervals[i][0]}, after {a[1]}. "
                "Since starts are sorted, nothing after it overlaps either.",
        )
    # Phase 3: append merged, then the rest.
    out.append(tuple(a))
    rest_from = len(out)
    out.extend(intervals[i:])
    assert out == expected
    yield Frame(
        *view(hl={k: "visited" for k in range(3)}, out=out,
              out_hl={rest_from - 1: "match", **{k: "active" for k in range(rest_from, len(out))}}, new=a),
        Text(f"phase 3: append [{a[0]},{a[1]}], then copy the rest",
             ("answer: " + " ".join(f"[{x},{y}]" for x, y in out), "match"), "O(n) time, O(n) space for the output"),
        say="Phase 3 appends the merged interval and copies the rest. The answer is [1,5] [6,9]. "
            "Each interval is looked at once: O(n) time, O(n) for the output.",
    )
