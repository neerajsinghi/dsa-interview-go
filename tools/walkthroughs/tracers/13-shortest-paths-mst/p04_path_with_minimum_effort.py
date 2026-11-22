import heapq

from walkthrough import Array, Frame, Grid, Row, Text

TITLE = "Path With Minimum Effort"
PATTERN = "Minimax Dijkstra"

INF = float("inf")
DIRS = [(1, 0), (-1, 0), (0, 1), (0, -1)]


def frames():
    heights = [[1, 2, 2], [3, 8, 2], [5, 3, 5]]
    m, n = len(heights), len(heights[0])

    def boards(effort, hl_h=None, hl_e=None):
        eff = [["∞" if effort[r * n + c] == INF else effort[r * n + c] for c in range(n)] for r in range(m)]
        return Row(Grid(heights, label="heights", hl=hl_h, cell=46),
                   Grid(eff, label="effort[cell] (best max-step so far)", hl=hl_e, cell=46), weights=[1, 1])

    def heap_arr(h):
        items = sorted(h)
        return Array([f"({i // n},{i % n}):{d}" for d, i in items] or ["—"], label="min-heap (cell:effort)",
                     show_index=False, cell=78)

    effort = [INF] * (m * n)
    yield Frame(
        boards(effort, hl_h={(0, 0): "special", (m - 1, n - 1): "special"}),
        Text("effort of a route = largest |height difference| between neighbours"),
        speaker="Interviewer",
        say="Walk from the top-left to the bottom-right, moving in four directions. "
            "A route's effort is its biggest single step. Find the minimum effort.",
    )
    yield Frame(
        boards(effort, hl_h={(0, 0): "special", (m - 1, n - 1): "special"}),
        Text("brute force: DFS every simple path, take the min of each path's max",
             ("exponential in the number of cells", "reject")),
        say="Trying every path is exponential. But a route's effort never goes down as it gets longer, "
            "and that's all Dijkstra needs. I relax with max instead of plus.",
    )

    effort[0] = 0
    h = [(0, 0)]
    yield Frame(
        boards(effort, hl_h={(0, 0): "active"}, hl_e={(0, 0): "active"}),
        heap_arr(h),
        Text("effort[(0,0)] = 0;  relax: next = max(effort[u], |h[u] - h[v]|)"),
        say="The start costs nothing. Each pop gives the unsettled cell with the smallest effort so far.",
    )
    shown = 0
    done = set()
    answer = None
    while h:
        d, u = heapq.heappop(h)
        if d > effort[u]:
            continue
        r, c = divmod(u, n)
        if r == m - 1 and c == n - 1:
            answer = d
            break
        done.add(u)
        lines = [f"pop ({r},{c}) effort {d}"]
        hl_h = {(r, c): "active"}
        hl_e = {(x // n, x % n): "visited" for x in done}
        hl_e[(r, c)] = "active"
        for dr, dc in DIRS:
            nr, nc = r + dr, c + dc
            if not (0 <= nr < m and 0 <= nc < n):
                continue
            step = abs(heights[nr][nc] - heights[r][c])
            nxt = max(d, step)
            v = nr * n + nc
            if nxt < effort[v]:
                lines.append((f"({nr},{nc}): max({d}, |{heights[nr][nc]}-{heights[r][c]}|={step}) = {nxt} < "
                              f"{'∞' if effort[v] == INF else effort[v]} → push", "match"))
                effort[v] = nxt
                heapq.heappush(h, (nxt, v))
                hl_h[(nr, nc)] = "match"
                hl_e[(nr, nc)] = "match"
            elif v not in done:
                hl_h[(nr, nc)] = "reject"
        shown += 1
        if shown == 1:
            say = "From the start, every neighbour beats infinity, so each one goes on the heap with its step size as its effort."
        elif len(lines) == 1:
            say = f"Pop ({r},{c}) at effort {d}. None of its neighbours improves, so nothing is pushed."
        else:
            say = (f"Pop ({r},{c}) at effort {d}; that value is now final. "
                   "A neighbour's new effort is the larger of my effort and the step to it.")
        if len(lines) == 1:
            lines.append(("no neighbour improves", "reject"))
        yield Frame(boards(effort, hl_h, hl_e), heap_arr(h), Text(*lines[:4]), say=say)

    assert answer == 2
    yield Frame(
        boards(effort, hl_h={(m - 1, n - 1): "match"}, hl_e={(m - 1, n - 1): "match"}),
        heap_arr(h),
        Text(f"pop ({m - 1},{n - 1}) effort {answer} → target, return", (f"answer = {answer}", "match"),
             "O(m·n·log(m·n)) time, O(m·n) space"),
        say=f"The target comes off the heap with effort {answer}, and the first time it's popped it is final. "
            "That's O(mn log mn) time and O(mn) space.",
    )
