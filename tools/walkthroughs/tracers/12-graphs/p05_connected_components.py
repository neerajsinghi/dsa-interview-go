import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from uf import UnionFind  # noqa: E402
from walkthrough import Array, Frame, Graph, Row, Text  # noqa: E402

TITLE = "Connected Components"
PATTERN = "Union-find"

POS = {0: (0.0, 0.15), 1: (0.3, 0.15), 2: (0.15, 0.9), 3: (0.75, 0.15), 4: (1.0, 0.9)}


def frames():
    n, edges = 5, [[0, 1], [1, 2], [3, 4]]
    expected = 2
    uf = UnionFind(n)
    nodes = list(range(n))

    def view(added, cur=None, roots_hl=None):
        he = {tuple(cur): "active"} if cur else {}
        g = Graph(nodes, [tuple(e) for e in added], label="graph (edges added so far)", pos=POS,
                  hl_edges=he, hl_nodes={v: "active" for v in (cur or [])}, h=150)
        forest = [(v, uf.parent[v]) for v in nodes if uf.parent[v] != v]
        roots = {v: "match" for v in nodes if uf.parent[v] == v}
        roots.update(roots_hl or {})
        f = Graph(nodes, forest, label="union-find forest (arrow = parent)", pos=POS, directed=True,
                  hl_nodes=roots, notes={v: f"r{uf.rank[v]}" for v in nodes if uf.parent[v] == v and uf.rank[v]},
                  h=150)
        return Row(g, f)

    def parent_row():
        return Array(uf.parent, label="parent[]", hl={v: "match" for v in nodes if uf.parent[v] == v})

    yield Frame(
        Graph(nodes, [tuple(e) for e in edges], label="n = 5, edges = [[0,1],[1,2],[3,4]]", pos=POS, h=150),
        speaker="Interviewer",
        say="Given n nodes and a list of undirected edges, count the connected components.",
    )
    yield Frame(
        Graph(nodes, [tuple(e) for e in edges], label="n = 5, edges = [[0,1],[1,2],[3,4]]", pos=POS, h=150),
        Text("brute force: BFS from every node, collect its reachable set, dedupe the sets",
             ("O(V · (V + E))", "reject")),
        say="A naive approach runs a search from every node and dedupes the reachable sets. That's "
            "O(V·(V+E)). One DFS per unvisited node is already O(V+E). Union-find also handles edges one at a time.",
    )
    yield Frame(
        view([]),
        parent_row(),
        Text((f"count = {uf.count}", "match"), "every node starts as its own group (its own root)"),
        say="Union-find starts with every node in its own group, so the count is 5. "
            "Each edge that joins two different groups lowers the count by one.",
    )

    added = []
    for a, b in edges:
        ra, rb = uf.find(a), uf.find(b)
        same_rank = uf.rank[ra] == uf.rank[rb]
        merged = uf.union(a, b)
        added.append([a, b])
        new_root = uf.find(a)
        child = rb if new_root == ra else ra
        if merged:
            say = (f"Edge {a}-{b}: find({a}) = {ra} and find({b}) = {rb} are different roots, so I merge. "
                   + (f"Ranks are equal, so {child} hangs under {new_root} and {new_root}'s rank grows."
                      if same_rank else
                      f"Root {child} has the lower rank, so it hangs under {new_root}."))
        else:
            say = f"Edge {a}-{b}: both ends already share root {ra}, so nothing changes."
        yield Frame(
            view(added, [a, b]),
            parent_row(),
            Text(f"union({a}, {b}): find({a}) = {ra}, find({b}) = {rb}",
                 (f"merged -> count = {uf.count}" if merged else "same group -> count unchanged", "match")),
            say=say,
        )

    assert uf.count == expected, uf.count
    yield Frame(
        view(added),
        parent_row(),
        Text((f"answer: {uf.count} components", "match"),
             "time: O(E · α(n)) ≈ O(E), plus O(n) to initialise",
             "space: O(n) for parent and rank"),
        say=f"Roots {' and '.join(str(v) for v in nodes if uf.parent[v] == v)} remain, so there are {uf.count} components. With path compression and union "
            "by rank each operation is nearly constant, so the total is about O(n + E).",
    )
