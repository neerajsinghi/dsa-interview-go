from walkthrough import Array, Frame, KV, Text

TITLE = "Longest Repeating Character Replacement"
PATTERN = "Sliding window"

# The first doc example ("ABAB", k = 2) never needs a slide, so it hides the key
# trick. The second example ("AABABBA", k = 1) shows growing AND sliding.


def frames():
    s, k = "AABABBA", 1
    chars = list(s)

    yield Frame(
        Array(chars, label="s"),
        Text(f'k = {k}', "longest substring that can become one repeated letter"),
        speaker="Interviewer",
        say="You may change at most k = 1 letter. Return the length of the longest substring "
            "that can be made of one repeated letter.",
    )
    yield Frame(
        Array(chars, label="s"),
        Text("a window is fixable when  length - maxFreq <= k", "brute force: check every substring → O(n²)"),
        say="A window is fixable if everything except its most common letter fits in k changes. "
            "Checking every substring is O(n²); one sliding window does it in O(n).",
    )

    count = {}
    max_freq, left = 0, 0
    slides = 0
    for right, c in enumerate(s):
        count[c] = count.get(c, 0) + 1
        max_freq = max(max_freq, count[c])
        length = right - left + 1
        need = length - max_freq
        hl = {i: "window" for i in range(left, right + 1)}
        hl[right] = "active"
        if need > k:
            out = s[left]
            count[out] -= 1
            hl[left] = "reject"
            slides += 1
            old_left = left
            left += 1
            line = (f"{length} - {max_freq} = {need} > k → slide: drop s[{old_left}] = '{out}', left++", "reject")
            if slides == 1:
                say = (f"Now {need} letters would need changing, more than k. Instead of shrinking, I slide: "
                       f"drop '{out}' on the left. The window keeps length {length - 1}, the best so far.")
            elif max(count.values()) < max_freq:
                say = (f"'{c}' enters, '{out}' leaves. maxFreq is now stale: no letter has {max_freq} any more. "
                       "That's fine, because a stale value only keeps the window sliding; it can't grow it wrongly.")
            else:
                say = (f"Still {need} changes needed, so I slide again: '{c}' enters, '{out}' leaves, "
                       f"and the length stays {length - 1}.")
            win = (old_left, right)
        else:
            line = (f"{length} - {max_freq} = {need} <= k → fixable, window grows to {length}", "match")
            say = (f"'{c}' enters. The most common letter appears {max_freq} times, so {need} change"
                   f"{'' if need == 1 else 's'} would do. Fixable, so the window grows to {length}.")
            win = (left, right)
            best_win = win
        yield Frame(
            Array(chars, label="s", hl=hl, window=win, pointers={"left": win[0], "right": right}),
            KV({ch: count.get(ch, 0) for ch in sorted(set(s))}, label="count", hl={c: "active"}),
            Text(f"maxFreq = {max_freq}", line),
            say=say,
        )

    ans = len(s) - left
    assert ans == 4  # doc comment: "AABABBA", k = 1 → 4
    bl, br = best_win
    assert br - bl + 1 == ans
    yield Frame(
        Array(chars, label="s", hl={i: "match" for i in range(bl, br + 1)}, window=(bl, br)),
        Text((f"answer: len(s) - left = {len(s)} - {left} = {ans}", "match"),
             f'that length was reached by "{s[bl:br + 1]}" → one change gives "{s[bl:br + 1].replace("B", "A")}"',
             "time O(n), space O(1): 26 counters"),
        say=f"The window never shrinks, so its final length, {len(s)} - {left} = {ans}, is the answer. "
            f"It was earned by \"{s[bl:br + 1]}\". One pass, 26 counters: O(n) time, O(1) space.",
    )
