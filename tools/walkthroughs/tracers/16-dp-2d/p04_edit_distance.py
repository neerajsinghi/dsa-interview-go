from walkthrough import Frame, Grid, Text

TITLE = "Edit Distance"
PATTERN = "2-D DP · two strings"


def frames():
    a, b = "horse", "ros"
    m, n = len(a), len(b)
    dp = [[None] * (n + 1) for _ in range(m + 1)]
    rl = ["\"\""] + list(a)
    cl = ["\"\""] + list(b)

    def g(hl=None):
        return Grid(dp, label="dp[i][j] = edits to turn a[:i] into b[:j]   (rows: horse, columns: ros)",
                    hl=hl, row_labels=rl, col_labels=cl, cell=40)

    yield Frame(
        g(),
        Text(f"a = \"{a}\", b = \"{b}\"", "insert, delete or replace one character: cost 1 each"),
        speaker="Interviewer",
        say="Turn \"horse\" into \"ros\" using inserts, deletes and replacements. Each costs 1. "
            "What's the fewest operations?",
    )
    yield Frame(
        g(),
        Text("on a mismatch try all three operations recursively", ("three branches each time → O(3^(m+n))", "reject")),
        say="The brute force tries all three operations at every mismatch: up to 3^(m+n) calls, "
            "mostly on the same pairs of prefixes.",
    )

    for i in range(m + 1):
        dp[i][0] = i
    for j in range(n + 1):
        dp[0][j] = j
    base = {(0, j): "match" for j in range(n + 1)}
    base.update({(i, 0): "match" for i in range(m + 1)})
    yield Frame(
        g(base),
        Text("dp[i][0] = i  (delete all i letters)", "dp[0][j] = j  (insert all j letters)"),
        say="State: dp[i][j] is the cost to turn the first i letters of a into the first j of b. "
            "Turning a prefix into empty costs its length, and so does building from empty.",
    )

    def cell(i, j):
        if a[i - 1] == b[j - 1]:
            dp[i][j] = dp[i - 1][j - 1]
            return (f"'{a[i-1]}' == '{b[j-1]}' → free, copy diag = {dp[i][j]}",
                    {(i, j): "active", (i - 1, j - 1): "compare"})
        up, left, diag = dp[i - 1][j], dp[i][j - 1], dp[i - 1][j - 1]
        dp[i][j] = 1 + min(up, left, diag)
        return (f"'{a[i-1]}' != '{b[j-1]}' → 1 + min(delete {up}, insert {left}, replace {diag}) = {dp[i][j]}",
                {(i, j): "active", (i - 1, j): "compare", (i, j - 1): "compare", (i - 1, j - 1): "compare"})

    says = {
        (1, 1): "'h' vs 'r' differ. Up means delete 'h', left means insert 'r', diagonal means replace. "
                "Replace is cheapest: 1 + 0 = 1.",
        (1, 2): "'h' vs 'o' differ. The best neighbour is 1, so 2.",
        (1, 3): "'h' vs 's' differ: 1 + 2 = 3.",
        (2, 1): "Row 'o': 'o' vs 'r' differ. Deleting 'o' after turning 'h' into 'r' costs 1 + 1 = 2.",
        (2, 2): "'o' matches 'o'. A match is free: copy the diagonal, the cost of turning 'h' into 'r', which is 1.",
        (2, 3): "'o' vs 's' differ. Inserting 's' after the 1 on the left gives 2.",
    }
    for i in (1, 2):
        for j in range(1, n + 1):
            line, h = cell(i, j)
            yield Frame(g(h), Text(line), say=says[(i, j)])

    row_says = {
        3: "Same rule for row 'r', left to right. 'r' matches 'r', so that cell copies the diagonal for free: 2.",
        4: "Row 's': 's' matches 's', so that cell copies the diagonal for free.",
        5: "Row 'e' matches nothing. The corner is 1 + min(2, 4, 3) = 3: delete the 'e'.",
    }
    for i in range(3, m + 1):
        h = {}
        lines = []
        for j in range(1, n + 1):
            line, _ = cell(i, j)
            h[(i, j)] = "active"
            lines.append(f"j={j}: " + line)
        yield Frame(g(h), Text(*lines, size=13), say=row_says[i])

    ans = dp[m][n]
    assert ans == 3
    yield Frame(
        g({(m, n): "match"}),
        Text(f"answer: dp[{m}][{n}] = {ans}  (horse → rorse → rose → ros)",
             ("O(m·n) time · O(n) space: one row plus a saved diagonal", "match")),
        say="The answer is 3. The Go rolling version keeps one row, and saves the old 'above' value "
            "before overwriting it, because the next cell needs it as its diagonal.",
    )
