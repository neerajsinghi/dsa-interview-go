import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from trieview import TNode, TrieView  # noqa: E402
from walkthrough import Frame, Grid, Text  # noqa: E402

TITLE = "Maximum XOR of Two Numbers"
PATTERN = "Bitwise trie"

MAX_BIT = 30  # same as the Go code: values below 2^31
SHOW = 5      # every example value is below 32, so only the low 5 levels differ


def insert_bits(root, x):
    cur = root
    for b in range(MAX_BIT, -1, -1):
        bit = (x >> b) & 1
        cur = cur.children.setdefault(bit, TNode())
    cur.end = True


def best_xor(root, x, steps=None):
    cur, xor, path = root, 0, ""
    for b in range(MAX_BIT, -1, -1):
        bit = (x >> b) & 1
        want = bit ^ 1
        if want in cur.children:
            xor |= 1 << b
            cur, took = cur.children[want], want
        else:
            cur, took = cur.children[bit], bit
        path += str(took)
        if steps is not None and b < SHOW:
            steps.append((b, bit, want, took, xor))
    return xor, path


def view_root(root):
    cur = root
    for _ in range(MAX_BIT + 1 - SHOW):
        cur = cur.children[0]
    return cur


def bits(x):
    return format(x, f"0{SHOW}b")


def leaf_notes(nums):
    return {bits(x): x for x in nums}


def path_hl(p, last="match"):
    hl = {p[:i]: "active" for i in range(len(p))}
    hl[p] = last
    return hl


def frames():
    nums = [3, 10, 5, 25, 2, 8]
    expected = 28
    table = Grid([list(bits(x)) for x in nums], label="numbers in binary (low 5 bits)",
                 row_labels=nums, col_labels=["16", "8", "4", "2", "1"], cell=30)

    yield Frame(
        table,
        speaker="Interviewer",
        say="Given non-negative numbers below 2^31, return the largest a XOR b over all pairs. "
            "For this list the answer is 28.",
    )
    yield Frame(
        table,
        Text("brute force: try every pair i <= j", ("O(n²) pairs", "reject")),
        say="Trying every pair is O(n²). XOR is decided bit by bit from the top, "
            "so I can do better by storing the numbers' bits in a trie.",
    )

    root = TNode()
    best = 0
    seen = []
    got = {}
    for x in nums:
        insert_bits(root, x)
        seen.append(x)
        steps = []
        v, path = best_xor(root, x, steps)
        vp = path[-SHOW:]
        old = best
        best = max(best, v)
        partner = int(vp, 2)
        got[x] = (partner, v)
        vr = view_root(root)
        if x == 25:
            yield Frame(
                TrieView(vr, label="bitwise trie (top 26 levels are all 0, not drawn)", notes=leaf_notes(seen),
                         hl={bits(25)[:i]: "compare" for i in range(SHOW + 1)}),
                Text(f"insert 25 = {bits(25)}, now ask: which stored y maximises 25 XOR y?"),
                say="Now 25 = 11001 goes in. To get the biggest XOR I walk from the top bit and "
                    "try to take the opposite bit at every level.",
            )
            walked = ""
            for b, bit, want, took, xor in steps:
                walked += str(took)
                ok = took == want
                hl = path_hl(walked, "match" if ok else "reject")
                yield Frame(
                    TrieView(vr, label="bitwise trie", notes=leaf_notes(seen), hl=hl),
                    Text(f"bit {b} (value {1 << b}): 25 has {bit}, want {want}",
                         (f"{'take ' + str(want) + ': XOR gains ' + str(1 << b) if ok else 'no ' + str(want) + ' child, must take ' + str(bit)}",
                          "match" if ok else "reject"),
                         f"XOR so far = {format(xor & 31, '05b')} = {xor}"),
                    say=(f"Bit {b}: 25 has a {bit}. A {want} child exists, so I take it and this bit "
                         f"is set in the XOR. Winning a higher bit beats all lower bits combined."
                         if ok and b == 4 else
                         f"Bit {b}: 25 has a {bit}, and a {want} child exists. I take it and gain {1 << b}."
                         if ok else
                         f"Bit {b}: 25 has a {bit}, but there's no {want} child below here. "
                         f"I'm forced to follow {bit}, and this bit adds nothing."),
                )
            yield Frame(
                TrieView(vr, label="bitwise trie", notes=leaf_notes(seen), hl=path_hl(vp)),
                Text(f"25 XOR {partner} = {bits(25)} ^ {bits(partner)} = {bits(v)} = {v}",
                     (f"best = max({old}, {v}) = {best}", "match")),
                say=f"The path ends at {partner}, so 25 XOR {partner} = {v}. That beats the old best of {old}.",
            )
            continue
        if x in (2, 8):
            if x == 8:
                yield Frame(
                    TrieView(vr, label="bitwise trie", notes=leaf_notes(seen)),
                    Text(f"insert 2, query 2: best partner {got[2][0]} -> {got[2][1]}",
                         f"insert 8, query 8: best partner {partner} -> {v}",
                         (f"best stays {best}", "match")),
                    say="The last two numbers go the same way: insert, then one greedy walk. "
                        f"2 pairs best with {got[2][0]} for {got[2][1]} and 8 with {partner} for {v}. Neither beats {best}.",
                )
            continue
        yield Frame(
            TrieView(vr, label="bitwise trie (top 26 levels are all 0, not drawn)", notes=leaf_notes(seen),
                     hl=path_hl(vp)),
            Text(f"insert {x} = {bits(x)}; greedy walk ends at {partner}",
                 (f"{x} XOR {partner} = {v}  ->  best = {best}", "match")),
            say=(f"I insert {x} first, then query it. It's the only number, so the walk ends at itself: XOR 0."
                 if len(seen) == 1 else
                 f"Insert {x}, then walk greedily. The walk ends at {partner}, giving {v}. "
                 "Each number is paired with everything inserted before it."),
        )

    assert best == expected, best
    yield Frame(
        TrieView(view_root(root), label="bitwise trie", notes=leaf_notes(nums), hl=path_hl(bits(5))),
        Text((f"answer: {best}  (5 XOR 25)", "match"),
             "time: O(n · 31)   brute force O(n²)",
             "space: O(n · 31) trie nodes"),
        say="The answer is 28. Each insert and each query walks 31 levels, so it's O(n·31), "
            "which is linear in n. The trie holds at most 31 nodes per number.",
    )
