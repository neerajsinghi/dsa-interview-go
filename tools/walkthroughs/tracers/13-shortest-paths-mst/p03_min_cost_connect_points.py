from walkthrough import Array, Frame, Graph, Text

TITLE = "Min Cost to Connect All Points"
PATTERN = "Prim O(n²)"

INF = float("inf")


def manhattan(a, b):
    return abs(a[0] - b[0]) + abs(a[1] - b[1])


def frames():
    points = [[0, 0], [2, 2], [3, 10], [5, 2], [7, 0]]
    n = len(points)
    nodes = list(range(n))
    # Lay the points out at their real coordinates (y grows upwards).
    pos = {i: (x / 7, 1 - y / 10) for i, (x, y) in enumerate(points)}
    names = [f"{i}:({x},{y})" for i, (x, y) in enumerate(points)]

    def graph(edges, hl_nodes=None, hl_edges=None, notes=None):
        return Graph(nodes, edges, label="points at their coordinates", pos=pos, hl_nodes=hl_nodes,
                     hl_edges=hl_edges, notes=notes, h=210)

    def md_arr(min_dist, in_tree, pick=None):
        hl = {v: "visited" for v in nodes if in_tree[v]}
        if pick is not None:
            hl[pick] = "match"
        return Array(min_dist, label="minDist[v]: cheapest edge from the tree to v", hl=hl, cell=60)

    yield Frame(
        graph([]),
        Text("points: " + "  ".join(names)),
        speaker="Interviewer",
        say="Connecting two points costs their Manhattan distance. "
            "What is the cheapest way to connect all of them so every point can reach every other?",
    )
    yield Frame(
        graph([]),
        Text(f"every pair is a candidate edge: {n * (n - 1) // 2} edges, O(n²) in general",
             "Kruskal: build + sort all edges → O(n² log n), O(n²) memory"),
        say="This is a minimum spanning tree on a complete graph. Kruskal would list and sort all n squared edges. "
            "With that many edges, array-based Prim is O(n squared) and needs no edge list.",
    )

    in_tree = [False] * n
    min_dist = [INF] * n
    parent = [None] * n  # only for drawing the chosen edge
    min_dist[0] = 0
    total = 0
    tree_edges = []
    yield Frame(
        graph([], hl_nodes={0: "active"}),
        md_arr(min_dist, in_tree),
        Text("minDist[0] = 0, the rest ∞"),
        say="I grow one tree from point 0. minDist[v] is the cheapest known edge from the tree to v, so point 0 starts at 0.",
    )
    for rnd in range(n):
        u = -1
        for v in nodes:
            if not in_tree[v] and (u == -1 or min_dist[v] < min_dist[u]):
                u = v
        in_tree[u] = True
        total += min_dist[u]
        if parent[u] is not None:
            tree_edges.append((parent[u], u, min_dist[u]))
        updates = []
        for v in nodes:
            if not in_tree[v]:
                d = manhattan(points[u], points[v])
                if d < min_dist[v]:
                    updates.append(f"{v}: {min_dist[v] if min_dist[v] != INF else '∞'}→{d}")
                    min_dist[v] = d
                    parent[v] = u
        hl_n = {v: "visited" for v in nodes if in_tree[v]}
        hl_n[u] = "match"
        hl_e = {(tree_edges[-1][0], tree_edges[-1][1]): "match"} if parent[u] is not None else {}
        added = f"add point {u} (cost {tree_edges[-1][2]} via {tree_edges[-1][0]})" if parent[u] is not None \
            else f"add point {u} (cost 0, the start)"
        if rnd == 0:
            say = "Round 1 picks point 0 for free. Then I compute its distance to every other point to fill minDist."
        elif rnd == n - 1:
            say = f"The last point, {u}, joins with an edge of {min_dist[u]}. Every point is now in the tree."
        else:
            say = (f"The closest point outside the tree is {u}, at distance {min_dist[u]}. By the cut property that edge is safe. "
                   f"Then {u} may offer cheaper edges to the rest.")
        yield Frame(
            graph(list(tree_edges), hl_nodes=hl_n, hl_edges=hl_e,
                  notes={v: ('∞' if min_dist[v] == INF else min_dist[v]) for v in nodes if not in_tree[v]}),
            md_arr(min_dist, in_tree, u),
            Text(added, ("updates: " + (", ".join(updates) if updates else "none"), "active" if updates else None),
                 f"total = {total}"),
            say=say,
        )

    assert total == 20
    yield Frame(
        graph(list(tree_edges), hl_nodes={v: "match" for v in nodes},
              hl_edges={(a, b): "match" for a, b, _ in tree_edges}),
        Text(" + ".join(str(w) for _, _, w in tree_edges) + f" = {total}", (f"answer = {total}", "match"),
             "Prim (array): O(n²) time, O(n) space"),
        say=f"The tree uses {n - 1} edges for a total of {total}. "
            "n rounds, each scanning n points, is O(n squared) time and O(n) extra space.",
    )
