import heapq

from walkthrough import Array, Frame, Grid, Row, Text

TITLE = "Swim in Rising Water"
PATTERN = "Minimax Dijkstra"

INF = float("inf")
DIRS = [(1, 0), (-1, 0), (0, 1), (0, -1)]


def frames():
    grid = [[0, 2], [1, 3]]
    n = len(grid)

    def boards(best, hl_g=None, hl_b=None):
        bb = [["∞" if best[r * n + c] == INF else best[r * n + c] for c in range(n)] for r in range(n)]
        return Row(Grid(grid, label="elevation", hl=hl_g, cell=56),
                   Grid(bb, label="best[cell] = earliest time to stand there", hl=hl_b, cell=56), weights=[1, 1.4])

    def heap_arr(h):
        items = sorted(h)
        return Array([f"({i // n},{i % n}):{t}" for t, i in items] or ["—"], label="min-heap (cell:time)",
                     show_index=False, cell=78)

    best = [INF] * (n * n)
    corners = {(0, 0): "special", (n - 1, n - 1): "special"}
    yield Frame(
        boards(best, hl_g=corners),
        Text("at time t you may swim between neighbours if both elevations ≤ t"),
        speaker="Interviewer",
        say="At time t the water is at level t everywhere. You can swim between adjacent cells when both are at most t. "
            "What's the earliest time you can reach the bottom-right?",
    )
    yield Frame(
        boards(best, hl_g=corners),
        Text("brute force: for t = 0, 1, 2, …: BFS using cells ≤ t", ("O(n²) BFS × n² values of t = O(n⁴)", "reject")),
        say="I could try every time t and run a BFS through cells at most t. That's n squared BFS runs. "
            "Instead, a route's cost is its highest cell, so minimax Dijkstra finds it directly.",
    )

    best[0] = grid[0][0]
    h = [(grid[0][0], 0)]
    yield Frame(
        boards(best, hl_g={(0, 0): "active"}, hl_b={(0, 0): "active"}),
        heap_arr(h),
        Text(f"best[(0,0)] = grid[0][0] = {grid[0][0]};  relax: t = max(best[u], grid[v])"),
        say="I can't start before the water covers my own cell, so the start costs its elevation. "
            "Entering a neighbour costs the larger of my time and its elevation.",
    )
    answer = None
    done = set()
    while h:
        t, u = heapq.heappop(h)
        if t > best[u]:
            continue
        r, c = divmod(u, n)
        if r == n - 1 and c == n - 1:
            answer = t
            break
        done.add(u)
        lines = [f"pop ({r},{c}) time {t}"]
        hl_g = {(r, c): "active"}
        hl_b = {(x // n, x % n): "visited" for x in done}
        hl_b[(r, c)] = "active"
        for dr, dc in DIRS:
            nr, nc = r + dr, c + dc
            if not (0 <= nr < n and 0 <= nc < n):
                continue
            v = nr * n + nc
            nt = max(t, grid[nr][nc])
            if nt < best[v]:
                lines.append((f"({nr},{nc}): max({t}, {grid[nr][nc]}) = {nt} < {'∞' if best[v] == INF else best[v]} → push", "match"))
                best[v] = nt
                heapq.heappush(h, (nt, v))
                hl_g[(nr, nc)] = "match"
                hl_b[(nr, nc)] = "match"
            else:
                lines.append((f"({nr},{nc}): max({t}, {grid[nr][nc]}) = {nt} ≥ {best[v]}, no change", "reject"))
                hl_g[(nr, nc)] = "reject"
        if any(l[1] == "match" for l in lines[1:]):
            say = (f"Pop ({r},{c}) at time {t}; nothing can reach it sooner. "
                   "Each neighbour's time is the higher of my time and its elevation.")
        else:
            say = f"Pop ({r},{c}) at time {t}. Its neighbours already have times at least as good, so nothing changes."
        yield Frame(boards(best, hl_g, hl_b), heap_arr(h), Text(*lines), say=say)

    assert answer == 3
    yield Frame(
        boards(best, hl_g={(n - 1, n - 1): "match"}, hl_b={(n - 1, n - 1): "match"}),
        heap_arr(h),
        Text(f"pop ({n - 1},{n - 1}) time {answer} → target", (f"answer = {answer}", "match"),
             "O(n² log n) time, O(n²) space"),
        say=f"The target pops at time {answer}: I must wait until the water covers the corner itself. "
            "n squared cells through a heap is O(n squared log n) time and O(n squared) space.",
    )
