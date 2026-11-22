from walkthrough import Array, Frame, KV, Text

TITLE = "Continuous Subarray Sum"
PATTERN = "Prefix mod k"


def mod(x, k):
    return ((x % k) + k) % k


def frames():
    nums, k = [23, 2, 4, 6, 7], 6
    want = True

    yield Frame(
        Array(nums, label="nums"),
        Text(f"k = {k}"),
        speaker="Interviewer",
        say="Is there a contiguous subarray of length at least 2 whose sum is a multiple of k? "
            "Zero counts as a multiple, and numbers can be negative.",
    )
    yield Frame(
        Array(nums, label="nums"),
        Text(f"k = {k}", "brute force: every subarray of length ≥ 2 → O(n²)"),
        say="Checking every subarray of length two or more is O(n²). I want one pass.",
    )
    first = {0: -1}
    yield Frame(
        Array(nums, label="nums"),
        KV(first, label="first  (remainder: first index)", hl={0: "special"}),
        Text("P[i] ≡ P[j] (mod k)  ⇒  sum between them is a multiple of k"),
        say="If two prefix sums leave the same remainder mod k, the part between them is a multiple "
            "of k. I store the first index of each remainder, seeded with 0 at −1.",
    )

    prefix = 0
    result = False
    for i, x in enumerate(nums):
        old = prefix
        prefix = mod(prefix + x, k)
        calc = f"prefix = ({old} + {x}) mod {k} = {prefix}"
        if prefix in first:
            j = first[prefix]
            if i - j >= 2:
                yield Frame(
                    Array(nums, label="nums", hl={t: "match" for t in range(j + 1, i + 1)},
                          window=(j + 1, i), pointers={"i": i}),
                    KV(first, label="first", hl={prefix: "match"}),
                    Text(calc, (f"remainder {prefix} first at {j}: span {i} − {j} = {i - j} ≥ 2 → true", "match")),
                    say=f"Remainder {prefix} was first seen at index {j}. The span is {i - j}, at least 2, "
                        f"so nums[{j + 1}..{i}] sums to a multiple of {k}. Return true.",
                )
                result = True
                break
            yield Frame(
                Array(nums, label="nums", hl={i: "reject"}, pointers={"i": i}),
                KV(first, label="first", hl={prefix: "reject"}),
                Text(calc, (f"span {i - j} < 2 → keep going", "reject")),
                say=f"Remainder {prefix} repeats, but only one element apart. That's a length-1 subarray, "
                    "so it doesn't count. I keep the older index.",
            )
        else:
            first[prefix] = i
            yield Frame(
                Array(nums, label="nums", hl={i: "active", **{t: "visited" for t in range(i)}}, pointers={"i": i}),
                KV(first, label="first", hl={prefix: "active"}),
                Text(calc, f"remainder {prefix} is new → first[{prefix}] = {i}"),
                say=f"I add {x} and reduce mod {k}, giving {prefix}. It's new, so I record index {i}. "
                    "Reducing as I go keeps the numbers small and never overflows.",
            )

    assert result == want
    yield Frame(
        Array(nums, label="nums", hl={1: "match", 2: "match"}, window=(1, 2)),
        Text("answer: true   (2 + 4 = 6)", ("time O(n) · space O(min(n, k))", "match")),
        say="The answer is true: 2 + 4 = 6. One pass, and the map holds at most min(n, k) remainders, "
            "so O(n) time and O(min(n, k)) space.",
    )
