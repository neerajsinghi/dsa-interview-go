from walkthrough import INK, Frame, Grid, Row, Text

TITLE = "Sudoku Solver"
PATTERN = "Backtracking + MRV"

PUZZLE = [
    "53..7....", "6..195...", ".98....6.",
    "8...6...3", "4..8.3..1", "7...2...6",
    ".6....28.", "...419..5", "....8..79",
]
ALL = (1 << 9) - 1


def popcount(x):
    return bin(x).count("1")


class Sudoku(Grid):
    """A 9x9 Grid with thick lines around the 3x3 boxes."""

    def draw(self, d, x, y, width):
        super().draw(d, x, y, width)
        c, lw = self._geom(width)
        y0 = y + (20 if self.label else 0)
        x0 = x + lw
        for k in range(0, 10, 3):
            d.line([(x0 + k * c, y0), (x0 + k * c, y0 + 9 * c)], fill=INK, width=2)
            d.line([(x0, y0 + k * c), (x0 + 9 * c, y0 + k * c)], fill=INK, width=2)


def solve_trace(board):
    """SolveSudoku with MRV, recording every placement, dead end and undo."""
    rows, cols, boxes = [0] * 9, [0] * 9, [0] * 9
    empty = []
    for r in range(9):
        for c in range(9):
            if board[r][c] == ".":
                empty.append([r, c])
                continue
            bit = 1 << (int(board[r][c]) - 1)
            b = (r // 3) * 3 + c // 3
            assert not (rows[r] & bit or cols[c] & bit or boxes[b] & bit)
            rows[r] |= bit
            cols[c] |= bit
            boxes[b] |= bit
    events = []

    def solve(remaining):
        if remaining == 0:
            return True
        best, best_count = -1, 10
        for i in range(remaining):
            r, c = empty[i]
            free = ALL & ~(rows[r] | cols[c] | boxes[(r // 3) * 3 + c // 3])
            k = popcount(free)
            if k < best_count:
                best, best_count = i, k
                if k == 0:
                    events.append(("dead", r, c, None, 0, [row[:] for row in board]))
                    return False
        last = remaining - 1
        empty[best], empty[last] = empty[last], empty[best]
        r, c = empty[last]
        b = (r // 3) * 3 + c // 3
        free = ALL & ~(rows[r] | cols[c] | boxes[b])
        while free:
            bit = free & -free
            free &= ~bit
            rows[r] |= bit
            cols[c] |= bit
            boxes[b] |= bit
            board[r][c] = str(bit.bit_length())
            events.append(("place", r, c, board[r][c], best_count, [row[:] for row in board]))
            if solve(remaining - 1):
                return True
            rows[r] &= ~bit
            cols[c] &= ~bit
            boxes[b] &= ~bit
            board[r][c] = "."
            events.append(("undo", r, c, None, best_count, [row[:] for row in board]))
        return False

    ok = solve(len(empty))
    return ok, events


def candidates(board, r, c):
    used = set(board[r]) | {board[k][c] for k in range(9)}
    br, bc = (r // 3) * 3, (c // 3) * 3
    used |= {board[br + i][bc + j] for i in range(3) for j in range(3)}
    return [str(d) for d in range(1, 10) if str(d) not in used]


def frames():
    clues = [list(r) for r in PUZZLE]
    board = [list(r) for r in PUZZLE]

    def grid(b, hl=None, label="board"):
        base = {(r, c): "visited" for r in range(9) for c in range(9) if PUZZLE[r][c] != "."}
        base.update(hl or {})
        shown = [[("" if v == "." else v) for v in row] for row in b]
        return Sudoku(shown, label=label, hl=base, cell=30)

    n_empty = sum(row.count(".") for row in PUZZLE)
    yield Frame(
        grid(clues),
        Text(f"{n_empty} empty cells", "each row, column and 3×3 box must hold 1-9 once"),
        speaker="Interviewer",
        say="Fill the empty cells so every row, column and 3 by 3 box contains each digit from 1 to 9 exactly once.",
    )
    yield Frame(
        grid(clues),
        Text("brute force: try 1-9 in every empty cell, check at the end",
             (f"up to 9^{n_empty} boards", "reject")),
        say="Trying every digit in every cell and checking at the end is hopeless. "
            "I check constraints as I go, and I choose carefully which cell to fill next.",
    )

    ok, events = solve_trace(board)
    assert ok

    # The first placement in detail: MRV picks a cell with one candidate.
    places = [e for e in events if e[0] == "place"]
    undos = [e for e in events if e[0] == "undo"]
    dead = [e for e in events if e[0] == "dead"]
    _, r0, c0, d0, k0, _ = places[0]
    cand = candidates(clues, r0, c0)
    assert len(cand) == k0
    yield Frame(
        Row(grid(clues, hl={(r0, c0): "active"}),
            Text("masks per row, column, box:", "rows[r] | cols[c] | boxes[b]", "= digits already used",
                 "", f"cell ({r0},{c0}):", f"candidates {cand}", (f"count = {k0} → fewest", "active"), size=14),
            weights=[1.1, 1]),
        say="I keep a 9-bit mask of used digits for each row, column and box. Each step I scan the empty cells "
            "and pick the one with the fewest candidates.",
    )

    shown = 0
    for ev, r, c, d, k, snap in events[:4]:
        if ev != "place":
            break
        shown += 1
        yield Frame(
            Row(grid(snap, hl={(r, c): "match"}),
                Text(f"place {d} at ({r},{c})", f"candidates: {k}", "", "forced: no choice" if k == 1 else "a guess",
                     size=14),
                weights=[1.1, 1]),
            say=(f"Cell ({r},{c}) has only one legal digit, {d}, so it's forced. Placing it updates three masks, "
                 "which can make other cells forced too.") if k == 1 else
                f"Cell ({r},{c}) has {k} candidates. I try the smallest, {d}, and undo it if the branch fails.",
        )

    guesses = sum(1 for e in places if e[4] > 1)
    yield Frame(
        Row(grid(places[len(places) // 2][5], hl={(places[len(places) // 2][1], places[len(places) // 2][2]): "match"}),
            Text(f"placements so far: {len(places) // 2 + 1}", "", "whole search:", f"guesses (> 1 candidate): {guesses}",
                 f"undos: {len(undos)}", f"dead ends: {len(dead)}", size=14),
            weights=[1.1, 1]),
        say=("Same again for the rest. On this puzzle, MRV always finds a cell with one candidate, "
             "so it never has to guess or backtrack.") if guesses == 0 else
            f"Same again for the rest. MRV had to guess {guesses} times and undo {len(undos)} placements.",
    )

    solved = ["".join(row) for row in board]
    assert solved[:3] == ["534678912", "672195348", "198342567"]
    yield Frame(
        Row(grid(board, hl={(r, c): "match" for r in range(9) for c in range(9) if PUZZLE[r][c] == "."},
                 label="solved"),
            Text(f"{len(places)} placements, {len(undos)} undos", ("solved → true", "match"),
                 "time: exponential worst case", "space: O(81)", size=14),
            weights=[1.1, 1]),
        say=f"Solved with {len(places)} placements. The worst case is exponential, but MRV keeps real puzzles to "
            "milliseconds. The masks and recursion use O(81) space.",
    )
