from walkthrough import Frame, Grid, Row, Text

TITLE = "Unique Paths (with obstacles)"
PATTERN = "2-D DP · grid"


def frames():
    grid = [[0, 0, 0], [0, 1, 0], [0, 0, 0]]
    m, n = len(grid), len(grid[0])
    dp = [[None] * n for _ in range(m)]
    rock = {(i, j): "reject" for i in range(m) for j in range(n) if grid[i][j] == 1}

    def widgets(hl=None, hl_board=None):
        board = [["#" if v else "·" for v in row] for row in grid]
        hb = dict(rock)
        hb.update(hl_board or {})
        h = dict(hl or {})
        return (Row(Grid(board, label="grid (# = blocked)", hl=hb, row_labels=range(m), col_labels=range(n), cell=46),
                    Grid(dp, label="ways[i][j] = routes from (0,0)", hl=h, row_labels=range(m),
                         col_labels=range(n), cell=46)),)

    yield Frame(
        *widgets(hl_board={(0, 0): "active", (2, 2): "match"}),
        Text("moves: right or down · blocked cells can't be entered"),
        speaker="Interviewer",
        say="A robot starts top-left and wants the bottom-right cell. It moves only right or down, "
            "and can't enter blocked cells. How many routes are there?",
    )
    yield Frame(
        *widgets(),
        Text("walk(i, j) = walk(i+1, j) + walk(i, j+1)", ("same cells revisited → O(2^(m+n))", "reject")),
        say="Recursing on every route is exponential: lots of routes pass through the same cell, "
            "and each one re-counts the rest of the way from there.",
    )

    def compute(i, j):
        if grid[i][j] == 1:
            dp[i][j] = 0
            return "blocked → 0", {(i, j): "reject"}
        if i == 0 and j == 0:
            dp[i][j] = 1
            return "start → 1", {(i, j): "match"}
        up = dp[i - 1][j] if i > 0 else 0
        left = dp[i][j - 1] if j > 0 else 0
        dp[i][j] = up + left
        hl = {(i, j): "active"}
        if i > 0:
            hl[(i - 1, j)] = "compare"
        if j > 0:
            hl[(i, j - 1)] = "compare"
        return f"ways[{i}][{j}] = up {up} + left {left} = {dp[i][j]}", hl

    # Row 0 in one go.
    hl = {}
    for j in range(n):
        compute(0, j)
        hl[(0, j)] = "match"
    yield Frame(
        *widgets(hl),
        Text("ways[0][0] = 1 · first row: only 'from the left' exists → 1, 1, 1"),
        say="I count routes into each cell instead. The only way into a cell is from above or from the left. "
            "The first row can only be reached from the left, so it's all 1s.",
    )
    says = {
        (1, 0): "Column 0 is the same: only reachable from above, so 1.",
        (1, 1): "This cell is blocked. No route can pass through it, so it holds 0. That's the obstacle rule.",
        (1, 2): "Here the left neighbour is the rock, which adds 0. Only the 1 route from above counts.",
        (2, 0): "Still only from above: 1.",
        (2, 1): "From above is the rock (0); from the left is 1. So 1 route.",
        (2, 2): "The target adds both neighbours: 1 from above (around the right side) plus 1 from the left.",
    }
    for i in range(1, m):
        for j in range(n):
            line, h = compute(i, j)
            yield Frame(*widgets(h), Text(line), say=says[(i, j)])

    ans = dp[m - 1][n - 1]
    assert ans == 2
    yield Frame(
        *widgets({(m - 1, n - 1): "match"}),
        Text(f"answer: ways[{m-1}][{n-1}] = {ans}", ("O(m·n) time · O(n) space with one rolling row", "match")),
        say=f"The answer is {ans}: around the centre either way. Each row only needs the row above, "
            "so the Go code keeps one row of n cells: O(m·n) time, O(n) space.",
    )
