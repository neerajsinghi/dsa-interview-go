from walkthrough import Frame, Grid, Row, Text

TITLE = "Pacific Atlantic Water Flow"
PATTERN = "Reverse multi-source BFS"

DIRS4 = [(1, 0), (-1, 0), (0, 1), (0, -1)]


def climb(heights, starts):
    """Port of the climb closure: BFS uphill. Also records the BFS depth of each cell."""
    m, n = len(heights), len(heights[0])
    seen = [[False] * n for _ in range(m)]
    depth = {}
    for r, c in starts:
        seen[r][c] = True
        depth[(r, c)] = 0
    queue = list(starts)
    while queue:
        r0, c0 = queue.pop(0)
        for dr, dc in DIRS4:
            r, c = r0 + dr, c0 + dc
            if r < 0 or r >= m or c < 0 or c >= n or seen[r][c]:
                continue
            if heights[r][c] >= heights[r0][c0]:
                seen[r][c] = True
                depth[(r, c)] = depth[(r0, c0)] + 1
                queue.append((r, c))
    return seen, depth


def frames():
    heights = [[1, 2, 2, 3, 5], [3, 2, 3, 4, 4], [2, 4, 5, 3, 1], [6, 7, 1, 4, 5], [5, 1, 1, 2, 4]]
    expected = [[0, 4], [1, 3], [1, 4], [2, 2], [3, 0], [3, 1], [4, 0]]
    m, n = len(heights), len(heights[0])

    def board(hl, label="heights"):
        return Grid(heights, label=label, hl=hl, cell=40)

    yield Frame(
        board({}),
        Text("Pacific: top and left edges.  Atlantic: bottom and right edges.",
             "water flows to a neighbour with height <= its own"),
        speaker="Interviewer",
        say="Rain flows to equal or lower neighbours and off the edges into the ocean. Return every cell "
            "whose water can reach both the Pacific and the Atlantic.",
    )
    yield Frame(
        board({}),
        Text("brute force: from every cell, search downhill and see which oceans it touches",
             ("m·n searches × O(m·n) each = O((m·n)²)", "reject")),
        say="The direct way runs a downhill search from every cell. That's one O(m·n) search per cell, "
            "so O((m·n)²), and neighbouring cells repeat almost the same work.",
    )

    pac = [(r, 0) for r in range(m)] + [(0, c) for c in range(n)]
    atl = [(r, n - 1) for r in range(m)] + [(m - 1, c) for c in range(n)]
    to_p, dp = climb(heights, pac)
    to_a, da = climb(heights, atl)

    for name, depth, style in (("Pacific", dp, "active"), ("Atlantic", da, "window")):
        levels = max(depth.values())
        for lvl in range(levels + 1):
            hl = {cell: style for cell, d in depth.items() if d < lvl}
            new = sorted(cell for cell, d in depth.items() if d == lvl)
            for cell in new:
                hl[cell] = "compare" if lvl else style
            if lvl == 0:
                say = (f"I reverse the flow. Every {name} edge cell drains into the {name}, so all of them "
                       "start in one queue. From there I climb uphill, to neighbours at least as high."
                       if name == "Pacific" else
                       "Same again for the Atlantic: start from the bottom and right edges and climb uphill.")
                lines = [f"{name} starts: {len(set(new))} edge cells (corners appear twice, seen stops repeats)"]
            else:
                cells = ", ".join(f"({r},{c})" for r, c in new)
                say = (f"BFS ring {lvl}: {cells} "
                       f"{'are' if len(new) > 1 else 'is'} at least as high as a cell already reached, so water "
                       f"from there can flow down to the {name}.")
                lines = [f"ring {lvl}: " + cells]
            if lvl == levels:
                lines.append((f"{name} reach: {sum(map(sum, to_p if name == 'Pacific' else to_a))} cells", style))
            yield Frame(board(hl, f"heights - {name} climb"), Text(*lines), say=say)

    res = [[r, c] for r in range(m) for c in range(n) if to_p[r][c] and to_a[r][c]]
    both = {(r, c): "match" for r, c in res}
    only = {}
    for r in range(m):
        for c in range(n):
            if to_p[r][c] and not to_a[r][c]:
                only[(r, c)] = "active"
            elif to_a[r][c] and not to_p[r][c]:
                only[(r, c)] = "window"
    yield Frame(
        board({**only, **both}, "heights - green = both oceans"),
        Text((f"both: {len(res)} cells", "match"), "blue = Pacific only, yellow = Atlantic only"),
        say="A cell reached by both climbs can drain into both oceans. So the answer is the intersection "
            "of the two seen grids, read in row-major order.",
    )

    assert res == expected, res
    yield Frame(
        board(both, "heights - answer"),
        Text((str(res).replace(" ", ""), "match"),
             "two multi-source BFS runs: O(m·n) time, O(m·n) space"),
        say="Two BFS runs, one per ocean, each visit every cell at most once. So it's O(m·n) time and space, "
            "instead of one search per cell.",
    )
