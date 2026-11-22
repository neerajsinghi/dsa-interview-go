from walkthrough import Frame, LinkedList, Text

TITLE = "Reorder List"
PATTERN = "Middle + reverse + weave"


def read(nxt, head):
    out, n = [], head
    while n is not None:
        out.append(n)
        n = nxt[n]
    return out


def frames():
    vals = [1, 2, 3, 4]
    N = len(vals)
    nxt = {i: (i + 1 if i + 1 < N else None) for i in range(N)}

    def view(ptrs=None, hl=None, label="nodes stay in place; arrows show Next"):
        return LinkedList(vals, label=label, hl=hl or {}, pointers=ptrs or {}, links=dict(nxt))

    def order(head):
        return " → ".join(str(vals[i]) for i in read(nxt, head)) + " → nil"

    yield Frame(
        view(label="head"),
        Text("L0 → L1 → L2 → L3   becomes   L0 → L3 → L1 → L2", "1 → 2 → 3 → 4   becomes   1 → 4 → 2 → 3"),
        speaker="Interviewer",
        say="Reorder the list to first, last, second, second-to-last, and so on. Relink the nodes in place; don't move values.",
    )
    yield Frame(
        view(label="head"),
        Text("put node pointers in a slice, relink by index i, j", "O(n) time, O(n) extra space"),
        say="With a slice of node pointers I can relink from both ends by index. That's O(n) extra space. "
            "Without it, I need three moves: find the middle, reverse the back half, weave.",
    )

    # 1. Middle with slow/fast.
    slow = fast = 0
    yield Frame(view({"slow": slow, "fast": fast}), Text("step 1: find the middle", "slow moves 1, fast moves 2"),
                say="Step 1: find the middle. slow moves one node, fast moves two.")
    while nxt[fast] is not None and nxt[nxt[fast]] is not None:
        slow, fast = nxt[slow], nxt[nxt[fast]]
    yield Frame(
        view({"slow": slow, "fast": fast}, {i: "active" for i in range(slow + 1)}),
        Text(f"fast.Next.Next is nil → stop", f"slow = {vals[slow]}: end of the first half"),
        say=f"fast can't take two more steps, so I stop. slow is on {vals[slow]}, the end of the first half. "
            "For odd lengths the first half gets the extra node.",
    )

    # 2. Reverse the second half and cut.
    prev, cur = None, nxt[slow]
    while cur is not None:
        nn = nxt[cur]
        nxt[cur] = prev
        prev, cur = cur, nn
    second = prev
    nxt[slow] = None
    first_half = read(nxt, 0)
    second_half = read(nxt, second)
    hl = {i: "active" for i in first_half}
    hl.update({i: "compare" for i in second_half})
    yield Frame(
        view({"head": 0, "second": second}, hl),
        Text("step 2: second = reverse(slow.Next); slow.Next = nil",
             f"first:  {order(0)}", f"second: {order(second)}"),
        say="Step 2: reverse the back half with the usual prev/cur loop, and cut slow.Next so the halves are separate. "
            "Now the back half runs from the end inward.",
    )

    # 3. Weave.
    first = 0
    k = 0
    while second is not None:
        f_next, s_next = nxt[first], nxt[second]
        nxt[first] = second
        nxt[second] = f_next
        k += 1
        a, b = first, second
        first, second = f_next, s_next
        yield Frame(
            view({"first": first if first is not None else N, "second": second if second is not None else N},
                 {a: "match", b: "match"}),
            Text(f"step 3: {vals[a]}.Next = {vals[b]},  {vals[b]}.Next = {vals[f_next] if f_next is not None else 'nil'}",
                 f"list now: {order(0)}"),
            say=(f"Step 3: weave. Save both nexts, point {vals[a]} at {vals[b]}, and {vals[b]} at the old next of {vals[a]}. "
                 "Then advance both." if k == 1 else
                 f"Same again: {vals[a]} points to {vals[b]}, {vals[b]} points to "
                 f"{vals[f_next] if f_next is not None else 'nil'}. second is now nil, so I stop."),
        )

    result = [vals[i] for i in read(nxt, 0)]
    assert result == [1, 4, 2, 3]
    yield Frame(
        LinkedList(result, label="head (reordered)", hl={0: "match", 1: "compare", 2: "match", 3: "compare"}),
        Text(("1 → 4 → 2 → 3", "match"), "O(n) time, O(1) extra space"),
        say="The list reads 1, 4, 2, 3. Each step is one linear pass with a few pointers: O(n) time, O(1) extra space.",
    )
