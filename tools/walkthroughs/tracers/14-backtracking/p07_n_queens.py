from walkthrough import Frame, Grid, KV, Row, Text

TITLE = "N-Queens I & II"
PATTERN = "Backtracking"


def render(queen_col):
    n = len(queen_col)
    return ["." * c + "Q" + "." * (n - c - 1) for c in queen_col]


def frames():
    n = 4

    # SolveNQueens, recording events.
    res = []
    queen_col = [0] * n
    cols, diag, anti = [False] * n, [False] * (2 * n), [False] * (2 * n)
    events = []

    def place(r):
        if r == n:
            res.append(render(queen_col))
            events.append(("solution", r, None, list(queen_col[:r]), render(queen_col)))
            return
        rejected = []
        placed_any = False
        for c in range(n):
            if cols[c] or diag[r - c + n] or anti[r + c]:
                rejected.append(c)
                continue
            queen_col[r] = c
            placed_any = True
            events.append(("place", r, c, list(queen_col[:r + 1]), list(rejected)))
            cols[c], diag[r - c + n], anti[r + c] = True, True, True
            place(r + 1)
            cols[c], diag[r - c + n], anti[r + c] = False, False, False
        if not placed_any:
            events.append(("dead", r, None, list(queen_col[:r]), list(rejected)))

    place(0)

    def board(queens, row=None, rejected=(), chosen=None, dead=False, sol=False):
        cells = [["" for _ in range(n)] for _ in range(n)]
        hl = {}
        for r, c in enumerate(queens):
            cells[r][c] = "Q"
            hl[(r, c)] = "match" if sol else "visited"
        if row is not None and row < n:
            for c in rejected:
                cells[row][c] = "×"
                hl[(row, c)] = "reject"
            if chosen is not None:
                hl[(row, chosen)] = "active"
        return Grid(cells, label=f"{n}×{n} board", hl=hl, row_labels=[f"r{r}" for r in range(n)],
                    col_labels=[f"c{c}" for c in range(n)], cell=50)

    def sets(queens):
        return Text(f"cols  = {{{', '.join(str(c) for c in queens)}}}",
                    f"r-c   = {{{', '.join(str(r - c) for r, c in enumerate(queens))}}}",
                    f"r+c   = {{{', '.join(str(r + c) for r, c in enumerate(queens))}}}")

    yield Frame(
        board([]),
        Text(f"n = {n}: place {n} queens, none attacking another", "attack = same row, column or diagonal"),
        speaker="Interviewer",
        say="Place 4 queens on a 4 by 4 board so that no two share a row, column or diagonal. Return every board.",
    )
    yield Frame(
        board([]),
        Text("brute force: choose 4 of 16 squares = 1820 boards, check each",
             ("most boards fail on the first two queens", "reject")),
        say="Trying every set of 4 squares means 1820 boards. Two queens can't share a row, "
            "so I place exactly one queen per row and check each square against the queens above in O(1).",
    )

    solutions_seen = 0
    for ev, r, c, queens, rejected in events:
        if solutions_seen >= 1:
            continue
        if ev == "place":
            if r == 0 and c == 0:
                say = ("Row 0, column 0. I record the column, the r-c diagonal and the r+c anti-diagonal in sets, "
                       "so later rows can check a square in O(1).")
            elif [x for x in range(c) if x not in rejected]:
                prev = [x for x in range(c) if x not in rejected][-1]
                say = (f"Row {r}: the queen in column {prev} led only to dead ends, so I unchoose it. "
                       f"Column {c} is the next free square.")
            elif rejected:
                say = (f"Row {r}: column{'s' if len(rejected) > 1 else ''} {', '.join(map(str, rejected))} "
                       f"{'are' if len(rejected) > 1 else 'is'} attacked. Column {c} is free, so I place a queen there.")
            else:
                say = f"Row {r}: column {c} is free, so the queen goes there."
            if r == 0 and c > 0:
                say = (f"Back up to row 0. Every board with the queen in column {c - 1} failed, "
                       f"so I unchoose it and try column {c}.")
            label = f"row {r}: rejected {rejected or 'none'}, place at c{c}"
            yield Frame(Row(board(queens, r, rejected, c), Text(label, mono=True, size=14), weights=[1, 1]),
                        sets(queens), say=say)
        elif ev == "dead":
            yield Frame(Row(board(queens, r, rejected), Text(f"row {r}: every column attacked", ("dead end → backtrack", "reject"),
                                                              size=14), weights=[1, 1]),
                        sets(queens),
                        say=f"Every square in row {r} is attacked, so this branch is dead. "
                            f"I return, and row {r - 1} unchooses its queen and tries the next column.")
        elif ev == "solution":
            solutions_seen += 1
            yield Frame(Row(board(queens, sol=True), Text(("r == n → record:", "match"), *rejected, size=14),
                            weights=[1, 1]),
                        sets(queens),
                        say="All 4 rows have a queen, so this board is a solution. I copy it into the results "
                            "and keep backtracking to find the rest.")

    # Summarise the rest of the search.
    idx = [k for k, e in enumerate(events) if e[0] == "solution"]
    later_places = sum(1 for e in events[idx[0] + 1:] if e[0] == "place")
    later_dead = sum(1 for e in events[idx[0] + 1:] if e[0] == "dead")
    expected = [[".Q..", "...Q", "Q...", "..Q."], ["..Q.", "Q...", "...Q", ".Q.."]]
    assert res == expected
    second = [row.index("Q") for row in res[1]]
    yield Frame(
        Row(board(second, sol=True), Text(("record:", "match"), *res[1], size=14), weights=[1, 1]),
        Text(f"after the first solution: {later_places} more placements, {later_dead} dead ends",
             ("second solution (the mirror image)", "match")),
        say="The search continues the same way. It finds the mirror image as the second solution, "
            "and every other branch hits a dead end.",
    )
    yield Frame(
        Row(board([row.index("Q") for row in res[0]], sol=True), board(second, sol=True), weights=[1, 1]),
        Text(f"answer: {len(res)} boards", "O(n!) time, O(n) space for the sets and recursion"),
        say="Two solutions for n = 4. Row r has at most n minus r safe columns, so the search is O(n factorial). "
            "The sets and recursion take O(n) space.",
    )
