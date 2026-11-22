from walkthrough import Frame, Grid, Text

TITLE = "Design Hit Counter"
PATTERN = "Ring of buckets"

W = 300
COLS = [0, 1, 2, 3, 299]  # the only buckets the example touches, plus the last


class HitCounter:
    def __init__(self):
        self.times = [0] * W
        self.counts = [0] * W

    def hit(self, ts):
        i = ts % W
        if self.times[i] != ts:
            self.times[i] = ts
            self.counts[i] = 0
        self.counts[i] += 1
        return i

    def get_hits(self, ts):
        total = 0
        live = []
        for i in range(W):
            if ts - self.times[i] < W:
                total += self.counts[i]
                if self.counts[i]:
                    live.append(i)
        return total, live


def frames():
    h = HitCounter()

    def grid(hl_cols=None):
        rows = [[h.times[i] for i in COLS], [h.counts[i] for i in COLS]]
        col_labels = [f"[{i}]" for i in COLS]
        hl = {}
        for c, st in (hl_cols or {}).items():
            k = COLS.index(c)
            hl[(0, k)] = st
            hl[(1, k)] = st
        return Grid(rows, label="buckets (index = ts % 300; [4..298] unused here)", row_labels=["time", "count"],
                    col_labels=col_labels, hl=hl, cell=60)

    yield Frame(
        grid(),
        Text("Hit(ts) records a hit · GetHits(ts) counts hits in (ts − 300, ts]"),
        speaker="Interviewer",
        say="Count hits in the last 5 minutes. Timestamps are in seconds and never go backwards. "
            "There may be a huge number of hits per second.",
    )
    yield Frame(
        grid(),
        Text("naive: queue of every hit timestamp, pop expired ones", ("memory grows with traffic", "reject")),
        say="A queue of every hit works, but a million hits per second means a million entries. "
            "I only need per-second counts, and only 300 seconds of them.",
    )

    def hit(ts):
        i = h.hit(ts)
        return i

    for ts in (1, 2, 3):
        i = hit(ts)
        yield Frame(
            grid({i: "active"}),
            Text(f"Hit({ts}) → bucket {ts} % 300 = {i}, time {h.times[i]}, count {h.counts[i]}"),
            say=f"Hit({ts}) goes to bucket {i}. It held an old second, so I reset it to time {ts} and count 1."
            if ts == 1 else f"Hit({ts}) fills bucket {i} the same way.",
        )

    total, live = h.get_hits(4)
    assert total == 3
    yield Frame(
        grid({i: "match" for i in live}),
        Text(("GetHits(4) → 3", "match"), "sum buckets with 4 − time < 300"),
        say="GetHits(4) sums every bucket whose time is within 300 seconds of 4. Buckets 1, 2 and 3 qualify: 3 hits.",
    )
    i = hit(300)
    yield Frame(
        grid({i: "active"}),
        Text(f"Hit(300) → bucket 300 % 300 = {i}"),
        say="Hit(300) wraps around to bucket 0. Its stored time was 0, not 300, so it's stale: reset and count 1.",
    )
    total, live = h.get_hits(300)
    assert total == 4
    yield Frame(
        grid({i: "match" for i in live}),
        Text(("GetHits(300) → 4", "match"), "300 − 1 = 299 < 300, so the hit at 1 still counts"),
        say="GetHits(300) counts buckets 0 to 3. The hit at second 1 is 299 seconds old, still inside the window.",
    )
    total, live = h.get_hits(301)
    assert total == 3
    yield Frame(
        grid({**{i: "match" for i in live}, 1: "reject"}),
        Text(("GetHits(301) → 3", "match"), ("bucket 1: 301 − 1 = 300, not < 300 → expired", "reject")),
        say="At 301 the hit at second 1 is exactly 300 seconds old, so it's outside (ts − 300, ts]. "
            "Bucket 1 is skipped on read and would be reused by the next Hit(301).",
    )
    yield Frame(
        grid({i: "match" for i in live}),
        Text("Hit O(1) · GetHits O(300) · space O(300)", "one sync.Mutex guards both arrays"),
        say="Hit is O(1) and GetHits is O(300), a constant. Memory is 300 buckets whatever the traffic. "
            "A mutex makes it safe for many goroutines.",
    )
