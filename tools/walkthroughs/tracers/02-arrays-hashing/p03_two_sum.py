from walkthrough import KV, Array, Frame, Text

TITLE = "Two Sum"
PATTERN = "Hash map"


def frames():
    nums, target = [2, 7, 11, 15], 9

    yield Frame(
        Array(nums, label="nums"),
        Text(f"target = {target}"),
        speaker="Interviewer",
        say="Return the indices of two different elements that add up to 9, or nil if there's no pair.",
    )
    yield Frame(
        Array(nums, label="nums", pointers={"i": 0, "j": 1}, hl={0: "compare", 1: "compare"}),
        Text(f"target = {target}", ("brute force: every pair → O(n²)", "reject")),
        say="Trying every pair is O(n²). For each x the inner loop is just searching for target − x. "
            "A map from value to index makes that search O(1).",
    )

    index_of = {}
    yield Frame(
        Array(nums, label="nums"),
        KV(index_of, label="indexOf (value: index)", empty="{ } empty"),
        Text(f"target = {target}"),
        say="I walk the slice once. For each x I look up target − x among the values seen so far, "
            "and only then insert x.",
    )

    ans = None
    for j, x in enumerate(nums):
        need = target - x
        if need in index_of:
            i = index_of[need]
            ans = [i, j]
            yield Frame(
                Array(nums, label="nums", hl={i: "match", j: "match"}, pointers={"i": i, "j": j}),
                KV(index_of, label="indexOf (value: index)", hl={need: "match"}),
                Text(f"x = {x}, need {target} − {x} = {need}", (f"indexOf[{need}] = {i} → return [{i}, {j}]", "match")),
                say=f"x is {x}, so I need {need}. {need} is in the map at index {i}, "
                    f"so the answer is [{i}, {j}].",
            )
            break
        index_of[x] = j
        yield Frame(
            Array(nums, label="nums", hl={j: "active"}, pointers={"j": j}),
            KV(index_of, label="indexOf (value: index)", hl={x: "active"}),
            Text(f"x = {x}, need {target} − {x} = {need}", f"{need} not in map → indexOf[{x}] = {j}"),
            say=f"x is {x}, so I need {need}. It isn't in the map yet. I store {x}: {j} after the lookup, "
                "so an element can never pair with itself.",
        )

    assert ans == [0, 1]
    yield Frame(
        Array(nums, label="nums", hl={0: "match", 1: "match"}),
        Text(("answer: [0, 1]", "match"), "TwoSumBrute: O(n²) time, O(1) space",
             "TwoSum:      O(n) time,  O(n) space"),
        say="The answer is [0, 1]. One pass with O(1) map lookups gives O(n) time, "
            "and the map costs O(n) space.",
    )
