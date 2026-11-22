from walkthrough import Array, Frame, KV, Text

TITLE = "Minimum Window Substring"
PATTERN = "Sliding window"


def frames():
    s, t = "ADOBECODEBANC", "ABC"

    need = {}
    for c in t:
        need[c] = need.get(c, 0) + 1

    def needkv(hl_key=None):
        shown = {c: need.get(c, 0) for c in sorted(set(t))}
        return KV(shown, label="need (count still required; < 0 = surplus)",
                  hl={hl_key: "active"} if hl_key in shown else {})

    def arr(left, right, best=None, hl_extra=None):
        hl = {k: "window" for k in range(left, right + 1)}
        if best:
            pass
        hl.update(hl_extra or {})
        ptr = {}
        if right >= 0:
            ptr = {"L": left, "R": right}
        return Array(list(s), label="s", hl=hl, pointers=ptr, window=(left, right) if right >= left else None)

    yield Frame(
        Array(list(s), label="s"),
        needkv(),
        Text(f's = "{s}", t = "{t}"'),
        speaker="Interviewer",
        say="Return the shortest substring of s that contains every character of t, duplicates included. "
            "If there are several, return the leftmost.",
    )
    yield Frame(
        Array(list(s), label="s"),
        needkv(),
        Text("brute force: every start, extend until valid → O(|s|²)"),
        say="For each start I could extend right until the window covers t. That's O(|s|²). "
            "A sliding window never moves backwards, so each index enters and leaves once.",
    )

    missing = len(t)
    best_l, best_len = 0, len(s) + 1
    left = 0
    shown_first_skip = False
    for right in range(len(s)):
        c = s[right]
        useful = need.get(c, 0) > 0
        if useful:
            missing -= 1
        need[c] = need.get(c, 0) - 1
        if useful and missing > 0:
            yield Frame(
                arr(left, right, hl_extra={right: "active"}),
                needkv(c),
                Text(f"R = {right}: '{c}' was needed → missing = {missing}"),
                say=f"I grow the window to the right. '{c}' is still needed, so missing drops to {missing}. "
                    "Letters t doesn't contain just go negative in need.",
            )
        elif not useful and missing > 0 and c in t and not shown_first_skip:
            shown_first_skip = True
            yield Frame(
                arr(left, right, hl_extra={right: "visited"}),
                needkv(c),
                Text(f"R = {right}: '{c}' is surplus → need['{c}'] = {need[c]}, missing stays {missing}"),
                say=f"Another '{c}' arrives, but the window doesn't need it. need goes to {need[c]} and missing "
                    "doesn't change. That surplus lets me drop the earlier copy later.",
            )
        silent = []
        while missing == 0:
            recorded = False
            if right - left + 1 < best_len:
                best_l, best_len = left, right - left + 1
                recorded = True
            d = s[left]
            need[d] = need.get(d, 0) + 1
            broke = need[d] > 0
            if broke:
                missing += 1
            if recorded or broke:
                lines = [f"window s[{left}..{right}] = \"{s[left:right + 1]}\" is valid"
                         + (f"  (already dropped {', '.join(silent)})" if silent else "")]
                if recorded:
                    lines.append((f"best = \"{s[best_l:best_l + best_len]}\" (length {best_len})", "match"))
                if broke:
                    lines.append((f"drop '{d}' → need['{d}'] = 1 → missing = 1, stop shrinking", "reject"))
                else:
                    lines.append(f"drop '{d}' (not needed), keep shrinking")
                hl = {k: "match" for k in range(left, right + 1)}
                hl[left] = "reject" if broke else "visited"
                say = ((f"'{c}' completes the window. I shrink past {', '.join(silent)}, none of them needed. "
                        if silent else f"missing is 0, so s[{left}..{right}] covers t. ")
                       + (f"Length {right - left + 1} is the best so far. " if recorded else "")
                       + (f"Dropping '{d}' breaks it, so I go back to growing R." if broke
                          else f"'{d}' isn't needed, so I drop it and try smaller."))
                yield Frame(
                    Array(list(s), label="s", hl=hl, pointers={"L": left, "R": right}, window=(left, right)),
                    needkv(d),
                    Text(*lines),
                    say=say,
                )
                silent = []
            else:
                silent.append(d)
            left += 1

    ans = s[best_l:best_l + best_len] if best_len <= len(s) else ""
    assert ans == "BANC"
    yield Frame(
        Array(list(s), label="s", hl={k: "match" for k in range(best_l, best_l + best_len)},
              window=(best_l, best_l + best_len - 1)),
        needkv(),
        Text((f'answer: "{ans}"', "match"), "time O(|s| + |t|) · space O(1) (256-entry array)"),
        say="The best window is \"BANC\". L and R each move at most |s| times, so it's O(|s| + |t|) time. "
            "need is a fixed 256-entry array: O(1) space.",
    )
