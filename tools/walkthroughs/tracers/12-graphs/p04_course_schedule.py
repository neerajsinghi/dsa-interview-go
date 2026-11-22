from walkthrough import Array, Frame, Graph, Row, Text

TITLE = "Course Schedule I & II"
PATTERN = "Topological sort (Kahn)"

POS = {0: (0.05, 0.5), 1: (0.5, 0.05), 2: (0.5, 0.95), 3: (0.95, 0.5)}


def prereq_graph(n, prereqs):
    adj = [[] for _ in range(n)]
    for a, b in prereqs:
        adj[b].append(a)
    return adj


def kahn(adj, on_step=None):
    n = len(adj)
    indeg = [0] * n
    for u in range(n):
        for v in adj[u]:
            indeg[v] += 1
    queue = [v for v in range(n) if indeg[v] == 0]
    order = []
    if on_step:
        on_step("init", None, list(indeg), list(queue), list(order), [])
    while queue:
        u = queue.pop(0)
        order.append(u)
        freed = []
        for v in adj[u]:
            indeg[v] -= 1
            if indeg[v] == 0:
                queue.append(v)
                freed.append(v)
        if on_step:
            on_step("pop", u, list(indeg), list(queue), list(order), freed)
    return order, len(order) == n


def frames():
    n, prereqs = 4, [[1, 0], [2, 0], [3, 1], [3, 2]]
    adj = prereq_graph(n, prereqs)
    edges = [(b, a) for a, b in prereqs]

    def graph(indeg, done, cur=None, freed=(), removed=()):
        hl = {v: "visited" for v in done}
        for v in freed:
            hl[v] = "match"
        if cur is not None:
            hl[cur] = "active"
        he = {e: "compare" for e in removed}
        return Graph(list(range(n)), edges, label="b -> a means take b before a", pos=POS, directed=True,
                     hl_nodes=hl, hl_edges=he, notes={v: f"in {indeg[v]}" for v in range(n)}, h=180)

    yield Frame(
        Graph(list(range(n)), edges, label="b -> a means take b before a", pos=POS, directed=True, h=180),
        Text("prerequisites = [[1,0],[2,0],[3,1],[3,2]]   ([a, b]: take b before a)"),
        speaker="Interviewer",
        say="There are 4 courses. Each pair [a, b] means take b before a. Can you finish them all, "
            "and if so, in what order?",
    )
    yield Frame(
        Graph(list(range(n)), edges, label="b -> a means take b before a", pos=POS, directed=True, h=180),
        Text("brute force: try orderings and check every prerequisite",
             ("n! orderings", "reject")),
        say="Trying every ordering and checking each pair is O(n!·E). I'll draw an edge b → a for each pair. "
            "A valid schedule is then a topological order of this graph.",
    )

    steps = []
    order, ok = kahn(adj, lambda *a: steps.append(a))
    for kind, u, indeg, queue, done, freed in steps:
        if kind == "init":
            yield Frame(
                Row(graph(indeg, []), Array(queue, label="queue (in-degree 0)", show_index=False, cell=46),
                    weights=[1.6, 1]),
                Text("in-degree = number of unmet prerequisites", "order = []"),
                say="Kahn's algorithm counts incoming edges per course. A course with in-degree 0 has no "
                    f"unmet prerequisites, so it can be taken now. Only course {queue[0]} qualifies.",
            )
            continue
        removed = [(u, v) for v in adj[u]]
        if freed:
            tail = (f"Courses {', '.join(map(str, freed))} drop to 0 and join the queue." if len(freed) > 1
                    else f"Course {freed[0]} drops to 0 and joins the queue.")
        elif adj[u]:
            tail = f"Course {adj[u][0]} drops to {indeg[adj[u][0]]}, still waiting on another prerequisite."
        else:
            tail = "It unlocks nothing. The queue is empty now."
        yield Frame(
            Row(graph(indeg, done[:-1], u, freed, removed),
                Array(queue or ["-"], label="queue", show_index=False, cell=46), weights=[1.6, 1]),
            Text(f"take course {u}; lower the in-degree of {adj[u] or 'nobody'}",
                 (f"order = {done}", "match")),
            say=(f"I take course {u} and remove its outgoing edges, which lowers the in-degree of its dependants. "
                 if not done[:-1] else f"Next I take course {u} from the queue. ") + tail,
        )

    assert ok and order == [0, 1, 2, 3], order
    # Part I's "false" case: the second example, a two-course cycle.
    cyc_adj = prereq_graph(2, [[1, 0], [0, 1]])
    cyc_in = [sum(v in cyc_adj[u] for u in range(2)) for v in range(2)]
    cyc_order, cyc_ok = kahn(cyc_adj)
    assert not cyc_ok and cyc_order == []
    yield Frame(
        Graph([0, 1], [(0, 1), (1, 0)], label="second example: [[1,0],[0,1]]",
              pos={0: (0.3, 0.5), 1: (0.7, 0.5)}, directed=True, hl_nodes={0: "reject", 1: "reject"},
              notes={v: f"in {cyc_in[v]}" for v in range(2)}, h=110),
        Text("no course has in-degree 0 -> queue starts empty",
             (f"order = {cyc_order}, {len(cyc_order)} of 2 courses -> CanFinish = false, FindOrder = []", "reject")),
        say="With a cycle, the courses wait on each other, so no in-degree ever reaches 0. "
            "Kahn emits fewer than n courses, which tells me it's impossible.",
    )
    yield Frame(
        graph(steps[-1][2], list(range(n))),
        Text((f"FindOrder = {order}, CanFinish = true", "match"),
             "emitted all n courses <=> no cycle",
             "time: O(V + E)   space: O(V + E)"),
        say="All 4 courses came out, so there's no cycle and the order is [0 1 2 3]. Each node and edge "
            "is handled once: O(V + E).",
    )
