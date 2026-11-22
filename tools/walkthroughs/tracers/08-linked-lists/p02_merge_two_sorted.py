from walkthrough import Frame, LinkedList, Text

TITLE = "Merge Two Sorted Lists"
PATTERN = "Dummy head"


def views(A, B, ia, ib, merged, used_a, used_b, pick=None):
    """A, B: value lists; ia/ib: index of the a/b pointer (len = nil);
    merged: [(src, idx)] spliced so far; used_*: indices already spliced."""
    ha = {i: "visited" for i in used_a}
    hb = {i: "visited" for i in used_b}
    if pick == "a":
        ha[ia] = "match"
        if ib < len(B):
            hb[ib] = "compare"
    elif pick == "b":
        hb[ib] = "match"
        if ia < len(A):
            ha[ia] = "compare"
    la = LinkedList(A, label="list a", hl=ha, pointers={"a": ia})
    lb = LinkedList(B, label="list b", hl=hb, pointers={"b": ib})
    mvals = ["D"] + [(A if s == "a" else B)[i] for s, i in merged]
    mh = {0: "special"}
    for k, (s, _) in enumerate(merged, start=1):
        mh[k] = "active" if s == "a" else "compare"
    lm = LinkedList(mvals, label="merged (D = dummy)", hl=mh, pointers={"tail": len(mvals) - 1})
    return la, lb, lm


def frames():
    A, B = [1, 2, 4], [1, 3, 4]

    yield Frame(
        LinkedList(A, label="list a"), LinkedList(B, label="list b"),
        Text("expected: 1 → 1 → 2 → 3 → 4 → 4"),
        speaker="Interviewer",
        say="Both lists are sorted. Splice their nodes into one sorted list, reusing the nodes, and return its head.",
    )
    yield Frame(
        LinkedList(A, label="list a"), LinkedList(B, label="list b"),
        Text("collect all values, sort, rebuild → O((m+n) log(m+n)) and new nodes"),
        say="I could dump every value into a slice, sort it and build new nodes. That's O((m+n) log(m+n)) and "
            "allocates. Both lists are already sorted, so I can merge them.",
    )

    ia = ib = 0
    merged, used_a, used_b = [], set(), set()
    yield Frame(
        *views(A, B, ia, ib, merged, used_a, used_b),
        Text("dummy := &ListNode{}, tail := dummy"),
        say="A dummy node in front means I never ask 'which list gives the head?'. tail always points at the last merged node.",
    )

    step = 0
    while ia < len(A) and ib < len(B):
        step += 1
        if A[ia] <= B[ib]:
            yield Frame(*views(A, B, ia, ib, merged, used_a, used_b, "a"),
                        Text(f"a = {A[ia]} <= b = {B[ib]} → tail.Next = a, a = a.Next"),
                        say=(f"{A[ia]} <= {B[ib]}, so a's node goes next. On a tie I take from a, which keeps the merge stable."
                             if A[ia] == B[ib] else
                             f"{A[ia]} is smaller than {B[ib]}, so I hang a's node on the tail and advance a."))
            merged.append(("a", ia))
            used_a.add(ia)
            ia += 1
        else:
            yield Frame(*views(A, B, ia, ib, merged, used_a, used_b, "b"),
                        Text(f"a = {A[ia]} > b = {B[ib]} → tail.Next = b, b = b.Next"),
                        say=f"{B[ib]} is smaller than {A[ia]}, so b's node goes next. Then tail moves onto it.")
            merged.append(("b", ib))
            used_b.add(ib)
            ib += 1

    rest = "a" if ia < len(A) else "b"
    left = [A[i] for i in range(ia, len(A))] if rest == "a" else [B[i] for i in range(ib, len(B))]
    if rest == "a":
        merged += [("a", i) for i in range(ia, len(A))]
        used_a |= set(range(ia, len(A)))
    else:
        merged += [("b", i) for i in range(ib, len(B))]
        used_b |= set(range(ib, len(B)))
    yield Frame(
        *views(A, B, len(A) if rest == "a" else ia, len(B) if rest == "b" else ib, merged, used_a, used_b),
        Text(f"{'a' if rest == 'b' else 'b'} is nil → tail.Next = {rest}  (attach {left})"),
        say=f"One list ran out. What's left of list {rest} is already sorted and bigger than everything merged, "
            "so I attach it in one step instead of looping.",
    )

    result = [(A if s == "a" else B)[i] for s, i in merged]
    assert result == [1, 1, 2, 3, 4, 4]
    yield Frame(
        LinkedList(result, label="return dummy.Next", hl={i: ("active" if s == "a" else "compare") for i, (s, _) in enumerate(merged)}),
        Text(("1 → 1 → 2 → 3 → 4 → 4", "match"), "blue = from a, orange = from b", "O(m + n) time, O(1) extra space"),
        say="I return dummy.Next, skipping the dummy. Each node was spliced once with no new nodes: O(m + n) time, O(1) extra space.",
    )
