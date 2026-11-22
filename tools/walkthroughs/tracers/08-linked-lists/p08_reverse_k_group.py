from walkthrough import Frame, LinkedList, Text

TITLE = "Reverse Nodes in k-Group"
PATTERN = "Sublist reversal"


def frames():
    vals, k = [1, 2, 3, 4, 5], 2
    values = ["D"] + vals  # slot 0 is the dummy
    N = len(values)
    nxt = {i: (i + 1 if i + 1 < N else None) for i in range(N)}
    done = set()

    def name(s):
        return "nil" if s is None else str(values[s])

    def order():
        out, s = [], nxt[0]
        while s is not None:
            out.append(str(values[s]))
            s = nxt[s]
        return " → ".join(out)

    def view(ptrs, hl=None):
        h = {0: "special"}
        h.update({s: "match" for s in done})
        h.update(hl or {})
        p = {key: (v if v is not None else N) for key, v in ptrs.items()}
        return LinkedList(values, label="nodes stay in place; arrows show Next", hl=h, pointers=p, links=dict(nxt))

    yield Frame(
        LinkedList(vals, label="head"),
        Text(f"k = {k}:  1 → 2 → 3 → 4 → 5   becomes   2 → 1 → 4 → 3 → 5", "a short last group stays as it is"),
        speaker="Interviewer",
        say="Reverse the nodes k at a time. If fewer than k nodes are left at the end, leave them. Relink nodes, O(1) extra memory.",
    )
    yield Frame(
        LinkedList(vals, label="head"),
        Text("copy values to a slice, reverse each chunk, rebuild", "O(n) time, O(n) extra space, and it swaps values"),
        say="Copying the values out, reversing chunks and rebuilding is easy, but it's O(n) space and doesn't relink nodes. "
            "I'll reverse each group in place between two anchors.",
    )

    group_prev = 0
    yield Frame(
        view({"groupPrev": group_prev}),
        Text("groupPrev = dummy: the node just before the current group"),
        say="A dummy in front gives the first group a 'node before' like every other group. groupPrev marks that node.",
    )

    g = 0
    while True:
        kth = group_prev
        i = 0
        while i < k and kth is not None:
            kth = nxt[kth]
            i += 1
        if kth is None:
            yield Frame(
                view({"groupPrev": group_prev}, {s: "visited" for s in range(group_prev + 1, N)}),
                Text(f"walk {k} nodes from groupPrev → reached nil", ("fewer than k nodes left → stop", "reject"),
                     f"list now: {order()}"),
                say=f"Only {name(nxt[group_prev])} is left, fewer than {k} nodes, so the walk hits nil. "
                    "I leave it as it is and return dummy.Next.",
            )
            break
        g += 1
        group_next = nxt[kth]
        group = []
        s = nxt[group_prev]
        while s != group_next:
            group.append(s)
            s = nxt[s]
        yield Frame(
            view({"groupPrev": group_prev, "kth": kth, "groupNext": group_next}, {s: "active" for s in group}),
            Text(f"group {g}: {' → '.join(name(s) for s in group)}   kth = {name(kth)}, groupNext = {name(group_next)}"),
            say=(f"Walk {k} nodes from groupPrev to find kth = {name(kth)}. If I hit nil first, the group is short. "
                 f"groupNext = {name(group_next)} is where the rest of the list starts." if g == 1 else
                 f"Next group: kth = {name(kth)}, groupNext = {name(group_next)}. There are {k} nodes, so I reverse them."),
        )

        prev, cur = group_next, nxt[group_prev]
        while cur != group_next:
            save = nxt[cur]
            nxt[cur] = prev
            prev, cur = cur, save
        first = nxt[group_prev]
        yield Frame(
            view({"groupPrev": group_prev, "prev": prev, "cur": cur}, {s: "active" for s in group}),
            Text(f"reverse with prev starting at groupNext ({name(group_next)})",
                 f"{name(first)}.Next = {name(nxt[first])},  {name(kth)}.Next = {name(nxt[kth])}"),
            say=(f"I reverse the group with the usual prev/cur loop, but prev starts at groupNext. "
                 f"So the old first node, {name(first)}, ends up pointing straight at {name(group_next)}."
                 if g == 1 else f"Same reversal: {name(first)} now points at {name(group_next)} and {name(kth)} points back at {name(first)}."),
        )

        nxt[group_prev] = kth
        done.update(group)
        yield Frame(
            view({"groupPrev": first}),
            Text(f"groupPrev.Next = kth  → {name(group_prev)} → {name(kth)}", f"groupPrev = {name(first)} (old first, now last)",
                 f"list now: {order()}"),
            say=f"Reconnect: {name(group_prev)} now points at {name(kth)}, the group's new head. The old first node "
                f"{name(first)} is now the group's tail, so it becomes groupPrev for the next group.",
        )
        group_prev = first

    result = order()
    assert result == "2 → 1 → 4 → 3 → 5"
    yield Frame(
        LinkedList([int(v) for v in result.split(" → ")], label="return dummy.Next", hl={0: "match", 1: "match", 2: "match", 3: "match"}),
        Text((result, "match"), "each node is walked twice (check + reverse): O(n) time, O(1) space"),
        say="The answer is 2, 1, 4, 3, 5. Each node is visited once to check the group and once to reverse it: "
            "O(n) time, O(1) extra space.",
    )
