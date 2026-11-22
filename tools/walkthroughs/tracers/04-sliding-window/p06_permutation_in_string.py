from walkthrough import Array, Frame, Grid, Text

TITLE = "Permutation in String"
PATTERN = "Fixed window"

LETTERS = "abdeio"  # the letters that occur in s1 or s2; the other 20 always match (0 = 0)


def counts(s):
    c = [0] * 26
    for ch in s:
        c[ord(ch) - 97] += 1
    return c


def table(need, have, changed=()):
    cols = list(LETTERS)
    rows = [[need[ord(ch) - 97] for ch in cols], [have[ord(ch) - 97] for ch in cols]]
    hl = {}
    for j, ch in enumerate(cols):
        i = ord(ch) - 97
        hl[(1, j)] = "match" if need[i] == have[i] else "reject"
        if ch in changed:
            hl[(1, j)] = "active" if need[i] != have[i] else "match"
    return Grid(rows, row_labels=["need", "have"], col_labels=cols, hl=hl, cell=40)


def frames():
    s1, s2 = "ab", "eidbaooo"
    m = len(s1)
    chars = list(s2)

    yield Frame(
        Array(chars, label="s2"),
        Text(f's1 = "{s1}"', "does some rearrangement of s1 appear in s2?"),
        speaker="Interviewer",
        say='Does some rearrangement of s1 = "ab" appear as a contiguous substring of s2?',
    )
    yield Frame(
        Array(chars, label="s2"),
        Text("permutation ⇔ same letter counts, same length", "brute force: recount all 26 letters per window → O(n·m)"),
        say="A substring is a permutation of s1 exactly when it has the same letter counts. "
            "Recounting each window is O(n·m). I'll slide a window of length 2 and update counts in O(1).",
    )

    need, have = counts(s1), counts(s2[:m])
    matches = sum(1 for i in range(26) if need[i] == have[i])
    off = [ch for ch in LETTERS if need[ord(ch) - 97] != have[ord(ch) - 97]]
    off_s = ", ".join(off[:-1]) + " and " + off[-1]
    yield Frame(
        Array(chars, label="s2", hl={0: "window", 1: "window"}, window=(0, m - 1)),
        table(need, have),
        Text(f"matches = {matches} of 26 letters", "window valid ⇔ matches == 26"),
        say=f'The first window is "{s2[:m]}". {matches} of 26 letters already have matching counts; '
            f"{off_s} don't. The window is a permutation when all 26 match.",
    )

    state = {"matches": matches}

    def update(c, delta):
        i = ord(c) - 97
        was = have[i] == need[i]
        if was:
            state["matches"] -= 1
        have[i] += delta
        now = have[i] == need[i]
        if now:
            state["matches"] += 1
        effect = "starts matching (+1)" if now and not was else "stops matching (-1)" if was and not now else "no change"
        return f"'{c}': have {have[i] - delta} → {have[i]}, need {need[i]} → {effect}"

    found = None
    for right in range(m, len(s2)):
        if state["matches"] == 26:
            found = right - m
            break
        before = state["matches"]
        note_in = update(s2[right], +1)
        note_out = update(s2[right - m], -1)
        lo = right - m + 1
        hl = {i: "window" for i in range(lo, right + 1)}
        hl[right] = "active"
        hl[right - m] = "reject"
        yield Frame(
            Array(chars, label="s2", hl=hl, window=(lo, right)),
            table(need, have, changed=(s2[right], s2[right - m])),
            Text(note_in, note_out, (f"matches {before} → {state['matches']}", "match" if state["matches"] == 26 else None)),
            say=f"Slide: '{s2[right]}' enters and '{s2[right - m]}' leaves. Only those two letters change, "
                + (f"so I update matches for just them: it goes from {before} to {state['matches']}."
                   if before != state["matches"] else f"so I update matches for just them. One loss and one gain: it stays {before}.")
                + (" All 26 now match." if state["matches"] == 26 else ""),
        )
    if found is None and state["matches"] == 26:
        found = len(s2) - m

    assert found is not None  # doc comment: s1 = "ab", s2 = "eidbaooo" → true
    yield Frame(
        Array(chars, label="s2", hl={i: "match" for i in range(found, found + m)}, window=(found, found + m - 1)),
        table(need, have),
        Text((f'matches == 26 → true  ("{s2[found:found + m]}")', "match"),
             "time O(n + m): O(1) work per slide", "space O(1): two arrays of 26 counts"),
        say=f'At the next check matches is 26, so "{s2[found:found + m]}" is a permutation of "{s1}": return true. '
            "Each slide touches two letters, so it's O(n + m) time and O(1) space.",
    )
