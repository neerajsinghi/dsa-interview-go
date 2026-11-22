from walkthrough import Frame, LinkedList, Text

TITLE = "Remove Nth Node From End"
PATTERN = "Two pointers, gap"


def frames():
    vals, n = [1, 2, 3, 4, 5], 2
    values = ["D"] + vals  # slot 0 is the dummy
    N = len(values)
    nxt = {i: (i + 1 if i + 1 < N else None) for i in range(N)}

    def view(slow, fast, hl=None, links=None):
        h = {0: "special"}
        h.update(hl or {})
        ptr = {"slow": slow, "fast": fast if fast is not None else N}
        return LinkedList(values, label="list (D = dummy)", hl=h, pointers=ptr, links=links or dict(nxt))

    def gap_text(slow, fast):
        f = "nil" if fast is None else values[fast]
        return f"slow = {values[slow]}, fast = {f}"

    yield Frame(
        LinkedList(vals, label="head"),
        Text(f"n = {n}  →  remove 4  →  1 → 2 → 3 → 5"),
        speaker="Interviewer",
        say="Delete the n-th node from the end, where n = 1 is the last node, and return the head. Can you do it in one pass?",
    )
    yield Frame(
        LinkedList(vals, label="head"),
        Text("pass 1: count length L = 5", "pass 2: walk L − n = 3 nodes, unlink the next one"),
        say="Two passes work: count the length, then walk to node L − n and unlink the one after it. "
            "To do one pass, I keep two pointers a fixed distance apart.",
    )

    slow = fast = 0
    yield Frame(
        view(slow, fast),
        Text("slow = fast = dummy", f"move fast n + 1 = {n + 1} steps first"),
        say="Both start at a dummy node in front of the head. The dummy means deleting the real head needs no special case.",
    )
    for i in range(n + 1):
        fast = nxt[fast]
    yield Frame(
        view(slow, fast, {fast: "compare"}),
        Text(gap_text(slow, fast), f"gap = n + 1 = {n + 1} nodes"),
        say=f"fast moves {n + 1} steps ahead. With a gap of n + 1, when fast falls off the end, slow stops one node "
            "before the target, which is what I need to unlink it.",
    )
    k = 0
    while fast is not None:
        slow, fast = nxt[slow], nxt[fast]
        k += 1
        yield Frame(
            view(slow, fast, {slow: "active"} | ({fast: "compare"} if fast is not None else {})),
            Text(gap_text(slow, fast), "move both one step" + ("  → fast is nil, stop" if fast is None else "")),
            say=("Now both move one step at a time, keeping the gap." if k == 1 else
                 "fast is nil, so slow is exactly n + 1 nodes from the end: the node before the one to delete."
                 if fast is None else "Same again: both move one step."),
        )

    target = nxt[slow]
    links = dict(nxt)
    links[slow] = nxt[target]
    links.pop(target)
    yield Frame(
        view(slow, None, {slow: "active", target: "reject"}, links),
        Text(f"slow.Next = slow.Next.Next  →  {values[slow]} now points to {values[nxt[target]] if nxt[target] else 'nil'}"),
        say=f"slow is on {values[slow]}. I set slow.Next = slow.Next.Next, which skips {values[target]}. "
            "Nothing points to it any more, so it's out of the list.",
    )

    out, i = [], links[0]
    while i is not None:
        out.append(values[i])
        i = links.get(i)
    assert out == [1, 2, 3, 5]
    yield Frame(
        LinkedList(out, label="return dummy.Next"),
        Text(("1 → 2 → 3 → 5", "match"), "one pass: O(L) time, O(1) space"),
        say="I return dummy.Next. If the head had been deleted, dummy.Next would already be the new head. "
            "One pass, O(L) time, O(1) space.",
    )
