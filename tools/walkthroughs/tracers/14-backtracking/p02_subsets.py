from walkthrough import Array, Frame, KV, Row, Text

TITLE = "Subsets I & II"
PATTERN = "Backtracking + dedup"


def fmt(path):
    return "[" + " ".join(map(str, path)) + "]"


def trace(nums):
    """SubsetsWithDup, yielding an event at every record and every skipped sibling."""
    s = sorted(nums)
    res, path = [], []

    def backtrack(start):
        res.append(list(path))  # every partial path is a subset
        yield "record", list(path), list(res), start, None
        for i in range(start, len(s)):
            if i > start and s[i] == s[i - 1]:
                yield "skip", list(path), list(res), start, i
                continue
            path.append(s[i])
            yield from backtrack(i + 1)
            path.pop()

    yield from backtrack(0)
    trace.result = res


def frames():
    nums = [1, 2, 2]
    s = sorted(nums)

    def view(path, res, start, i=None, new=None, style="active"):
        hl = {}
        if i is not None:
            hl[i] = style
            hl[i - 1] = "compare"
        return [
            Row(Array(s, label="sorted nums", hl=hl, pointers={"start": start} if start < len(s) else {}, cell=52),
                Array(path or ["—"], label="path", show_index=False, cell=52), weights=[1, 1]),
            KV([fmt(r) for r in res], label=f"results ({len(res)})", hl={new: "match"} if new else None),
        ]

    yield Frame(
        Array(nums, label="nums", cell=52),
        Text("return every subset exactly once ([1 2] from either 2 is the same subset)"),
        speaker="Interviewer",
        say="The list may contain repeated values. Return every distinct subset exactly once, in any order.",
    )
    yield Frame(
        Array(nums, label="nums", cell=52),
        Text("brute force: all 2^n bitmasks, then dedupe with a set of sorted keys",
             ("[1 2] appears twice (mask 011 and 101)", "reject")),
        say="I could build all 2 to the n subsets and dedupe with a hash set, but that does extra work and needs "
            "keys for slices. Better: never generate a duplicate in the first place.",
    )

    events = list(trace(nums))
    for ev, path, res, start, i in events:
        if ev == "record":
            f = fmt(path)
            if not path:
                say = ("I sort first so equal values sit together. Every call records its path, so the empty "
                       "subset is recorded at the root.")
            elif len(path) == len(s):
                say = f"Record {f}. start is past the end, so this branch is done; I unchoose and go back up."
            else:
                say = (f"Choose {path[-1]} and recurse with start = {start}. Every node in the tree is an answer, "
                       f"so I record {f} right away.")
            yield Frame(*view(path, res, start, new=f), Text(f"record {f}", f"loop i from start = {start}"), say=say)
        else:
            say = (f"At this depth I already tried a {s[i]} as the next element. Taking the second {s[i]} here "
                   f"would rebuild the same subsets, so I skip it: i > start and nums[i] == nums[i-1].")
            yield Frame(*view(path, res, start, i=i, style="reject"),
                        Text(f"i = {i} > start = {start} and nums[{i}] == nums[{i - 1}] = {s[i]}",
                             ("skip: same value already tried at this depth", "reject")), say=say)

    res = trace.result
    expected = [[], [1], [2], [1, 2], [2, 2], [1, 2, 2]]
    assert sorted(res) == sorted(expected)
    yield Frame(
        KV([fmt(r) for r in res], label=f"results ({len(res)})", hl={fmt(r): "match" for r in res}),
        Text(f"{len(res)} distinct subsets, no duplicates generated", "O(n · 2^n) time, O(n) recursion besides the output"),
        say="Six distinct subsets and no duplicates. There are at most 2 to the n subsets and copying each costs O(n), "
            "so it's O(n times 2 to the n) time and O(n) extra space.",
    )
