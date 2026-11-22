import heapq

from walkthrough import MUTED, font, Array, Frame, Intervals, Text


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


TITLE = "Meeting Rooms II"
PATTERN = "Sort + min-heap"

AXIS = (0, 30)


def heap_view(h, hl=None):
    # Show the heap as a sorted row so the minimum (the root) is always first.
    vals = sorted(h)
    return Array(vals if vals else [None], label="min-heap of end times (rooms in use)",
                 hl=hl or ({} if vals else {0: "empty"}), pointers={"min": 0} if vals else {}, show_index=False)


def frames():
    meetings = [(0, 30), (5, 10), (15, 20)]
    expected = 2

    yield Frame(
        MarkedIntervals(meetings, label="meetings [start, end)", axis=AXIS),
        heap_view([]),
        speaker="Interviewer",
        say="Meetings are half-open ranges, so one ending at 10 frees its room for one starting at 10. "
            "What's the minimum number of rooms?",
    )
    yield Frame(
        MarkedIntervals(meetings, label="meetings [start, end)", axis=AXIS),
        heap_view([]),
        Text("brute force: at each start, count meetings live then → O(n²)"),
        say="The answer is the most meetings live at once. I could count the live meetings at every start time, "
            "but that's O(n²).",
    )

    srt = sorted(meetings)  # sortedByStart
    h = []
    for k, (s, e) in enumerate(srt):
        hl = {j: "visited" for j in range(k)}
        hl[k] = "active"
        if h and h[0] <= s:
            freed = heapq.heappop(h)
            heapq.heappush(h, e)
            yield Frame(
                MarkedIntervals(srt, label="meetings sorted by start", axis=AXIS, hl=hl, marks={"now": s}),
                heap_view(h, hl={sorted(h).index(e): "match"}),
                Text(f"min end {freed} ≤ start {s} → reuse that room", (f"pop {freed}, push {e}", "match")),
                say=f"The earliest-ending room frees at {freed}, by the time [{s},{e}) starts. "
                    f"I reuse it: pop {freed}, push {e}. The room count stays {len(h)}.",
            )
        else:
            had = h[0] if h else None
            heapq.heappush(h, e)
            if had is None:
                say = f"No rooms are open yet, so [{s},{e}) opens one. I push its end time, {e}."
                why = "heap empty → open a room"
            else:
                say = (f"The earliest any room frees is {had}, but this meeting starts at {s}. "
                       f"No room is free, so I open a new one. Rooms: {len(h)}.")
                why = f"min end {had} > start {s} → all rooms busy"
            yield Frame(
                MarkedIntervals(srt, label="meetings sorted by start", axis=AXIS, hl=hl, marks={"now": s}),
                heap_view(h, hl={sorted(h).index(e): "active"}),
                Text(why, (f"push {e}", "active")),
                say=say,
            )

    rooms = len(h)
    assert rooms == expected
    yield Frame(
        MarkedIntervals(srt, label="meetings sorted by start", axis=AXIS, hl={j: "match" for j in range(len(srt))}),
        heap_view(h),
        Text((f"answer: {rooms} rooms (heap size)", "match"), "O(n log n) time, O(n) space"),
        say=f"The heap never shrinks, it only swaps, so its size is the peak: {rooms} rooms. "
            "Sorting and heap operations are O(n log n); the heap is O(n).",
    )
