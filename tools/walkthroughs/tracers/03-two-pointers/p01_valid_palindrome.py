from walkthrough import Array, Frame, Text

TITLE = "Valid Palindrome"
PATTERN = "Two pointers"

S = "A man, a plan, a canal: Panama"
SHOWN = 7  # events drawn one by one before summarising the rest


def is_alnum(c):
    return c.isascii() and c.isalnum()


def cells():
    return ["␣" if c == " " else c for c in S]


def arr(hl, l, r):
    return Array(cells(), label="s", hl=hl, pointers={"l": l, "r": r}, cell=23)


def frames():
    yield Frame(
        Array(cells(), label="s", cell=23),
        Text('s = "A man, a plan, a canal: Panama"', "ignore case and anything that isn't a letter or digit"),
        speaker="Interviewer",
        say="Ignoring case and every character that isn't a letter or digit, does this string read the same "
            "forwards and backwards?",
    )
    yield Frame(
        Array(cells(), label="s", cell=23),
        Text('brute force: clean = "amanaplanacanalpanama", compare with reverse',
             "O(n) time, but O(n) extra space for the copy"),
        say="The easy way builds a cleaned, lower-cased copy and compares it with its reverse. "
            "That's O(n) time but O(n) extra space. I can compare in place instead.",
    )

    # Run IsPalindrome, recording every loop iteration.
    events = []
    l, r = 0, len(S) - 1
    done = set()
    result = True
    while l < r:
        if not is_alnum(S[l]):
            events.append(("skipl", l, r))
            l += 1
            continue
        if not is_alnum(S[r]):
            events.append(("skipr", l, r))
            r -= 1
            continue
        if S[l].lower() != S[r].lower():
            events.append(("diff", l, r))
            result = False
            break
        events.append(("same", l, r))
        l += 1
        r -= 1
    assert result is True  # doc comment: "A man, a plan, a canal: Panama" → true
    final_l, final_r = l, r

    def show(c):
        return {" ": "space", ",": "comma", ":": "colon"}.get(c, repr(c))

    for k, (kind, l, r) in enumerate(events[:SHOWN]):
        hl = {i: "visited" for i in done}
        if kind == "same":
            hl.update({l: "match", r: "match"})
            line = (f"s[{l}] = {S[l]!r}, s[{r}] = {S[r]!r} → same letter ignoring case", "match")
            say = (f"{S[l]!r} and {S[r]!r} match once I lower-case them, so both ends are fine. "
                   "Move both pointers inward.")
            if k == 0:
                say = ("l starts at the front and r at the back. "
                       f"{S[l]!r} and {S[r]!r} match ignoring case, so I move both inward.")
            done.update({l, r})
        elif kind == "skipl":
            hl.update({l: "reject", r: "active"})
            line = (f"s[{l}] = {show(S[l])} is not a letter or digit → l++", "reject")
            say = f"s[{l}] is a {show(S[l])}, which doesn't count. Skip it by moving l right; r waits."
            done.add(l)
        else:
            hl.update({l: "active", r: "reject"})
            line = (f"s[{r}] = {show(S[r])} is not a letter or digit → r--", "reject")
            say = f"s[{r}] is a {show(S[r])}, which doesn't count. Skip it by moving r left; l waits."
            done.add(r)
        yield Frame(arr(hl, l, r), Text(f"l = {l}, r = {r}", line), say=say)

    rest = events[SHOWN:]
    compares = sum(1 for e in rest if e[0] == "same")
    skips = len(rest) - compares
    for kind, l, r in rest:
        done.update({l, r} if kind == "same" else ({l} if kind == "skipl" else {r}))
    fl, fr = final_l, final_r
    final_hl = {i: ("match" if is_alnum(S[i]) else "visited") for i in done}
    yield Frame(
        arr(final_hl, fl, fr),
        Text(f"{compares} more matching pairs, {skips} more skips", f"l = {fl}, r = {fr} → l >= r, loop ends"),
        say="The rest goes the same way: skip punctuation and spaces, compare letters, move inward. "
            "Every pair matched, and now the pointers have met in the middle.",
    )
    yield Frame(
        arr(final_hl, fl, fr),
        Text(("answer: true", "match"), "time O(n): each index is visited once", "space O(1): no cleaned copy"),
        say="No mismatch, so the answer is true. Each pointer only moves inward, so it's O(n) time, "
            "and I never built a copy, so it's O(1) extra space.",
    )
