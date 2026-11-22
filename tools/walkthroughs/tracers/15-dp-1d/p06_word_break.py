from walkthrough import Array, Frame, KV, Text

TITLE = "Word Break"
PATTERN = "1-D DP · prefix"


def frames():
    s, words = "leetcode", ["leet", "code"]
    d = set(words)
    max_len = max(len(w) for w in words)
    n = len(s)
    ok = [None] * (n + 1)

    def widgets(hl_s=None, hl_ok=None, hl_d=None):
        return (
            Array(list(s), label="s", hl=hl_s),
            Array(ok, label="ok[i] = the prefix s[:i] can be split into words", hl=hl_ok),
            KV(words, label="dict", hl=hl_d),
        )

    yield Frame(
        *widgets(),
        speaker="Interviewer",
        say="Can \"leetcode\" be split into a sequence of dictionary words? The dictionary is leet and code, "
            "and words may be reused.",
    )
    yield Frame(
        *widgets(),
        Text("try every word as the next piece, recurse on the rest", ("same suffixes re-checked → exponential", "reject")),
        say="The brute force tries each word as the next piece and recurses. With many overlapping words "
            "it re-checks the same suffix over and over: exponential.",
    )

    ok[0] = True
    yield Frame(
        *widgets(hl_ok={0: "match"}),
        Text("ok[0] = T  (the empty prefix)", f"look back at most {max_len} letters (longest word)"),
        say="State: ok[i] says whether the first i letters split into words. The empty prefix is fine. "
            "I only look back as far as the longest word, 4 letters.",
    )

    for i in range(1, n + 1):
        ok[i] = False
        tried = []
        hit = None
        j = i - 1
        while j >= 0 and i - j <= max_len:
            tried.append(j)
            if ok[j] and s[j:i] in d:
                ok[i] = True
                hit = j
                break
            j -= 1
        if i in (1, 4, 8):
            hl_ok = {i: "match" if ok[i] else "reject"}
            for t in tried:
                hl_ok.setdefault(t, "compare")
            if hit is not None:
                hl_s = {k: "match" for k in range(hit, i)}
                line = (f"ok[{hit}] = T and \"{s[hit:i]}\" ∈ dict  →  ok[{i}] = T", "match")
                if i == 4:
                    say = "At i = 4, the piece \"leet\" is a word and ok[0] is true. So the prefix \"leet\" splits."
                else:
                    say = ("At i = 8, the piece \"code\" is a word and ok[4] is true. "
                           "So the whole string splits: leet + code.")
                yield Frame(*widgets(hl_s, hl_ok, {s[hit:i]: "match"}),
                            Text("pieces tried: " + ", ".join(f"\"{s[t:i]}\"" for t in tried), line), say=say)
            else:
                yield Frame(*widgets({i - 1: "reject"}, hl_ok),
                            Text(f"\"{s[0:i]}\" is not a word", (f"ok[{i}] = F", "reject")),
                            say="For i = 1, the only piece is \"l\", which isn't a word. So ok[1] is false.")
        elif i in (3, 7):
            lo = 2 if i == 3 else 5
            yield Frame(
                *widgets(hl_ok={k: "reject" for k in range(lo, i + 1)}),
                Text(f"ok[{lo}..{i}] = F  (no word ends at those positions after a true ok[j])"),
                say=("Same for i = 2 and 3: \"le\", \"lee\" and their suffixes aren't words." if i == 3 else
                     "For i = 5 to 7, pieces like \"c\", \"co\", \"cod\" aren't words, and the longer ones aren't either."),
            )

    assert ok[n] is True
    yield Frame(
        *widgets(hl_ok={0: "match", 4: "match", 8: "match"}),
        Text(f"answer: ok[{n}] = T  (leet | code)", ("O(n · L) substring checks · O(n + dict) space", "match")),
        say="ok[8] is true, so the answer is true. Each position tries at most L lookbacks, L being the "
            "longest word: O(n · L) checks, O(n) extra space.",
    )
