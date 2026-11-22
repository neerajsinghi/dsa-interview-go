from walkthrough import Array, Frame, KV, Text

TITLE = "Subarray Sum Equals K"
PATTERN = "Prefix + hash map"


def frames():
    nums, k = [1, 1, 1], 2
    want = 2

    yield Frame(
        Array(nums, label="nums"),
        Text(f"k = {k}"),
        speaker="Interviewer",
        say="Count the contiguous, non-empty subarrays whose sum is exactly k. The numbers can be "
            "negative or zero.",
    )
    yield Frame(
        Array(nums, label="nums", hl={0: "compare", 1: "compare"}, window=(0, 1)),
        Text(f"k = {k}", "brute force: every start, extend a running sum → O(n²)",
             ("sliding window fails: negatives break the monotone sum", "reject")),
        say="Trying every start is O(n²). A sliding window won't work, because with negative numbers "
            "growing or shrinking the window doesn't move the sum in one direction.",
    )

    seen = {0: 1}
    yield Frame(
        Array(nums, label="nums"),
        KV(seen, label="seen  (prefix: count)", hl={0: "special"}),
        Text("subarray (i, j] sums to k  ⇔  P[i] = P[j] − k"),
        say="A subarray sums to k when two prefix sums differ by k. So at each step I ask how many "
            "earlier prefixes equal prefix − k. I seed the map with the empty prefix 0.",
    )

    count, prefix = 0, 0
    where = {0: [-1]}  # display only: index after which each prefix value occurred
    for j, x in enumerate(nums):
        prefix += x
        need = prefix - k
        found = seen.get(need, 0)
        before = dict(seen)
        count += found
        seen[prefix] = seen.get(prefix, 0) + 1
        start = where[need][-1] + 1 if found else None
        where.setdefault(prefix, []).append(j)
        hl = {need: "match"} if found else {}
        line = (f"seen[{need}] = {found} → count = {count}", "match" if found else "reject")
        if found:
            say = (f"prefix is {prefix}, so I need an earlier prefix of {need}. It appeared {'once' if found == 1 else f'{found} times'}, "
                   f"so count becomes {count}. Then I record prefix {prefix}.")
        else:
            say = (f"prefix is {prefix}, so I need {need}. It's not in the map, so no subarray ends here. "
                   f"I record prefix {prefix} after the lookup.")
        yield Frame(
            Array(nums, label="nums", hl={**{i: "visited" for i in range(j)}, j: "active"},
                  pointers={"j": j}, window=(start, j) if found else None),
            KV(before, label="seen before lookup", hl=hl),
            Text(f"prefix = {prefix}   need prefix − k = {prefix} − {k} = {need}", line),
            say=say,
        )

    assert count == want
    yield Frame(
        Array(nums, label="nums"),
        KV(seen, label="seen (final)"),
        Text(f"answer: {count}   ([1,1] at 0..1 and 1..2)", ("time O(n) · space O(n)", "match")),
        say="The answer is 2. One pass with one map lookup per element: O(n) time and O(n) space "
            "for the map.",
    )
