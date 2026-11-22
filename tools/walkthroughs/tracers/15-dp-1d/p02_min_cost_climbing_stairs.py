from walkthrough import Array, Frame, Text

TITLE = "Min Cost Climbing Stairs"
PATTERN = "1-D DP"


def frames():
    cost = [10, 15, 20]
    n = len(cost)
    dp = [None] * (n + 1)

    def widgets(hl_cost=None, hl_dp=None):
        return (
            Array(cost + [""], label="cost[i] (index 3 is the top, it's free)", hl=hl_cost),
            Array(dp, label="best[i] = min cost to stand on position i", hl=hl_dp),
        )

    yield Frame(
        *widgets(),
        Text("start on stair 0 or 1 · pay for a stair, then climb 1 or 2"),
        speaker="Interviewer",
        say="cost[i] is the price of stepping on stair i. You can start on stair 0 or 1. "
            "Return the cheapest way to reach the top, just past the last stair.",
    )
    yield Frame(
        *widgets(),
        Text("try every sequence of 1- and 2-steps", ("overlapping subproblems → O(2^n)", "reject")),
        say="Trying every sequence of moves is exponential, and it keeps re-solving the same question: "
            "what's the cheapest way to stand on stair i?",
    )

    dp[0], dp[1] = 0, 0
    yield Frame(
        *widgets(hl_dp={0: "match", 1: "match"}),
        Text("best[0] = 0, best[1] = 0  (free starting points)"),
        say="So that's my state: best[i], the cheapest way to stand on position i before paying for it. "
            "Positions 0 and 1 are free starts.",
    )
    for i in range(2, n + 1):
        a, b = dp[i - 1] + cost[i - 1], dp[i - 2] + cost[i - 2]
        dp[i] = min(a, b)
        win = i - 1 if a <= b else i - 2
        yield Frame(
            *widgets(hl_cost={i - 1: "compare", i - 2: "compare"},
                     hl_dp={i: "active", i - 1: "compare", i - 2: "compare"}),
            Text(f"from {i-1}: best[{i-1}] + cost[{i-1}] = {dp[i-1]} + {cost[i-1]} = {a}",
                 f"from {i-2}: best[{i-2}] + cost[{i-2}] = {dp[i-2]} + {cost[i-2]} = {b}",
                 (f"best[{i}] = min({a}, {b}) = {dp[i]}  (came from {win})", "active")),
            say=(f"To stand on position {i} I arrived from {i-1} or {i-2}, paying for the stair I left. "
                 f"I take the cheaper one: {dp[i]}." if i == 2 else
                 f"Position {i} is the top. Coming from stair 2 costs {a}; jumping two from stair 1 costs {b}. "
                 f"The jump wins."),
        )

    assert dp[n] == 15
    yield Frame(
        *widgets(hl_cost={1: "match"}, hl_dp={n: "match"}),
        Text(f"answer: best[{n}] = {dp[n]}  (start on stair 1, jump 2)", ("O(n) time · O(1) space", "match")),
        say=f"The answer is best[3] = {dp[n]}: start on stair 1 and jump to the top. "
            "Each cell reads only two earlier ones, so the Go code keeps two variables: O(n) time, O(1) space.",
    )
