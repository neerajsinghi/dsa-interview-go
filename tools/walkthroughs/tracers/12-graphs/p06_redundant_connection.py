import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from uf import UnionFind  # noqa: E402
from walkthrough import Array, Frame, Graph, Row, Text  # noqa: E402

TITLE = "Redundant Connection"
PATTERN = "Union-find cycle check"

POS = {1: (0.5, 0.05), 2: (0.1, 0.95), 3: (0.9, 0.95)}


def frames():
    edges = [[1, 2], [1, 3], [2, 3]]
    expected = [2, 3]
    nodes = [1, 2, 3]
    uf = UnionFind(len(edges) + 1)  # nodes are 1-based

    def view(added, cur=None, style="active"):
        he = {tuple(e): "visited" for e in added}
        if cur:
            he[tuple(cur)] = style
        g = Graph(nodes, [tuple(e) for e in edges if e in added or e == cur], label="edges processed",
                  pos=POS, hl_edges=he, hl_nodes={v: style for v in (cur or [])}, h=160)
        forest = [(v, uf.parent[v]) for v in nodes if uf.parent[v] != v]
        f = Graph(nodes, forest, label="union-find forest (arrow = parent)", pos=POS, directed=True,
                  hl_nodes={v: "match" for v in nodes if uf.parent[v] == v}, h=160)
        return Row(g, f)

    def parent_row():
        return Array(uf.parent[1:], label="parent[1..3]", show_index=False)

    all_g = Graph(nodes, [tuple(e) for e in edges], label="edges = [[1,2],[1,3],[2,3]]", pos=POS, h=160)
    yield Frame(
        all_g,
        speaker="Interviewer",
        say="A tree on nodes 1 to n got one extra edge, so there's exactly one cycle. Return an edge to "
            "remove that makes it a tree again. If several work, return the last one in the input.",
    )
    yield Frame(
        all_g,
        Text("brute force: remove each edge from last to first,",
             "then BFS to check the rest is still connected",
             ("n removals × O(n) BFS = O(n²)", "reject")),
        say="I could try removing each edge, starting from the last, and check with a BFS that the graph "
            "is still connected. That's O(n²). Union-find finds the cycle in one pass.",
    )
    yield Frame(
        view([]),
        parent_row(),
        Text("every node is its own group"),
        say="I process the edges in order. An edge between two different groups merges them. An edge "
            "whose ends are already in the same group closes a cycle.",
    )

    added, answer = [], None
    for e in edges:
        a, b = e
        ra, rb = uf.find(a), uf.find(b)
        if not uf.union(a, b):
            answer = e
            yield Frame(
                view(added, e, "reject"),
                parent_row(),
                Text(f"edge [{a},{b}]: find({a}) = {ra}, find({b}) = {rb}",
                     (f"same root -> [{a},{b}] closes the cycle", "reject")),
                say=f"Edge [{a},{b}]: both ends already have root {ra}. They're connected through earlier edges, "
                    "so this edge closes the cycle. I return it.",
            )
            break
        added.append(e)
        yield Frame(
            view(added, e),
            parent_row(),
            Text(f"edge [{a},{b}]: find({a}) = {ra}, find({b}) = {rb}",
                 ("different roots -> union, keep the edge", "match")),
            say=f"Edge [{a},{b}]: roots {ra} and {rb} differ, so the edge joins two separate groups. "
                "It can't be on a cycle yet. I union them.",
        )

    assert answer == expected, answer
    yield Frame(
        view(added, answer, "reject"),
        Text((f"answer: {answer}", "match"),
             "every other cycle edge appeared earlier, so this is also the last one",
             "time: O(n · α(n))   space: O(n)"),
        say="The first edge that closes the cycle is also the last cycle edge in the input, which is "
            "what the problem asks for. One pass of near-constant unions: O(n·α(n)).",
    )
