from walkthrough import KV, Array, Frame, Text

TITLE = "Top K Frequent Elements"
PATTERN = "Bucket sort"


def bucket_view(buckets, hl=None, pointers=None):
    vals = ["[" + ",".join(map(str, b)) + "]" if b else "·" for b in buckets]
    return Array(vals, label="buckets[f] = values seen exactly f times", hl=hl or {}, pointers=pointers or {}, cell=64)


def frames():
    nums, k = [1, 1, 1, 2, 2, 3], 2

    yield Frame(
        Array(nums, label="nums"),
        Text(f"k = {k}"),
        speaker="Interviewer",
        say="Return the k most frequent values, most frequent first. Break ties by the smaller value.",
    )

    # Counter(nums)
    counts = {}
    for x in nums:
        counts[x] = counts.get(x, 0) + 1
    yield Frame(
        Array(nums, label="nums", hl={i: "visited" for i in range(len(nums))}),
        KV(counts, label="counts (value: times)"),
        Text("sort the distinct values by count → O(n + d log d)", ("works, but the sort isn't needed", "compare")),
        say="First I count with a map: one pass. I could then sort the values by count, "
            "but a count can't exceed n, so I can bucket by count instead.",
    )

    buckets = [[] for _ in range(len(nums) + 1)]
    yield Frame(
        KV(counts, label="counts (value: times)"),
        bucket_view(buckets),
        Text(f"len(nums) = {len(nums)} → {len(nums) + 1} buckets, index = frequency"),
        say=f"I make n + 1 = {len(nums) + 1} buckets. Bucket f will hold every value that appears exactly f times.",
    )
    for v, c in counts.items():
        buckets[c].append(v)
    yield Frame(
        KV(counts, label="counts (value: times)", hl={v: "active" for v in counts}),
        bucket_view(buckets, hl={c: "active" for c in counts.values()}),
        Text("  ".join(f"{v} → buckets[{c}]" for v, c in counts.items())),
        say="Each value drops into the bucket for its count. No comparisons, so this is O(n).",
    )

    res = []
    f = len(buckets) - 1
    skipped = []
    while f > 0 and len(res) < k:
        b = sorted(buckets[f])
        if not b:
            skipped.append(f)
            f -= 1
            continue
        if skipped:
            yield Frame(
                bucket_view(buckets, hl={s: "visited" for s in skipped}, pointers={"f": skipped[-1]}),
                Text(f"buckets {skipped[0]}..{skipped[-1]} are empty → skip", f"res = {res}"),
                say="Now I read buckets from the highest frequency down. The top ones are empty, so I skip them.",
            )
            skipped = []
        took = []
        for v in b:
            if len(res) == k:
                break
            res.append(v)
            took.append(v)
        done = len(res) == k
        yield Frame(
            bucket_view(buckets, hl={f: "match"}, pointers={"f": f}),
            Text(f"buckets[{f}] = {b} → take {took}", (f"res = {res}" + ("  (len = k, stop)" if done else ""), "match")),
            say=(f"Bucket {f} holds {b[0]}: it appears {f} times, the most of any value. It goes first."
                 if len(res) == 1 and not done else
                 f"Bucket {f} gives {took[-1]}. Now I have k = {k} values, so I stop."
                 if done else f"Bucket {f} gives {took}."),
        )
        f -= 1

    assert res == [1, 2]
    yield Frame(
        Array(res, label="answer"),
        Text((f"answer: {res}", "match"),
             "TopKFrequentSort: O(n + d log d) time, O(n) space",
             "TopKFrequent:     O(n) time (+ sorting ties in a bucket), O(n) space"),
        say=f"The answer is {res}. Counting and bucketing are O(n), and so is the walk over n + 1 buckets. "
            "Space is O(n) for the map and buckets.",
    )
