from walkthrough import Frame, Grid, Text

TITLE = "Longest Palindromic Subsequence"
PATTERN = "2-D DP · interval"


def frames():
    s = "bbbab"
    n = len(s)
    dp = [[None] * n for _ in range(n)]
    lower = {(i, j): "empty" for i in range(n) for j in range(n) if i > j}
    labels = [f"{k}:{c}" for k, c in enumerate(s)]

    def g(hl=None):
        h = dict(lower)
        h.update(hl or {})
        return Grid(dp, label="dp[i][j] = longest palindromic subsequence of s[i..j]   (row i, column j)",
                    hl=h, row_labels=labels, col_labels=labels, cell=44)

    yield Frame(
        g(),
        Text(f"s = \"{s}\""),
        speaker="Interviewer",
        say="Return the length of the longest subsequence of \"bbbab\" that reads the same both ways.",
    )
    yield Frame(
        g(),
        Text("try every subsequence and check it", ("2^n subsequences", "reject")),
        say="Checking all 2^n subsequences is too slow. A palindrome is decided by its two ends, "
            "so I'll solve every interval s[i..j] instead.",
    )

    for i in range(n):
        dp[i][i] = 1
    yield Frame(
        g({(i, i): "match" for i in range(n)}),
        Text("dp[i][i] = 1  (one letter is a palindrome)", "dp[i][j] needs row i+1 → fill i from the bottom up"),
        say="State: dp[i][j] is the answer for s[i..j]. A single letter is a palindrome of length 1. "
            "Each cell reads the row below, so I fill rows bottom-up.",
    )

    def cell(i, j):
        if s[i] == s[j]:
            inner = dp[i + 1][j - 1] if i + 1 <= j - 1 else 0
            dp[i][j] = inner + 2
            h = {(i, j): "active"}
            if i + 1 <= j - 1:
                h[(i + 1, j - 1)] = "compare"
            return f"s[{i}] == s[{j}] = '{s[i]}' → inner + 2 = {inner} + 2 = {dp[i][j]}", h
        dp[i][j] = max(dp[i + 1][j], dp[i][j - 1])
        return (f"s[{i}]='{s[i]}' != s[{j}]='{s[j]}' → max(below {dp[i+1][j]}, left {dp[i][j-1]}) = {dp[i][j]}",
                {(i, j): "active", (i + 1, j): "compare", (i, j - 1): "compare"})

    says = {
        (3, 4): "s[3..4] is \"ab\". The ends differ, so one of them is left out. "
                "I take the better of dropping the left end (below) or the right end (left).",
        (2, 3): "\"ba\": the ends differ again, so 1.",
        (2, 4): "\"bab\": the ends match. Wrap them around the best of the inside, \"a\": 1 + 2 = 3.",
    }
    for i, j in [(3, 4), (2, 3), (2, 4)]:
        line, h = cell(i, j)
        yield Frame(g(h), Text(line), say=says[(i, j)])

    for i, say in [(1, "Same rule for row 1, left to right. Its matching ends wrap the inside answer and add 2."),
                   (0, None)]:
        lines, h = [], {}
        for j in range(i + 1, n):
            line, _ = cell(i, j)
            h[(i, j)] = "active"
            lines.append(line)
        if say:
            yield Frame(g(h), Text(*lines, size=13), say=say)
        else:
            last = lines[-1]
            yield Frame(
                g({**h, (0, n - 1): "match", (1, n - 2): "compare"}),
                Text(*lines[:-1], (last, "match"), size=13),
                say="Row 0 the same way. The corner is the whole string: both ends are 'b', "
                    "so it wraps \"bba\", worth 2, and adds 2.",
            )

    ans = dp[0][n - 1]
    assert ans == 4
    yield Frame(
        g({(0, n - 1): "match"}),
        Text(f"answer: dp[0][{n-1}] = {ans}  (\"bbbb\")", ("O(n²) time · O(n) space with one row", "match")),
        say="The answer is 4, \"bbbb\". The Go row version keeps one row plus the saved diagonal: "
            "O(n²) time, O(n) space.",
    )
