from walkthrough import Array, Frame, Text

TITLE = "Rate Limiter (token bucket)"
PATTERN = "Lazy refill"


class TokenBucket:
    def __init__(self, capacity, rate, now):
        self.capacity = float(capacity)
        self.rate = rate
        self.tokens = float(capacity)
        self.now = now
        self.last = now()

    def refill(self):
        now = self.now()
        elapsed = now - self.last
        added = 0.0
        if elapsed > 0:
            before = self.tokens
            self.tokens = min(self.capacity, self.tokens + elapsed * self.rate)
            added = self.tokens - before
            self.last = now
        return elapsed, added

    def allow_n(self, n):
        elapsed, added = self.refill()
        if self.tokens >= n:
            self.tokens -= n
            return True, elapsed, added
        return False, elapsed, added


def frames():
    clock = [0.0]
    b = TokenBucket(3, 1.0, lambda: clock[0])

    def bucket(hl=None):
        full = int(b.tokens)
        cells = ["●" if i < full else "" for i in range(3)]
        h = {i: "match" for i in range(full)}
        h.update({i: "empty" for i in range(full, 3)})
        h.update(hl or {})
        return Array(cells, label="bucket (capacity 3, refills 1 token/s)", hl=h, cell=56)

    def view(call, say, ok=None, hl=None, extra=None):
        lines = [f"t = {clock[0]:g}s   tokens = {b.tokens:g}   last = {b.last:g}s"]
        if extra:
            lines.append(extra)
        lines.append((call, "match" if ok else "reject") if ok is not None else call)
        return Frame(bucket(hl), Text(*lines), say=say)

    yield Frame(
        bucket(),
        Text("capacity 3, rate 1/s, starts full", "Allow() / AllowN(n): answer now, never block"),
        speaker="Interviewer",
        say="Build a rate limiter that allows bursts of up to 3 requests, then 1 per second. Allow answers "
            "immediately. Many goroutines will call it.",
    )
    yield Frame(
        bucket(),
        Text("naive: a ticker goroutine adds a token every second", ("an extra goroutine per limiter, hard to test", "reject")),
        say="A ticker goroutine that adds tokens works, but it costs a goroutine per limiter and makes tests "
            "depend on real time. I'll refill lazily: on each call, add tokens for the time that passed.",
    )

    results = []
    for k in range(4):
        before = b.tokens
        ok, _, _ = b.allow_n(1)
        results.append(ok)
        if ok:
            say = (f"t = 0. Refill adds nothing, since no time passed. {before:g} ≥ 1, so I spend a token and allow."
                   if k == 0 else f"Same instant, {before:g} tokens left before this call. Spend one, allow.")
            yield view("Allow() → true", say, True, {int(b.tokens): "active"})
        else:
            yield view("Allow() → false", "The burst is used up: 0 tokens, and no time has passed. "
                                          "Reject without blocking.", False)
    assert results == [True, True, True, False]

    clock[0] = 1.0
    ok, elapsed, added = b.allow_n(1)
    assert ok
    yield view("Allow() → true", "One second later. Refill adds elapsed × rate = 1 token, then the request "
                                 "spends it. last moves to 1s.", True, {0: "active"},
               f"refill: elapsed {elapsed:g}s × 1/s = +{added:g} token")

    clock[0] = 10.0
    ok, elapsed, added = b.allow_n(5)
    assert not ok
    yield view("AllowN(5) → false", "Nine seconds later, refill would add 9, but the bucket caps at 3. "
                                    "5 > 3, and AllowN is all or nothing, so nothing is spent.", False, None,
               f"refill: min(3, 0 + {elapsed:g} × 1) = {b.tokens:g}")
    yield Frame(
        bucket(),
        Text("tokens = min(capacity, tokens + elapsed·rate); spend if tokens ≥ n",
             "O(1) per call · O(1) space · one sync.Mutex"),
        say="Each call is O(1) time and space. Refill, check and spend must happen as one step, "
            "so they run under a single mutex. The clock is injected, so tests need no sleeps.",
    )
