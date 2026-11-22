from walkthrough import Array, Frame, Text, Tree

TITLE = "Fibonacci: naive vs memo vs iterative"
PATTERN = "Recursion tree"


def level_order(root):
    """Turn {id: (value, [child ids])} rooted at root into a level-order list."""
    nodes, kids = root
    out, queue = [nodes[0][0]], [0]
    while queue:
        p = queue.pop(0)
        for c in kids.get(p, [None, None]):
            out.append(None if c is None else nodes[c][0])
            if c is not None:
                queue.append(c)
    while out and out[-1] is None:
        out.pop()
    return out


def trace_naive(n):
    """Run FibNaive, recording the call tree. Returns (result, tree, call count)."""
    nodes, kids = [], {}

    def fib(k):
        me = len(nodes)
        nodes.append((k,))
        if k < 2:
            return k, me
        a, ca = fib(k - 1)
        b, cb = fib(k - 2)
        kids[me] = [ca, cb]
        return a + b, me

    res, _ = fib(n)
    return res, (nodes, kids), len(nodes)


def trace_memo(n):
    """Run FibMemo, recording the call tree and which calls were cache hits."""
    nodes, kids, hits, memo = [], {}, set(), {}

    def fib(k):
        me = len(nodes)
        nodes.append((k,))
        if k < 2:
            return k, me
        if k in memo:
            hits.add(me)
            return memo[k], me
        a, ca = fib(k - 1)
        b, cb = fib(k - 2)
        kids[me] = [ca, cb]
        memo[k] = a + b
        return memo[k], me

    res, _ = fib(n)
    return res, (nodes, kids), hits


def bfs_index(tree):
    """Map call id → its position in the level-order list (for highlights)."""
    nodes, kids = tree
    pos, order, queue, i = {0: 0}, [], [0], 1
    while queue:
        p = queue.pop(0)
        for c in kids.get(p, [None, None]):
            if c is not None:
                pos[c] = i
                queue.append(c)
            i += 1
    return pos


def frames():
    n = 10
    small = 5

    yield Frame(
        Text("F(0) = 0, F(1) = 1, F(n) = F(n-1) + F(n-2)", "n = 10  →  expected 55"),
        speaker="Interviewer",
        say="Return the n-th Fibonacci number. For n = 10 the answer is 55. What does each version cost?",
    )

    # FibNaive: draw the real call tree for F(5) and mark repeated calls.
    res5, tree, calls5 = trace_naive(small)
    res10, _, calls10 = trace_naive(n)
    pos = bfs_index(tree)
    seen, hl = set(), {}
    nodes, kids = tree

    def mark(c, repeated):
        k = nodes[c][0]
        rep = repeated or (k >= 2 and k in seen)
        if k >= 2:
            seen.add(k)
        if rep:
            hl[pos[c]] = "reject"
        for ch in kids.get(c, []):
            mark(ch, rep)

    mark(0, False)
    yield Frame(
        Tree(level_order(tree), label="FibNaive(5) call tree (node k = a call to F(k))", hl=hl),
        Text(f"F(5): {calls5} calls     F(10): {calls10} calls     ~1.618ⁿ growth", ("red = repeated work", "reject")),
        say=f"The naive recursion follows the definition. Each call branches twice, and whole subtrees repeat. "
            f"F(5) makes {calls5} calls and F(10) makes {calls10}.",
    )

    # FibMemo: the same tree, with cache hits cut off.
    m5, mtree, hits = trace_memo(small)
    mpos = bfs_index(mtree)
    yield Frame(
        Tree(level_order(mtree), label="FibMemo(5) call tree", hl={mpos[h]: "match" for h in hits}),
        Text(f"{len(mtree[0])} calls instead of {calls5}", ("green = answered from the memo, no recursion", "match")),
        say="If I store each F(k) the first time, a repeat is a map lookup. The tree collapses to a path: "
            "O(n) time, but O(n) memo plus O(n) stack.",
    )

    # FibIter: keep only the last two values.
    fib = [None] * (n + 1)
    prev, cur = 0, 1
    fib[0], fib[1] = prev, cur
    yield Frame(
        Array(fib, label="F(0..10)", pointers={"prev": 0, "cur": 1}, hl={0: "active", 1: "active"}),
        Text("prev, cur = 0, 1   (F(0), F(1))"),
        say="I can go bottom-up instead. Each value needs only the two before it, so I keep just prev and cur.",
    )
    for t in range(1, n + 1):
        prev, cur = cur, prev + cur
        if t + 1 <= n:
            fib[t + 1] = cur
        if t <= 3 or t == n:
            ptr = {"prev": t}
            if t + 1 <= n:
                ptr["cur"] = t + 1
            hl = {t: "active"}
            if t + 1 <= n:
                hl[t + 1] = "match"
            line = f"step {t}: prev, cur = cur, prev+cur = {prev}, {cur}"
            if t == 1:
                say = "Each step slides the pair one place right: the new cur is prev plus cur."
            elif t == n:
                say = f"After {n} steps prev holds F({n}) = {prev}. cur already ran one ahead, to {cur}."
            else:
                say = f"Step {t}: prev is F({t}) = {prev}, cur is F({t + 1}) = {cur}. Two variables, no matter how big n is."
            yield Frame(
                Array(fib, label="F(0..10)", pointers=ptr, hl=hl),
                Text(line, (f"cur = {cur} past the end (not stored)", "visited") if t == n else ""),
                say=say,
            )
        elif t == 4:
            yield Frame(
                Array(fib[:], label="F(0..10)", pointers={"prev": t, "cur": t + 1},
                      hl={t: "active", t + 1: "match"}),
                Text(f"step {t}: prev, cur = {prev}, {cur}", "… same for steps 5 to 9"),
                say="The same update repeats for every step, one addition each time.",
            )

    assert prev == 55 and res10 == 55 and res5 == 5 and m5 == 5
    yield Frame(
        Text(("FibIter(10) = 55", "match"),
             f"FibNaive  O(φⁿ) time   O(n) stack   ({calls10} calls for n=10)",
             "FibMemo   O(n) time    O(n) memo + O(n) stack",
             "FibIter   O(n) time    O(1) space"),
        say="The answer is 55. Naive is exponential, memo is linear with O(n) memory, "
            "and the loop is linear with O(1) space.",
    )
