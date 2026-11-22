from walkthrough import Frame, Graph, KV, Row, Text, stack_widget

TITLE = "Clone Graph"
PATTERN = "DFS + map"

POS = {1: (0.15, 0.1), 2: (0.85, 0.1), 3: (0.85, 0.9), 4: (0.15, 0.9)}


class GNode:
    def __init__(self, val):
        self.val = val
        self.neighbors = []


def frames():
    # 1 -- 2, 2 -- 3, 3 -- 4, 4 -- 1 (adjacency in the usual LeetCode order)
    adj = {1: [2, 4], 2: [1, 3], 3: [2, 4], 4: [1, 3]}
    orig = {v: GNode(v) for v in adj}
    for v, nbs in adj.items():
        orig[v].neighbors = [orig[u] for u in nbs]
    und_edges = [(1, 2), (2, 3), (3, 4), (4, 1)]

    copies = {}
    calls = []

    def view(hl_orig=None, hl_clone=None, edge=None):
        vals = sorted(o.val for o in copies)
        cn = [f"{v}'" for v in vals]
        ce = []
        for o, c in copies.items():
            for nb in c.neighbors:
                ce.append((f"{o.val}'", f"{nb.val}'"))
        g1 = Graph([1, 2, 3, 4], und_edges, label="original", pos=POS, hl_nodes=hl_orig or {},
                   hl_edges={edge: "compare"} if edge else {}, h=170)
        g2 = Graph(cn, ce, label="clone (arrows = neighbour lists filled so far)",
                   pos={f"{v}'": POS[v] for v in vals}, hl_nodes=hl_clone or {}, directed=True, h=170)
        return Row(g1, g2)

    def side():
        return Row(KV({o.val: f"{o.val}'" for o in copies}, label="copies map", empty="{ }"),
                   stack_widget([f"clone({c})" for c in calls], label="recursion"), weights=[1, 1.4])

    yield Frame(
        Row(Graph([1, 2, 3, 4], und_edges, label="original", pos=POS, h=170),
            Graph([], [], label="clone", h=170)),
        speaker="Interviewer",
        say="You get node 1 of a connected undirected graph. Return a deep copy: new nodes with the "
            "same values and the same connections, sharing no pointers.",
    )
    yield Frame(
        Row(Graph([1, 2, 3, 4], und_edges, label="original", pos=POS, h=170),
            Graph([], [], label="clone", h=170)),
        Text("naive: copy a node, then copy each neighbour recursively",
             ("1 -> 2 -> 1 -> 2 -> ... the cycle never ends", "reject")),
        say="Copying a node and then blindly copying its neighbours loops forever, because the graph "
            "has a cycle. I need to remember which nodes already have a copy.",
    )

    out = []
    hits = [0]

    def clone(n, frm):
        calls.append(n.val)
        if n in copies:
            hits[0] += 1
            if hits[0] == 1:
                say = (f"Node {n.val} is already in the map, so I return {n.val}' instead of recursing. "
                       "This lookup is what stops the cycle.")
            elif len(copies) < 4:
                say = f"{n.val} is in the map again, so I return {n.val}' and append it to the list of {frm}'."
            else:
                say = (f"All four copies exist now, so every remaining neighbour is a map hit. "
                       f"Here {n.val}' goes into the neighbour list of {frm}'.")
            out.append(Frame(
                view({n.val: "match", frm: "active"}, {f"{n.val}'": "match", f"{frm}'": "active"}, (frm, n.val)),
                side(),
                Text((f"{n.val} is already in the map -> return {n.val}'", "match")),
                say=say,
            ))
            calls.pop()
            return copies[n]
        c = GNode(n.val)
        copies[n] = c  # register first: neighbours may point back at n
        out.append(Frame(
            view({n.val: "active"} | ({frm: "visited"} if frm else {}), {f"{n.val}'": "active"},
                 (frm, n.val) if frm else None),
            side(),
            Text(f"create {n.val}', store copies[{n.val}] = {n.val}'",
                 f"then clone its neighbours {[nb.val for nb in n.neighbors]}"),
            say=(f"Node {n.val} has no copy yet. I create {n.val}' and put it in the map before touching "
                 "any neighbour, so a cycle that leads back here finds it."
                 if len(copies) == 1 else
                 f"Node {n.val} is new too: create {n.val}', register it, then recurse into its neighbours."),
        ))
        for nb in n.neighbors:
            c.neighbors.append(clone(nb, n.val))
        calls.pop()
        return c

    root = clone(orig[1], None)
    yield from out

    # Check the clone: same shape, no shared pointers.
    seen, stack = set(), [root]
    while stack:
        c = stack.pop()
        if id(c) in seen:
            continue
        seen.add(id(c))
        assert all(c is not o for o in orig.values())
        assert [nb.val for nb in c.neighbors] == adj[c.val]
        stack.extend(c.neighbors)
    assert len(seen) == 4

    yield Frame(
        view(),
        Text(("clone(1) returns 1': a new 4-cycle 1'-2'-3'-4'", "match"),
             "time: O(V + E), each node created once, each edge copied once",
             "space: O(V) for the map and the recursion"),
        say="Every node is copied once and every neighbour entry once, so it's O(V + E). "
            "The map holds V entries, and it doubles as the visited set.",
    )
