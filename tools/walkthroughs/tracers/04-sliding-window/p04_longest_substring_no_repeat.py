from walkthrough import Array, Frame, KV, Text

TITLE = "Longest Substring Without Repeating Characters"
PATTERN = "Variable window"


def frames():
    s = "abcabcbb"
    chars = list(s)

    yield Frame(
        Array(chars, label="s"),
        Text('s = "abcabcbb"', "length of the longest substring with no repeated character"),
        speaker="Interviewer",
        say="Return the length of the longest substring in which no character appears twice.",
    )
    yield Frame(
        Array(chars, label="s"),
        Text("brute force: from every start, extend until a repeat", "O(n²)"),
        say="From every start I could extend until I hit a repeat. That's O(n²). Instead I'll keep one window "
            "with no repeats and remember where each character was last seen.",
    )

    next_after = {}  # Go uses [256]int; 0 means "never seen"
    best, left = 0, 0
    best_range = (0, 0)
    jumps = 0
    for right, c in enumerate(s):
        prev = next_after.get(c, 0)
        old_left = left
        left = max(left, prev)
        next_after[c] = right + 1
        length = right - left + 1
        improved = length > best
        if improved:
            best, best_range = length, (left, right)
        hl = {i: "window" for i in range(left, right + 1)}
        hl[right] = "active"
        for i in range(old_left, left):
            hl[i] = "reject"
        if left != old_left:
            jump = (f"'{c}' last seen at {prev - 1}, inside the window → left = {left}", "reject")
            jumps += 1
            say = (f"'{c}' repeats too: it was last at {prev - 1}, so left jumps to {left}. "
                   f"The window is now length {length}.") if jumps > 1 else (f"'{c}' already appears at index {prev - 1}, inside the window. Rather than shrinking one step "
                   f"at a time, I jump left straight to {left}, just past it.")
        elif prev:
            jump = (f"'{c}' last seen at {prev - 1}, left of the window → left stays {left}", None)
            say = (f"'{c}' was seen at {prev - 1}, but that's already outside the window. "
                   "That's why I take the max: left must never move backwards.")
        else:
            jump = (f"'{c}' not seen before → left stays {left}", None)
            say = (f"'{c}' is new, so the window just grows. Length {length}"
                   + (", a new best." if improved else "."))
        if right == 0:
            say = "I record one past each character's last index, so 0 means never seen. 'a' is new; the window is just 'a'."
        yield Frame(
            Array(chars, label="s", hl=hl, window=(left, right), pointers={"left": left, "right": right}),
            KV({k: v for k, v in sorted(next_after.items())}, label="nextAfter (last index + 1)", hl={c: "active"}),
            Text(jump, (f"length = {length}   best = {best}", "match" if improved else None)),
            say=say,
        )

    assert best == 3  # doc comment: "abcabcbb" → 3
    bl, br = best_range
    yield Frame(
        Array(chars, label="s", hl={i: "match" for i in range(bl, br + 1)}, window=(bl, br)),
        Text((f'answer: {best}  ("{s[bl:br + 1]}")', "match"), "time O(n): one pass, left only moves forward",
             "space O(σ): one slot per possible byte (256)"),
        say=f"The answer is {best}, first reached with \"{s[bl:br + 1]}\". It's one pass, so O(n) time, "
            "plus a table of 256 entries, which is O(σ) space for the alphabet.",
    )
