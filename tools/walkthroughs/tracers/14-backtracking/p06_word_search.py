from collections import Counter

from walkthrough import Array, Frame, Grid, KV, Row, Text

TITLE = "Word Search"
PATTERN = "Grid backtracking"


def frames():
    board_rows = ["ABCE", "SFCS", "ADEE"]
    original_word = "ABCCED"
    board = [list(r) for r in board_rows]
    m, n = len(board), len(board[0])

    def grid(hl=None, label="board"):
        shown = [[(ch if ch else "·") for ch in row] for row in board]
        return Grid(shown, label=label, hl=hl, cell=48)

    yield Frame(
        grid(),
        Text(f'word = "{original_word}"', "adjacent = up/down/left/right; each cell used at most once"),
        speaker="Interviewer",
        say="Can the word be spelled by a path of neighbouring cells, using each cell at most once?",
    )
    yield Frame(
        grid(),
        Text("brute force: every path of length L from every cell",
             ("m·n starts × 4^L paths, most fail on letter 2 or 3", "reject")),
        say="Every path from every cell is far too many, but most die on the first wrong letter. "
            "So I DFS from each matching cell and backtrack as soon as a letter doesn't match.",
    )

    # Pruning exactly as Exist does it.
    word = original_word
    have = Counter(ch for row in board for ch in row)
    need = Counter(word)
    assert len(word) <= m * n and all(need[c] <= have[c] for c in need)
    reversed_word = have[word[-1]] < have[word[0]]
    if reversed_word:
        word = word[::-1]
    yield Frame(
        grid(hl={(r, c): ("special" if board[r][c] == word[0] else None) for r in range(m) for c in range(n)
                 if board[r][c] == word[0]}),
        Row(KV({ch: f"{need[ch]}/{have[ch]}" for ch in sorted(need)}, label="letters needed / in grid"),
            Text(f"have['{original_word[-1]}'] = {have[original_word[-1]]} < have['{original_word[0]}'] = {have[original_word[0]]}",
                 (f'search "{word}" instead', "active")), weights=[1, 1]),
        say=f"Two cheap checks first. The grid has enough of every letter. And '{original_word[-1]}' is rarer than "
            f"'{original_word[0]}', so I search the reversed word: fewer starting cells.",
    )

    events = []

    def dfs(r, c, i):
        if r < 0 or r >= m or c < 0 or c >= n or board[r][c] != word[i]:
            if 0 <= r < m and 0 <= c < n and i > 0 and board[r][c]:  # skip cells already on the path
                events.append(("fail", r, c, i))
            return False
        if i == len(word) - 1:
            events.append(("found", r, c, i))
            return True
        events.append(("step", r, c, i))
        ch = board[r][c]
        board[r][c] = ""  # choose: mark used
        found = dfs(r + 1, c, i + 1) or dfs(r - 1, c, i + 1) or dfs(r, c + 1, i + 1) or dfs(r, c - 1, i + 1)
        board[r][c] = ch  # unchoose: restore
        if not found:
            events.append(("back", r, c, i))
        return found

    result = False
    starts_tried = 0
    for r in range(m):
        for c in range(n):
            starts_tried += 1
            if dfs(r, c, 0):
                result = True
                break
        if result:
            break

    # Replay the events on a copy of the board to draw each step.
    path = []
    fails = []
    for ev, r, c, i in events:
        if ev == "fail":
            fails.append((r, c))
            continue
        if ev == "back":
            path.pop()
            continue
        path.append((r, c))
        hl = {p: "visited" for p in path[:-1]}
        for f in fails:
            hl[f] = "reject"
        hl[(r, c)] = "match"
        prefix = word[:i + 1]
        lines = [f'matched "{prefix}"  (index {i})']
        if fails:
            lines.append((f"rejected neighbours: " + ", ".join(f"({a},{b})={board_rows[a][b]}" for a, b in fails), "reject"))
        if ev == "found":
            say = f"The last letter matches at ({r},{c}), so the whole word is on the board. Return true right away."
        elif i == 0:
            say = (f"The scan reaches the only '{word[0]}' at ({r},{c}). I mark it used by overwriting it, "
                   "then try down, up, right, left for the next letter.")
        else:
            say = (f"'{word[i]}' matches at ({r},{c}). Mark it, then look for '{word[i + 1]}'. "
                   + ("Red cells were tried and didn't match." if fails else ""))
        yield Frame(
            Row(Grid(board_rows, label="board (grey = on the path, marked used)", hl=hl, cell=48),
                Array(list(word), label="word being searched", hl={k: ("match" if k <= i else None) for k in range(len(word))},
                      cell=40), weights=[1, 1.1]),
            Text(*lines),
            say=say,
        )
        fails = []

    assert result is True
    assert ["".join(row) for row in board] == board_rows  # board restored
    yield Frame(
        Grid(board_rows, label="board (restored)", hl={p: "match" for p in path}, cell=48),
        Text(f'"{original_word}" found (as "{word}")', ("answer = true", "match"),
             "O(m·n·3^L) time, O(L) recursion"),
        say="True. Every mark is undone on the way back, so the board is restored. "
            "Each start explores up to 3 directions per letter: O(mn times 3 to the L) time, O(L) recursion.",
    )
