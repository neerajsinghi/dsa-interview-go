import bisect

from walkthrough import Array, Frame, Text

TITLE = "Time-Based Key-Value Store"
PATTERN = "Map + binary search"


class TimeMap:
    def __init__(self):
        self.data = {}

    def set(self, key, value, ts):
        entries = self.data.setdefault(key, [])
        if not entries or entries[-1][0] < ts:
            entries.append((ts, value))  # fast path: increasing timestamps
            return
        i = bisect.bisect_left([e[0] for e in entries], ts)
        if i < len(entries) and entries[i][0] == ts:
            entries[i] = (ts, value)
            return
        entries.insert(i, (ts, value))

    def get(self, key, ts):
        entries = self.data.get(key, [])
        i = bisect.bisect_left([e[0] for e in entries], ts)  # first index with ts' >= ts
        found = i < len(entries) and entries[i][0] == ts
        if found:
            return entries[i][1], i, True
        if i == 0:
            return "", i, False
        return entries[i - 1][1], i, False


def frames():
    m = TimeMap()

    def view(hl=None, ptr=None, *lines):
        cells = [f"{t}:{v}" for t, v in m.data.get("foo", [])]
        return (Array(cells, label='data["foo"]  (timestamp:value, sorted by time)', hl=hl or {},
                      pointers=ptr or {}, cell=96),
                Text(*lines))

    yield Frame(
        *view(None, None, "Set(key, value, ts) · Get(key, ts) → value with largest timestamp ≤ ts"),
        speaker="Interviewer",
        say="Store values per key with a timestamp. Get(key, ts) returns the value with the largest "
            "timestamp at or before ts, or an empty string if there isn't one.",
    )
    yield Frame(
        *view(None, None, "naive: scan all versions of the key on every Get → O(n)"),
        say="I could keep every version in a list and scan it on each Get. That's O(n) per Get. "
            "Timestamps usually arrive in order, so the list is already sorted, and I can binary search it.",
    )

    m.set("foo", "bar", 1)
    yield Frame(
        *view({0: "active"}, None, ('Set("foo", "bar", 1)', "active")),
        say="The first Set creates the slice for \"foo\" and appends (1, bar). Appending keeps it sorted "
            "while timestamps increase.",
    )

    def get_frame(ts, expect, say):
        v, i, found = m.get("foo", ts)
        assert v == expect
        n = len(m.data["foo"])
        ans_idx = i if found else i - 1
        hl = {}
        if 0 <= ans_idx < n:
            hl[ans_idx] = "match"
        ptr = {"i": i}
        shown = f'"{v}"'
        return Frame(
            *view(hl, ptr, f"search: first index with timestamp ≥ {ts} → i = {i}"
                           + (" (exact hit)" if found else ""),
                  (f'Get("foo", {ts}) → {shown}', "match" if v else "reject")),
            say=say,
        )

    yield get_frame(1, "bar", "Get at 1: the search finds timestamp 1 exactly at i = 0, so I return \"bar\".")
    yield get_frame(3, "bar", "Get at 3: the first timestamp ≥ 3 would be at i = 1, past the end. "
                              "The floor is the entry just before it: \"bar\".")
    m.set("foo", "bar2", 4)
    yield Frame(
        *view({1: "active"}, None, ('Set("foo", "bar2", 4)', "active")),
        say="Set at 4 is newer than the last entry, so the fast path appends. Out-of-order writes would "
            "binary search and insert instead.",
    )
    yield get_frame(4, "bar2", "Get at 4 is an exact hit at i = 1: \"bar2\".")
    yield get_frame(5, "bar2", "Get at 5: i lands at 2, past the end, so the answer is entry 1: \"bar2\".")
    yield get_frame(0, "", "Get at 0: i = 0, so every stored timestamp is newer than 0. "
                           "There's no value that old, so I return an empty string.")
    yield Frame(
        *view(None, None, "Set O(1) amortized (append) · Get O(log n) · space O(entries)",
              "sync.RWMutex: Gets read in parallel, Sets take the write lock"),
        say="Set is an O(1) append in the usual case and Get is O(log n). Get only reads, so an RWMutex "
            "lets many goroutines call Get at the same time.",
    )
