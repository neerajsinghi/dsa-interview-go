from walkthrough import Array, Frame, Text

TITLE = "Valid Palindrome II"
PATTERN = "Two pointers"


def is_pal_range(s, l, r):
    while l < r:
        if s[l] != s[r]:
            return False
        l, r = l + 1, r - 1
    return True


def frames():
    # "aba" (the first doc example) never mismatches; "abca" shows the skip.
    s = "abca"
    chars = list(s)

    yield Frame(
        Array(chars, label="s"),
        Text('s = "abca"'),
        speaker="Interviewer",
        say="Return true if s can be a palindrome after deleting at most one character. Here s is \"abca\".",
    )
    yield Frame(
        Array(chars, label="s"),
        Text("brute force: delete each index, re-check → n checks × O(n) = O(n²)"),
        say="I could try deleting every index and re-check the whole string each time. "
            "That's O(n²). Most of those deletions are pointless, because only a mismatch tells me what to delete.",
    )

    l, r = 0, len(s) - 1
    answer = True
    while l < r:
        if s[l] != s[r]:
            yield Frame(
                Array(chars, label="s", hl={l: "reject", r: "reject"}, pointers={"l": l, "r": r}),
                Text(f"s[{l}] = '{s[l]}' ≠ s[{r}] = '{s[r]}'", ("one of these two must be deleted", "reject")),
                say=f"'{s[l]}' and '{s[r]}' differ. Everything outside already matches, so the deleted "
                    "character has to be one of these two. I get exactly one branch point.",
            )
            left_ok = is_pal_range(s, l + 1, r)
            hl = {i: ("match" if left_ok else "reject") for i in range(l + 1, r + 1)}
            hl[l] = "visited"
            yield Frame(
                Array(chars, label="s", hl=hl, pointers={"l+1": l + 1, "r": r}),
                Text(f"skip s[{l}]: is s[{l + 1}..{r}] = \"{s[l + 1:r + 1]}\" a palindrome?",
                     (f"→ {str(left_ok).lower()}", "match" if left_ok else "reject")),
                say=f"First I skip the left one and check s[{l + 1}..{r}]. "
                    + ("That is a single character, so it's a palindrome and I can stop."
                       if left_ok else "It isn't a palindrome, so I try the other side."),
            )
            if left_ok:
                answer = True
            else:
                right_ok = is_pal_range(s, l, r - 1)
                hl = {i: ("match" if right_ok else "reject") for i in range(l, r)}
                hl[r] = "visited"
                yield Frame(
                    Array(chars, label="s", hl=hl, pointers={"l": l, "r-1": r - 1}),
                    Text(f"skip s[{r}]: is s[{l}..{r - 1}] a palindrome? {right_ok}"),
                    say=f"Now I skip the right one and check s[{l}..{r - 1}].",
                )
                answer = right_ok
            break
        yield Frame(
            Array(chars, label="s", hl={l: "match", r: "match"}, pointers={"l": l, "r": r}),
            Text(f"s[{l}] = '{s[l]}' = s[{r}] = '{s[r]}'", ("match → move both inward", "match")),
            say=f"Both ends are '{s[l]}', so they can stay. I move l right and r left.",
        )
        l, r = l + 1, r - 1

    assert answer is True  # doc: "abca" -> true
    yield Frame(
        Array(chars, label="s", hl={0: "match", 2: "match", 3: "match", 1: "visited"}),
        Text(("answer: true  (delete 'b' → \"aca\")", "match"), "time O(n) · space O(1)"),
        say="So the answer is true. The pointers walk inward once, and the one extra check after the mismatch "
            "is also linear. That's O(n) time and O(1) space.",
    )
