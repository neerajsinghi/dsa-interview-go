from walkthrough import Frame, KV, LinkedList, Text

TITLE = "LRU Cache"
PATTERN = "Map + doubly linked list"


class LRU:
    """Port of LRUCache: dict for lookup, list order[0] = most recent (front)."""

    def __init__(self, cap):
        self.cap = cap
        self.items = {}
        self.order = []  # stands in for the doubly linked list; front = most recent

    def get(self, k):
        if k not in self.items:
            return None
        self.order.remove(k)
        self.order.insert(0, k)  # moveToFront
        return self.items[k]

    def put(self, k, v):
        if self.cap <= 0:
            return None
        if k in self.items:
            self.items[k] = v
            self.order.remove(k)
            self.order.insert(0, k)
            return None
        evicted = None
        if len(self.order) == self.cap:
            evicted = self.order.pop()  # back()
            del self.items[evicted]
        self.order.insert(0, k)
        self.items[k] = v
        return evicted


def frames():
    c = LRU(2)

    def view(call, say, hl_key=None, style="active", extra=None):
        lst = [f"{k}:{c.items[k]}" for k in c.order]
        ll_hl = {c.order.index(hl_key): style} if hl_key in c.order else {}
        lines = [(call, style if style != "active" else None)]
        if extra:
            lines.append(extra)
        return Frame(
            KV({k: f"node({k})" for k in sorted(c.items)}, label="items: map[key] to node",
               hl={hl_key: style} if hl_key in c.items else {}),
            LinkedList(lst or ["—"], label="order (front = most recent, back = least recent)", hl=ll_hl,
                       pointers=({"front": 0, "back": len(lst) - 1} if len(lst) > 1 else
                                 ({"front": 0} if lst else {}))),
            Text(*lines),
            say=say,
        )

    yield Frame(
        KV({}, label="items: map[key] to node"),
        LinkedList(["—"], label="order (front = most recent, back = least recent)"),
        Text("capacity = 2 · Get(key) · Put(key, value) · both O(1)"),
        speaker="Interviewer",
        say="Build a cache with capacity 2. Get and Put must be O(1). When it's full, Put evicts the "
            "least recently used key.",
    )
    yield Frame(
        KV({}, label="items: map[key] to node"),
        LinkedList(["—"], label="order (front = most recent, back = least recent)"),
        Text("map + timestamps: eviction scans for the oldest → O(n)",
             "slice in recency order: moving a key → O(n)"),
        say="A map alone can't tell me the oldest key without a scan. A slice keeps order, but moving a key is O(n). "
            "A map pointing into a doubly linked list does both in O(1).",
    )

    ops = [("put", 1, 1), ("put", 2, 2), ("get", 1, 1), ("put", 3, 3), ("get", 2, None),
           ("put", 4, 4), ("get", 1, None), ("get", 3, 3), ("get", 4, 4)]
    says = {
        0: "Put(1,1): new node at the front, and the map points to it.",
        1: "Put(2,2): another new node at the front. Key 1 slides toward the back.",
        2: "Get(1) finds the node through the map and moves it to the front in O(1). Now 2 is least recent.",
        3: "Put(3,3) when full: pop the back node, 2, and delete it from the map. "
           "That's why each node stores its key. Then insert 3 at the front.",
        4: "Get(2) misses: 2 was evicted. A miss doesn't change the order.",
        5: "Put(4,4) is full again. The back is 1, so 1 is evicted and 4 goes to the front.",
        6: "Get(1) misses, as expected.",
        7: "Get(3) hits and moves 3 to the front.",
        8: "Get(4) hits and moves 4 to the front. Every call was a map lookup plus a few pointer changes.",
    }
    for k, (op, key, expect) in enumerate(ops):
        if op == "put":
            evicted = c.put(key, expect)
            extra = (f"evicted {evicted} from the back", "reject") if evicted is not None else None
            yield view(f"Put({key}, {expect})", says[k], key, "active", extra)
        else:
            r = c.get(key)
            assert r == expect, (op, key, r, expect)
            if r is None:
                yield view(f"Get({key}) → miss", says[k], None, "reject")
            else:
                yield view(f"Get({key}) → {r}", says[k], key, "match")

    yield Frame(
        KV({k: f"node({k})" for k in sorted(c.items)}, label="items: map[key] to node"),
        LinkedList([f"{k}:{c.items[k]}" for k in c.order], label="order (front = most recent, back = least recent)",
                   pointers={"front": 0, "back": 1}),
        Text("Get O(1) · Put O(1) · space O(capacity)", "sync.Mutex, not RWMutex: Get moves nodes"),
        say="Both operations are O(1) and space is O(capacity). For goroutines I use a plain Mutex: "
            "Get mutates the list, so a read lock would be a data race.",
    )
