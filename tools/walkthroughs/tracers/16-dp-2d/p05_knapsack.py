from walkthrough import Frame, Grid, Text

TITLE = "0/1 Knapsack"
PATTERN = "2-D DP · take or skip"


def frames():
    weights, values, cap = [1, 3, 4, 5], [1, 4, 5, 7], 7
    n = len(weights)
    dp = [[None] * (cap + 1) for _ in range(n + 1)]
    rl = ["none"] + [f"w{w} v{v}" for w, v in zip(weights, values)]

    def g(hl=None):
        return Grid(dp, label="dp[i][c] = best value from the first i items with capacity c",
                    hl=hl, row_labels=rl, col_labels=[f"c={c}" for c in range(cap + 1)], cell=52)

    yield Frame(
        g(),
        Text(f"weights = {weights}, values = {values}, capacity = {cap}", "each item: take it whole or leave it"),
        speaker="Interviewer",
        say="Each item has a weight and a value. The bag holds at most 7. Each item is taken whole, at most once. "
            "What's the best total value?",
    )
    yield Frame(
        g(),
        Text("try every subset of items", ("2^n subsets", "reject")),
        say="Trying every subset is 2^n. But after deciding the first few items, all that matters is the "
            "capacity left. That's a small state I can tabulate.",
    )

    dp[0] = [0] * (cap + 1)
    yield Frame(
        g({(0, c): "match" for c in range(cap + 1)}),
        Text("dp[0][c] = 0  (no items, no value)"),
        say="State: dp[i][c] is the best value using the first i items with capacity c. "
            "With no items, every capacity is worth 0.",
    )

    key = {1: 1, 2: 4, 3: 7, 4: 7}
    says = {
        1: "Item 1 weighs 1. For each capacity I either skip it, copying the cell above, or take it: "
           "its value plus the row above at c minus 1.",
        2: "Item 2 weighs 3, worth 4. At c = 4, taking it leaves capacity 1, worth 1 above. 4 + 1 = 5 beats skipping.",
        3: "Item 3 weighs 4, worth 5. At c = 7, taking it leaves 3, worth 4 above: 5 + 4 = 9 beats 5.",
        4: "Item 4 weighs 5, worth 7. At c = 7, taking it leaves 2, worth only 1: 8. Skipping keeps 9.",
    }
    for i in range(1, n + 1):
        w, v = weights[i - 1], values[i - 1]
        dp[i] = [0] * (cap + 1)
        for c in range(cap + 1):
            dp[i][c] = dp[i - 1][c]
            if w <= c:
                dp[i][c] = max(dp[i][c], dp[i - 1][c - w] + v)
        c = key[i]
        skip = dp[i - 1][c]
        take = dp[i - 1][c - w] + v
        hl = {(i, cc): "active" for cc in range(cap + 1)}
        hl[(i, c)] = "match"
        hl[(i - 1, c)] = "compare"
        hl[(i - 1, c - w)] = "compare"
        yield Frame(
            g(hl),
            Text(f"row {i}: dp[{i}][c] = max(dp[{i-1}][c], dp[{i-1}][c-{w}] + {v})",
                 (f"c={c}: max(skip {skip}, take {dp[i-1][c-w]} + {v} = {take}) = {dp[i][c]}", "match")),
            say=says[i],
        )

    ans = dp[n][cap]
    assert ans == 9
    yield Frame(
        g({(n, cap): "match"}),
        Text(f"answer: dp[{n}][{cap}] = {ans}  (items weighing 3 and 4)",
             ("O(n·W) time · O(W) space: one row, capacity high → low", "match")),
        say="The answer is 9. Each row reads only the row above, so the Go code keeps one array and loops capacity "
            "downward. Upward would reuse an item. O(n·W) time, O(W) space.",
    )
