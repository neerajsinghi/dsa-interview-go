from walkthrough import KV, Array, Frame, Text

TITLE = "Contains Duplicate, three ways"
PATTERN = "Time vs space"


def frames():
    nums = [1, 2, 3, 1]

    yield Frame(
        Array(nums, label="nums"),
        speaker="Interviewer",
        say="Report whether any value appears at least twice. Show me more than one way and compare them.",
    )

    # ContainsDuplicateBrute: count the comparisons it really makes.
    comparisons, brute = 0, False
    pair = None
    for i in range(len(nums)):
        for j in range(i + 1, len(nums)):
            comparisons += 1
            if nums[i] == nums[j]:
                brute, pair = True, (i, j)
                break
        if brute:
            break
    n = len(nums)
    yield Frame(
        Array(nums, label="nums", hl={pair[0]: "match", pair[1]: "match"}, pointers={"i": pair[0], "j": pair[1]}),
        Text(f"brute force: compare every pair, up to n(n-1)/2 = {n * (n - 1) // 2}",
             f"here it found nums[{pair[0]}] == nums[{pair[1]}] after {comparisons} comparisons",
             ("O(n²) time, O(1) space", "reject")),
        say="Way one compares every pair. No extra memory, but the pairs grow as n squared, "
            "so 10⁵ values means about 5 billion comparisons.",
    )

    # ContainsDuplicateSort: clone, sort, check neighbours.
    s = sorted(nums)
    sort_ans, at = False, None
    for i in range(1, len(s)):
        if s[i] == s[i - 1]:
            sort_ans, at = True, i
            break
    yield Frame(
        Array(nums, label="nums (untouched)"),
        Array(s, label="sorted copy", hl={at - 1: "match", at: "match"}, pointers={"i-1": at - 1, "i": at}),
        Text(f"s[{at - 1}] == s[{at}] == {s[at]} → true", ("O(n log n) time, O(n) space for the copy", "compare")),
        say="Way two sorts a copy. Equal values become neighbours, so one pass over adjacent pairs finds them. "
            "Sorting costs O(n log n).",
    )

    # ContainsDuplicate: hash set, stop at the first repeat.
    seen = {}
    yield Frame(
        Array(nums, label="nums"),
        KV([], label="seen (hash set)", empty="{ } empty"),
        say="Way three trades memory for speed. I keep a set of values seen so far; "
            "a lookup in a Go map is O(1) on average.",
    )
    ans = False
    for i, x in enumerate(nums):
        if x in seen:
            ans = True
            yield Frame(
                Array(nums, label="nums", hl={i: "match", nums.index(x): "match"}, pointers={"x": i}),
                KV(list(seen), label="seen (hash set)", hl={x: "match"}),
                Text((f"{x} is already in seen → return true", "match")),
                say=f"{x} is already in the set, so this is a repeat. I return true right away "
                    "without looking at anything else.",
            )
            break
        seen[x] = True
        yield Frame(
            Array(nums, label="nums", hl={i: "active"} | {k: "visited" for k in range(i)}, pointers={"x": i}),
            KV(list(seen), label="seen (hash set)", hl={x: "active"}),
            Text(f"{x} not in seen → add {x}"),
            say=f"{x} isn't in the set yet, so it's new. I add it and move on." if i == 0
            else f"{x} is new too. Add it and keep going.",
        )

    assert ans is True and brute is True and sort_ans is True
    yield Frame(
        Text(("answer: true", "match"),
             "ContainsDuplicateBrute  O(n²)       O(1)",
             "ContainsDuplicateSort   O(n log n)  O(n)  copy",
             "ContainsDuplicate       O(n)        O(n)  set"),
        say="All three say true. The hash set is fastest at O(n) time but uses O(n) memory. "
            "If memory is tight, sorting is the middle ground.",
    )
