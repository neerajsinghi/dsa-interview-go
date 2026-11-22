from walkthrough import Array, Frame, Text

TITLE = "Climbing Stairs"
PATTERN = "1-D DP"


def frames():
    n = 5
    dp = [None] * (n + 1)

    def arr(hl=None):
        return Array(dp, label="ways[i] = move sequences that land on step i", hl=hl)

    yield Frame(
        arr(),
        Text(f"n = {n}", "each move climbs 1 or 2 steps"),
        speaker="Interviewer",
        say="A staircase has 5 steps. Each move climbs 1 or 2 steps. How many different sequences of moves reach the top?",
    )
    yield Frame(
        arr(),
        Text("ways(n) = ways(n-1) + ways(n-2), by plain recursion", ("the same subproblems repeat → O(2^n)", "reject")),
        say="The brute force recurses on the last move. But ways(3) gets computed again inside ways(5) and ways(4), "
            "so the call tree doubles each level: O(2^n).",
    )

    dp[0], dp[1] = 1, 1
    yield Frame(
        arr({0: "match", 1: "match"}),
        Text("ways[0] = 1  (the empty sequence)", "ways[1] = 1  (one single step)"),
        say="I'll save each answer instead. The base cases: one way to stand on step 0, doing nothing, "
            "and one way to reach step 1.",
    )
    for i in range(2, n + 1):
        dp[i] = dp[i - 1] + dp[i - 2]
        if i == 2:
            say = ("Every way to reach step 2 ends with a 1-step from step 1 or a 2-step from step 0. "
                   "Those groups don't overlap, so I add them.")
        elif i == 3:
            say = "Same rule for step 3: the ways ending with a 1-step plus the ways ending with a 2-step."
        else:
            say = f"Step {i} only reads the two cells to its left, which are already filled. That's why I go left to right."
        yield Frame(
            arr({i: "active", i - 1: "compare", i - 2: "compare"}),
            Text(f"ways[{i}] = ways[{i-1}] + ways[{i-2}] = {dp[i-1]} + {dp[i-2]} = {dp[i]}"),
            say=say,
        )

    assert dp[n] == 8
    yield Frame(
        arr({n: "match"}),
        Text(f"answer: ways[{n}] = {dp[n]}", ("O(n) time · O(1) space (keep only the last two)", "match")),
        say=f"The answer is ways[5] = {dp[n]}. It's Fibonacci. Each cell needs only the previous two, "
            "so the Go code keeps two variables: O(n) time, O(1) space.",
    )
