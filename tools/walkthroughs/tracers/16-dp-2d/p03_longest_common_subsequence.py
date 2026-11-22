from walkthrough import Frame, Grid, Text

TITLE = "Longest Common Subsequence"
PATTERN = "2-D DP · two strings"


def frames():
    a, b = "abcde", "ace"
    m, n = len(a), len(b)
    dp = [[None] * (n + 1) for _ in range(m + 1)]
    rl = ["\"\""] + list(a)
    cl = ["\"\""] + list(b)

    def g(hl=None):
        return Grid(dp, label="dp[i][j] = LCS of a[:i] and b[:j]   (rows: a = abcde, columns: b = ace)",
                    hl=hl, row_labels=rl, col_labels=cl, cell=40)

    yield Frame(
        g(),
        Text(f"a = \"{a}\", b = \"{b}\""),
        speaker="Interviewer",
        say="Given \"abcde\" and \"ace\", return the length of the longest subsequence they share. "
            "A subsequence keeps order but may skip characters.",
    )
    yield Frame(
        g(),
        Text("if first chars match, take both; else drop one from either string", ("two branches each time → O(2^(m+n))", "reject")),
        say="The brute force branches twice whenever the first characters differ. That's exponential, "
            "and it keeps re-solving the same pair of suffixes.",
    )

    for i in range(m + 1):
        dp[i][0] = 0
    for j in range(n + 1):
        dp[0][j] = 0
    base = {(0, j): "match" for j in range(n + 1)}
    base.update({(i, 0): "match" for i in range(m + 1)})
    yield Frame(
        g(base),
        Text("dp[0][*] = dp[*][0] = 0  (an empty prefix shares nothing)"),
        say="State: dp[i][j] is the LCS of the first i letters of a and the first j of b. "
            "If either prefix is empty, the answer is 0.",
    )

    def cell(i, j):
        if a[i - 1] == b[j - 1]:
            dp[i][j] = dp[i - 1][j - 1] + 1
            return (f"'{a[i-1]}' == '{b[j-1]}' → diag + 1 = {dp[i-1][j-1]} + 1 = {dp[i][j]}",
                    {(i, j): "active", (i - 1, j - 1): "compare"})
        dp[i][j] = max(dp[i - 1][j], dp[i][j - 1])
        return (f"'{a[i-1]}' != '{b[j-1]}' → max(up {dp[i-1][j]}, left {dp[i][j-1]}) = {dp[i][j]}",
                {(i, j): "active", (i - 1, j): "compare", (i, j - 1): "compare"})

    says = {
        (1, 1): "'a' matches 'a'. Both letters can end the common subsequence, so I take the diagonal and add 1.",
        (1, 2): "'a' vs 'c' differ. One of them isn't used, so I take the better of dropping either: up or left.",
        (1, 3): "Same for 'e': the 1 from the left carries over.",
        (2, 1): "Row 'b': 'b' vs 'a' differ, so max(up, left) = 1.",
        (2, 2): "'b' vs 'c' differ; the 1 carries over.",
        (2, 3): "'b' vs 'e' differ; still 1.",
    }
    for i in (1, 2):
        for j in range(1, n + 1):
            line, h = cell(i, j)
            yield Frame(g(h), Text(line), say=says[(i, j)])

    row_says = {
        3: "Same rule for row 'c'. 'c' matches 'c', so that cell is the diagonal 1 plus 1 = 2.",
        4: "Row 'd' matches nothing, so every cell copies the best of up and left.",
        5: "Row 'e': 'e' matches 'e', so the corner is the diagonal 2 plus 1 = 3.",
    }
    for i in range(3, m + 1):
        h = {}
        lines = []
        for j in range(1, n + 1):
            line, _ = cell(i, j)
            h[(i, j)] = "active"
            if a[i - 1] == b[j - 1]:
                lines.append(f"j={j}: " + line)
        if not lines:
            lines.append("no match in this row → every cell = max(up, left)")
        yield Frame(g(h), Text(*lines), say=row_says[i])

    ans = dp[m][n]
    assert ans == 3
    # backtrack like LCSString
    i, j, path, out = m, n, {}, []
    while i > 0 and j > 0:
        if a[i - 1] == b[j - 1]:
            path[(i, j)] = "match"
            out.append(a[i - 1])
            i, j = i - 1, j - 1
        elif dp[i - 1][j] >= dp[i][j - 1]:
            path[(i, j)] = "visited"
            i -= 1
        else:
            path[(i, j)] = "visited"
            j -= 1
    word = "".join(reversed(out))
    assert word == "ace"
    yield Frame(
        g(path),
        Text(f"answer: dp[{m}][{n}] = {ans}  · walk back from the corner → \"{word}\"",
             ("O(m·n) time · O(min(m,n)) space with two rows", "match")),
        say="The answer is 3. Walking back from the corner, every diagonal step on a match gives a letter: \"ace\". "
            "Each row reads only the one above, so two rows are enough.",
    )
