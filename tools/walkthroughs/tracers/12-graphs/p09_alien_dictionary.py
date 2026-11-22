from walkthrough import Array, Frame, Graph, KV, Row, Text

TITLE = "Alien Dictionary"
PATTERN = "Graph + Kahn"

POS = {"w": (0.0, 0.15), "e": (0.25, 0.85), "r": (0.5, 0.15), "t": (0.75, 0.85), "f": (1.0, 0.15)}


def frames():
    words = ["wrt", "wrf", "er", "ett", "rftt"]
    expected = "wertf"

    present = [False] * 26
    for w in words:
        for ch in w:
            present[ord(ch) - 97] = True
    letters = [chr(97 + c) for c in range(26) if present[c]]
    edge = [[False] * 26 for _ in range(26)]
    indeg = [0] * 26
    facts = []

    def graph(hl=None, he=None, notes=None):
        es = [(chr(97 + u), chr(97 + v)) for u in range(26) for v in range(26) if edge[u][v]]
        return Graph(letters, es, label="letter graph (u -> v means u comes before v)", pos=POS, directed=True,
                     hl_nodes=hl or {}, hl_edges=he or {}, notes=notes or {}, h=150)

    def wordlist(i=None):
        hl = {words[i]: "active", words[i + 1]: "active"} if i is not None else {}
        return KV(words, label="dictionary (sorted in alien order)", hl=hl)

    yield Frame(
        wordlist(),
        speaker="Interviewer",
        say="These words are sorted by an unknown alien alphabet. Return any letter order consistent with "
            "them, or the empty string if none exists.",
    )
    yield Frame(
        wordlist(),
        Text(f"letters: {', '.join(letters)}",
             "brute force: try every ordering of the letters, check the list is sorted",
             (f"{len(letters)}! orderings, up to 26! in general", "reject")),
        say="I could try every permutation of the letters and check if the list is sorted under it. "
            "That's factorial. The list itself tells me which letters come first.",
    )

    for i in range(len(words) - 1):
        a, b = words[i], words[i + 1]
        j = 0
        while j < len(a) and j < len(b) and a[j] == b[j]:
            j += 1
        if j < len(a) and j < len(b):
            u, v = ord(a[j]) - 97, ord(b[j]) - 97
            if not edge[u][v]:
                edge[u][v] = True
                indeg[v] += 1
            facts.append(f"{a[j]} < {b[j]}")
            same = f"they share \"{a[:j]}\", then " if j else ""
            yield Frame(
                wordlist(i),
                graph({a[j]: "active", b[j]: "active"}, {(a[j], b[j]): "active"}),
                Text(f"{a} / {b}: first difference at index {j}", (f"fact: {a[j]} < {b[j]}", "match")),
                say=(f"Compare neighbours \"{a}\" and \"{b}\": {same}'{a[j]}' meets '{b[j]}'. Since \"{a}\" comes first, "
                     f"{a[j]} comes before {b[j]}." + (" Letters after that tell me nothing." if i == 0 else "")),
            )
        elif len(a) > len(b):
            raise AssertionError("example has no prefix conflict")

    queue = [c for c in range(26) if present[c] and indeg[c] == 0]
    notes = {chr(97 + c): f"in {indeg[c]}" for c in range(26) if present[c]}
    yield Frame(
        graph(notes=notes, hl={chr(97 + c): "match" for c in queue}),
        Row(Array([chr(97 + c) for c in queue], label="queue (in-degree 0)", show_index=False, cell=46),
            Text("facts: " + ", ".join(facts), "order = \"\""), weights=[1, 2.5]),
        say="Now it's a topological sort. I count incoming edges for each letter. Only "
            f"'{chr(97 + queue[0])}' has none, so it's the only letter that can come first.",
    )

    order = []
    while queue:
        u = queue.pop(0)
        order.append(chr(97 + u))
        freed = []
        for v in range(26):
            if edge[u][v]:
                indeg[v] -= 1
                if indeg[v] == 0:
                    queue.append(v)
                    freed.append(chr(97 + v))
        notes = {chr(97 + c): f"in {indeg[c]}" for c in range(26) if present[c]}
        hl = {ch: "visited" for ch in order}
        hl[chr(97 + u)] = "active"
        for ch in freed:
            hl[ch] = "match"
        if len(order) == 1:
            say = (f"I emit '{order[-1]}' and remove its edges. '{freed[0]}' drops to in-degree 0, so it joins the queue. "
                   "Each step repeats this.")
        elif freed:
            say = f"Emit '{order[-1]}'. That frees '{freed[0]}'."
        else:
            say = f"Emit '{order[-1]}'. It has no outgoing edges and the queue is empty, so Kahn stops."
        yield Frame(
            graph(hl, {(chr(97 + u), f): "compare" for f in freed}, notes),
            Row(Array([chr(97 + c) for c in queue] or ["-"], label="queue", show_index=False, cell=46),
                Text((f'order = "{"".join(order)}"', "match")), weights=[1, 2.5]),
            say=say,
        )

    result = "".join(order) if len(order) == len(letters) else ""
    assert result == expected, result
    yield Frame(
        graph({ch: "match" for ch in order}),
        Text((f'answer: "{result}"  ({len(order)} of {len(letters)} letters emitted -> no cycle)', "match"),
             "leftover letters would mean a cycle -> \"\"",
             "time: O(C) for C total characters; the graph has at most 26 nodes"),
        say="All five letters came out, so there's no cycle and the answer is \"wertf\". Building facts reads "
            "each character once, and the graph is at most 26 letters: O(C).",
    )
