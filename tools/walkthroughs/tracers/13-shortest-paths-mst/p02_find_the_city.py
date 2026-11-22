from walkthrough import Frame, Graph, Grid, Row, Text

TITLE = "Find the City (threshold distance)"
PATTERN = "Floyd-Warshall"

INF = float("inf")
POS = {0: (0.05, 0.12), 1: (0.45, 0.12), 2: (0.92, 0.12), 3: (0.92, 0.9)}


def frames():
    n, roads, threshold = 4, [[0, 1, 3], [1, 2, 1], [1, 3, 4], [2, 3, 1]], 4
    nodes = list(range(n))
    edges = [tuple(r) for r in roads]

    def graph(hl_nodes=None, hl_edges=None, notes=None):
        return Graph(nodes, edges, label="roads (undirected)", pos=POS, hl_nodes=hl_nodes,
                     hl_edges=hl_edges, notes=notes, h=190)

    def table(dist, hl=None, label="dist[i][j]"):
        return Grid(dist, label=label, hl=hl, row_labels=[f"i={i}" for i in nodes],
                    col_labels=[f"j={j}" for j in nodes], cell=44)

    yield Frame(
        graph(),
        Text(f"n = {n}, threshold = {threshold}"),
        speaker="Interviewer",
        say="For each city, count the other cities within distance 4. Return the city with the fewest; "
            "on a tie, return the larger label.",
    )

    # FindTheCity: undirected road = two directed edges, then Floyd-Warshall.
    dist = [[0 if i == j else INF for j in nodes] for i in nodes]
    for u, v, w in roads:
        dist[u][v] = min(dist[u][v], w)
        dist[v][u] = min(dist[v][u], w)

    yield Frame(
        Row(graph(), table(dist), weights=[1, 1]),
        Text("one BFS/DFS per city? weights differ, so BFS is wrong",
             "every pair needs a shortest distance → all-pairs"),
        say="I need the shortest distance between every pair, not from one source. "
            "Exploring every path from each city is exponential. With n this small, Floyd-Warshall gives all pairs in O(n cubed).",
    )
    yield Frame(
        Row(graph(), table(dist, hl={(i, j): "active" for i in nodes for j in nodes if dist[i][j] not in (0, INF)}),
            weights=[1, 1]),
        Text("dist[i][j] = road weight, 0 on the diagonal, ∞ otherwise"),
        say="I start the table with the direct roads in both directions, 0 on the diagonal and infinity elsewhere.",
    )

    for k in nodes:
        changed = {}
        lines = []
        for i in nodes:
            if dist[i][k] == INF:
                continue
            for j in nodes:
                if dist[k][j] != INF and dist[i][k] + dist[k][j] < dist[i][j]:
                    old = dist[i][j]
                    dist[i][j] = dist[i][k] + dist[k][j]
                    changed[(i, j)] = "match"
                    if i < j:
                        lines.append((f"dist[{i}][{j}]: {old if old != INF else '∞'} → {dist[i][k]} + {dist[k][j]} = {dist[i][j]}", "match"))
        hl = dict(changed)
        for t in nodes:
            hl.setdefault((k, t), "special")
            hl.setdefault((t, k), "special")
        if not lines:
            lines = [("no pair gets shorter through this city", "reject")]
        if k == 0 and not changed:
            say = ("Round k = 0: may any trip stop at city 0? City 0 is a dead end with one road, so nothing improves.")
        elif changed:
            say = (f"Round k = {k}: for every pair I ask whether going from i to {k} and then to j beats the current best. "
                   f"{len(changed) // 2} pair(s) get shorter.")
        else:
            say = f"Round k = {k}: no pair gets shorter by passing through city {k}."
        yield Frame(
            Row(graph(hl_nodes={k: "special"}), table(dist, hl=hl, label=f"after k = {k}"), weights=[1, 1]),
            Text(f"dist[i][j] = min(dist[i][j], dist[i][{k}] + dist[{k}][j])", *lines[:3]),
            say=say,
        )

    # pickCity
    best, best_count = -1, n + 1
    counts = {}
    for i in nodes:
        count = sum(1 for j in nodes if i != j and dist[i][j] <= threshold)
        counts[i] = count
        if count <= best_count:
            best, best_count = i, count
    within = {(i, j): ("match" if dist[i][j] <= threshold else "reject") for i in nodes for j in nodes if i != j}
    yield Frame(
        Row(graph(notes={i: f"{counts[i]} near" for i in nodes}), table(dist, hl=within, label="≤ 4 green, > 4 red"),
            weights=[1, 1]),
        Text(*(f"city {i}: {counts[i]} reachable" for i in nodes)),
        say="Now I count, for each row, the cells at most 4 away. Cities 0 and 3 both reach only two others.",
    )
    assert best == 3
    yield Frame(
        Row(graph(hl_nodes={best: "match"}, notes={i: f"{counts[i]} near" for i in nodes}),
            table(dist, hl={(best, j): "match" for j in nodes}), weights=[1, 1]),
        Text(f"tie between 0 and 3 → larger label", (f"answer = {best}", "match"), "O(n³) time, O(n²) space"),
        say="I scan cities in order and update on count <= best, so the later city wins the tie. "
            "The answer is 3. Floyd-Warshall is O(n cubed) time and O(n squared) space.",
    )
