from walkthrough import Array, Frame, Graph, KV, Row, Text

TITLE = "Word Ladder"
PATTERN = "BFS on implicit graph"

POS = {"hit": (0.0, 0.5), "hot": (0.18, 0.5), "dot": (0.42, 0.05), "lot": (0.42, 0.95),
       "dog": (0.7, 0.05), "log": (0.7, 0.95), "cog": (0.95, 0.5)}


def one_apart(a, b):
    return sum(x != y for x, y in zip(a, b)) == 1


def frames():
    begin, end = "hit", "cog"
    word_list = ["hot", "dot", "dog", "lot", "log", "cog"]
    expected = 5
    nodes = [begin] + word_list
    edges = [(a, b) for i, a in enumerate(nodes) for b in nodes[i + 1:] if one_apart(a, b)]

    def graph(hl, he=None):
        return Graph(nodes, edges, label="implicit graph: words one letter apart", pos=POS,
                     hl_nodes=hl, hl_edges=he or {}, h=170)

    yield Frame(
        graph({begin: "active", end: "special"}),
        Text(f'begin = "{begin}", end = "{end}"', f"wordList = {word_list}".replace("'", "")),
        speaker="Interviewer",
        say="Change hit into cog one letter at a time. Every word on the way must be in the list. "
            "Return the number of words in the shortest sequence, or 0.",
    )
    yield Frame(
        graph({begin: "active", end: "special"}),
        Text("brute force: compare every pair of words to build the graph: O(N² · L)",
             ("or DFS over all paths: exponential", "reject")),
        say="Words are nodes and one-letter changes are edges. Building that graph by comparing all "
            "pairs is O(N²·L), and a DFS over paths is exponential. Fewest steps means BFS.",
    )

    # Port of LadderLength.
    words = set(word_list)
    assert end in words
    words.discard(begin)
    queue = [begin]
    visited = {begin: 0}
    parent = {}
    steps = 1
    answer = 0
    yield Frame(
        Row(graph({begin: "active"}), weights=[1]),
        Row(Array(queue, label=f"queue (level {steps})", show_index=False, cell=52),
            KV(sorted(words), label="unvisited word set"), weights=[1, 2.4]),
        say="I start the BFS at hit with steps = 1. I won't build the edges up front. For each word I try all "
            "25 other letters at each position and keep candidates that are still in the set.",
    )
    while queue:
        nxt = []
        tried = 0
        for w in queue:
            if w == end:
                answer = steps
                break
            b = list(w)
            for i in range(len(b)):
                orig = b[i]
                for ch in "abcdefghijklmnopqrstuvwxyz":
                    if ch == orig:
                        continue
                    tried += 1
                    b[i] = ch
                    cand = "".join(b)
                    if cand in words:
                        words.discard(cand)  # deleting = marking visited
                        nxt.append(cand)
                        visited[cand] = steps
                        parent[cand] = w
                b[i] = orig
        if answer:
            break
        hl = {v: "visited" for v in visited}
        for w in queue:
            hl[w] = "active"
        for w in nxt:
            hl[w] = "match"
        he = {(a, c): "match" for a in queue for c in nxt if one_apart(a, c)}
        found = ", ".join(nxt)
        if steps == 1:
            say = (f"Level {steps}: from hit I generate {tried} candidates. Only hot is in the set. I delete it "
                   "from the set right away, which marks it visited.")
        elif end in nxt:
            say = f"Level {steps}: {', '.join(queue)} reach {found}. The end word is in the next level."
        else:
            say = f"Level {steps}: expanding {', '.join(queue)} finds {found}. Each is removed from the set as it is queued."
        steps += 1
        queue = nxt
        yield Frame(
            graph(hl, he),
            Row(Array(queue, label=f"queue (level {steps})", show_index=False, cell=52),
                KV(sorted(words), label="unvisited word set"), weights=[1, 2.4]),
            Text(f"tried {tried} one-letter changes, kept: {found}"),
            say=say,
        )

    assert answer == expected, answer
    path = [end]
    while path[-1] != begin:
        path.append(parent[path[-1]])
    path.reverse()
    yield Frame(
        graph({w: "match" for w in path}, {(a, b): "match" for a, b in zip(path, path[1:])}),
        Text((f'dequeued "{end}" at level {answer} -> answer {answer}', "match"),
             "one shortest ladder: " + " -> ".join(path),
             "time: O(N · L · 26)   space: O(N · L)"),
        say=f"cog comes out of the queue at level {answer}, so the answer is {answer}. Each word is queued once "
            "and makes L·25 candidates: O(N·L·26). Bidirectional BFS shrinks this in practice.",
    )
