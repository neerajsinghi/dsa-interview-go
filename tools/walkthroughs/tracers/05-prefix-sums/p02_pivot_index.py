from walkthrough import Array, Frame, Text

TITLE = "Find Pivot Index"
PATTERN = "Prefix sums"


def frames():
    nums = [1, 7, 3, 6, 5, 6]
    want = 3

    yield Frame(
        Array(nums, label="nums"),
        speaker="Interviewer",
        say="Find the leftmost index where the sum strictly to its left equals the sum strictly "
            "to its right. An empty side counts as 0. Return -1 if there's none.",
    )
    yield Frame(
        Array(nums, label="nums", hl={3: "active", 0: "compare", 1: "compare", 2: "compare",
                                      4: "compare", 5: "compare"}),
        Text("brute force: for each i, add up both sides", "→ O(n) per index, O(n²) total"),
        say="I could add up both sides for every index. That's O(n) work per index and O(n²) "
            "overall, and most of those sums repeat.",
    )

    total = sum(nums)
    yield Frame(
        Array(nums, label="nums"),
        Text(f"total = {' + '.join(map(str, nums))} = {total}", "right = total − left − nums[i]"),
        say=f"If I know the total, {total}, I only need a running left sum. The right side is "
            "whatever is left over: total minus left minus the current number.",
    )

    left = 0
    answer = -1
    for i, x in enumerate(nums):
        right = total - left - x
        hl = {j: "visited" for j in range(i)}
        hl.update({j: "compare" for j in range(i + 1, len(nums))})
        if left == right:
            hl[i] = "match"
            yield Frame(
                Array(nums, label="nums", hl=hl, pointers={"i": i}),
                Text(f"left = {left}   right = {total} − {left} − {x} = {right}",
                     (f"{left} == {right} → pivot at index {i}", "match")),
                say=f"Left is {left} and right is {total} − {left} − {x} = {right}. They're equal, "
                    f"so {i} is the pivot. I return straight away, so it's the leftmost one.",
            )
            answer = i
            break
        hl[i] = "active"
        yield Frame(
            Array(nums, label="nums", hl=hl, pointers={"i": i}),
            Text(f"left = {left}   right = {total} − {left} − {x} = {right}",
                 (f"{left} ≠ {right} → left += {x}", "reject")),
            say=f"At index {i}, left is {left} and right is {right}. Not equal, so I add {x} to the "
                "left sum and move on." if i else
                f"At index 0 the left side is empty, so left is 0. Right is {right}. Not equal, "
                f"so I add {x} to left and move on.",
        )
        left += x

    assert answer == want
    yield Frame(
        Array(nums, label="nums", hl={answer: "match"}, pointers={"pivot": answer}),
        Text(f"answer: {answer}   (1+7+3 = 11 = 5+6)", ("time O(n) · space O(1)", "match")),
        say="The answer is 3. One pass for the total, one pass for the running left sum: O(n) time "
            "and O(1) extra space.",
    )
