from walkthrough import Array, Frame, Text

TITLE = "Jump Game II"
PATTERN = "Greedy · BFS levels"


def frames():
    nums = [2, 3, 1, 1, 4]
    expected = 2

    yield Frame(
        Array(nums, label="nums (max jump from each index)"),
        speaker="Interviewer",
        say="Same board, but the last index is always reachable. Return the minimum number of jumps to get there.",
    )
    yield Frame(
        Array(nums, label="nums (max jump from each index)"),
        Text("DP: dp[j] = min(dp[j], dp[i] + 1) for every jump i → j → O(n²)"),
        say="A DP keeps the fewest jumps to each index and relaxes every jump. That's O(n²). "
            "I think of it as BFS instead.",
    )

    jumps, end, farthest = 0, 0, 0
    level_start = 0
    yield Frame(
        Array(nums, label="nums", hl={0: "window"}, window=(0, 0), pointers={"i": 0, "end": 0}),
        Text("jumps = 0   end = 0   farthest = 0"),
        say="Indices reachable in k jumps form one window. Window 0 is just index 0. "
            "While I scan a window, I track the farthest index one more jump can reach.",
    )
    for i in range(len(nums) - 1):
        old_far = farthest
        farthest = max(farthest, i + nums[i])
        hl = {k: "window" for k in range(level_start, end + 1)}
        hl[i] = "active"
        if i == end:
            jumps += 1
            new_start, end = end + 1, farthest
            yield Frame(
                Array(nums, label="nums", hl=hl, window=(level_start, i), pointers={"i": i, "end": i}),
                Text(f"farthest = max({old_far}, {i} + {nums[i]}) = {farthest}",
                     (f"i == end → jump: jumps = {jumps}, end = {end}", "match")),
                say=f"I've scanned the whole window. One more jump reaches up to index {farthest}, "
                    f"so jumps becomes {jumps} and the next window ends at {end}.",
            )
            level_start = new_start
        else:
            yield Frame(
                Array(nums, label="nums", hl=hl, window=(level_start, end), pointers={"i": i, "end": end}),
                Text(f"farthest = max({old_far}, {i} + {nums[i]}) = {farthest}", f"i < end = {end}: keep scanning"),
                say=f"Index {i} is still inside the window. From here I can reach {i + nums[i]}, "
                    f"so farthest is {farthest}. I don't jump yet.",
            )

    assert jumps == expected
    last = len(nums) - 1
    yield Frame(
        Array(nums, label="nums", hl={0: "match", 1: "match", last: "match"}, window=(level_start, end)),
        Text((f"answer: {jumps} jumps  (0 → 1 → 4)", "match"), "loop stops before the last index",
             "O(n) time, O(1) space"),
        say=f"The last index is in window {jumps}, so the answer is {jumps}. "
            "I never jump from the last index. One pass: O(n) time, O(1) space.",
    )
