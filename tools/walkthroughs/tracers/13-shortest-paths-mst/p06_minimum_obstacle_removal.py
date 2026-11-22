from walkthrough import Array, Frame, Grid, Row, Text

TITLE = "Minimum Obstacle Removal"
PATTERN = "0-1 BFS"

INF = float("inf")
DIRS = [(1, 0), (-1, 0), (0, 1), (0, -1)]


class Deque:
    """Same layout as the Go deque: front is a stack, back is a queue read from head."""

    def __init__(self):
        self.front, self.back, self.head = [], [], 0

    def push_front(self, x):
        self.front.append(x)

    def push_back(self, x):
        self.back.append(x)

    def __len__(self):
        return len(self.front) + len(self.back) - self.head

    def pop_front(self):
        if self.front:
            return self.front.pop()
        x = self.back[self.head]
        self.head += 1
        return x

    def items(self):
        return list(reversed(self.front)) + self.back[self.head:]


def frames():
    grid = [[0, 1, 1], [1, 1, 0], [1, 1, 0]]
    m, n = len(grid), len(grid[0])
    N = m * n
    # gridGraph: the edge into cell (r, c) weighs grid[r][c].
    adj = [[] for _ in range(N)]
    for r in range(m):
        for c in range(n):
            for dr, dc in DIRS:
                nr, nc = r + dr, c + dc
                if 0 <= nr < m and 0 <= nc < n:
                    adj[r * n + c].append((nr * n + nc, grid[nr][nc]))

    def cell(i):
        return f"({i // n},{i % n})"

    def boards(dist, hl_g=None, hl_d=None):
        dd = [["∞" if dist[r * n + c] == INF else dist[r * n + c] for c in range(n)] for r in range(m)]
        base = {(r, c): "compare" for r in range(m) for c in range(n) if grid[r][c]}
        base.update(hl_g or {})
        return Row(Grid(grid, label="grid (1 = obstacle, costs 1 to enter)", hl=base, cell=46),
                   Grid(dd, label="dist = obstacles removed so far", hl=hl_d, cell=46), weights=[1, 1])

    def dq_arr(dq):
        items = dq.items()
        return Array([cell(i) for i in items] or ["—"], label="deque (front on the left)", show_index=False, cell=64)

    dist = [INF] * N
    yield Frame(
        boards(dist, hl_g={(0, 0): "special", (m - 1, n - 1): "special"}),
        Text("move up/down/left/right; entering an obstacle removes it (cost 1)"),
        speaker="Interviewer",
        say="Go from the top-left to the bottom-right corner. You may remove obstacles. What's the fewest you must remove?",
    )
    yield Frame(
        boards(dist),
        Text("each move costs 0 (empty) or 1 (obstacle) → shortest path",
             "Dijkstra works: O(mn log mn);  BFS alone ignores the weights"),
        say="It's a shortest path where entering a cell costs its value, 0 or 1. Dijkstra works, "
            "but with only two weights a deque can replace the heap and drop the log factor.",
    )

    dist[0] = 0
    dq = Deque()
    dq.push_back(0)
    yield Frame(
        boards(dist, hl_d={(0, 0): "active"}),
        dq_arr(dq),
        Text("dist[(0,0)] = 0; push it"),
        say="A 0-cost move goes to the front of the deque, a 1-cost move to the back. "
            "The deque then holds at most two distances, in sorted order, just like Dijkstra's heap.",
    )
    pops = 0
    shown = 0
    while len(dq):
        u = dq.pop_front()
        pops += 1
        r, c = divmod(u, n)
        lines = [f"pop {cell(u)}  dist {dist[u]}"]
        hl_g = {(r, c): "active"}
        hl_d = {(r, c): "active"}
        for v, w in adj[u]:
            nd = dist[u] + w
            if nd < dist[v]:
                old = dist[v]
                dist[v] = nd
                if w == 0:
                    dq.push_front(v)
                    side = "front"
                else:
                    dq.push_back(v)
                    side = "back"
                lines.append((f"{cell(v)}: {dist[u]} + {w} = {nd} < {'∞' if old == INF else old} → push {side}", "match"))
                hl_g[(v // n, v % n)] = "match"
                hl_d[(v // n, v % n)] = "match"
        if len(lines) == 1:
            lines.append(("no neighbour improves", "reject"))
        shown += 1
        if shown == 1:
            say = "From the start both neighbours are obstacles, so both cost 1 and go to the back."
        elif any("front" in l[0] for l in lines[1:] if isinstance(l, tuple)):
            say = (f"Pop {cell(u)} at distance {dist[u]}. A move into an empty cell costs nothing, "
                   "so that cell goes to the FRONT and is processed next at the same distance.")
        elif len(lines) == 2 and lines[1][1] == "reject":
            say = f"Pop {cell(u)}. Every neighbour already has a distance at least as good, so nothing changes."
        else:
            say = (f"Pop {cell(u)} at distance {dist[u]}. Each obstacle neighbour costs one more removal, "
                   "so it goes to the back of the deque.")
        if shown <= 7:
            yield Frame(boards(dist, hl_g, hl_d), dq_arr(dq), Text(*lines[:4]), say=say)
        else:
            assert len(lines) == 2 and lines[1][1] == "reject"  # the summary says later pops change nothing

    ans = dist[N - 1]
    assert ans == 2
    yield Frame(
        boards(dist, hl_d={(m - 1, n - 1): "match"}),
        Text(f"{pops} pops in total; the later pops change nothing", (f"answer = dist[(2,2)] = {ans}", "match"),
             "O(m·n) time, O(m·n) space"),
        say=f"The remaining pops improve nothing, and the deque empties. The corner needs {ans} removals. "
            "Each cell enters the deque a constant number of times: O(mn) time and space.",
    )
