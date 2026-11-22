from walkthrough import Array, Frame, KV, Text

TITLE = "Repeated DNA Sequences"
PATTERN = "2-bit rolling hash"

CODE = {"A": 0, "C": 1, "G": 2, "T": 3}
L = 10
MASK = (1 << (2 * L)) - 1


def bits(w):
    b = format(w, "020b")
    return " ".join(b[k:k + 2] for k in range(0, 20, 2))


def frames():
    s = "AAAAACCCCCAAAAACCCCCCAAAAAGGGTTT"
    cell = 22

    def arr(hl=None, ptr=None, window=None):
        return Array(list(s), label="s", hl=hl or {}, pointers=ptr or {}, window=window, cell=cell)

    yield Frame(
        arr(),
        Text("return every 10-letter substring that occurs more than once"),
        KV([], label="answer", empty="[ ]"),
        speaker="Interviewer",
        say="s uses only A, C, G and T. Return every 10-letter substring that appears more than once, "
            "each listed once.",
    )
    yield Frame(
        arr(),
        Text("map[string]int of every window → n keys of 10 bytes each",
             "each key is hashed byte by byte"),
        KV([], label="answer", empty="[ ]"),
        say="A map from each 10-letter substring to its count works, and it's linear. But every key is a "
            "10-byte string that the map has to hash. I can make each window a small integer instead.",
    )
    yield Frame(
        arr(),
        Text("A=00  C=01  G=10  T=11   → 10 letters = 20 bits",
             ("window = (window<<2 | code(c)) & (1<<20 − 1)", "special")),
        KV([], label="answer", empty="[ ]"),
        say="Four letters fit in 2 bits, so a 10-letter window is a 20-bit number. Different windows always "
            "get different numbers, so it's a perfect hash with no collisions to check.",
    )

    count = {}
    out = []
    window = 0
    for i in range(len(s)):
        prev = window
        window = ((window << 2) | CODE[s[i]]) & MASK
        if i < L - 1:
            continue
        count[window] = count.get(window, 0) + 1
        lo = i - L + 1
        show = i in (L - 1, L) or count[window] == 2
        if not show:
            continue
        hl = {k: "window" for k in range(lo, i + 1)}
        hl[i] = "active"
        if i == L - 1:
            say = ("After 10 letters the window is full. \"AAAAACCCCC\" becomes five 00 pairs then five 01 pairs. "
                   "I count it in a map[uint32]int.")
            lines = [f"i = {i}: window = {bits(window)} = {window}", f"count[{window}] = 1"]
        elif i == L:
            say = ("Sliding by one is O(1): shift left 2 bits, OR in the new letter, and the mask drops "
                   "the oldest letter's 2 bits.")
            lines = [f"before: {bits(prev)}", f"after:  {bits(window)} = {window}  ('{s[i]}' in, '{s[lo - 1]}' out)"]
        else:
            sub = s[lo:i + 1]
            out.append(sub)
            hl = {k: "match" for k in range(lo, i + 1)}
            say = (f"At i = {i} window {window} is seen for the second time, so \"{sub}\" repeats. "
                   "I add it only when the count hits exactly 2, so it's listed once.")
            lines = [f"i = {i}: window = {bits(window)} = {window}", (f"count[{window}] = 2 → add \"{sub}\"", "match")]
        yield Frame(
            arr(hl, {"i": i}, window=(lo, i)),
            Text(*lines),
            KV(out, label="answer", hl={x: "match" for x in out}, empty="[ ]"),
            say=say,
        )

    out.sort()
    assert out == ["AAAAACCCCC", "CCCCCAAAAA"]
    yield Frame(
        arr(),
        Text((f"answer: {out}", "match"), f"{len(count)} distinct windows counted",
             "time O(n) · space O(n)"),
        KV(out, label="answer", hl={x: "match" for x in out}),
        say="After the scan I sort the two hits for a stable order. Each step is O(1), so it's O(n) time, "
            "and the map holds at most n windows: O(n) space.",
    )
