from walkthrough import Array, Frame, KV, Text, Widget

TITLE = "LFU Cache"
PATTERN = "Freq lists + minFreq"


class LFU:
    """Port of LFUCache. freqs[f] is a list with the front (index 0) = most recent."""

    def __init__(self, cap):
        self.cap = cap
        self.items = {}  # key -> [val, freq]
        self.freqs = {}
        self.min_freq = 0

    def _touch(self, k):
        f = self.items[k][1]
        old = self.freqs[f]
        old.remove(k)
        if not old:
            del self.freqs[f]
            if self.min_freq == f:
                self.min_freq += 1
        self.items[k][1] = f + 1
        self.freqs.setdefault(f + 1, []).insert(0, k)

    def get(self, k):
        if k not in self.items:
            return None
        self._touch(k)
        return self.items[k][0]

    def put(self, k, v):
        if self.cap <= 0:
            return None
        if k in self.items:
            self.items[k][0] = v
            self._touch(k)
            return None
        victim = None
        if len(self.items) == self.cap:
            lst = self.freqs[self.min_freq]
            victim = lst.pop()  # back = least recent among least frequent
            if not lst:
                del self.freqs[self.min_freq]
            del self.items[victim]
        self.items[k] = [v, 1]
        self.freqs.setdefault(1, []).insert(0, k)
        self.min_freq = 1
        return victim


class FreqLists(Widget):
    """One row per frequency: 'freq f' label, then its keys front to back; minFreq is marked."""

    def __init__(self, freqs, min_freq, hl=None, rows=3):
        self.freqs = freqs
        self.min_freq = min_freq
        self.hl = hl or {}
        self.rows = rows

    def height(self, width):
        return 20 + self.rows * 50

    def draw(self, d, x, y, width):
        from walkthrough import font, MUTED, STYLES, centered
        d.text((x, y), "freqs: freq to LRU list (front = most recent, evict from back)", font=font("bold", 14), fill=MUTED)
        y += 22
        keys = sorted(self.freqs)
        if not keys:
            d.text((x, y + 10), "(empty)", font=font("mono", 14), fill=MUTED)
            return
        for r, f in enumerate(keys):
            yy = y + r * 50
            label = f"freq {f}"
            d.text((x, yy + 20), label, font=font("mono", 15), fill="#1f2328", anchor="lm")
            for k, key in enumerate(self.freqs[f]):
                bx = x + 90 + k * 70
                fill, border, ink = STYLES[self.hl.get(key)]
                d.rounded_rectangle([bx, yy + 2, bx + 56, yy + 38], 8, fill=fill, outline=border,
                                    width=2 if key in self.hl else 1)
                centered(d, bx + 28, yy + 20, f"k{key}", font("mono", 16), ink)
                if k + 1 < len(self.freqs[f]):
                    d.text((bx + 58, yy + 20), "⇄", font=font("mono", 14), fill=MUTED, anchor="lm")
            if f == self.min_freq:
                tx = x + 90 + len(self.freqs[f]) * 70 + 6
                d.text((tx, yy + 20), "← minFreq", font=font("mono", 14), fill="#8250df", anchor="lm")


def frames():
    c = LFU(2)

    def view(call, say, key=None, style="active", extra=None):
        kv = {k: f"v{v[0]} f{v[1]}" for k, v in sorted(c.items.items())}
        lines = [(call, None if style == "active" else style)]
        if extra:
            lines.append(extra)
        return Frame(
            KV(kv, label="items: key to (value, freq)", hl={key: style} if key in kv else {}),
            FreqLists(c.freqs, c.min_freq, hl={key: style} if key in c.items else {}),
            Text(*lines),
            say=say,
        )

    yield Frame(
        KV({}, label="items: key to (value, freq)"),
        FreqLists({}, 0),
        Text("capacity = 2 · evict least frequently used · ties: least recently used"),
        speaker="Interviewer",
        say="Now evict the least frequently used key. Every Get hit and every Put counts as a use. "
            "Break ties by evicting the least recently used. Get and Put must be O(1).",
    )
    yield Frame(
        KV({}, label="items: key to (value, freq)"),
        FreqLists({}, 0),
        Text("naive: store a count per key, scan for the minimum on eviction → O(n)",
             "or a heap by (freq, lastUse) → O(log n)"),
        say="Counting uses and scanning for the minimum is O(n) per eviction, and a heap is O(log n). "
            "For O(1) I group keys by frequency, each group an LRU list, and track the smallest frequency.",
    )

    victim = c.put(1, 1)
    yield view("Put(1, 1)", "Put(1,1): a new key always has freq 1, so it goes to the front of list 1 "
                            "and minFreq becomes 1.", 1)
    c.put(2, 2)
    yield view("Put(2, 2)", "Put(2,2) also goes to the front of list 1. Key 1 is now the older one there.", 2)
    r = c.get(1)
    assert r == 1
    yield view("Get(1) → 1", "Get(1) is a use: key 1 moves from list 1 to the front of list 2. "
                             "List 1 still has key 2, so minFreq stays 1.", 1, "match")
    victim = c.put(3, 3)
    assert victim == 2
    yield view("Put(3, 3)", "Full. minFreq is 1, so I evict from the back of list 1: key 2. "
                            "Then 3 enters list 1 and minFreq is reset to 1.", 3, "active",
               ("evicted 2 (freq 1)", "reject"))
    r = c.get(3)
    assert r == 3
    yield view("Get(3) → 3", "Get(3) moves key 3 to list 2. List 1 is now empty and was minFreq, "
                             "so minFreq becomes 2. No search needed.", 3, "match")
    victim = c.put(4, 4)
    assert victim == 1
    yield view("Put(4, 4)", "Full again. Keys 1 and 3 both have freq 2. Key 1 is at the back of list 2, "
                            "so it's the least recent of the two and gets evicted.", 4, "active",
               ("evicted 1 (tie at freq 2, older use)", "reject"))
    yield Frame(
        KV({k: f"v{v[0]} f{v[1]}" for k, v in sorted(c.items.items())}, label="items: key to (value, freq)"),
        FreqLists(c.freqs, c.min_freq),
        Text("Get O(1) · Put O(1) · space O(capacity)", "sync.Mutex: Get changes frequencies"),
        say="Every step is a few map lookups and list moves, so Get and Put are O(1), with O(capacity) space. "
            "A plain Mutex guards it, since Get changes frequencies.",
    )
