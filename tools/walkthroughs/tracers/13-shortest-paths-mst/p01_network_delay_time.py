import heapq

from walkthrough import Array, Frame, Graph, Row, Text

TITLE = "Network Delay Time"
PATTERN = "Dijkstra"

INF = float("inf")
POS = {1: (0.06, 0.85), 2: (0.3, 0.12), 3: (0.56, 0.85), 4: (0.94, 0.85)}


def frames():
    times, n, k = [[2, 1, 1], [2, 3, 1], [3, 4, 1]], 4, 2
    nodes = list(range(1, n + 1))
    edges = [tuple(t) for t in times]
    adj = {u: [] for u in nodes}
    for u, v, w in times:
        adj[u].append((v, w))

    def notes(dist):
        return {v: f"d={'∞' if dist[v] == INF else dist[v]}" for v in nodes}

    def heap_arr(h):
        items = sorted(h)
        return Array([f"{u}:{d}" for d, u in items] or ["—"], label="min-heap (node:dist)",
                     show_index=False, cell=64)

    def graph(dist, hl_nodes=None, hl_edges=None):
        return Graph(nodes, edges, label="links (directed, weight = delay)", pos=POS, directed=True,
                     notes=notes(dist), hl_nodes=hl_nodes, hl_edges=hl_edges, h=200)

    dist0 = {v: INF for v in nodes}
    yield Frame(
        graph(dist0, hl_nodes={k: "special"}),
        Text(f"n = {n}, start k = {k}"),
        speaker="Interviewer",
        say="A signal starts at node 2 and travels along these directed links. "
            "How long until every node has heard it? Return -1 if some node never does.",
    )
    yield Frame(
        graph(dist0, hl_nodes={k: "special"}),
        Text("brute force: DFS every path, keep the best time per node", "worst case: exponential number of paths"),
        say="Each node hears the signal at its shortest-path time, and the answer is the largest of those. "
            "Trying every path is exponential, but weights are non-negative, so Dijkstra fits.",
    )

    dist = {v: INF for v in nodes}
    dist[k] = 0
    h = [(0, k)]
    settled = set()
    yield Frame(
        graph(dist, hl_nodes={k: "active"}),
        heap_arr(h),
        Text(f"dist[{k}] = 0, everything else ∞"),
        say="I set the start's distance to 0 and push it on a min-heap. "
            "The heap always hands me the closest node I haven't settled yet.",
    )
    while h:
        d, u = heapq.heappop(h)
        if d > dist[u]:
            yield Frame(graph(dist), heap_arr(h), Text(f"pop {u}:{d}, but dist[{u}] = {dist[u]}", ("stale, skip", "reject")),
                        say=f"This entry for {u} is stale, a shorter path was already found. Skip it.")
            continue
        settled.add(u)
        lines = [f"pop {u}:{d}  → dist[{u}] = {d} is final"]
        hl_e = {}
        for v, w in adj[u]:
            nd = d + w
            if nd < dist[v]:
                lines.append((f"relax {u}→{v}: {d} + {w} = {nd} < {'∞' if dist[v] == INF else dist[v]}  → push {v}:{nd}", "match"))
                dist[v] = nd
                heapq.heappush(h, (nd, v))
                hl_e[(u, v)] = "match"
            else:
                lines.append((f"{u}→{v}: {d} + {w} = {nd} ≥ {dist[v]}, no change", "reject"))
                hl_e[(u, v)] = "reject"
        if not adj[u]:
            lines.append("no outgoing links")
        hl_n = {x: "visited" for x in settled}
        hl_n[u] = "active"
        if u == k:
            say = "Pop node 2 at time 0. Both of its links improve their targets from infinity to 1, so I push nodes 1 and 3."
        elif adj[u]:
            say = (f"Pop node {u} at time {d}. Nothing can beat it later, because any other route passes through a node "
                   f"that is at least as far. Relax its links.")
        else:
            say = f"Pop node {u} at time {d}. It's final, and it has no outgoing links, so nothing else changes."
        yield Frame(graph(dist, hl_nodes=hl_n, hl_edges=hl_e), heap_arr(h), Text(*lines), say=say)

    reached = [dist[v] for v in nodes]
    ans = -1 if INF in reached else max(reached)
    assert ans == 2
    slow = max(nodes, key=lambda v: dist[v])
    yield Frame(
        graph(dist, hl_nodes={v: ("match" if v == slow else "visited") for v in nodes}),
        Text(f"dist = {[dist[v] for v in nodes]}  (nodes 1..4)", (f"answer = max = {ans}", "match"),
             "Dijkstra O(E log V) time, O(V + E) space"),
        say=f"Every node is reachable, so the answer is the largest distance: {ans}, at node {slow}. "
            "Dijkstra with a binary heap is O(E log V) time and O(V + E) space.",
    )
