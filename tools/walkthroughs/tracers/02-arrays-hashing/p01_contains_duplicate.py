from walkthrough import KV, Array, Frame, Text

TITLE = "Contains Duplicate"
PATTERN = "Hash set"


def frames():
    nums = [1, 2, 3, 1]

    yield Frame(
        Array(nums, label="nums"),
        speaker="Interviewer",
        say="Given a slice of integers, report whether any value occurs more than once.",
    )

    n = len(nums)
    yield Frame(
        Array(nums, label="nums", pointers={"i": 0, "j": n - 1}, hl={0: "compare", n - 1: "compare"}),
        Text(f"brute force: compare every pair → n(n-1)/2 = {n * (n - 1) // 2} pairs here",
             ("O(n²) time, O(1) space", "reject")),
        say="The brute force compares every pair. That's O(n²). The inner loop is really asking "
            "'have I seen this value before?', and a set answers that in O(1).",
    )

    seen = []
    yield Frame(
        Array(nums, label="nums"),
        KV(seen, label="seen", empty="{ } empty"),
        Text("seen := make(Set[int], len(nums))"),
        say="So I keep a set of values seen so far. For each x, I check the set first, then add x.",
    )

    ans = False
    for i, x in enumerate(nums):
        if x in seen:
            ans = True
            first = nums.index(x)
            pair = (first, i)
            yield Frame(
                Array(nums, label="nums", hl={first: "match", i: "match"}, pointers={"x": i}),
                KV(seen, label="seen", hl={x: "match"}),
                Text(f"seen.Has({x}) = true", (f"return true (index {first} and {i})", "match")),
                say=f"{x} is already in the set, so it's a repeat. I return true at once; "
                    "nothing after this needs checking.",
            )
            break
        seen.append(x)
        yield Frame(
            Array(nums, label="nums", hl={k: "visited" for k in range(i)} | {i: "active"}, pointers={"x": i}),
            KV(seen, label="seen", hl={x: "active"}),
            Text(f"seen.Has({x}) = false → seen.Add({x})"),
            say=f"{x} isn't in the set, so it's the first time I've seen it. Add it and move on."
            if i == 0 else f"{x} is new as well. Add it.",
        )

    assert ans is True
    yield Frame(
        Array(nums, label="nums", hl={pair[0]: "match", pair[1]: "match"}),
        Text(("answer: true", "match"),
             "ContainsDuplicateBrute: O(n²) time, O(1) space",
             "ContainsDuplicate:      O(n) time,  O(n) space"),
        say="The answer is true. One pass with O(1) set operations is O(n) time, "
            "paid for with O(n) memory for the set.",
    )
