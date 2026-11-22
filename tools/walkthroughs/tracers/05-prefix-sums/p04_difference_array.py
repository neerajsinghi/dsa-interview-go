from walkthrough import Array, Frame, Text

TITLE = "Range Addition (difference array)"
PATTERN = "Difference array"


def frames():
    n = 5
    updates = [(1, 3, 2), (2, 4, 3), (0, 2, -2)]
    want = [-2, 0, 3, 5, 3]
    upd_txt = "updates: " + "  ".join(f"[{s},{e},{d:+d}]" for s, e, d in updates)

    yield Frame(
        Array([0] * n, label="out (n = 5)"),
        Text(upd_txt),
        speaker="Interviewer",
        say="Start with five zeros. Each update [start, end, delta] adds delta to every position "
            "from start to end. Return the array after all updates.",
    )
    yield Frame(
        Array([0, 2, 2, 2, 0], label="out", hl={1: "compare", 2: "compare", 3: "compare"}, window=(1, 3)),
        Text(upd_txt, "brute force: loop over every position of every update → O(n·q)"),
        say="Looping over each range touches up to n cells per update, so O(n·q). With long ranges "
            "and many updates that's slow.",
    )

    diff = [0] * (n + 1)
    yield Frame(
        Array(diff, label="diff (n + 1 slots)"),
        Text(upd_txt),
        say="Instead I keep a difference array with one extra slot. An update only marks where its "
            "delta starts and where it stops, two O(1) writes.",
    )
    for k, (s, e, d) in enumerate(updates):
        diff[s] += d
        diff[e + 1] -= d
        yield Frame(
            Array(diff, label="diff", hl={s: "match", e + 1: "reject"}),
            Text(f"update [{s},{e},{d:+d}]:  diff[{s}] += {d},  diff[{e + 1}] −= {d}"),
            say=f"Update {k + 1}: +{d} at index {s} turns the delta on, and −{d} at index {e + 1} turns "
                "it off just after the range ends." if d > 0 else
                f"Update {k + 1} has a negative delta: {d} at index {s} and +{-d} at index {e + 1}. "
                "Same two marks, signs flipped.",
        )

    out = [None] * n
    running = 0
    for i in range(n):
        running += diff[i]
        out[i] = running
        if i < 2:
            yield Frame(
                Array(diff, label="diff", hl={i: "active", **{j: "visited" for j in range(i)}}),
                Array(out, label="out", hl={i: "match"}),
                Text(f"running += diff[{i}] = {running}  →  out[{i}] = {running}"),
                say=f"Now one prefix-sum pass. The running total at index {i} is {running}: every delta "
                    "that has started and not yet stopped." if i == 0 else
                    f"running = {running - diff[i]} + {diff[i]} = {running}. A prefix sum spreads each "
                    "mark over the cells after it.",
            )
    yield Frame(
        Array(diff, label="diff", hl={j: "visited" for j in range(n)}),
        Array(out, label="out", hl={2: "match", 3: "match", 4: "match"}),
        Text(f"out[2..4]: running = {out[2]}, {out[3]}, {out[4]}"),
        say="Same again for indices 2 to 4. The last slot of diff is never read; it only exists so "
            "end + 1 never goes out of range.",
    )

    assert out == want, out
    yield Frame(
        Array(out, label="out", hl={i: "match" for i in range(n)}),
        Text(f"answer: {out}", ("time O(n + q) · space O(n)", "match")),
        say="The answer is [-2, 0, 3, 5, 3]. Each update is O(1) and the final pass is O(n), "
            "so O(n + q) time and O(n) space.",
    )
