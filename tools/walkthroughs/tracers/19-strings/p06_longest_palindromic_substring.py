from walkthrough import Array, Frame, Text

TITLE = "Longest Palindromic Substring"
PATTERN = "Manacher"


def frames():
    s = "babad"
    n = len(s)
    m = 2 * n + 1
    t = ["#" if i % 2 == 0 else s[i // 2] for i in range(m)]

    def match(i, j):
        return i % 2 == 0 or s[i // 2] == s[j // 2]

    yield Frame(
        Array(list(s), label="s"),
        Text('"babad"'),
        speaker="Interviewer",
        say="Return the longest palindromic substring of s. If there's a tie, return the leftmost one.",
    )
    yield Frame(
        Array(list(s), label="s"),
        Text("brute force: every substring × palindrome check → O(n³)",
             "expand around 2n−1 centres → O(n²)"),
        say="Checking every substring is O(n³). Expanding around each of the 2n−1 centres is O(n²) and fine for "
            "small n. Manacher gets O(n) by reusing work from mirrored positions.",
    )
    yield Frame(
        Array(t, label="t (virtual: # between letters)"),
        Array([""] * m, label="rad", show_index=False),
        Text("t[i] = '#' for even i, s[i/2] for odd i", "rad[i] in t = palindrome length in s"),
        say="I put a separator between letters, so every palindrome has odd length and one centre. "
            "The Go code never builds t; it uses i/2 to read s.",
    )

    rad = [None] * m
    center = right = 0
    best = 0
    for i in range(m):
        r = 0
        old_right = right
        mirror = None
        if i < right:
            mirror = 2 * center - i
            r = min(rad[mirror], right - i)
        free = r
        while i - r - 1 >= 0 and i + r + 1 < m and match(i - r - 1, i + r + 1):
            r += 1
        rad[i] = r
        moved = i + r > right
        if moved:
            center, right = i, i + r
        if r > (rad[best] or 0):
            best = i
        # Only show the steps that teach something: real letters and any reuse.
        if i % 2 == 1 or free > 0:
            hl = {k: "window" for k in range(i - r, i + r + 1)}
            hl[i] = "active"
            if mirror is not None and mirror != i:
                hl[mirror] = "special"
            ptrs = {"i": i}
            if mirror is not None and mirror != i:
                ptrs["mirror"] = mirror
            lines = []
            if mirror is not None:
                lines.append(f"i < right = {old_right}: start at min(rad[{mirror}] = {rad[mirror]}, "
                             f"{old_right} − {i} = {old_right - i}) = {free}")
            else:
                lines.append(f"i = {i} ≥ right = {old_right}: start from radius 0")
            lines.append((f"expanded to rad[{i}] = {r}" + (f"  → \"{s[(i - r) // 2:(i - r) // 2 + r]}\"" if r else ""),
                          "match" if r else None))
            say = (f"At i = {i} I can start with radius {free} for free, copied from its mirror {mirror}. "
                   f"Then I expand only past what's known, reaching {r}."
                   if free > 0 else
                   f"At i = {i}, '{t[i]}' has no free radius, so I expand outward while both sides match. "
                   f"Radius {r}" + (f" gives \"{s[(i - r) // 2:(i - r) // 2 + r]}\"." if r else "."))
            if moved and r > 1:
                say += " The right edge moves forward."
            yield Frame(
                Array(t, label="t (virtual)", hl=hl, pointers=ptrs),
                Array(rad, label="rad", show_index=False, hl={i: "active"}),
                Text(*lines),
                say=say,
            )

    start = (best - rad[best]) // 2
    ans = s[start:start + rad[best]]
    assert ans == "bab"
    hl = {k: "match" for k in range(best - rad[best], best + rad[best] + 1)}
    yield Frame(
        Array(t, label="t (virtual)", hl=hl, pointers={"best": best}),
        Array(rad, label="rad", show_index=False, hl={best: "match"}),
        Text((f'answer: s[{start}:{start + rad[best]}] = "{ans}"', "match"),
             "time O(n) · space O(n)"),
        say="The largest radius is 3 at i = 3, so the answer is \"bab\". \"aba\" ties, but I only replace on a "
            "strictly larger radius, so the leftmost wins. The right edge only moves forward: O(n).",
    )
