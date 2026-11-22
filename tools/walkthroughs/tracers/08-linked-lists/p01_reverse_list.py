from walkthrough import Frame, LinkedList, Text

TITLE = "Reverse Linked List"
PATTERN = "Pointer rewiring"

NIL = 0  # slot 0 is drawn as a "nil" box so prev has somewhere to point


def view(vals, nxt, prev, cur, nxt_ptr=None, flipped=None):
    values = ["nil"] + vals
    hl = {NIL: "empty"}
    for i in range(1, len(values)):
        if flipped and i in flipped:
            hl[i] = "match"
    if cur is not None:
        hl[cur] = "active"
    ptr = {"prev": prev}
    ptr["cur"] = cur if cur is not None else len(values)
    if nxt_ptr is not None:
        ptr["next"] = nxt_ptr
    links = {i: (j if j is not None else None) for i, j in nxt.items()}
    return LinkedList(values, label="list (slot 0 = nil)", hl=hl, pointers=ptr, links=links)


def read(nxt, head):
    out, n = [], head
    while n is not None and n != NIL:
        out.append(n)
        n = nxt.get(n)
    return out


def frames():
    vals = [1, 2, 3, 4]
    # Node k (1-based slot) holds vals[k-1]; nxt maps slot -> slot, None = nil.
    nxt = {i: i + 1 for i in range(1, len(vals))}
    nxt[len(vals)] = None

    yield Frame(
        LinkedList(vals, label="head"),
        Text("1 → 2 → 3 → 4 → nil   becomes   4 → 3 → 2 → 1 → nil"),
        speaker="Interviewer",
        say="Reverse a singly linked list in place and return the new head.",
    )
    yield Frame(
        LinkedList(vals, label="head"),
        Text("copy values to a slice, rebuild backwards → O(n) extra space",
             "recursion → O(n) call stack"),
        say="I could copy the values out and rebuild, or recurse, but both use O(n) extra memory. "
            "I'll flip each arrow in place with three pointers.",
    )

    prev, cur = NIL, 1
    yield Frame(
        view(vals, nxt, prev, cur),
        Text("prev = nil, cur = head"),
        say="prev starts at nil, since the old head becomes the new tail. cur starts at the head.",
    )

    flipped = set()
    step = 0
    while cur is not None:
        save = nxt[cur]  # next := cur.Next
        nxt[cur] = prev  # cur.Next = prev (slot 0 stands for nil)
        flipped.add(cur)
        step += 1
        before = (prev, cur)
        prev, cur = cur, save
        if step <= 2:
            say = (f"Save next = {vals[save - 1] if save else 'nil'} first, or I'd lose the rest of the list. "
                   f"Then point {vals[before[1] - 1]} back at {'nil' if before[0] == NIL else vals[before[0] - 1]}, "
                   "and slide prev and cur forward.")
        else:
            say = (f"Same three moves: save next, flip {vals[before[1] - 1]} to point at "
                   f"{vals[before[0] - 1]}, slide forward." + (" cur is now nil, so the loop ends." if cur is None else ""))
        yield Frame(
            view(vals, nxt, prev, cur, flipped=flipped),
            Text(f"next = cur.Next  → {vals[save - 1] if save else 'nil'}",
                 f"cur.Next = prev  → {vals[before[1] - 1]} now points to {'nil' if before[0] == NIL else vals[before[0] - 1]}",
                 "prev, cur = cur, next"),
            say=say,
        )

    order = [vals[k - 1] for k in read(nxt, prev)]
    assert order == [4, 3, 2, 1]
    yield Frame(
        LinkedList(order, label="return prev (new head)", hl={0: "match"}),
        Text(("4 → 3 → 2 → 1 → nil", "match"), "O(n) time, O(1) space"),
        say="prev sits on 4, the old tail, so I return it. Every node was visited once and I used three pointers: "
            "O(n) time, O(1) space.",
    )
