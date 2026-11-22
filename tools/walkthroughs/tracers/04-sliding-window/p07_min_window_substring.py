from walkthrough import Array, Frame, KV, Text

TITLE = "Minimum Window Substring"
PATTERN = "Variable window"


def frames():
    s, t = "ADOBECODEBANC", "ABC"
    chars = list(s)

    yield Frame(
        Array(chars, label="s"),
        Text(f't = "{t}"', "shortest substring of s containing every byte of t (with counts)"),
        speaker="Interviewer",
        say='Return the shortest substring of s that contains every character of t = "ABC", counting repeats. '
            'If none exists, return "".',
    )
    yield Frame(
        Array(chars, label="s"),
        Text("brute force: from every start, extend until it covers t", "O(n²·σ) with a full coverage check"),
        say="From every start I could extend until the window covers t, which is O(n²) windows with a check each. "
            "Instead: expand right until valid, then shrink left while it stays valid.",
    )

    need = {}
    for c in t:
        need[c] = need.get(c, 0) + 1
    missing = len(t)
    best_start, best_len = 0, -1
    left = 0

    def view(lo, hi, extra=None, ptr=True):
        hl = {i: "window" for i in range(lo, hi + 1)}
        for i in range(lo, hi + 1):
            if s[i] in t:
                hl[i] = "active"
        if extra:
            hl.update(extra)
        return Array(chars, label="s", hl=hl, window=(lo, hi) if lo <= hi else None,
                     pointers={"left": lo, "right": hi} if ptr else None)

    def needs():
        return KV({c: need.get(c, 0) for c in t}, label="need (still required; negative = surplus)")

    def best_line():
        if best_len == -1:
            return "best = none yet"
        return f'best = "{s[best_start:best_start + best_len]}" (len {best_len})'

    expand_from = 0
    for right in range(len(s)):
        c = s[right]
        if need.get(c, 0) > 0:
            missing -= 1
        need[c] = need.get(c, 0) - 1
        if missing:
            continue
        # The window just became valid: show the expansion that got us here.
        added = s[expand_from:right + 1]
        first_valid = best_len == -1
        yield Frame(
            view(left, right),
            needs(),
            Text(f'expand right over "{added}" → missing = 0, window valid', best_line()),
            say=("I grow the right edge and count down what t still needs. At 'C' nothing is missing, "
                 f'so "{s[left:right + 1]}" is my first valid window.') if first_valid else
                (f"I expand again until '{c}' fills the gap. missing is 0, so "
                 f'"{s[left:right + 1]}" is valid. Now I shrink from the left.'),
        )
        skipped = []
        while missing == 0:
            length = right - left + 1
            improved = best_len == -1 or length < best_len
            if improved:
                best_start, best_len = left, length
            d = s[left]
            need[d] = need.get(d, 0) + 1
            if need[d] > 0:
                missing += 1
            if missing:
                note = (f"Dropping {', '.join(skipped)} kept it valid but never beat the best. " if skipped else "")
                yield Frame(
                    view(left, right, {**({i: "match" for i in range(best_start, best_start + best_len)}
                                          if improved else {}), left: "reject"}),
                    needs(),
                    Text((f"len {length}" + ("  new best" if improved else "") + f'   {best_line()}',
                          "match" if improved else None),
                         (f"drop s[{left}] = '{d}' → need['{d}'] = 1, missing = 1 → invalid", "reject")),
                    say=(note + (f"Length {length} is a new best. " if improved else "")
                         + f"Dropping '{d}' loses a required character, so the window is invalid again. Back to expanding."),
                )
            elif improved:
                note = (f"Dropping {', '.join(skipped)} kept it valid. " if skipped else "")
                skipped = []
                yield Frame(
                    view(left, right, {left: "reject"}),
                    needs(),
                    Text((f"len {length}  new best   {best_line()}", "match"),
                         f"drop s[{left}] = '{d}' → still valid, keep shrinking"),
                    say=note + f"Length {length} is a new best. Dropping '{d}' doesn't break coverage, so I keep shrinking.",
                )
            else:
                skipped.append(f"'{d}'")
            left += 1
        expand_from = right + 1

    ans = "" if best_len == -1 else s[best_start:best_start + best_len]
    assert ans == "BANC"  # doc comment: s = "ADOBECODEBANC", t = "ABC" → "BANC"
    yield Frame(
        view(best_start, best_start + best_len - 1, ptr=False,
             extra={i: "match" for i in range(best_start, best_start + best_len)}),
        Text((f'answer: "{ans}"', "match"), "time O(n + m): right and left each pass over s once",
             "space O(σ): one counter per byte value (256)"),
        say=f'right reached the end, so the answer is "{ans}". left and right each move forward at most n times, '
            "so it's O(n + m) time with a 256-entry table.",
    )
