from walkthrough import Frame, LinkedList, Text

TITLE = "Linked List Cycle I & II"
PATTERN = "Fast & slow pointers"

VALS = [3, 2, 0, -4]
NXT = {0: 1, 1: 2, 2: 3, 3: 1}  # -4 links back to 2


def view(ptrs, hl=None, label="list (tail links back)"):
    return LinkedList(VALS, label=label, hl=hl or {}, pointers=ptrs, links=dict(NXT))


def frames():
    yield Frame(
        view({}),
        Text("3 → 2 → 0 → -4 → back to 2"),
        speaker="Interviewer",
        say="A list may loop back on itself. First tell me if there's a cycle. Then return the node where the cycle starts. O(1) memory.",
    )
    yield Frame(
        view({}),
        Text("hash set of visited nodes → first repeat is the start", "O(n) time, O(n) space"),
        say="I could store every node pointer in a set; the first node I see twice is the cycle start. "
            "That's O(n) space. Floyd's two pointers do it in O(1).",
    )

    slow = fast = 0
    yield Frame(
        view({"slow": slow, "fast": fast}, {0: "active"}),
        Text("slow moves 1 step, fast moves 2 steps"),
        say="slow moves one node per step, fast moves two. If there's no cycle, fast hits nil. If there is, "
            "fast gains one node per step and must land on slow.",
    )

    meet = None
    step = 0
    while True:
        # Port of: for fast != nil && fast.Next != nil (always true in this cyclic example)
        slow, fast = NXT[slow], NXT[NXT[fast]]
        step += 1
        if slow == fast:
            meet = slow
            yield Frame(
                view({"slow": slow, "fast": fast}, {slow: "match"}),
                Text(f"step {step}: slow = {VALS[slow]}, fast = {VALS[fast]}", ("slow == fast → cycle exists", "match")),
                say=f"After {step} steps both pointers sit on {VALS[slow]}. They met, so there is a cycle. "
                    "That answers part 1.",
            )
            break
        yield Frame(
            view({"slow": slow, "fast": fast}, {slow: "active", fast: "compare"}),
            Text(f"step {step}: slow = {VALS[slow]}, fast = {VALS[fast]}  (not equal)"),
            say=f"slow moves to {VALS[slow]}, fast jumps two nodes to {VALS[fast]}"
                + (", following the back link." if step > 1 else ".") + " Not equal yet, keep going.",
        )
    assert meet == 3

    p = 0
    yield Frame(
        view({"p": p, "slow": slow}, {p: "active", slow: "compare"}),
        Text("tail length a, meet b steps into a cycle of length c", "2(a + b) = a + b + kc  →  a = kc − b",
             "p = head; move p and slow one step at a time"),
        say="For the start: the distance from head to the start equals the distance from the meeting point "
            "to the start, mod the cycle length. So I put p at head and walk both one step at a time.",
    )
    walked = 0
    while p != slow:
        p, slow = NXT[p], NXT[slow]
        walked += 1
    yield Frame(
        view({"p": p, "slow": slow}, {p: "match"}),
        Text(f"p → {VALS[p]}, slow → {VALS[slow]}", (f"p == slow → cycle starts at node {VALS[p]}", "match"),
             "O(n) time, O(1) space"),
        say=f"After {walked} step{'s' if walked != 1 else ''} both are on {VALS[p]}, so the cycle starts there. Each phase is at most a few laps: "
            "O(n) time and only two pointers, O(1) space.",
    )
    assert VALS[p] == 2
