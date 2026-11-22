from walkthrough import Array, Frame, Grid, Text

TITLE = "Burst Balloons"
PATTERN = "2-D DP · interval"


def frames():
    nums = [3, 1, 5, 8]
    vals = [1] + nums + [1]
    N = len(vals)
    dp = [[None] * N for _ in range(N)]
    blank = {(l, r): "empty" for l in range(N) for r in range(N) if r <= l}
    labels = [f"{k}:{v}" for k, v in enumerate(vals)]

    def widgets(hl=None, hl_v=None):
        h = dict(blank)
        h.update(hl or {})
        return (
            Array(vals, label="vals = [1] + nums + [1]", hl=hl_v, cell=44),
            Grid(dp, label="dp[l][r] = best coins from bursting everything strictly between l and r",
                 hl=h, row_labels=labels, col_labels=labels, cell=46),
        )

    yield Frame(
        *widgets(),
        Text(f"nums = {nums} · bursting i earns left · nums[i] · right (current neighbours)"),
        speaker="Interviewer",
        say="Bursting a balloon earns its value times its current left and right neighbours. "
            "Missing neighbours count as 1. Burst all of [3, 1, 5, 8] for the most coins.",
    )
    yield Frame(
        *widgets(),
        Text("try every burst order", ("n! orders", "reject")),
        say="Trying every order is n!. Choosing the first balloon doesn't split the problem, because its "
            "neighbours change. But the last balloon in a range always has fixed neighbours.",
    )

    for l in range(N - 1):
        dp[l][l + 1] = 0
    yield Frame(
        *widgets({(l, l + 1): "match" for l in range(N - 1)}, {0: "special", N - 1: "special"}),
        Text("pad with 1 at both ends", "dp[l][l+1] = 0  (nothing strictly between)"),
        say="I pad both ends with 1. dp[l][r] is the best for the balloons strictly between l and r. "
            "Neighbouring l and r have nothing between them: 0.",
    )

    def solve(l, r):
        best, cands = 0, []
        for k in range(l + 1, r):
            got = dp[l][k] + vals[l] * vals[k] * vals[r] + dp[k][r]
            cands.append((k, got))
            best = max(best, got)
        dp[l][r] = best
        return cands

    def fmt(l, r, cands, show=True):
        parts = [f"k={k}: {dp[l][k]} + {vals[l]}·{vals[k]}·{vals[r]} + {dp[k][r]} = {got}" for k, got in cands]
        return parts

    # width 2: all at once
    lines, h = [], {}
    for l in range(N - 2):
        r = l + 2
        cands = solve(l, r)
        h[(l, r)] = "active"
        k = l + 1
        lines.append(f"dp[{l}][{r}] = {vals[l]}·{vals[k]}·{vals[r]} = {dp[l][r]}")
    yield Frame(*widgets(h), Text(*lines, size=13),
                say="Width 2: one balloon between l and r, so it is the last one. Its neighbours are l and r: "
                    "the product of three values.")

    # width 3: first interval in detail, the rest at once
    l, r = 0, 3
    cands = solve(l, r)
    h = {(l, r): "active"}
    for k, _ in cands:
        h[(l, k)] = "compare"
        h[(k, r)] = "compare"
    yield Frame(*widgets(h, {l: "window", r: "window"}),
                Text(*fmt(l, r, cands), (f"dp[0][3] = {dp[0][3]}", "active"), size=13),
                say="Width 3, (0, 3): try each balloon k as the last one. Its neighbours are then 0 and 3, "
                    "and the two sides were already solved. The best is k = 1.")
    h = {}
    lines = []
    for l in (1, 2):
        cands = solve(l, l + 3)
        h[(l, l + 3)] = "active"
        lines.append(f"dp[{l}][{l+3}] = max({', '.join(str(g) for _, g in cands)}) = {dp[l][l+3]}")
    yield Frame(*widgets(h), Text(*lines, size=13),
                say="Same rule for the rest of width 3. Shorter intervals are always ready, because I fill by width.")

    h, lines = {}, []
    for l in (0, 1):
        cands = solve(l, l + 4)
        h[(l, l + 4)] = "active"
        lines.append(f"dp[{l}][{l+4}] = max({', '.join(str(g) for _, g in cands)}) = {dp[l][l+4]}")
    yield Frame(*widgets(h), Text(*lines, size=13), say="Width 4, same rule: three choices of last balloon each.")

    l, r = 0, N - 1
    cands = solve(l, r)
    bestk = max(cands, key=lambda t: t[1])[0]
    h = {(l, r): "active", (l, bestk): "compare", (bestk, r): "compare"}
    yield Frame(*widgets(h, {bestk: "match"}),
                Text(*fmt(l, r, cands), (f"dp[0][5] = {dp[0][5]}  (k = {bestk}: balloon 8 last)", "active"), size=13),
                say="The whole row, (0, 5). Bursting 8 last is best: 1·8·1 = 8, plus the best for 3, 1, 5 on its left, 159.")

    ans = dp[0][N - 1]
    assert ans == 167
    yield Frame(
        *widgets({(0, N - 1): "match"}),
        Text(f"answer: dp[0][{N-1}] = {ans}", ("O(n³) time · O(n²) space", "match")),
        say="The answer is 167. There are O(n²) intervals and each tries O(n) last balloons: "
            "O(n³) time, O(n²) space.",
    )
