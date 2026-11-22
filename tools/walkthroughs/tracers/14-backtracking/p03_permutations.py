from walkthrough import Array, Frame, KV, Row, Text

TITLE = "Permutations I & II"
PATTERN = "Backtracking + used[]"


def fmt(path):
    return "[" + " ".join(map(str, path)) + "]"


def trace(nums):
    """PermuteUnique, yielding at every choose, leaf and duplicate skip."""
    s = sorted(nums)
    res, path, used = [], [], [False] * len(s)

    def backtrack():
        if len(path) == len(s):
            res.append(list(path))
            yield "leaf", list(path), list(used), list(res), None
            return
        for i, x in enumerate(s):
            if used[i]:
                continue
            if i > 0 and s[i] == s[i - 1] and not used[i - 1]:
                yield "skip", list(path), list(used), list(res), i
                continue
            used[i] = True
            path.append(x)
            yield "choose", list(path), list(used), list(res), i
            yield from backtrack()
            path.pop()
            used[i] = False

    yield from backtrack()
    trace.result = res


def frames():
    nums = [1, 1, 2]
    s = sorted(nums)

    def view(path, used, res, i=None, style="active", new=None):
        hl = {k: "visited" for k, u in enumerate(used) if u}
        if i is not None:
            hl[i] = style
        return [
            Row(Array(s, label="sorted nums (grey = used)", hl=hl, cell=52),
                Array(path or ["—"], label="path", show_index=False, cell=52), weights=[1, 1]),
            Array(used, label="used[]", hl={k: "active" for k, u in enumerate(used) if u}, cell=52),
            KV([fmt(r) for r in res], label=f"results ({len(res)})", hl={new: "match"} if new else None, empty="[ ]"),
        ]

    yield Frame(
        Array(nums, label="nums", cell=52),
        Text("return every DISTINCT ordering once"),
        speaker="Interviewer",
        say="The list may have repeated values. Return every distinct ordering exactly once.",
    )
    yield Frame(
        Array(nums, label="nums", cell=52),
        Text("brute force: all n! orderings, dedupe with a set",
             ("3! = 6 orderings, but only 3 distinct; the two 1s swap places", "reject")),
        say="Generating all n factorial orderings and deduping with a set works, but half the work here is wasted. "
            "I'd rather never build a duplicate.",
    )

    events = list(trace(nums))
    first_choose = True
    for ev, path, used, res, i in events:
        if ev == "choose" and first_choose:
            first_choose = False
            yield Frame(*view(path, used, res, i=i),
                        Text(f"choose nums[{i}] = {s[i]}, used[{i}] = true"),
                        say="I sort, then fill one position at a time with any unused element. "
                            "used[] marks what the path already holds. First choice: the first 1.")
        elif ev == "leaf":
            f = fmt(path)
            yield Frame(*view(path, used, res, new=f),
                        Text(f"len(path) == n → record {f}"),
                        say=f"The path has all {len(s)} elements, so {f} is an answer. "
                            "Then I unchoose: pop the last value and clear its used flag.")
        elif ev == "skip":
            yield Frame(*view(path, used, res, i=i, style="reject"),
                        Text(f"nums[{i}] == nums[{i - 1}] and used[{i - 1}] is false",
                             ("skip: equal copies must be used left to right", "reject")),
                        say=f"The second 1 equals the first, and the first is not in the path. "
                            "Choosing it here would repeat a branch I already explored, so I skip it.")

    res = trace.result
    assert sorted(res) == [[1, 1, 2], [1, 2, 1], [2, 1, 1]]
    yield Frame(
        KV([fmt(r) for r in res], label=f"results ({len(res)})", hl={fmt(r): "match" for r in res}),
        Text("3 distinct permutations, no duplicates generated", "O(n · n!) time, O(n) for path + used[]"),
        say="Three distinct orderings, each produced once. There are at most n factorial leaves and each copy costs O(n): "
            "O(n times n factorial) time, O(n) extra space.",
    )
