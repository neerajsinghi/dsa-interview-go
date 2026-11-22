from walkthrough import Array, Frame, Grid, Text

TITLE = "Interleaving String"
PATTERN = "2-D DP · two strings"


def frames():
    s1, s2, s3 = "aabcc", "dbbca", "aadbbcbcac"
    m, n = len(s1), len(s2)
    dp = [[None] * (n + 1) for _ in range(m + 1)]
    rl = ["\"\""] + list(s1)
    cl = ["\"\""] + list(s2)

    def widgets(hl=None, k=None):
        return (
            Array(list(s3), label="s3", hl={k: "active"} if k is not None else None, cell=40),
            Grid(dp, label="dp[i][j] = s1[:i] and s2[:j] interleave into s3[:i+j]  (rows: s1, cols: s2)",
                 hl=hl, row_labels=rl, col_labels=cl, cell=36),
        )

    yield Frame(
        *widgets(),
        Text(f"s1 = \"{s1}\", s2 = \"{s2}\""),
        speaker="Interviewer",
        say="Can s3 be formed by riffling s1 and s2 together, keeping each one's letters in order?",
    )
    yield Frame(
        *widgets(),
        Text("greedy fails when both strings offer the next letter", ("trying both branches → O(2^(m+n))", "reject")),
        say="Greedy fails: when both strings offer the next letter, I have to try both. "
            "Plain branching is exponential, but the state is just how much of each string I've used.",
    )

    def cell(i, j):
        k = i + j - 1
        if i == 0 and j == 0:
            dp[i][j] = True
            return
        up = i > 0 and dp[i - 1][j] and s1[i - 1] == s3[k]
        left = j > 0 and dp[i][j - 1] and s2[j - 1] == s3[k]
        dp[i][j] = bool(up or left)

    cell(0, 0)
    for j in range(1, n + 1):
        cell(0, j)
    yield Frame(
        *widgets({(0, j): "match" if dp[0][j] else "reject" for j in range(n + 1)}),
        Text("dp[0][0] = T · row 0 uses only s2: \"d\" != s3[0] = \"a\" → F from there on"),
        say="dp[i][j] means: using i letters of s1 and j of s2 builds the first i+j of s3. "
            "Empty plus empty works. Row 0 uses only s2, and 'd' doesn't match 'a'.",
    )

    # Row 1, first two cells individually.
    cell(1, 0)
    yield Frame(
        *widgets({(1, 0): "active", (0, 0): "compare"}, k=0),
        Text("from above: dp[0][0] = T and s1[0] = 'a' == s3[0] = 'a'", ("dp[1][0] = T", "match")),
        say="A cell's last letter came from s1, moving down, or from s2, moving right. "
            "Here s1's 'a' matches s3[0], and the cell above is true.",
    )
    cell(1, 1)
    yield Frame(
        *widgets({(1, 1): "active", (0, 1): "compare", (1, 0): "compare"}, k=1),
        Text("from above: dp[0][1] = F", "from left:  dp[1][0] = T but s2[0] = 'd' != s3[1] = 'a'",
             ("dp[1][1] = F", "reject")),
        say="For dp[1][1] the next letter of s3 is 'a'. Above is false, and coming from the left would need 'd'. "
            "So it's false.",
    )
    for j in range(2, n + 1):
        cell(1, j)
    row_says = {
        1: "Same rule for the rest of row 1. Only the cells where s3 can keep going stay true.",
        2: "Row 2 adds the second 'a' of s1. Now s2's 'd', 'b', 'b', 'c' can follow it, one by one.",
        3: "Row 3 adds s1's 'b'. A cell is true only if a true neighbour leads in and the letter it adds is the next letter of s3.",
        4: "Row 4 adds s1's first 'c'. The true cells move right, because more of s2 has to be used by now.",
        5: "Row 5 adds the last 'c'. The corner says all of s1 and s2 make all of s3.",
    }
    for i in range(1, m + 1):
        if i > 1:
            for j in range(n + 1):
                cell(i, j)
        h = {(i, j): "match" if dp[i][j] else "active" for j in range(n + 1)}
        yield Frame(*widgets(h), Text(f"row {i}: " + " ".join("T" if v else "F" for v in dp[i])), say=row_says[i])

    assert dp[m][n] is True
    # Walk back to show one interleaving.
    path, i, j = {}, m, n
    while i or j:
        path[(i, j)] = "match"
        k = i + j - 1
        if i > 0 and dp[i - 1][j] and s1[i - 1] == s3[k]:
            i -= 1
        else:
            j -= 1
    path[(0, 0)] = "match"
    yield Frame(
        *widgets(path),
        Text(f"answer: dp[{m}][{n}] = T", ("O(m·n) time · O(n) space with one rolling row", "match")),
        say="The answer is true. The green path is one valid riffle: down takes from s1, right from s2. "
            "Each row needs only the row above, so one row suffices: O(m·n) time, O(n) space.",
    )
