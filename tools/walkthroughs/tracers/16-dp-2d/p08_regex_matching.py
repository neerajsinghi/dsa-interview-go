from walkthrough import Frame, Grid, Text

TITLE = "Regular Expression Matching"
PATTERN = "2-D DP · two strings"


def frames():
    s, p = "aab", "c*a*b"
    m, n = len(s), len(p)
    dp = [[None] * (n + 1) for _ in range(m + 1)]
    rl = [f"\"{s[i:]}\"" for i in range(m + 1)]
    cl = [f"{p[j:]}" if j < n else "\"\"" for j in range(n + 1)]

    def g(hl=None):
        return Grid(dp, label="dp[i][j] = s[i:] matches p[j:]   (rows: suffixes of s, cols: suffixes of p)",
                    hl=hl, row_labels=rl, col_labels=cl, cell=50)

    yield Frame(
        g(),
        Text(f"s = \"{s}\", p = \"{p}\"", "'.' = any one char · 'x*' = zero or more x · match the whole string"),
        speaker="Interviewer",
        say="Implement matching with '.' and '*'. The pattern must match the whole string. "
            "Does \"c*a*b\" match \"aab\"?",
    )
    yield Frame(
        g(),
        Text("x* → use it zero times, or consume one char and keep x*", ("plain recursion is exponential", "reject")),
        say="The recursion branches at every star: skip it, or eat one letter and stay. "
            "Inputs like a*a*a*a*c blow up, but there are only (m+1)(n+1) suffix pairs.",
    )

    for i in range(m + 1):
        dp[i][n] = i == m
    yield Frame(
        g({(i, n): "match" if dp[i][n] else "reject" for i in range(m + 1)}),
        Text("empty pattern matches only the empty string", "fill i from m down to 0, j from n-1 down to 0"),
        say="State: dp[i][j] says whether s from i matches p from j. The empty pattern matches only the empty "
            "string. Cells read i+1 and j+1 or j+2, so I fill from the bottom-right.",
    )

    def cell(i, j):
        first = i < m and (p[j] == s[i] or p[j] == ".")
        if j + 1 < n and p[j + 1] == "*":
            skip = dp[i][j + 2]
            eat = first and dp[i + 1][j]
            dp[i][j] = bool(skip or eat)
            src = {(i, j + 2): "compare"}
            if first:
                src[(i + 1, j)] = "compare"
            return (f"p[{j}:]=\"{p[j:]}\": skip x* → dp[{i}][{j+2}]={'T' if skip else 'F'}"
                    f", eat one → {'T' if eat else 'F'}  ⇒ {'T' if dp[i][j] else 'F'}"), src
        dp[i][j] = bool(first and dp[i + 1][j + 1])
        if p[j] == "*":
            return None, {}
        return (f"p[{j}]='{p[j]}' vs \"{s[i:]}\": first {'✓' if first else '✗'}"
                + (f" and dp[{i+1}][{j+1}]" if first else "") + f"  ⇒ {'T' if dp[i][j] else 'F'}",
                {(i + 1, j + 1): "compare"} if first else {})

    row_says = {
        3: "Row for the empty string. 'b' needs a letter, so false. 'a*' and 'c*' can only be skipped, "
           "and what's after them fails. Columns starting at '*' are never true.",
        2: "For \"b\": 'b' matches and the rest is empty, so true. Then 'a*' and 'c*' are skipped, "
           "and the true cell two to the right carries left.",
        1: "For \"ab\": 'b' fails on 'a'. But 'a*' can eat the 'a' and stay, landing on \"b\" vs \"a*b\", "
           "which is true.",
        0: "For \"aab\": 'a*' eats one 'a' again and reads the row below. Then 'c*' is skipped zero times, "
           "carrying the true across.",
    }
    for i in range(m, -1, -1):
        lines, h = [], {}
        for j in range(n - 1, -1, -1):
            line, src = cell(i, j)
            h[(i, j)] = "match" if dp[i][j] else "active"
            if line:
                lines.append(line)
                if i < m and j == (2 if i <= 1 else 4):
                    for c in src:
                        h.setdefault(c, "compare")
        yield Frame(g(h), Text(*lines, size=13), say=row_says[i])

    assert dp[0][0] is True
    yield Frame(
        g({(0, 0): "match"}),
        Text("answer: dp[0][0] = T  (c* matches nothing, a* matches \"aa\", b matches \"b\")",
             ("O(m·n) time · O(m·n) space", "match")),
        say="dp[0][0] is true, so the pattern matches: c* takes nothing, a* takes \"aa\", b takes \"b\". "
            "Each cell is O(1) work: O(m·n) time and space.",
    )
