from walkthrough import Array, Frame, Grid, Row, Text

TITLE = "Rotting Oranges"
PATTERN = "Multi-source BFS"

DIRS4 = [(1, 0), (-1, 0), (0, 1), (0, -1)]


def frames():
    grid = [[2, 1, 1], [1, 1, 0], [0, 1, 1]]
    expected = 4
    m, n = len(grid), len(grid[0])
    g = [row[:] for row in grid]

    def board(new=()):
        hl = {}
        for r in range(m):
            for c in range(n):
                hl[(r, c)] = {0: "empty", 1: "match", 2: "reject"}[g[r][c]]
        for cell in new:
            hl[cell] = "compare"
        return Grid(g, label="grid (0 empty, 1 fresh, 2 rotten)", hl=hl, cell=46)

    def qview(q, label):
        return Array([f"{r},{c}" for r, c in q] or ["-"], label=label, show_index=False, cell=52)

    yield Frame(
        board(),
        speaker="Interviewer",
        say="Each minute, every fresh orange next to a rotten one goes rotten. "
            "How many minutes until no fresh orange is left? Return -1 if that never happens.",
    )
    yield Frame(
        board(),
        Text("brute force: simulate minute by minute",
             "each minute, rescan the whole grid for fresh oranges next to rotten ones",
             ("up to m·n minutes × m·n cells = O((m·n)²)", "reject")),
        say="I could simulate: each minute, scan the whole grid and rot the fresh oranges that touch "
            "a rotten one. That rescans everything each minute, O((m·n)²) in the worst case.",
    )

    queue, fresh = [], 0
    for r in range(m):
        for c in range(n):
            if g[r][c] == 2:
                queue.append((r, c))
            elif g[r][c] == 1:
                fresh += 1
    yield Frame(
        Row(board(), qview(queue, "queue (minute 0)"), weights=[1, 1]),
        Text(f"fresh = {fresh}", "every rotten orange starts in the queue"),
        say="Rot spreads like ripples from every rotten orange at once. So I put all rotten oranges "
            f"in the queue at minute 0 and count the fresh ones: {fresh}.",
    )

    minutes = 0
    while queue and fresh > 0:
        nxt = []
        for r0, c0 in queue:
            for dr, dc in DIRS4:
                r, c = r0 + dr, c0 + dc
                if 0 <= r < m and 0 <= c < n and g[r][c] == 1:
                    g[r][c] = 2  # rot it now so it is enqueued only once
                    fresh -= 1
                    nxt.append((r, c))
        prev = queue
        queue = nxt
        minutes += 1
        cells = ", ".join(f"({r},{c})" for r, c in nxt)
        if minutes == 1:
            say = (f"Minute 1: I process the whole queue as one level. Each fresh neighbour is rotted right "
                   f"away, so it can't be queued twice. {len(nxt)} oranges rot: {cells}.")
        elif fresh == 0:
            say = f"Minute {minutes}: {cells} rots. That was the last fresh orange, so the loop stops."
        else:
            say = f"Minute {minutes}: the oranges that rotted last minute spread to {cells}. {fresh} fresh left."
        yield Frame(
            Row(board(nxt), qview(queue, f"queue after minute {minutes}"), weights=[1, 1]),
            Text(f"processed level: {', '.join(f'({r},{c})' for r, c in prev)}",
                 (f"minute {minutes}: rotted {cells}", "compare"), f"fresh = {fresh}"),
            say=say,
        )

    result = -1 if fresh > 0 else minutes
    assert result == expected, result
    yield Frame(
        board(),
        Text((f"fresh = 0 -> answer: {result} minutes", "match"),
             "if fresh > 0 at the end -> -1 (some orange is walled off)",
             "time: O(m·n)   space: O(m·n) for the queue"),
        say=f"No fresh oranges remain, so the answer is {result}. Each cell is queued at most once, "
            "so it's O(m·n) time and space. Leftover fresh oranges would mean -1.",
    )
