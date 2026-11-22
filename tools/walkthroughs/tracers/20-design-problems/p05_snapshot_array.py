import bisect

from walkthrough import Frame, Grid, Text

TITLE = "Snapshot Array"
PATTERN = "Copy-on-write history"


class SnapshotArray:
    def __init__(self, length):
        self.history = [[] for _ in range(length)]
        self.snap_id = 0

    def set(self, i, val):
        h = self.history[i]
        if h and h[-1][0] == self.snap_id:
            h[-1] = (self.snap_id, val)
            return
        h.append((self.snap_id, val))

    def snap(self):
        self.snap_id += 1
        return self.snap_id - 1

    def get(self, i, snap_id):
        h = self.history[i]
        j = bisect.bisect_right([r[0] for r in h], snap_id)  # first record with snap > snap_id
        return (0 if j == 0 else h[j - 1][1]), j


def frames():
    a = SnapshotArray(3)
    COLS = 3

    def grid(hl=None):
        rows = []
        for h in a.history:
            row = [f"s{sn}:{v}" for sn, v in h]
            rows.append(row + [""] * (COLS - len(row)))
        hl = dict(hl or {})
        for r, h in enumerate(a.history):
            for c in range(len(h), COLS):
                hl.setdefault((r, c), "empty")
        return Grid(rows, label="history[i]: (snapID:value) records", row_labels=["i=0", "i=1", "i=2"],
                    hl=hl, cell=72)

    yield Frame(
        grid(),
        Text("Set(i, v) · Snap() → id · Get(i, snapID)", "array of 3 zeros"),
        speaker="Interviewer",
        say="An array starts as zeros. Set writes a value, Snap returns a snapshot id, and Get(i, snapID) "
            "returns what index i held at that snapshot.",
    )
    yield Frame(
        grid(),
        Text("naive: copy the whole array on every Snap", ("O(length) per Snap, even if nothing changed", "reject")),
        say="Copying the array on each Snap is O(length) time and memory per call, even if one index changed. "
            "Instead each index keeps its own history, and only changed indices pay.",
    )
    yield Frame(
        grid(),
        Text(("NewSnapshotArray(3)", "active"), f"snapID = {a.snap_id} (the id the next Snap returns)"),
        say="The constructor makes three empty histories. An empty history means the value has always been 0. "
            "snapID is the version currently being written.",
    )
    a.set(0, 5)
    yield Frame(
        grid({(0, 0): "active"}),
        Text(("Set(0, 5)", "active"), f"snapID = {a.snap_id}"),
        say="Set(0, 5) appends the record (snap 0, value 5) to index 0. Indices 1 and 2 cost nothing.",
    )
    r = a.snap()
    assert r == 0
    yield Frame(
        grid(),
        Text((f"Snap() → {r}", "match"), f"snapID = {a.snap_id}"),
        say="Snap is O(1): it returns 0 and bumps snapID to 1. Nothing is copied. Records tagged 0 are now frozen.",
    )
    a.set(0, 6)
    yield Frame(
        grid({(0, 1): "active"}),
        Text(("Set(0, 6)", "active"), f"snapID = {a.snap_id}"),
        say="Set(0, 6) happens in version 1, so it appends a new record instead of overwriting. "
            "A second Set before the next Snap would overwrite this record.",
    )
    v, j = a.get(0, 0)
    assert v == 5
    yield Frame(
        grid({(0, j - 1): "match", (0, j): "visited"}),
        Text(("Get(0, 0) → 5", "match"), f"first record with snap > 0 is at j = {j}; answer is record j − 1"),
        say="Get(0, 0) binary searches index 0's history for the first record after snapshot 0. "
            "That's j = 1, so the answer is the record before it: 5.",
    )
    yield Frame(
        grid(),
        Text("Set O(1) · Snap O(1) · Get O(log s) · space O(length + sets)"),
        say="Set and Snap are O(1), and Get is O(log s) for s records at that index. "
            "Memory grows with the number of Sets, not with length times snapshots.",
    )
