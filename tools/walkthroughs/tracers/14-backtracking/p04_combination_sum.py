from walkthrough import Array, Frame, KV, Row, Text

TITLE = "Combination Sum"
PATTERN = "Backtracking + prune"


def fmt(path):
    return "[" + " ".join(map(str, path)) + "]"


def trace(candidates, target):
    """CombinationSum, yielding at every choose, every answer and every sorted break."""
    cands = sorted(candidates)
    res, path = [], []

    def backtrack(start, remaining):
        if remaining == 0:
            res.append(list(path))
            yield "found", list(path), list(res), start, remaining, None
            return
        for i in range(start, len(cands)):
            if cands[i] > remaining:
                yield "break", list(path), list(res), start, remaining, i
                break
            path.append(cands[i])
            yield "choose", list(path), list(res), start, remaining - cands[i], i
            yield from backtrack(i, remaining - cands[i])  # i, not i+1: reuse allowed
            path.pop()

    yield from backtrack(0, target)
    trace.result = res


def frames():
    candidates, target = [2, 3, 6, 7], 7
    cands = sorted(candidates)

    def view(path, res, rem, i=None, style="active", new=None):
        hl = {i: style} if i is not None else {}
        if style == "reject" and i is not None:
            for j in range(i + 1, len(cands)):
                hl[j] = "reject"
        return [
            Row(Array(cands, label="sorted candidates", hl=hl, pointers={"i": i} if i is not None else {}, cell=52),
                Array(path or ["—"], label=f"path   (remaining = {rem})", show_index=False, cell=52),
                weights=[1, 1.2]),
            KV([fmt(r) for r in res], label=f"results ({len(res)})", hl={new: "match"} if new else None, empty="[ ]"),
        ]

    yield Frame(
        Array(candidates, label="candidates", cell=52),
        Text(f"target = {target}; each candidate may be reused"),
        speaker="Interviewer",
        say="Return every combination of candidates that adds up to 7. You can reuse a number, "
            "and [2 2 3] and [3 2 2] count as the same answer.",
    )
    yield Frame(
        Array(candidates, label="candidates", cell=52),
        Text("brute force: try every sequence of picks up to depth T/min",
             ("generates [2 2 3], [2 3 2], [3 2 2] … then must dedupe", "reject")),
        say="Trying every sequence of picks finds each combination many times, once per ordering. "
            "If I only ever pick in non-decreasing order, each combination comes out once.",
    )

    events = list(trace(candidates, target))
    found = 0
    skipped_choose = skipped_break = 0
    for ev, path, res, start, rem, i in events:
        if found >= 1 and ev != "found" and not (path == [7] or (ev == "choose" and len(path) == 1)):
            # After the first answer, show only top-level choices and answers.
            skipped_choose += ev == "choose"
            skipped_break += ev == "break"
            continue
        if ev == "choose":
            if len(path) == 1 and found:
                say = (f"Back at the root I move on to {path[0]}. Its subtree only uses {path[0]} or larger, "
                       f"so it can't rebuild an answer that starts with a smaller number.")
                extra = f"(not shown since last step: {skipped_choose} pick(s), {skipped_break} break(s), no answers)"
                skipped_choose = skipped_break = 0
            else:
                say = (f"Choose {cands[i]}, so {rem} remains. I recurse with start = i, not i + 1, "
                       "so the same number may be picked again.")
                extra = f"recurse with start = {i}"
            yield Frame(*view(path, res, rem, i=i), Text(f"choose {cands[i]}: remaining {rem + cands[i]} - {cands[i]} = {rem}",
                                                         extra), say=say)
        elif ev == "break":
            yield Frame(*view(path, res, rem, i=i, style="reject"),
                        Text(f"{cands[i]} > remaining {rem}", ("break: every later candidate is even bigger", "reject")),
                        say=f"{cands[i]} is bigger than the {rem} that remains. The list is sorted, so everything after it "
                            "is bigger too. I break instead of continue, and unchoose.")
        else:
            found += 1
            f = fmt(path)
            yield Frame(*view(path, res, 0, new=f), Text(f"remaining == 0 → record {f}"),
                        say=f"Remaining hit 0, so {f} is an answer. I copy the path into the results and unchoose.")

    res = trace.result
    assert res == [[2, 2, 3], [7]]
    yield Frame(
        KV([fmt(r) for r in res], label=f"results ({len(res)})", hl={fmt(r): "match" for r in res}),
        Text("answer = [2 2 3] [7]", "O(n^(T/m)) worst case time, O(T/m) recursion  (m = smallest candidate)"),
        say="Two answers. The tree is at most T over m levels deep, where m is the smallest candidate, "
            "so the worst case is n to the T over m. The sorted break cuts most of it.",
    )
