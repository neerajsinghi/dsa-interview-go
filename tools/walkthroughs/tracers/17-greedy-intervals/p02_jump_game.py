from walkthrough import Array, Frame, Text

TITLE = "Jump Game"
PATTERN = "Greedy"


def frames():
    nums = [2, 3, 1, 1, 4]
    expected = True
    last = len(nums) - 1

    yield Frame(
        Array(nums, label="nums (max jump from each index)"),
        speaker="Interviewer",
        say="You start at index 0. nums[i] is the longest jump you can make from i. Can you reach the last index?",
    )
    yield Frame(
        Array(nums, label="nums (max jump from each index)"),
        Text("DP: from each reachable i, mark i+1 .. i+nums[i] → O(n²)"),
        say="I could mark every index each reachable index can jump to. That's O(n²) in the worst case, "
            "because one index can mark up to n others.",
    )

    reach = 0
    yield Frame(
        Array(nums, label="nums", hl={0: "window"}, window=(0, 0), pointers={"reach": 0}),
        Text("reach = 0"),
        say="Instead I keep one number: reach, the farthest index I can get to so far. "
            "Every index up to reach is reachable, because I can always jump shorter.",
    )
    answer = None
    for i, n in enumerate(nums):
        if i > reach:
            answer = False
            break
        old = reach
        reach = max(reach, i + n)
        done = reach >= last
        hl = {k: "window" for k in range(0, min(reach, last) + 1)}
        hl[i] = "active"
        if done:
            hl[last] = "match"
        yield Frame(
            Array(nums, label="nums", hl=hl, window=(0, min(reach, last)),
                  pointers={"i": i, "reach": min(reach, last)} if reach != i else {"i": i}),
            Text(f"i = {i} ≤ reach = {old}, so i is reachable",
                 f"reach = max({old}, {i} + {n}) = {reach}",
                 (f"reach ≥ last index {last} → true", "match") if done else ""),
            say=(f"From index {i} I can get to {i + n}. That covers the last index, so I can stop and return true."
                 if done else
                 f"Index {i} is within reach. From here I can get to {i + n}, so reach becomes {reach}."),
        )
        if done:
            answer = True
            break
    if answer is None:
        answer = True

    assert answer == expected
    yield Frame(
        Array(nums, label="nums", hl={k: "match" for k in range(len(nums))}),
        Text(("answer: true", "match"), "fails only if the scan reaches i > reach", "O(n) time, O(1) space"),
        say="The answer is true. If the scan ever reached an index past reach, I'd return false. "
            "One pass, one variable: O(n) time, O(1) space.",
    )
