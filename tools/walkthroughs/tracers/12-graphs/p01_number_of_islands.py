from walkthrough import Frame, Grid, Row, Text, stack_widget

TITLE = "Number of Islands"
PATTERN = "Grid DFS flood fill"

DIRS4 = [(1, 0), (-1, 0), (0, 1), (0, -1)]
ISLAND_STYLES = ["match", "special", "window"]


def frames():
    rows = ["11000", "11000", "00100", "00011"]
    expected = 3
    grid = [list(r) for r in rows]
    m, n = len(grid), len(grid[0])
    seen = [[False] * n for _ in range(m)]
    owner = {}

    def board(cur=None, cur_style="active"):
        hl = {}
        for r in range(m):
            for c in range(n):
                if grid[r][c] == "0":
                    hl[(r, c)] = "empty"
        for cell, k in owner.items():
            hl[cell] = ISLAND_STYLES[k % len(ISLAND_STYLES)]
        if cur is not None:
            hl[cur] = cur_style
        return Grid(grid, label="grid", hl=hl, cell=44)

    yield Frame(
        board(),
        speaker="Interviewer",
        say="1 is land and 0 is water. Land cells that touch up, down, left or right form one island. "
            "How many islands are there?",
    )
    yield Frame(
        board(),
        Text("brute force: for each land cell, search its whole island",
             "and compare against islands already found",
             ("the same island is explored once per cell: O((m·n)²)", "reject")),
        say="A naive approach explores the island around every land cell and dedupes the results. "
            "The same island gets explored once per cell, so it's O((m·n)²).",
    )

    islands = 0
    frames_out = []
    detail = [True]  # show every DFS step for the first island only

    def sink(r, c, stack, k):
        if r < 0 or r >= m or c < 0 or c >= n or grid[r][c] != "1" or seen[r][c]:
            return
        seen[r][c] = True
        owner[(r, c)] = k
        stack.append(f"{r},{c}")
        if detail[0]:
            frames_out.append(Frame(
                Row(board((r, c)), stack_widget(list(stack), label="DFS call stack (r,c)"), weights=[1, 1]),
                Text(f"mark ({r},{c}) seen, then try down, up, right, left"),
                say=(f"({r},{c}) is unseen land, so islands becomes {k + 1}. I mark it seen before recursing, "
                     "then try its neighbours: down, up, right, left."
                     if len(stack) == 1 else
                     f"({r},{c}) is unseen land next to the cell below it on the stack, so the DFS "
                     "steps in and marks it. Water and seen cells return at once."),
            ))
        for dr, dc in DIRS4:
            sink(r + dr, c + dc, stack, k)
        stack.pop()

    skipped_note = False
    for r in range(m):
        for c in range(n):
            if grid[r][c] == "1" and not seen[r][c]:
                islands += 1
                sink(r, c, [], islands - 1)
                if frames_out:
                    yield from frames_out
                    frames_out.clear()
                    yield Frame(
                        board(),
                        Text((f"island {islands} fully marked", "match"), f"islands = {islands}"),
                        say="The DFS has returned, so every cell of this island is marked. "
                            "I count it once and go back to scanning row by row.",
                    )
                    detail[0] = False
                else:
                    size = sum(1 for v in owner.values() if v == islands - 1)
                    yield Frame(
                        board((r, c), "compare"),
                        Text(f"({r},{c}) is unseen land -> new island", (f"islands = {islands}", "match"),
                             f"flood fill marked {size} cell{'s' if size > 1 else ''}"),
                        say=(f"Scanning on, ({r},{c}) is land I haven't seen, so it starts island {islands}. "
                             + ("It has no land neighbours, so the fill marks just this cell."
                                if size == 1 else
                                f"The flood fill marks all {size} of its cells.")),
                    )
            elif grid[r][c] == "1" and not skipped_note and (r, c) != (0, 0):
                skipped_note = True
                yield Frame(
                    board((r, c), "visited"),
                    Text(f"({r},{c}) is land but already seen -> skip", f"islands = {islands}"),
                    say=f"The scan reaches ({r},{c}). It's land, but the flood fill already marked it, "
                        "so it belongs to an island I've counted. Skip.",
                )

    assert islands == expected, islands
    yield Frame(
        board(),
        Text((f"answer: {islands}", "match"),
             "time: O(m·n), every cell is marked once",
             "space: O(m·n) for seen and the recursion"),
        say="Three islands. Each cell is marked at most once and checks four neighbours, "
            "so it's O(m·n) time and O(m·n) space for the seen grid and the stack.",
    )
