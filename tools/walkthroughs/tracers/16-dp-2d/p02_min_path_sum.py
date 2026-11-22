from walkthrough import Frame, Grid, Row, Text

TITLE = "Minimum Path Sum"
PATTERN = "2-D DP · grid"


def frames():
    grid = [[1, 3, 1], [1, 5, 1], [4, 2, 1]]
    m, n = len(grid), len(grid[0])
    dp = [[None] * n for _ in range(m)]

    def widgets(hl=None, hl_g=None):
        return (Row(Grid(grid, label="grid (cost of each cell)", hl=hl_g, row_labels=range(m), col_labels=range(n), cell=46),
                    Grid(dp, label="best[i][j] = cheapest route to (i,j)", hl=hl, row_labels=range(m),
                         col_labels=range(n), cell=46)),)

    yield Frame(
        *widgets(),
        Text("move right or down · total = sum of the cells you step on"),
        speaker="Interviewer",
        say="Each cell has a cost. Moving only right or down, go from the top-left to the bottom-right. "
            "Return the smallest total cost, both ends included.",
    )
    yield Frame(
        *widgets(),
        Text("try every route", ("routes share cells → O(2^(m+n))", "reject")),
        say="Trying every route is exponential, and routes share most of their cells. "
            "The cheapest way to reach a cell doesn't depend on where you go next, so I can save it.",
    )

    def compute(i, j):
        c = grid[i][j]
        hl = {(i, j): "active"}
        if i == 0 and j == 0:
            dp[i][j] = c
            return f"best[0][0] = grid[0][0] = {c}", hl
        if i == 0:
            dp[i][j] = dp[i][j - 1] + c
            hl[(i, j - 1)] = "compare"
            return f"best[{i}][{j}] = left {dp[i][j-1]} + {c} = {dp[i][j]}", hl
        if j == 0:
            dp[i][j] = dp[i - 1][j] + c
            hl[(i - 1, j)] = "compare"
            return f"best[{i}][{j}] = up {dp[i-1][j]} + {c} = {dp[i][j]}", hl
        up, left = dp[i - 1][j], dp[i][j - 1]
        dp[i][j] = min(up, left) + c
        hl[(i - 1, j)] = "compare"
        hl[(i, j - 1)] = "compare"
        return f"best[{i}][{j}] = min(up {up}, left {left}) + {c} = {dp[i][j]}", hl

    says = {
        (0, 0): "State: best[i][j] is the cheapest route from the start to (i, j). The start costs its own value, 1.",
        (0, 1): "The first row has one parent, the cell to the left, so I just add.",
        (0, 2): "Same for the rest of the first row: 4 + 1 = 5.",
        (1, 0): "The first column also has one parent, the cell above.",
        (1, 1): "Inner cells come from above or the left. I take the cheaper one and add this cell's cost.",
        (1, 2): "Coming down from 5 costs less than coming from 7 on the left: 5 + 1 = 6.",
        (2, 0): "First column again: 2 + 4 = 6.",
        (2, 1): "Above is 7, left is 6. The left wins: 6 + 2 = 8.",
        (2, 2): "The target: above is 6, left is 8. Taking 6 gives 6 + 1 = 7.",
    }
    for i in range(m):
        for j in range(n):
            line, h = compute(i, j)
            yield Frame(*widgets(h, {(i, j): "active"}), Text(line), say=says[(i, j)])

    ans = dp[m - 1][n - 1]
    assert ans == 7
    path = {(0, 0), (0, 1), (0, 2), (1, 2), (2, 2)}
    yield Frame(
        *widgets({p: "match" for p in path}, {p: "match" for p in path}),
        Text(f"answer: best[{m-1}][{n-1}] = {ans}  (1→3→1→1→1)", ("O(m·n) time · O(n) space with one row", "match")),
        say=f"The answer is {ans}, along 1, 3, 1, 1, 1. The Go rolling version keeps one row: "
            "row[j] is 'above' before I overwrite it. O(m·n) time, O(n) space.",
    )
