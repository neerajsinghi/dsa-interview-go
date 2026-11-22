from walkthrough import Array, Frame, Grid, Text

TITLE = "Valid Anagram"
PATTERN = "Count array"

LETTERS = [chr(ord("a") + k) for k in range(26)]


def counts_grid(cs, ct, hl=None, label="letterCounts ([26]int)"):
    return Grid([cs, ct], label=label, row_labels=["s", "t"], col_labels=LETTERS, hl=hl or {}, cell=26)


def frames():
    s, t = "listen", "silent"

    yield Frame(
        Array(list(s), label="s"),
        Array(list(t), label="t"),
        speaker="Interviewer",
        say="Given two lowercase strings, tell me whether one is a rearrangement of the other's letters.",
    )
    yield Frame(
        Array(sorted(s), label="sorted(s)"),
        Array(sorted(t), label="sorted(t)"),
        Text(("equal → anagrams, but sorting costs O(n log n)", "compare")),
        say="The simple way sorts both strings and compares them. That's IsAnagramSort, O(n log n). "
            "Order doesn't matter here, only how many of each letter.",
    )

    same_len = len(s) == len(t)
    yield Frame(
        Array(list(s), label="s"),
        Array(list(t), label="t"),
        Text(f"len(s) = {len(s)}, len(t) = {len(t)} → {'continue' if same_len else 'return false'}"),
        say=f"First a cheap check: different lengths can't be anagrams. Both have {len(s)} letters, so I go on "
            "and count letters with a [26]int, one counter per letter.",
    )

    cs, ct = [0] * 26, [0] * 26
    for i, ch in enumerate(s):
        cs[ord(ch) - 97] += 1
        if i < 2 or i == len(s) - 1:
            say = (f"'{ch}' is letter {ord(ch) - 97}, so counter {ord(ch) - 97} goes up by one. "
                   "Array indexing, no hashing." if i == 0 else
                   f"'{ch}' bumps its own counter." if i == 1 else
                   "After the rest of s, the row holds how many of each letter s has.")
            yield Frame(
                Array(list(s), label="s", hl={k: "visited" for k in range(i)} | {i: "active"}, pointers={"i": i}),
                counts_grid(cs, ct, {(0, ord(ch) - 97): "active"}),
                Text(f"s[{i}] = '{ch}' → counts['{ch}'-'a'] = {cs[ord(ch) - 97]}"),
                say=say,
            )
    for i, ch in enumerate(t):
        ct[ord(ch) - 97] += 1
    yield Frame(
        Array(list(t), label="t", hl={k: "visited" for k in range(len(t))}),
        counts_grid(cs, ct, {(1, k): "active" for k in range(26) if ct[k]}),
        Text("same pass over t fills the second row"),
        say="I count t the same way into its own array. Each string is read once: O(n).",
    )

    ans = same_len and cs == ct
    used = [k for k in range(26) if cs[k] or ct[k]]
    yield Frame(
        counts_grid(cs, ct, {(r, k): "match" for k in used for r in range(2)}),
        Text(f"letterCounts(s) == letterCounts(t) → {str(ans).lower()}", ("arrays compare with == in Go", "match")),
        say="Go arrays are comparable, so I compare the two [26]int with ==. "
            "Every column matches, so they're anagrams.",
    )

    assert ans is True
    yield Frame(
        Text(("answer: true", "match"),
             "IsAnagramSort:    O(n log n) time, O(n) space",
             "IsAnagram:        O(n) time, O(1) space (26 counters)",
             "IsAnagramUnicode: O(n) time, O(k) space (map[rune]int)"),
        say="The answer is true. Counting is O(n) time and O(1) space because the alphabet is fixed. "
            "For Unicode I'd count runes in a map instead.",
    )
