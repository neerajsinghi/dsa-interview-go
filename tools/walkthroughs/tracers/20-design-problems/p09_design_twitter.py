from walkthrough import Array, Frame, KV, Row, Text

TITLE = "Design Twitter"
PATTERN = "Heap k-way merge"

FEED_SIZE = 10


class Twitter:
    def __init__(self):
        self.clock = 0
        self.tweets = {}     # user -> [(id, time)] oldest first
        self.following = {}  # follower -> set

    def post(self, user, tid):
        self.clock += 1
        self.tweets.setdefault(user, []).append((tid, self.clock))

    def follow(self, a, b):
        if a == b:
            return
        self.following.setdefault(a, set()).add(b)

    def unfollow(self, a, b):
        self.following.get(a, set()).discard(b)

    def feed_steps(self, user):
        """Port of GetNewsFeed. Yields (event, heap snapshot, feed) and returns the feed."""
        h = []  # cursors [user, list, pos]; max-heap by list[pos].time

        def key(c):
            return c[1][c[2]][1]

        def sift_down(i):
            n = len(h)
            while True:
                l, r, big = 2 * i + 1, 2 * i + 2, i
                if l < n and key(h[l]) > key(h[big]):
                    big = l
                if r < n and key(h[r]) > key(h[big]):
                    big = r
                if big == i:
                    return
                h[i], h[big] = h[big], h[i]
                i = big

        for u in [user] + sorted(self.following.get(user, ())):
            lst = self.tweets.get(u, [])
            if lst:
                h.append([u, lst, len(lst) - 1])
        for i in range(len(h) // 2 - 1, -1, -1):
            sift_down(i)
        feed = []
        steps = [("init", [list(c) for c in h], list(feed))]
        while h and len(feed) < FEED_SIZE:
            top = h[0]
            feed.append(top[1][top[2]][0])
            if top[2] == 0:
                h[0] = h[-1]
                h.pop()
                if h:
                    sift_down(0)
            else:
                top[2] -= 1
                sift_down(0)
            steps.append(("pop", [list(c) for c in h], list(feed)))
        return feed, steps


def frames():
    tw = Twitter()

    def state(heap=None, hl_user=None, feed=None, call=None, call_style=None, lines=()):
        users = [1, 2]
        arrays = []
        for u in users:
            lst = tw.tweets.get(u, [])
            hl = {}
            if heap is not None:
                for c in heap:
                    if c[0] == u:
                        hl[c[2]] = "active"
            arrays.append(Array([f"{tid}@t{t}" for tid, t in lst] or ["—"], label=f"tweets[{u}] (oldest first)",
                                hl=hl, cell=70, show_index=False))
        heap_cells = [f"u{c[0]}:{c[1][c[2]][0]}@t{c[1][c[2]][1]}" for c in heap] if heap is not None else []
        follows = {f"{a}": "{" + ",".join(str(x) for x in sorted(s)) + "}" for a, s in tw.following.items()}
        out = [
            Row(*arrays, KV(follows, label="following", empty="{ }"), weights=[1, 1, 1]),
            Array(heap_cells or ["empty"], label="max-heap of cursors (by time)", hl={0: "special"} if heap_cells else {},
                  cell=80, show_index=False),
        ]
        txt = list(lines)
        if call:
            txt.append((call, call_style))
        out.append(Text(*txt) if txt else Text(""))
        return out

    yield Frame(
        *state(lines=["PostTweet · Follow · Unfollow · GetNewsFeed (10 newest, newest first)"]),
        speaker="Interviewer",
        say="Users post tweets and follow each other. GetNewsFeed returns the 10 most recent tweet ids from "
            "the user and everyone they follow, newest first.",
    )
    yield Frame(
        *state(lines=["naive: gather every followee's tweets, sort all → O(T log T)"]),
        say="I could gather all tweets from everyone the user follows and sort them. That's O(T log T) for T tweets, "
            "but I only need 10. Each user's list is already sorted by time.",
    )

    tw.post(1, 5)
    yield Frame(*state(call="PostTweet(1, 5)", call_style="active"),
                say="PostTweet bumps a global clock and appends (5, t1) to user 1's list. Lists stay sorted for free.")

    feed, steps = tw.feed_steps(1)
    assert feed == [5]
    yield Frame(*state(heap=steps[-1][1], call="GetNewsFeed(1) → [5]", call_style="match"),
                say="User 1 follows nobody yet, so the merge has one list. The feed is [5].")

    tw.follow(1, 2)
    yield Frame(*state(call="Follow(1, 2)", call_style="active"),
                say="Follow adds 2 to user 1's followee set. Nothing is copied: this is fan-out on read.")
    tw.post(2, 6)
    yield Frame(*state(call="PostTweet(2, 6)", call_style="active"),
                say="User 2 posts tweet 6 at time 2, so it's newer than tweet 5.")

    feed, steps = tw.feed_steps(1)
    assert feed == [6, 5]
    for k, (ev, heap, f) in enumerate(steps):
        if ev == "init":
            yield Frame(*state(heap=heap, lines=["one cursor per user, at their newest tweet; heap.Init is O(f)"],
                               call="GetNewsFeed(1)", call_style="active"),
                        say="GetNewsFeed puts one cursor per user, at their newest tweet, into a max-heap by time. "
                            "The top is the newest tweet overall.")
        else:
            last = k == len(steps) - 1
            yield Frame(*state(heap=heap, lines=[f"pop newest → feed = {f}"],
                               call=("GetNewsFeed(1) → [6 5]" if last else None),
                               call_style="match"),
                        say=(f"Pop the top: tweet {f[-1]}. That user has no older tweet, so the cursor leaves the heap. "
                             + ("The heap is empty, so the feed is [6, 5]." if last else
                                "Otherwise I'd step it back one and sift down.")))

    tw.unfollow(1, 2)
    feed, steps = tw.feed_steps(1)
    assert feed == [5]
    yield Frame(*state(heap=steps[-1][1], call="Unfollow(1, 2); GetNewsFeed(1) → [5]", call_style="match"),
                say="Unfollow removes 2 from the set. User 2's tweets still exist, but the next feed only merges "
                    "user 1's list: [5].")
    yield Frame(*state(lines=["Post/Follow/Unfollow O(1) · GetNewsFeed O(f + 10·log f)",
                              "sync.RWMutex: feeds are built under a read lock"]),
                say="Posting and following are O(1). A feed costs O(f) to build the heap plus 10 pops at O(log f). "
                    "Feeds only read, so they run in parallel under an RWMutex.")
