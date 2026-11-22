from walkthrough import Array, Frame, Graph, Row, Text

TITLE = "Cheapest Flights Within K Stops"
PATTERN = "Bellman-Ford, k+1 rounds"

INF = float("inf")
POS = {0: (0.04, 0.5), 1: (0.4, 0.06), 2: (0.4, 0.94), 3: (0.96, 0.5)}


def fmt(x):
    return "∞" if x == INF else x


def frames():
    n = 4
    flights = [[0, 1, 100], [1, 2, 100], [2, 0, 100], [1, 3, 600], [2, 3, 200]]
    src, dst, k = 0, 3, 1
    nodes = list(range(n))
    edges = [tuple(f) for f in flights]

    def graph(dist=None, hl_nodes=None, hl_edges=None):
        notes = {v: fmt(dist[v]) for v in nodes} if dist else None
        return Graph(nodes, edges, label="flights (directed, price)", pos=POS, directed=True,
                     notes=notes, hl_nodes=hl_nodes, hl_edges=hl_edges, h=190)

    def arr(vals, label, hl=None):
        return Array([fmt(x) for x in vals], label=label, hl=hl, cell=58)

    yield Frame(
        graph(hl_nodes={src: "special", dst: "special"}),
        Text(f"src = {src}, dst = {dst}, k = {k} stop (so at most {k + 1} flights)"),
        speaker="Interviewer",
        say="Find the cheapest price from city 0 to city 3 with at most one stop in between. Return -1 if there's no such route.",
    )
    yield Frame(
        graph(hl_nodes={src: "special", dst: "special"}),
        Text("brute force: DFS every route of ≤ k+1 flights → O(outdegree^(k+1))",
             "plain Dijkstra on price: wrong, cheapest route to a city may use too many stops"),
        say="Trying every route is exponential in k. Plain Dijkstra on price is wrong: the cheapest way into city 2 "
            "uses a stop, and then it can't continue. Bellman-Ford counts edges per round.",
    )

    # The trap: relax in place, one round. Run it to show what goes wrong.
    bad = [INF] * n
    bad[src] = 0
    for f, t, p in flights:
        if bad[f] != INF and bad[f] + p < bad[t]:
            bad[t] = bad[f] + p
    yield Frame(
        graph(bad, hl_edges={(0, 1): "reject", (1, 2): "reject", (2, 3): "reject"}),
        arr(bad, "dist after ONE in-place round (wrong)", hl={dst: "reject"}),
        Text(f"in place, round 1 already gives dist[3] = {bad[3]}", ("that uses 3 flights: 0→1→2→3", "reject")),
        say="One trap first. If I relax in place, one round can chain several flights. "
            f"Here round 1 already reaches city 3 for {bad[3]}, which uses three flights. So I relax from a snapshot.",
    )

    dist = [INF] * n
    dist[src] = 0
    yield Frame(
        graph(dist, hl_nodes={src: "active"}),
        arr(dist, "dist", hl={src: "active"}),
        Text("dist[src] = 0;  run k + 1 = 2 rounds"),
        say="I start with dist[0] = 0 and everything else infinite. After round i, dist holds the cheapest price using at most i flights.",
    )
    for rnd in range(k + 1):
        prev = list(dist)
        lines, hl_e = [], {}
        for f, t, p in flights:
            if prev[f] != INF and prev[f] + p < dist[t]:
                lines.append((f"{f}→{t}: prev[{f}] + {p} = {prev[f] + p} < {fmt(dist[t])}  → dist[{t}] = {prev[f] + p}", "match"))
                dist[t] = prev[f] + p
                hl_e[(f, t)] = "match"
            elif prev[f] == INF:
                hl_e.setdefault((f, t), None)
            else:
                lines.append((f"{f}→{t}: prev[{f}] + {p} = {prev[f] + p} ≥ {fmt(dist[t])}, no change", "reject"))
                hl_e[(f, t)] = "reject"
        hl_e = {e: s for e, s in hl_e.items() if s}
        changed = {i: "match" for i in nodes if dist[i] != prev[i]}
        skipped = [f"{f}→{t}" for f, t, _ in flights if prev[f] == INF]
        if skipped:
            lines.append(("skip " + ", ".join(skipped) + " (prev[from] = ∞)", None))
        if rnd == 0:
            say = ("Round 1 reads only from the snapshot, so it can extend paths by one flight. "
                   "Only flights leaving city 0 count, so city 1 becomes 100.")
        else:
            say = ("Round 2 extends the round-1 paths by one more flight. "
                   f"City 3 becomes {dist[3]} via route 0-1-3, and city 2 becomes {dist[2]}. The 2 to 3 flight is not used yet: city 2 wasn't reached in the snapshot.")
        yield Frame(
            graph(dist, hl_nodes=changed, hl_edges=hl_e),
            Row(arr(prev, f"prev (snapshot before round {rnd + 1})"), arr(dist, f"dist after round {rnd + 1}", hl=changed)),
            Text(*lines),
            say=say,
        )

    ans = -1 if dist[dst] == INF else dist[dst]
    assert ans == 700
    yield Frame(
        graph(dist, hl_nodes={dst: "match"}, hl_edges={(0, 1): "match", (1, 3): "match"}),
        Text(f"after k + 1 = {k + 1} rounds: dist[{dst}] = {dist[dst]}", (f"answer = {ans}  (0→1→3)", "match"),
             "O(k·E) time, O(n) space"),
        say=f"I stop after k + 1 rounds, so dist[3] = {ans} only counts routes with at most two flights. "
            "Each round scans every flight: O(k times E) time and O(n) space.",
    )
